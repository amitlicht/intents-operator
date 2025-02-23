package intents_reconcilers

import (
	"context"
	"fmt"
	otterizev1alpha3 "github.com/otterize/intents-operator/src/operator/api/v1alpha3"
	mocks "github.com/otterize/intents-operator/src/operator/controllers/intents_reconcilers/mocks"
	"github.com/otterize/intents-operator/src/shared/testbase"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"testing"
	"time"
)

const (
	testNamespace = "test-namespace"
)

type PodLabelReconcilerTestSuite struct {
	testbase.MocksSuiteBase
	Reconciler *PodLabelReconciler
}

func (s *PodLabelReconcilerTestSuite) SetupTest() {
	s.MocksSuiteBase.SetupTest()
	s.Client = mocks.NewMockClient(s.Controller)
	s.Reconciler = NewPodLabelReconciler(s.Client, nil)
	s.Reconciler.Recorder = s.Recorder
}

func (s *PodLabelReconcilerTestSuite) TearDownTest() {
	s.Reconciler = nil
	s.MocksSuiteBase.TearDownTest()
}

func (s *PodLabelReconcilerTestSuite) TestClientAccessLabelAdded() {
	clientIntentsName := "client-intents"
	serviceName := "test-client"

	namespacedName := types.NamespacedName{
		Namespace: testNamespace,
		Name:      clientIntentsName,
	}
	req := ctrl.Request{
		NamespacedName: namespacedName,
	}

	serverName := "test-server"
	intentsSpec := otterizev1alpha3.IntentsSpec{
		Service: otterizev1alpha3.Service{Name: serviceName},
		Calls: []otterizev1alpha3.Intent{
			{
				Name: serverName,
			},
		},
	}

	emptyIntents := &otterizev1alpha3.ClientIntents{}
	s.Client.EXPECT().Get(gomock.Any(), req.NamespacedName, gomock.Eq(emptyIntents)).DoAndReturn(
		func(ctx context.Context, name types.NamespacedName, intents *otterizev1alpha3.ClientIntents, options ...client.ListOption) error {
			intents.Spec = &intentsSpec
			intents.Namespace = testNamespace
			return nil
		})

	var intents otterizev1alpha3.ClientIntents
	intents.Spec = &intentsSpec

	listOption := &client.ListOptions{Namespace: testNamespace}
	labelSelector := labels.SelectorFromSet(map[string]string{
		"intents.otterize.com/server": "test-client-test-namespace-537e87",
	})

	labelMatcher := client.MatchingLabelsSelector{Selector: labelSelector}
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
			Labels:    make(map[string]string),
		},
	}

	s.Client.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Eq(listOption), gomock.Eq(labelMatcher)).DoAndReturn(
		func(ctx context.Context, pds *v1.PodList, opts ...client.ListOption) error {
			pds.Items = append(pds.Items, pod)
			return nil
		})

	updatedPod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
			Labels: map[string]string{
				"intents.otterize.com/access-test-server-test-namespace-8ddecb": "true",
				"intents.otterize.com/client":                                   "test-client-test-namespace-537e87",
			},
		},
		Spec: v1.PodSpec{},
	}

	s.Client.EXPECT().Patch(gomock.Any(), gomock.Eq(&updatedPod), gomock.Any()).Return(nil)

	res, err := s.Reconciler.Reconcile(context.Background(), req)
	s.NoError(err)
	s.Empty(res)
}

