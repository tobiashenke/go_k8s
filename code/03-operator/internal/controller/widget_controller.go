package controller

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	logf "sigs.k8s.io/controller-runtime/pkg/log"

	appsv1alpha1 "github.com/tobiashenke/go_k8s/operator/api/v1alpha1"
)

type WidgetReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=apps.henkebyte.dev,resources=widgets,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=apps.henkebyte.dev,resources=widgets/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=apps.henkebyte.dev,resources=widgets/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=configmaps,verbs=get;list;watch;create;update;patch;delete

func (r *WidgetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := logf.FromContext(ctx)
	log.Info("Reconciling Widget", "name", req.Name, "namespace", req.Namespace)

	var widget appsv1alpha1.Widget
	err := r.Get(ctx, req.NamespacedName, &widget)
	if errors.IsNotFound(err) {
		log.Info("Widget not found, likely deleted", "name", req.Name)
		return ctrl.Result{}, nil
	}
	if err != nil {
		return ctrl.Result{}, err
	}

	configMap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      widget.Name,
			Namespace: widget.Namespace,
		},
		Data: map[string]string{
			"message": widget.Spec.Message,
		},
	}

	err = ctrl.SetControllerReference(&widget, &configMap, r.Scheme)
	if err != nil {
		return ctrl.Result{}, err
	}

	_, err = controllerutil.CreateOrUpdate(ctx, r.Client, &configMap, func() error {
		configMap.Data = map[string]string{
			"message": widget.Spec.Message,
		}
		return nil
	})
	if err != nil {
		return ctrl.Result{}, err
	}

	log.Info("ConfigMap reconciled", "name", configMap.Name, "message", widget.Spec.Message)
	return ctrl.Result{}, nil
}

func (r *WidgetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Widget{}).
		Owns(&corev1.ConfigMap{}).
		Named("widget").
		Complete(r)
}
