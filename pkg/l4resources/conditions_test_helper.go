package l4resources

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
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
func assertConditionsMatch(t *testing.T, want []conditionMatcher, got []metav1.Condition) {
	t.Helper()

	var wantConditions []metav1.Condition
	for _, c := range want {
		wantConditions = append(wantConditions, metav1.Condition{
			Type:    c.Type,
			Status:  c.Status,
			Reason:  c.Reason,
			Message: c.Message,
		})
	}

	sortFn := func(a, b metav1.Condition) bool {
		return a.Type < b.Type
	}

	diff := cmp.Diff(wantConditions, got,
		cmpopts.SortSlices(sortFn),
		cmpopts.IgnoreFields(metav1.Condition{}, "LastTransitionTime", "ObservedGeneration"),
	)

	if diff != "" {
		t.Errorf("Conditions mismatch (-want +got):\n%s", diff)
	}
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