func (s *PodLabelReconcilerTestSuite) TestClientAccessLabelAddedTruncatedNameAndNamespace() {
	clientIntentsName := "client-intents"
	serviceName := "test-client-with-a-very-long-name-more-than-20-characters"
	longNamespace := "test-namespace-with-a-very-long-name-more-than-20-characters"

	namespacedName := types.NamespacedName{
		Namespace: longNamespace,
		Name:      clientIntentsName,
	}
	req := ctrl.Request{
		NamespacedName: namespacedName,
	}

	serverName := "test-server"
	intentsSpec := otterizev1alpha3.IntentsSpec{
		Service: otterizev1alpha3.Service{Name: serviceName},
		Calls: []otterizev1alpha3.Intent{
			{
				Name: serverName,
			},
		},
	}

	emptyIntents := &otterizev1alpha3.ClientIntents{}
	s.Client.EXPECT().Get(gomock.Any(), req.NamespacedName, gomock.Eq(emptyIntents)).DoAndReturn(
		func(ctx context.Context, name types.NamespacedName, intents *otterizev1alpha3.ClientIntents, options ...client.ListOption) error {
			intents.Spec = &intentsSpec
			intents.Namespace = testNamespace
			return nil
		})

	var intents otterizev1alpha3.ClientIntents
	intents.Spec = &intentsSpec

	listOption := &client.ListOptions{Namespace: longNamespace}
	labelSelector := labels.SelectorFromSet(map[string]string{
		"intents.otterize.com/server": "test-client-with-a-v-test-namespace-14e99d",
	})

	labelMatcher := client.MatchingLabelsSelector{Selector: labelSelector}
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
			Labels:    make(map[string]string),
		},
	}

	s.Client.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Eq(listOption), gomock.Eq(labelMatcher)).DoAndReturn(
		func(ctx context.Context, pds *v1.PodList, opts ...client.ListOption) error {
			pds.Items = append(pds.Items, pod)
			return nil
		})

	updatedPod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
			Labels: map[string]string{
				"intents.otterize.com/access-test-server-test-namespace-with--a1ac14": "true",
				"intents.otterize.com/client":                                         "test-client-with-a-v-test-namespace-14e99d",
			},
		},
		Spec: v1.PodSpec{},
	}

	s.Client.EXPECT().Patch(gomock.Any(), gomock.Eq(&updatedPod), gomock.Any()).Return(nil)

	res, err := s.Reconciler.Reconcile(context.Background(), req)
	s.NoError(err)
	s.Empty(res)
}

func (s *PodLabelReconcilerTestSuite) TestClientAccessLabelRemoved() {
	s.testClientAccessLabelRemovedWithParams(map[string]string{"a": "b"})
}

func (s *PodLabelReconcilerTestSuite) TestClientAccessLabelRemovedNoPodAnnotations() {
	s.testClientAccessLabelRemovedWithParams(nil)
}

func (s *PodLabelReconcilerTestSuite) testClientAccessLabelRemovedWithParams(podAnnotations map[string]string) {
	clientIntentsName := "client-intents"
	serviceName := "test-client"

	namespacedName := types.NamespacedName{
		Namespace: testNamespace,
		Name:      clientIntentsName,
	}
	req := ctrl.Request{
		NamespacedName: namespacedName,
	}

	serverName := "test-server"
	intentsSpec := otterizev1alpha3.IntentsSpec{
		Service: otterizev1alpha3.Service{Name: serviceName},
		Calls: []otterizev1alpha3.Intent{
			{
				Name: serverName,
			},
		},
	}

	emptyIntents := &otterizev1alpha3.ClientIntents{}

	var deletedIntents otterizev1alpha3.ClientIntents
	deletedIntents.Spec = &intentsSpec
	deletedIntents.Namespace = testNamespace
	deletedIntents.SetDeletionTimestamp(&metav1.Time{Time: time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)})

	s.Client.EXPECT().Get(gomock.Any(), req.NamespacedName, gomock.Eq(emptyIntents)).DoAndReturn(
		func(ctx context.Context, name types.NamespacedName, intents *otterizev1alpha3.ClientIntents, options ...client.ListOption) error {
			*intents = deletedIntents
			return nil
		})

	// Now the reconciler should handle the deletion of the client intents
	listOption := &client.ListOptions{Namespace: testNamespace}
	labelSelector := labels.SelectorFromSet(map[string]string{
		"intents.otterize.com/server": "test-client-test-namespace-537e87",
	})

	labelMatcher := client.MatchingLabelsSelector{Selector: labelSelector}
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
			Labels: map[string]string{
				"intents.otterize.com/access-test-server-test-namespace-8ddecb": "true",
				otterizev1alpha3.OtterizeClientLabelKey:                         "true",
			},
			Annotations: podAnnotations,
		},
		Spec: v1.PodSpec{},
	}

	s.Client.EXPECT().List(gomock.Any(), gomock.AssignableToTypeOf(&v1.PodList{}), gomock.Eq(listOption), gomock.Eq(labelMatcher)).DoAndReturn(
		func(ctx context.Context, pds *v1.PodList, opts ...client.ListOption) error {
			pds.Items = append(pds.Items, pod)
			return nil
		})

	if podAnnotations == nil {
		podAnnotations = make(map[string]string)
	}

	podAnnotations[otterizev1alpha3.AllIntentsRemovedAnnotation] = "true"
	updatedPod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test-pod",
			Labels: map[string]string{
				otterizev1alpha3.OtterizeClientLabelKey: "true",
			},
			Annotations: podAnnotations,
		},
		Spec: v1.PodSpec{},
	}

	s.Client.EXPECT().Patch(gomock.Any(), gomock.Eq(&updatedPod), MatchPatch(client.MergeFrom(&pod))).Return(nil)

	res, err := s.Reconciler.Reconcile(context.Background(), req)
	s.NoError(err)
	s.Empty(res)
}

