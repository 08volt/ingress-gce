package l4lb

import (
	"errors"
	"testing"
	"time"

	svcLBStatus "github.com/GoogleCloudPlatform/gke-networking-api/apis/serviceloadbalancerstatus/v1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/ingress-gce/pkg/context"
	"k8s.io/ingress-gce/pkg/utils/common"
	"k8s.io/klog/v2"

	fakesvclbstatusclientset "github.com/GoogleCloudPlatform/gke-networking-api/client/serviceloadbalancerstatus/clientset/versioned/fake"

	"github.com/google/go-cmp/cmp"
	"k8s.io/apimachinery/pkg/runtime"
	k8stesting "k8s.io/client-go/testing"
)

func TestFinalizerWasRemovedUnexpectedly(t *testing.T) {
	testCases := []struct {
		desc           string
		oldService     *v1.Service
		newService     *v1.Service
		finalizerName  string
		expectedResult bool
	}{
		{
			desc:           "Clean service",
			oldService:     &v1.Service{},
			newService:     &v1.Service{},
			finalizerName:  "random",
			expectedResult: false,
		},
		{
			desc: "Empty finalizers",
			oldService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{},
				},
			},
			newService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{},
				},
			},
			finalizerName:  "random",
			expectedResult: false,
		},
		{
			desc: "Changed L4 Finalizer",
			oldService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{common.LegacyILBFinalizer, "random"},
				},
			},
			newService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{"random", "gke.networking.io/l4-ilb-v1-fake"},
				},
			},
			finalizerName:  common.LegacyILBFinalizer,
			expectedResult: true,
		},
		{
			desc: "Removed L4 Finalizer",
			oldService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{common.LegacyILBFinalizer, "random"},
				},
			},
			newService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{"random"},
				},
			},
			finalizerName:  common.LegacyILBFinalizer,
			expectedResult: true,
		},
		{
			desc: "Added L4 ILB v2 Finalizer",
			oldService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{"random"},
				},
			},
			newService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{"random", common.ILBFinalizerV2},
				},
			},
			finalizerName:  common.ILBFinalizerV2,
			expectedResult: false,
		},
		{
			desc: "Service with NetLB Finalizer hasn't changed",
			oldService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{common.NetLBFinalizerV2, "random"},
				},
			},
			newService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{"random", common.NetLBFinalizerV2},
				},
			},
			finalizerName:  common.NetLBFinalizerV2,
			expectedResult: false,
		},
		{
			desc: "Finalizer was removed but given name is wrong",
			oldService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{common.NetLBFinalizerV2, "random"},
				},
			},
			newService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{"random"},
				},
			},
			finalizerName:  common.ILBFinalizerV2,
			expectedResult: false,
		},
		{
			desc: "Finalizer was removed and service to be deleted",
			oldService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					Finalizers: []string{common.NetLBFinalizerV2, "random"},
				},
			},
			newService: &v1.Service{
				ObjectMeta: metav1.ObjectMeta{
					DeletionTimestamp: &metav1.Time{Time: time.Date(2024, 12, 30, 0, 0, 0, 0, time.Local)},
					Finalizers:        []string{common.ILBFinalizerV2},
				},
			},
			finalizerName:  common.ILBFinalizerV2,
			expectedResult: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			gotResult := finalizerWasRemovedUnexpectedly(tc.oldService, tc.newService, tc.finalizerName)
			if gotResult != tc.expectedResult {
				t.Errorf("finalizerWasRemoved(oldSvc=%v, newSvc=%v, finalizer=%s) returned %v, but expected %v", tc.oldService, tc.newService, tc.finalizerName, gotResult, tc.expectedResult)
			}
		})
	}
}

func TestServiceLoadBalancerStatusResourcesEqual(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		desc string
		a    svcLBStatus.ServiceLoadBalancerStatusStatus
		b    svcLBStatus.ServiceLoadBalancerStatusStatus
		want bool
	}{
		{
			desc: "Identical resources",
			a:    svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: []string{"res-a", "res-b"}},
			b:    svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: []string{"res-a", "res-b"}},
			want: true,
		},
		{
			desc: "Equal resources, different order",
			a:    svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: []string{"res-b", "res-a"}},
			b:    svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: []string{"res-a", "res-b"}},
			want: true,
		},
		{
			desc: "Different resource list length",
			a:    svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: []string{"res-a"}},
			b:    svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: []string{"res-a", "res-b"}},
			want: false,
		},
		{
			desc: "Both lists are nil/empty",
			a:    svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: nil},
			b:    svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: []string{}},
			want: true,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			got := serviceLoadBalancerStatusResourcesEqual(tc.a, tc.b)
			if got != tc.want {
				t.Errorf("serviceLoadBalancerStatusResourcesEqual(%+v, %+v) = %v, want %v", tc.a, tc.b, got, tc.want)
			}
		})
	}
}

