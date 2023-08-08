/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhooks

import (
	"context"
	"fmt"
	"github.com/asaskevich/govalidator"
	otterizev1alpha2 "github.com/otterize/intents-operator/src/operator/api/v1alpha2"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"strings"
)

type ProtectedServiceValidator struct {
	client.Client
}

func (v *ProtectedServiceValidator) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&otterizev1alpha2.ProtectedService{}).
		WithValidator(v).
		Complete()
}

func NewProtectedServiceValidator(c client.Client) *ProtectedServiceValidator {
	return &ProtectedServiceValidator{
		Client: c,
	}
}

//+kubebuilder:webhook:path=/validate-k8s-otterize-com-v1alpha2-protectedservices,mutating=false,failurePolicy=fail,sideEffects=None,groups=k8s.otterize.com,resources=protectedservices,verbs=create;update,versions=v1alpha2,name=protectedservice.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &ProtectedServiceValidator{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (v *ProtectedServiceValidator) ValidateCreate(ctx context.Context, obj runtime.Object) error {
	var allErrs field.ErrorList
	protectedService := obj.(*otterizev1alpha2.ProtectedService)

	protectedServicesList := &otterizev1alpha2.ProtectedServiceList{}
	if err := v.List(ctx, protectedServicesList, &client.ListOptions{Namespace: protectedService.Namespace}); err != nil {
		return err
	}

	if err := v.validateNoDuplicateClients(protectedService, protectedServicesList); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := v.validateSpec(protectedService); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	gvk := protectedService.GroupVersionKind()
	return errors.NewInvalid(
		schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind},
		protectedService.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (v *ProtectedServiceValidator) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) error {
	var allErrs field.ErrorList
	protectedService := newObj.(*otterizev1alpha2.ProtectedService)

	protectedServicesList := &otterizev1alpha2.ProtectedServiceList{}
	if err := v.List(ctx, protectedServicesList, &client.ListOptions{Namespace: protectedService.Namespace}); err != nil {
		return err
	}

	if err := v.validateNoDuplicateClients(protectedService, protectedServicesList); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := v.validateSpec(protectedService); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil
	}

	gvk := protectedService.GroupVersionKind()
	return errors.NewInvalid(
		schema.GroupKind{Group: gvk.Group, Kind: gvk.Kind},
		protectedService.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (v *ProtectedServiceValidator) ValidateDelete(ctx context.Context, obj runtime.Object) error {
	return nil
}

func (v *ProtectedServiceValidator) validateNoDuplicateClients(
	protectedService *otterizev1alpha2.ProtectedService, protectedServicesList *otterizev1alpha2.ProtectedServiceList) *field.Error {

	protectedServiceName := protectedService.Spec.Name
	for _, protectedServiceFromList := range protectedServicesList.Items {
		// Deny admission if intents already exist for this client, and it's not the same object being updated
		if protectedServiceFromList.Spec.Name == protectedServiceName && protectedServiceFromList.Name != protectedService.Spec.Name {
			return &field.Error{
				Type:     field.ErrorTypeDuplicate,
				Field:    "name",
				BadValue: protectedServiceName,
				Detail: fmt.Sprintf(
					"Protected service for service %s already exist in resource %s", protectedServiceName, protectedServiceFromList.Name),
			}
		}
	}
	return nil
}

// validateSpec
func (v *ProtectedServiceValidator) validateSpec(protectedService *otterizev1alpha2.ProtectedService) *field.Error {
	serviceName := strings.ReplaceAll(protectedService.Spec.Name, "-", "")
	serviceName = strings.ReplaceAll(serviceName, "_", "")
	// Validate Service Name contains only lowercase alphanumeric characters
	// Service name should be a valid RFC 1123 subdomain name
	// It's a namespaced resource, we do not expect resources in other namespaces
	if !govalidator.IsAlphanumeric(serviceName) || !govalidator.IsLowerCase(serviceName) {
		message := fmt.Sprintf("Invalid Name: %s. Service name must contain only lowercase alphanumeric characters, '-' or '_'", protectedService.Spec.Name)
		return &field.Error{
			Type:   field.ErrorTypeForbidden,
			Field:  "Name",
			Detail: message,
		}
	}

	return nil
}
