package loadbalancers

import (
	"testing"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/ingress-gce/pkg/utils"
)

// conditionMatcher represents the expected values for a condition
type conditionMatcher struct {
	Type    string
	Status  metav1.ConditionStatus
	Reason  string
	Message string // if empty, will not be checked
}

// assertConditionsMatch verifies that the actual conditions match the expected conditions
func assertConditionsMatch(t *testing.T, expected []conditionMatcher, actual []metav1.Condition) {
	t.Helper()

	// Filter out removed conditions (Status=False) from actual conditions
	var filteredActual []metav1.Condition
	for _, cond := range actual {
		if cond.Status != metav1.ConditionFalse {
			filteredActual = append(filteredActual, cond)
		}
	}

	if len(expected) != len(filteredActual) {
		t.Errorf("Expected %d conditions, got %d.\nExpected: %+v\nActual: %+v",
			len(expected), len(filteredActual), expected, filteredActual)
		return
	}

	// Build a map of actual conditions by type for easier lookup
	actualMap := make(map[string]metav1.Condition)
	for _, cond := range filteredActual {
		actualMap[cond.Type] = cond
	}

	// Check each expected condition
	for _, exp := range expected {
		act, found := actualMap[exp.Type]
		if !found {
			t.Errorf("Expected condition type %q not found in actual conditions.\nActual conditions: %+v",
				exp.Type, actual)
			continue
		}

		if act.Status != exp.Status {
			t.Errorf("Condition %q: expected status %q, got %q", exp.Type, exp.Status, act.Status)
		}

		if act.Reason != exp.Reason {
			t.Errorf("Condition %q: expected reason %q, got %q", exp.Type, exp.Reason, act.Reason)
		}

		// Only check message if it's specified in the matcher
		if exp.Message != "" && act.Message != exp.Message {
			t.Errorf("Condition %q: expected message %q, got %q", exp.Type, exp.Message, act.Message)
		}
	}
}

// assertConditionExists verifies that a specific condition exists with the expected values
func assertConditionExists(t *testing.T, conditions []metav1.Condition, condType string, expectedStatus metav1.ConditionStatus, expectedReason string) {
	t.Helper()

	for _, cond := range conditions {
		if cond.Type == condType {
			if cond.Status != expectedStatus {
				t.Errorf("Condition %q: expected status %q, got %q", condType, expectedStatus, cond.Status)
			}
			if cond.Reason != expectedReason {
				t.Errorf("Condition %q: expected reason %q, got %q", condType, expectedReason, cond.Reason)
			}
			return
		}
	}

	t.Errorf("Condition type %q not found in conditions: %+v", condType, conditions)
}

// assertConditionNotExists verifies that a specific condition type does not exist
func assertConditionNotExists(t *testing.T, conditions []metav1.Condition, condType string) {
	t.Helper()

	for _, cond := range conditions {
		if cond.Type == condType {
			t.Errorf("Condition type %q should not exist but was found: %+v", condType, cond)
			return
		}
	}
}

// assertConditionMessage verifies that a condition exists and has a specific message
func assertConditionMessage(t *testing.T, conditions []metav1.Condition, condType string, expectedMessage string) {
	t.Helper()

	for _, cond := range conditions {
		if cond.Type == condType {
			if cond.Message != expectedMessage {
				t.Errorf("Condition %q: expected message %q, got %q", condType, expectedMessage, cond.Message)
			}
			return
		}
	}

	t.Errorf("Condition type %q not found in conditions: %+v", condType, conditions)
}

// Helper to build expected  IPv4 conditions
func expectedL4LBIPv4Conditions(protocol string, bsName, frName, hcName, fwName, hcFwName string) []conditionMatcher {
	conditions := []conditionMatcher{
		{
			Type:    utils.BackendServiceConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: bsName,
		},
		{
			Type:    utils.HealthCheckConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: hcName,
		},
		{
			Type:    utils.FirewallRuleConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: fwName,
		},
		{
			Type:    utils.FirewallHealthCheckConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: hcFwName,
		},
	}

	// Add protocol-specific forwarding rule condition
	switch protocol {
	case "TCP":
		conditions = append(conditions, conditionMatcher{
			Type:    utils.TCPForwardingRuleConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: frName,
		})
	case "UDP":
		conditions = append(conditions, conditionMatcher{
			Type:    utils.UDPForwardingRuleConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: frName,
		})
	case "L3":
		conditions = append(conditions, conditionMatcher{
			Type:    utils.L3ForwardingRuleConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: frName,
		})
	}

	return conditions
}

