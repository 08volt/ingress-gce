package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/ingress-gce/pkg/annotations"
)

const (
	BackendServiceConditionType = "ServiceLoadBalancerBackendService"

	TCPForwardingRuleConditionType = "ServiceLoadBalancerTCPForwardingRule"
	UDPForwardingRuleConditionType = "ServiceLoadBalancerUDPForwardingRule"
	L3ForwardingRuleConditionType  = "ServiceLoadBalancerL3ForwardingRule"

	TCPIPv6ForwardingRuleConditionType = "ServiceLoadBalancerTCPIPv6ForwardingRule"
	UDPIPv6ForwardingRuleConditionType = "ServiceLoadBalancerUDPIPv6ForwardingRule"
	L3IPv6ForwardingRuleConditionType  = "ServiceLoadBalancerL3IPv6ForwardingRule"

	// ForwardingRule Conditions without protocol prefix for error reporting
	IPv6ForwardingRuleConditionType     = "ServiceLoadBalancerIPv6ForwardingRule"
	ForwardingRuleResourceConditionType = "ServiceLoadBalancerForwardingRule"

	HealthCheckConditionType             = "ServiceLoadBalancerHealthCheck"
	FirewallRuleConditionType            = "ServiceLoadBalancerFirewallRule"
	IPv6FirewallRuleConditionType        = "ServiceLoadBalancerIPv6FirewallRule"
	FirewallHealthCheckConditionType     = "ServiceLoadBalancerFirewallRuleForHealthCheck"
	FirewallHealthCheckIPv6ConditionType = "ServiceLoadBalancerFirewallRuleForHealthCheckIPv6"

	ConditionReason                         = "GCEResourceAllocated"
	ConditionResourceAllocationFailedReason = "GCEResourceAllocationFailed"
)

func resourceAnnotationKeyToConditionType(annotationKey string) string {
	switch annotationKey {
	case annotations.BackendServiceResource:
		return BackendServiceConditionType

	case annotations.TCPForwardingRuleResource:
		return TCPForwardingRuleConditionType
	case annotations.UDPForwardingRuleResource:
		return UDPForwardingRuleConditionType
	case annotations.L3ForwardingRuleResource:
		return L3ForwardingRuleConditionType

	case annotations.TCPIPv6ForwardingRuleResource:
		return TCPIPv6ForwardingRuleConditionType
	case annotations.UDPIPv6ForwardingRuleResource:
		return UDPIPv6ForwardingRuleConditionType
	case annotations.L3IPv6ForwardingRuleResource:
		return L3IPv6ForwardingRuleConditionType

	case annotations.ForwardingRuleIPv6Resource:
		return IPv6ForwardingRuleConditionType
	case annotations.ForwardingRuleResource:
		return ForwardingRuleResourceConditionType

	case annotations.HealthcheckResource:
		return HealthCheckConditionType

	case annotations.FirewallRuleResource:
		return FirewallRuleConditionType
	case annotations.FirewallRuleIPv6Resource:
		return IPv6FirewallRuleConditionType

	case annotations.FirewallForHealthcheckResource:
		return FirewallHealthCheckConditionType
	case annotations.FirewallForHealthcheckIPv6Resource:
		return FirewallHealthCheckIPv6ConditionType
	default:
		return ""
	}
}

func NewConditionResourceAllocated(conditionType string, resourceName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               conditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            resourceName,
	}
}

func NewConditionResourceAllocationFailed(conditionType string, errMsg string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               conditionType,
		Status:             metav1.ConditionFalse,
		Reason:             ConditionResourceAllocationFailedReason,
		Message:            errMsg,
	}
}

func NewConditionResourceAllocationFailedFromAnnotation(resourceAnnotationKey string, errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(
		resourceAnnotationKeyToConditionType(resourceAnnotationKey),
		errMsg,
	)
}

// NewBackendServiceCondition creates a condition for the backend service.
func NewBackendServiceCondition(bsName string) metav1.Condition {
	return NewConditionResourceAllocated(BackendServiceConditionType, bsName)
}

// NewBackendServiceFailedCondition creates a failed condition for the backend service.
func NewBackendServiceFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(BackendServiceConditionType, errMsg)
}

// NewTCPForwardingRuleCondition creates a condition for the TCP forwarding rule.
func NewTCPForwardingRuleCondition(frName string) metav1.Condition {
	return NewConditionResourceAllocated(TCPForwardingRuleConditionType, frName)
}

// NewTCPForwardingRuleFailedCondition creates a failed condition for the TCP forwarding rule.
func NewTCPForwardingRuleFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(TCPForwardingRuleConditionType, errMsg)
}

