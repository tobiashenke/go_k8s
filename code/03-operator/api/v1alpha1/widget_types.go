package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// WidgetSpec defines the desired state of Widget
type WidgetSpec struct {
	Message string `json:"message"`
}

// WidgetStatus defines the observed state of Widget.
type WidgetStatus struct {
	// +listType=map
	// +listMapKey=type
	// +optional
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Widget is the Schema for the widgets API
type Widget struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              WidgetSpec   `json:"spec"`
	Status            WidgetStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WidgetList contains a list of Widget
type WidgetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []Widget `json:"items"`
}

func init() {
	SchemeBuilder.Register(func(s *runtime.Scheme) error {
		s.AddKnownTypes(SchemeGroupVersion, &Widget{}, &WidgetList{})
		return nil
	})
}