func (s *PodLabelReconcilerTestSuite) TestAccessLabelChangedOnIntentsEdit() {
	clientIntentsName := "client-intents"
	serviceName := "test-client"

	namespacedName := types.NamespacedName{
		Namespace: testNamespace,
		Name:      clientIntentsName,
	}
	req := ctrl.Request{
		NamespacedName: namespacedName,
	}

	serverName := "test-server"
	intentsSpec := otterizev1alpha3.IntentsSpec{
		Service: otterizev1alpha3.Service{Name: serviceName},
		Calls: []otterizev1alpha3.Intent{
			{
				Name: serverName,
			},
		},
	}

	emptyIntents := &otterizev1alpha3.ClientIntents{}
	s.Client.EXPECT().Get(gomock.Any(), req.NamespacedName, gomock.Eq(emptyIntents)).DoAndReturn(
		func(ctx context.Context, name types.NamespacedName, intents *otterizev1alpha3.ClientIntents, options ...client.ListOption) error {
			intents.Spec = &intentsSpec
			intents.Namespace = testNamespace
			return nil
		})

	var intents otterizev1alpha3.ClientIntents
	intents.Spec = &intentsSpec
	intents.Namespace = testNamespace

	listOption := &client.ListOptions{Namespace: testNamespace}
	labelSelector := labels.SelectorFromSet(map[string]string{
		"intents.otterize.com/server": "test-client-test-namespace-537e87",
	})

	labelMatcher := client.MatchingLabelsSelector{Selector: labelSelector}
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
			Labels:    make(map[string]string),
		},
	}

	s.Client.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Eq(listOption), gomock.Eq(labelMatcher)).DoAndReturn(
		func(ctx context.Context, pds *v1.PodList, opts ...client.ListOption) error {
			pds.Items = append(pds.Items, pod)
			return nil
		})

	updatedPod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
			Labels: map[string]string{
				"intents.otterize.com/access-test-server-test-namespace-8ddecb": "true",
				"intents.otterize.com/client":                                   "test-client-test-namespace-537e87",
			},
		},
		Spec: v1.PodSpec{},
	}

	s.Client.EXPECT().Patch(gomock.Any(), gomock.Eq(&updatedPod), gomock.Any()).Return(nil)

	res, err := s.Reconciler.Reconcile(context.Background(), req)
	s.NoError(err)
	s.Empty(res)

	// Now all the way through again, but with a different server name

	intentsSpec.Calls[0].Name = "test-server-2"

	emptyIntents = &otterizev1alpha3.ClientIntents{}
	s.Client.EXPECT().Get(gomock.Any(), req.NamespacedName, gomock.Eq(emptyIntents)).DoAndReturn(
		func(ctx context.Context, name types.NamespacedName, intents *otterizev1alpha3.ClientIntents, options ...client.ListOption) error {
			intents.Spec = &intentsSpec
			intents.Namespace = testNamespace
			return nil
		})

	s.Client.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Eq(listOption), gomock.Eq(labelMatcher)).DoAndReturn(
		func(ctx context.Context, pds *v1.PodList, opts ...client.ListOption) error {
			pds.Items = append(pds.Items, pod)
			return nil
		})

	updatedPod = v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
			Labels: map[string]string{
				"intents.otterize.com/access-test-server-2-test-namespace-e4423b": "true",
				"intents.otterize.com/client":                                     "test-client-test-namespace-537e87",
			},
		},
		Spec: v1.PodSpec{},
	}

	s.Client.EXPECT().Patch(gomock.Any(), gomock.Eq(&updatedPod), gomock.Any()).Return(nil)

	res, err = s.Reconciler.Reconcile(context.Background(), req)
	s.NoError(err)
	s.Empty(res)

}

