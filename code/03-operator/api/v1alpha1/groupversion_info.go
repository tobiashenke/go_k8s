// Package v1alpha1 contains API Schema definitions for the apps v1alpha1 API group.
// +kubebuilder:object:generate=true
// +groupName=apps.henkebyte.dev
package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	SchemeGroupVersion = schema.GroupVersion{Group: "apps.henkebyte.dev", Version: "v1alpha1"}
	GroupVersion       = SchemeGroupVersion
	SchemeBuilder      = runtime.NewSchemeBuilder(func(scheme *runtime.Scheme) error {
		metav1.AddToGroupVersion(scheme, SchemeGroupVersion)
		return nil
	})
	AddToScheme = SchemeBuilder.AddToScheme
)
