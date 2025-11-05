package utils

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	L4LBBackendServiceConditionType          = "ServiceLoadBalancerBackendService"
	L4LBTCPForwardingRuleConditionType       = "ServiceLoadBalancerTCPForwardingRule"
	L4LBUDPForwardingRuleConditionType       = "ServiceLoadBalancerUDPForwardingRule"
	L4LBL3ForwardingRuleConditionType        = "ServiceLoadBalancerL3ForwardingRule"
	L4LBTCPIPv6ForwardingRuleConditionType   = "ServiceLoadBalancerTCPIPv6ForwardingRule"
	L4LBUDPIPv6ForwardingRuleConditionType   = "ServiceLoadBalancerUDPIPv6ForwardingRule"
	L4LBL3IPv6ForwardingRuleConditionType    = "ServiceLoadBalancerL3IPv6ForwardingRule"
	L4LBHealthCheckConditionType             = "ServiceLoadBalancerHealthCheck"
	L4LBFirewallConditionType                = "ServiceLoadBalancerFirewall"
	L4LBIPv6FirewallConditionType            = "ServiceLoadBalancerIPv6Firewall"
	L4LBFirewallHealthCheckConditionType     = "ServiceLoadBalancerFirewallHealthCheck"
	L4LBFirewallHealthCheckIPv6ConditionType = "ServiceLoadBalancerFirewallHealthCheckIPv6"
	L4LBNetworkEndpointGroupConditionType    = "ServiceLoadBalancerNetworkEndpointGroup"

	L4LBConditionReason = "GCEResourceAllocated"
)

// NewBackendServiceCondition creates a condition for the backend service.
func NewBackendServiceCondition(bsName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBBackendServiceConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            bsName,
	}
}

// NewTCPForwardingRuleCondition creates a condition for the TCP forwarding rule.
func NewTCPForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBTCPForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            frName,
	}
}

// NewUDPForwardingRuleCondition creates a condition for the UDP forwarding rule.
func NewUDPForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBUDPForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            frName,
	}
}

// NewL3ForwardingRuleCondition creates a condition for the L3 forwarding rule.
func NewL3ForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBL3ForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            frName,
	}
}

// NewTCPIPv6ForwardingRuleCondition creates a condition for the TCP forwarding rule.
func NewTCPIPv6ForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBTCPIPv6ForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            frName,
	}
}

// NewUDPIPv6ForwardingRuleCondition creates a condition for the UDP forwarding rule.
func NewUDPIPv6ForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBUDPIPv6ForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            frName,
	}
}

// NewL3IPv6ForwardingRuleCondition creates a condition for the L3 forwarding rule.
func NewL3IPv6ForwardingRuleCondition(frName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBL3IPv6ForwardingRuleConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            frName,
	}
}

// NewHealthCheckCondition creates a condition for the health check.
func NewHealthCheckCondition(hcName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBHealthCheckConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            hcName,
	}
}

// NewFirewallCondition creates a condition for the firewall.
func NewFirewallCondition(fwName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBFirewallConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            fwName,
	}
}

// NewIPv6FirewallCondition creates a condition for the IPv6 firewall.
func NewIPv6FirewallCondition(fwName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBIPv6FirewallConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            fwName,
	}
}

// NewFirewallHealthCheckCondition creates a condition for the firewall health check.
func NewFirewallHealthCheckCondition(fwName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBFirewallHealthCheckConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            fwName,
	}
}

// NewFirewallHealthCheckIPv6Condition creates a condition for the firewall health check IPv6.
func NewFirewallHealthCheckIPv6Condition(fwName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBFirewallHealthCheckIPv6ConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            fwName,
	}
}

// NewNetworkEndpointGroupCondition creates a condition for the network endpoint group.
func NewNetworkEndpointGroupCondition(negName string) metav1.Condition {
	return metav1.Condition{
		LastTransitionTime: metav1.Now(),
		Type:               L4LBNetworkEndpointGroupConditionType,
		Status:             metav1.ConditionTrue,
		Reason:             L4LBConditionReason,
		Message:            negName,
	}
}
