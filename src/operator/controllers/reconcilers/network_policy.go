package reconcilers

import (
	"context"
	"fmt"
	otterizev1alpha1 "github.com/otterize/intents-operator/shared/api/v1alpha1"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/networking/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const OtterizeNetworkPolicyNameTemplate = "access-to-%s-from-%s"
const OtterizeNetworkPolicy = "otterize/network-policy"
const NetworkPolicyFinalizerName = "otterize-intents.policies/finalizer"

type NetworkPolicyReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

func (r *NetworkPolicyReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	intents := &otterizev1alpha1.Intents{}
	err := r.Get(ctx, req.NamespacedName, intents)
	if k8serrors.IsNotFound(err) {
		return ctrl.Result{}, nil
	}

	if err != nil {
		return ctrl.Result{}, err
	}

	if intents.Spec == nil {
		return ctrl.Result{}, nil
	}

	logrus.Infof("Reconciling network policies for service %s in namespace %s",
		intents.Spec.Service.Name, req.Namespace)

	// Object is deleted, handle finalizer and network policy clean up
	if !intents.DeletionTimestamp.IsZero() {
		err := r.cleanFinalizerAndPolicies(ctx, intents)
		if err != nil {
			if k8serrors.IsConflict(err) {
				return ctrl.Result{Requeue: true}, nil
			}
			return ctrl.Result{}, err
		}
		return ctrl.Result{}, nil
	}

	if !controllerutil.ContainsFinalizer(intents, NetworkPolicyFinalizerName) {
		logrus.Infof("Adding finalizer %s", NetworkPolicyFinalizerName)
		controllerutil.AddFinalizer(intents, NetworkPolicyFinalizerName)
		if err := r.Update(ctx, intents); err != nil {
			return ctrl.Result{}, err
		}
	}
	for _, intent := range intents.GetCallsList() {
		if intent.Namespace == "" {
			// We never actually update the intent, we just set it here, so we can to access it later
			intent.Namespace = req.Namespace
		}
		err := r.handleNetworkPolicyCreation(ctx, intent, req.Namespace)
		if err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *NetworkPolicyReconciler) handleNetworkPolicyCreation(
	ctx context.Context, intent otterizev1alpha1.Intent, intentsObjNamespace string) error {

	policyName := fmt.Sprintf(OtterizeNetworkPolicyNameTemplate, intent.Server, intentsObjNamespace)
	policy := &v1.NetworkPolicy{}
	err := r.Get(ctx, types.NamespacedName{Name: policyName, Namespace: intent.Namespace}, policy)

	// No matching network policy found, create one
	if k8serrors.IsNotFound(err) {
		logrus.Infof(
			"Creating network policy to enable access from namespace %s to %s", intentsObjNamespace, intent.Server)
		policy := r.buildNetworkPolicyObjectForIntent(intent, policyName, intentsObjNamespace)
		err := r.Create(ctx, policy)
		if err != nil {
			return err
		}
		return nil

	} else if err != nil {
		return err
	}

	// Found network policy, check for diff
	// TODO: Add code to compensate for changes in intents vs existing policies
	return nil
}

func (r *NetworkPolicyReconciler) cleanFinalizerAndPolicies(
	ctx context.Context, intents *otterizev1alpha1.Intents) error {
	if controllerutil.ContainsFinalizer(intents, NetworkPolicyFinalizerName) {
		logrus.Infof("Removing network policies for deleted intents for service: %s", intents.Spec.Service.Name)
		for _, intent := range intents.GetCallsList() {
			if intent.Namespace == "" {
				intent.Namespace = intents.Namespace
			}

			var intentsList otterizev1alpha1.IntentsList
			err := r.List(
				ctx, &intentsList,
				&client.MatchingFields{otterizev1alpha1.OtterizeTargetServerIndexField: intent.Server},
				&client.ListOptions{Namespace: intents.Namespace})

			if err != nil {
				return err
			}

			if len(intentsList.Items) == 1 {
				// We have only 1 intents resource that has this server as its target - and it's the current one
				// We need to delete the network policy that allows access from this namespace, as there are no other
				// clients in that namespace that need to access the target server
				logrus.Infof("No other intents in the namespace reference target server: %s", intent.Server)
				logrus.Infoln("Removing matching network policy for server")
				if err = r.removeNetworkPolicy(ctx, intent, intents.Namespace); err != nil {
					return err
				}
			}
		}

		controllerutil.RemoveFinalizer(intents, NetworkPolicyFinalizerName)
		if err := r.Update(ctx, intents); err != nil {
			return err
		}

	}

	return nil
}

func (r *NetworkPolicyReconciler) removeNetworkPolicy(
	ctx context.Context,
	intent otterizev1alpha1.Intent,
	intentsObjNamespace string) error {

	policyName := fmt.Sprintf(OtterizeNetworkPolicyNameTemplate, intent.Server, intentsObjNamespace)
	policy := &v1.NetworkPolicy{}
	err := r.Get(ctx, types.NamespacedName{Name: policyName, Namespace: intent.Namespace}, policy)

	if k8serrors.IsNotFound(err) {
		return nil
	} else {
		err := r.Delete(ctx, policy)
		if err != nil {
			return nil
		}
	}
	return nil
}

// buildNetworkPolicyObjectForIntent builds the network policy that represents the intent from the parameter
func (r *NetworkPolicyReconciler) buildNetworkPolicyObjectForIntent(
	intent otterizev1alpha1.Intent, policyName, intentsObjNamespace string) *v1.NetworkPolicy {
	formattedTargetServer := otterizev1alpha1.GetFormattedOtterizeIdentity(intent.Server, intent.Namespace)

	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      policyName,
			Namespace: intent.Namespace,
			Labels: map[string]string{
				OtterizeNetworkPolicy: "true",
			},
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					otterizev1alpha1.OtterizeServerLabelKey: formattedTargetServer,
				},
			},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					From: []v1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									fmt.Sprintf(
										otterizev1alpha1.OtterizeAccessLabelKey, formattedTargetServer): "true",
								},
							},
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									otterizev1alpha1.OtterizeNamespaceLabelKey: intentsObjNamespace,
								},
							},
						},
					},
				},
			},
		},
	}
}