func (s *PodLabelReconcilerTestSuite) TestPodLabelFinalizerAdded() {
	clientIntentsName := "client-intents"
	serviceName := "test-client"

	namespacedName := types.NamespacedName{
		Namespace: testNamespace,
		Name:      clientIntentsName,
	}
	req := ctrl.Request{
		NamespacedName: namespacedName,
	}

	serverName := "test-server"
	intentsSpec := otterizev1alpha3.IntentsSpec{
		Service: otterizev1alpha3.Service{Name: serviceName},
		Calls: []otterizev1alpha3.Intent{
			{
				Name: serverName,
			},
		},
	}

	emptyIntents := &otterizev1alpha3.ClientIntents{}
	s.Client.EXPECT().Get(gomock.Any(), req.NamespacedName, gomock.Eq(emptyIntents)).DoAndReturn(
		func(ctx context.Context, name types.NamespacedName, intents *otterizev1alpha3.ClientIntents, options ...client.ListOption) error {
			intents.Spec = &intentsSpec
			return nil
		})

	s.Client.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	res, err := s.Reconciler.Reconcile(context.Background(), req)
	s.NoError(err)
	s.Empty(res)
}

func (s *PodLabelReconcilerTestSuite) TestPodLabelFinalizerRemoved() {
	clientIntentsName := "client-intents"
	serviceName := "test-client"

	namespacedName := types.NamespacedName{
		Namespace: testNamespace,
		Name:      clientIntentsName,
	}
	req := ctrl.Request{
		NamespacedName: namespacedName,
	}

	serverName := "test-server"
	intentsSpec := otterizev1alpha3.IntentsSpec{
		Service: otterizev1alpha3.Service{Name: serviceName},
		Calls: []otterizev1alpha3.Intent{
			{
				Name: serverName,
			},
		},
	}

	emptyIntents := &otterizev1alpha3.ClientIntents{}
	deletionTimestamp := &metav1.Time{Time: time.Now()}
	s.Client.EXPECT().Get(gomock.Any(), req.NamespacedName, gomock.Eq(emptyIntents)).DoAndReturn(
		func(ctx context.Context, name types.NamespacedName, intents *otterizev1alpha3.ClientIntents, options ...client.ListOption) error {
			intents.Spec = &intentsSpec
			intents.DeletionTimestamp = deletionTimestamp
			return nil
		})

	s.Client.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil)

	res, err := s.Reconciler.Reconcile(context.Background(), req)
	s.NoError(err)
	s.Empty(res)
}

func (s *PodLabelReconcilerTestSuite) TestClientAccessLabelAddFailedPatch() {
	clientIntentsName := "client-intents"
	serviceName := "test-client"

	namespacedName := types.NamespacedName{
		Namespace: testNamespace,
		Name:      clientIntentsName,
	}
	req := ctrl.Request{
		NamespacedName: namespacedName,
	}

	serverName := "test-server"
	intentsSpec := otterizev1alpha3.IntentsSpec{
		Service: otterizev1alpha3.Service{Name: serviceName},
		Calls: []otterizev1alpha3.Intent{
			{
				Name: serverName,
			},
		},
	}

	emptyIntents := &otterizev1alpha3.ClientIntents{}
	s.Client.EXPECT().Get(gomock.Any(), req.NamespacedName, gomock.Eq(emptyIntents)).DoAndReturn(
		func(ctx context.Context, name types.NamespacedName, intents *otterizev1alpha3.ClientIntents, options ...client.ListOption) error {
			intents.Spec = &intentsSpec
			intents.Namespace = testNamespace
			return nil
		})

	var intents otterizev1alpha3.ClientIntents
	intents.Spec = &intentsSpec

	listOption := &client.ListOptions{Namespace: testNamespace}
	labelSelector := labels.SelectorFromSet(map[string]string{
		"intents.otterize.com/server": "test-client-test-namespace-537e87",
	})

	labelMatcher := client.MatchingLabelsSelector{Selector: labelSelector}
	pod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
			Labels:    make(map[string]string),
		},
	}

	s.Client.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Eq(listOption), gomock.Eq(labelMatcher)).DoAndReturn(
		func(ctx context.Context, pds *v1.PodList, opts ...client.ListOption) error {
			pds.Items = append(pds.Items, pod)
			return nil
		})

	updatedPod := v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "test-pod",
			Namespace: testNamespace,
			Labels: map[string]string{
				"intents.otterize.com/access-test-server-test-namespace-8ddecb": "true",
				"intents.otterize.com/client":                                   "test-client-test-namespace-537e87",
			},
		},
		Spec: v1.PodSpec{},
	}

	s.Client.EXPECT().Patch(gomock.Any(), gomock.Eq(&updatedPod), gomock.Any()).Return(fmt.Errorf("Patch failed"))

	_, err := s.Reconciler.Reconcile(context.Background(), req)
	s.Error(err)
	s.ExpectEvent(ReasonUpdatePodFailed)
}

func TestPodLabelReconcilerTestSuite(t *testing.T) {
	suite.Run(t, new(PodLabelReconcilerTestSuite))
}