// Helper to build expected  dual-stack conditions
func expectedL4LBDualStackConditions(protocol string, bsName, ipv4FrName, ipv6FrName, hcName, ipv4FwName, ipv6FwName, hcFwName, hcFwIPv6Name string) []conditionMatcher {
	conditions := []conditionMatcher{
		{
			Type:    utils.BackendServiceConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: bsName,
		},
		{
			Type:    utils.HealthCheckConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: hcName,
		},
		{
			Type:    utils.FirewallRuleConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: ipv4FwName,
		},
		{
			Type:    utils.IPv6FirewallRuleConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: ipv6FwName,
		},
		{
			Type:    utils.FirewallHealthCheckConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: hcFwName,
		},
		{
			Type:    utils.FirewallHealthCheckIPv6ConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: hcFwIPv6Name,
		},
	}

	// Add protocol-specific forwarding rule conditions for IPv4 and IPv6
	switch protocol {
	case "TCP":
		conditions = append(conditions,
			conditionMatcher{
				Type:    utils.TCPForwardingRuleConditionType,
				Status:  metav1.ConditionTrue,
				Reason:  utils.ConditionReason,
				Message: ipv4FrName,
			},
			conditionMatcher{
				Type:    utils.TCPIPv6ForwardingRuleConditionType,
				Status:  metav1.ConditionTrue,
				Reason:  utils.ConditionReason,
				Message: ipv6FrName,
			},
		)
	case "UDP":
		conditions = append(conditions,
			conditionMatcher{
				Type:    utils.UDPForwardingRuleConditionType,
				Status:  metav1.ConditionTrue,
				Reason:  utils.ConditionReason,
				Message: ipv4FrName,
			},
			conditionMatcher{
				Type:    utils.UDPIPv6ForwardingRuleConditionType,
				Status:  metav1.ConditionTrue,
				Reason:  utils.ConditionReason,
				Message: ipv6FrName,
			},
		)
	case "L3":
		conditions = append(conditions,
			conditionMatcher{
				Type:    utils.L3ForwardingRuleConditionType,
				Status:  metav1.ConditionTrue,
				Reason:  utils.ConditionReason,
				Message: ipv4FrName,
			},
			conditionMatcher{
				Type:    utils.L3IPv6ForwardingRuleConditionType,
				Status:  metav1.ConditionTrue,
				Reason:  utils.ConditionReason,
				Message: ipv6FrName,
			},
		)
	}

	return conditions
}

// Helper to build expected L4 LB IPv6 Only conditions
func expectedL4LBIPv6Conditions(protocol string, bsName, ipv6FrName, hcName, ipv6FwName, hcFwIPv6Name string) []conditionMatcher {
	conditions := []conditionMatcher{
		{
			Type:    utils.BackendServiceConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: bsName,
		},
		{
			Type:    utils.HealthCheckConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: hcName,
		},
		{
			Type:    utils.IPv6FirewallRuleConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: ipv6FwName,
		},
		{
			Type:    utils.FirewallHealthCheckIPv6ConditionType,
			Status:  metav1.ConditionTrue,
			Reason:  utils.ConditionReason,
			Message: hcFwIPv6Name,
		},
	}

	// Add protocol-specific forwarding rule conditions for IPv4 and IPv6
	switch protocol {
	case "TCP":
		conditions = append(conditions,
			conditionMatcher{
				Type:    utils.TCPIPv6ForwardingRuleConditionType,
				Status:  metav1.ConditionTrue,
				Reason:  utils.ConditionReason,
				Message: ipv6FrName,
			},
		)
	case "UDP":
		conditions = append(conditions,
			conditionMatcher{
				Type:    utils.UDPIPv6ForwardingRuleConditionType,
				Status:  metav1.ConditionTrue,
				Reason:  utils.ConditionReason,
				Message: ipv6FrName,
			},
		)
	case "L3":
		conditions = append(conditions,
			conditionMatcher{
				Type:    utils.L3IPv6ForwardingRuleConditionType,
				Status:  metav1.ConditionTrue,
				Reason:  utils.ConditionReason,
				Message: ipv6FrName,
			},
		)
	}

	return conditions
}