func TestGenerateServiceLoadBalancerStatus(t *testing.T) {
	t.Parallel()
	testService := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "test-service", Namespace: "test-ns", UID: "test-uid"},
	}
	wantCR := &svcLBStatus.ServiceLoadBalancerStatus{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-service-status",
			Namespace: "test-ns",
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(testService, v1.SchemeGroupVersion.WithKind("Service")),
			},
		},
		Spec: svcLBStatus.ServiceLoadBalancerStatusSpec{},
		Status: svcLBStatus.ServiceLoadBalancerStatusStatus{
			GceResources: []string{"res-1", "res-2"},
		},
	}

	gotCR := GenerateServiceLoadBalancerStatus(testService, []string{"res-1", "res-2"})
	if diff := cmp.Diff(wantCR, gotCR); diff != "" {
		t.Errorf("GenerateServiceLoadBalancerStatus() returned diff (-want +got):\n%s", diff)
	}
}

func TestEnsureServiceLoadBalancerStatusCR(t *testing.T) {
	t.Parallel()
	testSvc := &v1.Service{
		ObjectMeta: metav1.ObjectMeta{Name: "test-svc", Namespace: "default", UID: "uid-1"},
	}
	crName := testSvc.Name + "-status"

	testCases := []struct {
		desc              string
		existingObjects   []runtime.Object
		gceResourceURLs   []string
		expectAction      string // "create", "update"
		expectSubresource string // "status" or ""
		expectError       bool
	}{
		{
			desc:            "CR does not exist, should create it",
			existingObjects: nil,
			gceResourceURLs: []string{"res-a", "res-b"},
			expectAction:    "create",
		},
		{
			desc: "CR exists and is up-to-date, should do nothing",
			existingObjects: []runtime.Object{
				&svcLBStatus.ServiceLoadBalancerStatus{
					ObjectMeta: metav1.ObjectMeta{Name: crName, Namespace: testSvc.Namespace},
					Spec:       svcLBStatus.ServiceLoadBalancerStatusSpec{},
					Status:     svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: []string{"res-b", "res-a"}},
				},
			},
			gceResourceURLs: []string{"res-a", "res-b"},
			expectAction:    "", // No action
		},
		{
			desc: "CR exists, status is outdated, should update status subresource",
			existingObjects: []runtime.Object{
				&svcLBStatus.ServiceLoadBalancerStatus{
					ObjectMeta: metav1.ObjectMeta{Name: crName, Namespace: testSvc.Namespace, ResourceVersion: "1"},
					Spec:       svcLBStatus.ServiceLoadBalancerStatusSpec{},
					Status:     svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: []string{"old-res"}},
				},
			},
			gceResourceURLs:   []string{"new-res-1", "new-res-2"},
			expectAction:      "update",
			expectSubresource: "status",
		},
		{
			desc: "CR exists and needs to be updated to have empty resources",
			existingObjects: []runtime.Object{
				&svcLBStatus.ServiceLoadBalancerStatus{
					ObjectMeta: metav1.ObjectMeta{Name: crName, Namespace: testSvc.Namespace, ResourceVersion: "1"},
					Spec:       svcLBStatus.ServiceLoadBalancerStatusSpec{},
					Status:     svcLBStatus.ServiceLoadBalancerStatusStatus{GceResources: []string{"res-a"}},
				},
			},
			gceResourceURLs:   []string{},
			expectAction:      "update",
			expectSubresource: "status",
		},
		{
			desc:            "Error on initial get, should return error",
			existingObjects: nil,
			gceResourceURLs: []string{"res-a"},
			expectError:     true, // Injected via reactor below
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			fakeSvcLBStatusClient := fakesvclbstatusclientset.NewSimpleClientset(tc.existingObjects...)
			if tc.expectError {
				fakeSvcLBStatusClient.PrependReactor("get", "serviceloadbalancerstatuses", func(action k8stesting.Action) (handled bool, ret runtime.Object, err error) {
					return true, nil, errors.New("API server is down")
				})
			}
			testContext := &context.ControllerContext{
				KubeClient:        fake.NewSimpleClientset(),
				SvcLBStatusClient: fakeSvcLBStatusClient,
			}
			logger := klog.TODO()

			err := ensureServiceLoadBalancerStatusCR(testContext, testSvc, tc.gceResourceURLs, logger)
			if (err != nil) != tc.expectError {
				t.Fatalf("ensureServiceLoadBalancerStatusCR() error = %v, wantErr %v", err, tc.expectError)
			}

			var triggeredAction k8stesting.Action
			for _, action := range fakeSvcLBStatusClient.Actions() {
				verb := action.GetVerb()
				if verb == "create" || verb == "update" {
					triggeredAction = action
					break
				}
			}

			if tc.expectAction == "" {
				if triggeredAction != nil {
					t.Errorf("Expected no action, but got %q", triggeredAction.GetVerb())
				}
				return
			}
			if triggeredAction == nil {
				t.Fatalf("Expected action %q, but got none", tc.expectAction)
			}
			if gotVerb := triggeredAction.GetVerb(); gotVerb != tc.expectAction {
				t.Errorf("Expected verb %q, got %q", tc.expectAction, gotVerb)
			}
			updateAction, isUpdate := triggeredAction.(k8stesting.UpdateAction)
			if isUpdate {
				if gotSub := updateAction.GetSubresource(); gotSub != tc.expectSubresource {
					t.Errorf("Expected subresource %q, got %q", tc.expectSubresource, gotSub)
				}
			}
		})
	}
}