// NewUDPForwardingRuleCondition creates a condition for the UDP forwarding rule.
func NewUDPForwardingRuleCondition(frName string) metav1.Condition {
	return NewConditionResourceAllocated(UDPForwardingRuleConditionType, frName)
}

// NewUDPForwardingRuleFailedCondition creates a failed condition for the UDP forwarding rule.
func NewUDPForwardingRuleFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(UDPForwardingRuleConditionType, errMsg)
}

// NewL3ForwardingRuleCondition creates a condition for the L3 forwarding rule.
func NewL3ForwardingRuleCondition(frName string) metav1.Condition {
	return NewConditionResourceAllocated(L3ForwardingRuleConditionType, frName)
}

// NewL3ForwardingRuleFailedCondition creates a failed condition for the L3 forwarding rule.
func NewL3ForwardingRuleFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(L3ForwardingRuleConditionType, errMsg)
}

// NewTCPIPv6ForwardingRuleCondition creates a condition for the TCP forwarding rule.
func NewTCPIPv6ForwardingRuleCondition(frName string) metav1.Condition {
	return NewConditionResourceAllocated(TCPIPv6ForwardingRuleConditionType, frName)
}

// NewTCPIPv6ForwardingRuleFailedCondition creates a failed condition for the TCP forwarding rule.
func NewTCPIPv6ForwardingRuleFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(TCPIPv6ForwardingRuleConditionType, errMsg)
}

// NewUDPIPv6ForwardingRuleCondition creates a condition for the UDP forwarding rule.
func NewUDPIPv6ForwardingRuleCondition(frName string) metav1.Condition {
	return NewConditionResourceAllocated(UDPIPv6ForwardingRuleConditionType, frName)
}

// NewUDPIPv6ForwardingRuleFailedCondition creates a failed condition for the UDP forwarding rule.
func NewUDPIPv6ForwardingRuleFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(UDPIPv6ForwardingRuleConditionType, errMsg)
}

// NewL3IPv6ForwardingRuleCondition creates a condition for the L3 forwarding rule.
func NewL3IPv6ForwardingRuleCondition(frName string) metav1.Condition {
	return NewConditionResourceAllocated(L3IPv6ForwardingRuleConditionType, frName)
}

// NewL3IPv6ForwardingRuleFailedCondition creates a failed condition for the L3 forwarding rule.
func NewL3IPv6ForwardingRuleFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(L3IPv6ForwardingRuleConditionType, errMsg)
}

// NewHealthCheckCondition creates a condition for the health check.
func NewHealthCheckCondition(hcName string) metav1.Condition {
	return NewConditionResourceAllocated(HealthCheckConditionType, hcName)
}

// NewHealthCheckFailedCondition creates a failed condition for the health check.
func NewHealthCheckFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(HealthCheckConditionType, errMsg)
}

// NewFirewallCondition creates a condition for the firewall.
func NewFirewallCondition(fwName string) metav1.Condition {
	return NewConditionResourceAllocated(FirewallRuleConditionType, fwName)
}

// NewFirewallFailedCondition creates a failed condition for the firewall.
func NewFirewallFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(FirewallRuleConditionType, errMsg)
}

// NewIPv6FirewallCondition creates a condition for the IPv6 firewall.
func NewIPv6FirewallCondition(fwName string) metav1.Condition {
	return NewConditionResourceAllocated(IPv6FirewallRuleConditionType, fwName)
}

// NewIPv6FirewallFailedCondition creates a failed condition for the IPv6 firewall.
func NewIPv6FirewallFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(IPv6FirewallRuleConditionType, errMsg)
}

// NewFirewallHealthCheckCondition creates a condition for the firewall health check.
func NewFirewallHealthCheckCondition(fwName string) metav1.Condition {
	return NewConditionResourceAllocated(FirewallHealthCheckConditionType, fwName)
}

// NewFirewallHealthCheckFailedCondition creates a failed condition for the firewall health check.
func NewFirewallHealthCheckFailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(FirewallHealthCheckConditionType, errMsg)
}

// NewFirewallHealthCheckIPv6Condition creates a condition for the firewall health check IPv6.
func NewFirewallHealthCheckIPv6Condition(fwName string) metav1.Condition {
	return NewConditionResourceAllocated(FirewallHealthCheckIPv6ConditionType, fwName)
}

// NewFirewallHealthCheckIPv6FailedCondition creates a failed condition for the firewall health check IPv6.
func NewFirewallHealthCheckIPv6FailedCondition(errMsg string) metav1.Condition {
	return NewConditionResourceAllocationFailed(FirewallHealthCheckIPv6ConditionType, errMsg)
}
