package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

func (r *WidgetReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = logf.FromContext(ctx)

	// TODO: your logic here

	return ctrl.Result{}, nil
}

func (r *WidgetReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&appsv1alpha1.Widget{}).
		Named("widget").
		Complete(r)
}
