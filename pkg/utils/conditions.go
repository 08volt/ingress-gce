package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	BackendServiceConditionType          = "ServiceLoadBalancerBackendService"
	TCPForwardingRuleConditionType       = "ServiceLoadBalancerTCPForwardingRule"
	UDPForwardingRuleConditionType       = "ServiceLoadBalancerUDPForwardingRule"
	L3ForwardingRuleConditionType        = "ServiceLoadBalancerL3ForwardingRule"
	TCPIPv6ForwardingRuleConditionType   = "ServiceLoadBalancerTCPIPv6ForwardingRule"
	UDPIPv6ForwardingRuleConditionType   = "ServiceLoadBalancerUDPIPv6ForwardingRule"
	L3IPv6ForwardingRuleConditionType    = "ServiceLoadBalancerL3IPv6ForwardingRule"
	HealthCheckConditionType             = "ServiceLoadBalancerHealthCheck"
	FirewallConditionType                = "ServiceLoadBalancerFirewall"
	IPv6FirewallConditionType            = "ServiceLoadBalancerIPv6Firewall"
	FirewallHealthCheckConditionType     = "ServiceLoadBalancerFirewallHealthCheck"
	FirewallHealthCheckIPv6ConditionType = "ServiceLoadBalancerFirewallHealthCheckIPv6"

	ConditionReason        = "GCEResourceAllocated"
	ConditionReasonRemoved = "GCEResourceRemoved"

	MessageResourceRemoved = "GCE resource has been removed"
)

func SetRemovedResourceCondition(cond *metav1.Condition) {
	cond.LastTransitionTime = metav1.Now()
	cond.Status = metav1.ConditionFalse
	cond.Reason = ConditionReasonRemoved
	cond.Message = MessageResourceRemoved
}

// NewBackendServiceCondition creates a condition for the backend service.
func NewBackendServiceCondition(bsName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               BackendServiceConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            bsName,
	}
}

// NewTCPForwardingRuleCondition creates a condition for the TCP forwarding rule.
func NewTCPForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               TCPForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            frName,
	}
}

// NewUDPForwardingRuleCondition creates a condition for the UDP forwarding rule.
func NewUDPForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               UDPForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            frName,
	}
}

// NewL3ForwardingRuleCondition creates a condition for the L3 forwarding rule.
func NewL3ForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L3ForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            frName,
	}
}

// NewTCPIPv6ForwardingRuleCondition creates a condition for the TCP forwarding rule.
func NewTCPIPv6ForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               TCPIPv6ForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            frName,
	}
}

// NewUDPIPv6ForwardingRuleCondition creates a condition for the UDP forwarding rule.
func NewUDPIPv6ForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               UDPIPv6ForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            frName,
	}
}

// NewL3IPv6ForwardingRuleCondition creates a condition for the L3 forwarding rule.
func NewL3IPv6ForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L3IPv6ForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            frName,
	}
}

// NewHealthCheckCondition creates a condition for the health check.
func NewHealthCheckCondition(hcName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               HealthCheckConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            hcName,
	}
}

// NewFirewallCondition creates a condition for the firewall.
func NewFirewallCondition(fwName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               FirewallConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            fwName,
	}
}

// NewIPv6FirewallCondition creates a condition for the IPv6 firewall.
func NewIPv6FirewallCondition(fwName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               IPv6FirewallConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            fwName,
	}
}

// NewFirewallHealthCheckCondition creates a condition for the firewall health check.
func NewFirewallHealthCheckCondition(fwName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               FirewallHealthCheckConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            fwName,
	}
}

// NewFirewallHealthCheckIPv6Condition creates a condition for the firewall health check IPv6.
func NewFirewallHealthCheckIPv6Condition(fwName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               FirewallHealthCheckIPv6ConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             ConditionReason,
		Message:            fwName,
	}
}

// NewFirewallHealthCheckIPv6ConditionRemoved creates a condition for the firewall health check IPv6.
func NewFirewallHealthCheckIPv6ConditionRemoved() metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               FirewallHealthCheckIPv6ConditionType,
		Status:             metav1.ConditionFalse,
		Reason:             ConditionReasonRemoved,
		Message:            MessageResourceRemoved,
	}
}
