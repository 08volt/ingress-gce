package readonly

import (
	"context"

	compute "google.golang.org/api/compute/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	corev1 "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/cloud-provider-gcp/providers/gce"
	"k8s.io/klog/v2"
)

// ReadOnlyKubeClient is a wrapper around the real Kubernetes client
// that blocks mutating API calls.
type ReadOnlyKubeClient struct {
	kubernetes.Interface
	logger klog.Logger
}

// NewReadOnlyKubeClient creates a new read-only client.
func NewReadOnlyKubeClient(realClient kubernetes.Interface, logger klog.Logger) kubernetes.Interface {
	return &ReadOnlyKubeClient{
		Interface: realClient,
		logger:    logger.WithName("ReadOnlyKubeClient"),
	}
}

// CoreV1 intercepts calls to the CoreV1 interface to return our read-only version.
func (c *ReadOnlyKubeClient) CoreV1() corev1.CoreV1Interface {
	return &readOnlyCoreV1{
		CoreV1Interface: c.Interface.CoreV1(),
		logger:          c.logger,
	}
}

// --- readOnlyCoreV1 and its methods ---

type readOnlyCoreV1 struct {
	corev1.CoreV1Interface
	logger klog.Logger
}

func (c *readOnlyCoreV1) Services(namespace string) corev1.ServiceInterface {
	return &readOnlyServices{
		ServiceInterface: c.CoreV1Interface.Services(namespace),
		logger:           c.logger,
	}
}

// --- readOnlyServices blocks updates ---

type readOnlyServices struct {
	corev1.ServiceInterface
	logger klog.Logger
}

// Update is a no-op in read-only mode.
func (s *readOnlyServices) Update(ctx context.Context, service *v1.Service, opts metav1.UpdateOptions) (*v1.Service, error) {
	s.logger.Info("[READ-ONLY] Blocked Service Update", "service", klog.KObj(service))
	// Return the original service to allow controller logic to proceed.
	return service, nil
}

// UpdateStatus is a no-op in read-only mode.
func (s *readOnlyServices) UpdateStatus(ctx context.Context, service *v1.Service, opts metav1.UpdateOptions) (*v1.Service, error) {
	s.logger.Info("[READ-ONLY] Blocked Service Status Update", "service", klog.KObj(service), "status", service.Status.LoadBalancer)
	// Return the original service.
	return service, nil
}

// Patch is a no-op in read-only mode.
func (s *readOnlyServices) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Service, err error) {
	s.logger.Info("[READ-ONLY] Blocked Service Patch", "serviceName", name, "patchType", pt, "subresources", subresources)
	// We must get the original service to return it, so the caller doesn't crash on a nil pointer.
	return s.ServiceInterface.Get(ctx, name, metav1.GetOptions{})
}

// ReadOnlyGCECloud wraps the real GCE cloud provider to block mutating calls.
type ReadOnlyGCECloud struct {
	*gce.Cloud
	logger klog.Logger
}

// NewReadOnlyGCECloud creates a new read-only GCE provider.
func NewReadOnlyGCECloud(realCloud *gce.Cloud, logger klog.Logger) *ReadOnlyGCECloud {
	return &ReadOnlyGCECloud{
		Cloud:  realCloud,
		logger: logger.WithName("ReadOnlyGCECloud"),
	}
}

// Override a few example methods. A full implementation would cover all mutating functions.

func (r *ReadOnlyGCECloud) CreateForwardingRule(fr *compute.ForwardingRule) error {
	r.logger.Info("[READ-ONLY] Blocked CreateForwardingRule", "name", fr.Name, "ipAddress", fr.IPAddress)
	return nil // Return success
}

func (r *ReadOnlyGCECloud) DeleteForwardingRule(name string, region string) error {
	r.logger.Info("[READ-ONLY] Blocked DeleteForwardingRule", "name", name, "region", region)
	return nil // Return success
}

func (r *ReadOnlyGCECloud) CreateBackendService(bs *compute.BackendService) error {
	r.logger.Info("[READ-ONLY] Blocked CreateBackendService", "name", bs.Name)
	return nil
}

// For operations that return a created resource, we might need to return a dummy object.
func (r *ReadOnlyGCECloud) CreateInstanceGroup(ig *compute.InstanceGroup, zone string) error {
	r.logger.Info("[READ-ONLY] Blocked CreateInstanceGroup", "name", ig.Name, "zone", zone)
	return nil
}
