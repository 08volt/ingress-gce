/*
Copyright 2026 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package syncers

import (
	"fmt"
	"sync"

	nodetopologyv1 "github.com/GoogleCloudPlatform/gke-networking-api/apis/nodetopology/v1"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud"
	"github.com/GoogleCloudPlatform/k8s-cloud-provider/pkg/cloud/meta"
	"k8s.io/apimachinery/pkg/util/sets"
	negbindingv1beta1 "k8s.io/ingress-gce/pkg/apis/negbinding/v1beta1"
	"k8s.io/ingress-gce/pkg/neg/types/shared"
	"k8s.io/ingress-gce/pkg/network"
	"k8s.io/ingress-gce/pkg/utils/zonegetter"
	"k8s.io/klog/v2"
)

// negOwnershipRegistry is a pure data structure tracking exclusive ownership of
// NEG names across NEGBinding CR syncers. It performs no callbacks or goroutine
// spawning; ReleaseAllOwnedExcept returns the released names so the caller can
// re-enqueue any bindings waiting to acquire them.
type negOwnershipRegistry interface {
	Acquire(negName string, owner string) (bool, string)
	ReleaseAllOwnedExcept(owner string, keep sets.Set[string]) []string
}

// NEGBindingTopologyProvider provides subnets and zones where NEGs should be managed
// based on NetworkEndpointGroupBinding CR.
//
// Ownership is refreshed out-of-band via RefreshOwnership (called by the manager on
// its serialized worker), which caches the owned NEG refs. The List* getters are
// read-only projections of that cache and perform no ownership mutation.
type NEGBindingTopologyProvider struct {
	ownerKey        string
	defaultSubnetID *cloud.ResourceID
	registry        negOwnershipRegistry

	// mu guards ownedRefs, which is written by RefreshOwnership (manager goroutine)
	// and read by the List* getters (syncer goroutine).
	mu        sync.Mutex
	ownedRefs []negbindingv1beta1.SpecNegRef
}

// NewNEGBindingTopologyProvider constructs a new NEGBindingTopologyProvider
func NewNEGBindingTopologyProvider(namespace, negBindingName, defaultSubnetURL string, registry negOwnershipRegistry) (*NEGBindingTopologyProvider, error) {
	defaultSubnetID, err := cloud.ParseResourceURL(defaultSubnetURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse default subnetwork URL %q: %w", defaultSubnetURL, err)
	}

	return &NEGBindingTopologyProvider{
		ownerKey:        fmt.Sprintf("%s/%s", namespace, negBindingName),
		defaultSubnetID: defaultSubnetID,
		registry:        registry,
	}, nil
}

// RefreshOwnership releases the NEG names no longer desired by the binding, acquires
// the currently desired ones, and caches the successfully owned refs. It returns the
// names of NEGs that were released so the caller can re-enqueue any bindings waiting
// to acquire them.
func (p *NEGBindingTopologyProvider) RefreshOwnership(binding *negbindingv1beta1.NetworkEndpointGroupBinding, logger klog.Logger) []string {
	keepNEGs := sets.New[string]()
	for _, ref := range binding.Spec.NetworkEndpointGroups {
		keepNEGs.Insert(ref.Name)
	}

	released := p.registry.ReleaseAllOwnedExcept(p.ownerKey, keepNEGs)

	var acquiredRefs []negbindingv1beta1.SpecNegRef
	for _, ref := range binding.Spec.NetworkEndpointGroups {
		acquired, owner := p.registry.Acquire(ref.Name, p.ownerKey)
		if acquired {
			acquiredRefs = append(acquiredRefs, ref)
		} else {
			logger.Info("NEG name is owned by another binding, skipping", "negName", ref.Name, "owner", p.ownerKey, "currentOwner", owner)
		}
	}

	p.mu.Lock()
	p.ownedRefs = acquiredRefs
	p.mu.Unlock()
	return released
}

// ownedRefsSnapshot returns the currently owned NEG refs computed by the last
// RefreshOwnership call.
func (p *NEGBindingTopologyProvider) ownedRefsSnapshot() []negbindingv1beta1.SpecNegRef {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.ownedRefs
}

// ListSubnetsInDefaultNetwork returns the list of subnets of the owned NEGs declared
// inside the NegBinding CR Spec.
func (p *NEGBindingTopologyProvider) ListSubnetsInDefaultNetwork(logger klog.Logger) []nodetopologyv1.SubnetConfig {
	// Return only subnets where NEGs are owned.
	subnets := sets.New[string]()
	for _, ref := range p.ownedRefsSnapshot() {
		subnets.Insert(ref.Subnet)
	}

	configs := []nodetopologyv1.SubnetConfig{}
	for subnet := range subnets {
		key := &meta.Key{
			Name:   subnet,
			Region: p.defaultSubnetID.Key.Region,
		}
		subnetPath := cloud.SelfLink(meta.VersionGA, p.defaultSubnetID.ProjectID, p.defaultSubnetID.Resource, key)
		configs = append(configs, nodetopologyv1.SubnetConfig{
			Name:       subnet,
			SubnetPath: subnetPath,
		})
	}
	return configs
}

// ListZonesPerSubnet returns a map of subnet to zones of the owned NEGs defined inside
// the NegBinding CR Spec.
// NEGBinding contains explicit locations (subnet + zone pairs), where NEG controller is expected to
// manage NEGs ignoring if any endpoints available there. Therefore ignoring filtering.
func (p *NEGBindingTopologyProvider) ListZonesPerSubnet(_ zonegetter.Filter, networkInfo network.NetworkInfo, logger klog.Logger) (shared.ZonesPerSubnetMap, error) {
	if !networkInfo.IsDefault {
		return nil, fmt.Errorf("NEGBinding does not support multi-network mode")
	}

	// Return only zones of subnets, where NEGs are owned.
	zonesPerSubnet := make(shared.ZonesPerSubnetMap)
	for _, ref := range p.ownedRefsSnapshot() {
		zonesPerSubnet[ref.Subnet] = sets.New(ref.Zones...)
	}
	return zonesPerSubnet, nil
}
