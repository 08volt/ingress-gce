package operator

import (
	"fmt"

	"k8s.io/ingress-gce/pkg/annotations"
	backendconfigv1 "k8s.io/ingress-gce/pkg/apis/backendconfig/v1"
	"k8s.io/ingress-gce/pkg/utils"
	"k8s.io/klog/v2"

	api_v1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/networking/v1"
)

// Services returns the wrapper
func Services(s []*api_v1.Service, logger klog.Logger) *ServicesOperator {
	return &ServicesOperator{s: s, logger: logger.WithName("ServicesOperator")}
}

// ServicesOperator is an operator wrapper for a list of Services.
type ServicesOperator struct {
	s []*api_v1.Service

	logger klog.Logger
}

// AsList returns the underlying list of Services
func (op *ServicesOperator) AsList() []*api_v1.Service {
	if op.s == nil {
		return []*api_v1.Service{}
	}
	return op.s
}

// ReferencesBackendConfig returns the Services that reference the given BackendConfig.
func (op *ServicesOperator) ReferencesBackendConfig(beConfig *backendconfigv1.BackendConfig) *ServicesOperator {
	dupes := map[string]bool{}

	var s []*api_v1.Service
	for _, svc := range op.s {
		key := fmt.Sprintf("%s/%s", svc.Namespace, svc.Name)
		if doesServiceReferenceBackendConfig(svc, beConfig, op.logger) && !dupes[key] {
			s = append(s, svc)
			dupes[key] = true
		}
	}
	return Services(s, op.logger)
}

// ReferencesL4LoggingConfigMap returns the Services that reference the given ConfigMap.
func (op *ServicesOperator) ReferencesL4LoggingConfigMap(configMap *api_v1.ConfigMap) *ServicesOperator {
	var s []*api_v1.Service
	if configMap == nil {
		return Services(s, op.logger)
	}

	dupes := map[string]bool{}
	for _, svc := range op.s {
		key := fmt.Sprintf("%s/%s", svc.Namespace, svc.Name)
		loggingConfigMapName, ok := annotations.FromService(svc).GetL4LoggingConfigMapAnnotation()
		if !dupes[key] && ok && loggingConfigMapName == configMap.Name && svc.Namespace == configMap.Namespace {
			s = append(s, svc)
			dupes[key] = true
		}
	}

	return Services(s, op.logger)
}

// ReferencedByIngress returns the Services that are referenced by the passed in Ingress.
func (op *ServicesOperator) ReferencedByIngress(ing *v1.Ingress) *ServicesOperator {
	dupes := map[string]bool{}

	var s []*api_v1.Service
	for _, svc := range op.s {
		key := fmt.Sprintf("%s/%s", svc.Namespace, svc.Name)
		if doesIngressReferenceService(ing, svc) && !dupes[key] {
			s = append(s, svc)
			dupes[key] = true
		}
	}
	return Services(s, op.logger)
}

// doesIngressReferenceService returns true if the passed in Ingress directly references
// the passed in Service.
func doesIngressReferenceService(ing *v1.Ingress, svc *api_v1.Service) bool {
	if ing.Namespace != svc.Namespace {
		return false
	}

	doesReference := false
	utils.TraverseIngressBackends(ing, func(id utils.ServicePortID) bool {
		if id.Service.Name == svc.Name {
			doesReference = true
			return true
		}
		return false
	})
	return doesReference
}
