/*
Copyright 2023 The Kubernetes Authors.

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

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
	v1alpha1 "k8s.io/cloud-provider-gcp/crd/apis/network/v1alpha1"
)

// NetworkLister helps list Networks.
// All objects returned here must be treated as read-only.
type NetworkLister interface {
	// List lists all Networks in the indexer.
	// Objects returned here must be treated as read-only.
	List(selector labels.Selector) (ret []*v1alpha1.Network, err error)
	// Get retrieves the Network from the index for a given name.
	// Objects returned here must be treated as read-only.
	Get(name string) (*v1alpha1.Network, error)
	NetworkListerExpansion
}

// networkLister implements the NetworkLister interface.
type networkLister struct {
	indexer cache.Indexer
}

// NewNetworkLister returns a new NetworkLister.
func NewNetworkLister(indexer cache.Indexer) NetworkLister {
	return &networkLister{indexer: indexer}
}

// List lists all Networks in the indexer.
func (s *networkLister) List(selector labels.Selector) (ret []*v1alpha1.Network, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.Network))
	})
	return ret, err
}

// Get retrieves the Network from the index for a given name.
func (s *networkLister) Get(name string) (*v1alpha1.Network, error) {
	obj, exists, err := s.indexer.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("network"), name)
	}
	return obj.(*v1alpha1.Network), nil
}