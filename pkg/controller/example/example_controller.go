package example

import (
	"context"

	examplev1alpha1 "github.com/christianalexander/example-operator/pkg/apis/example/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_example")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Example Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileExample{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("example-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Example
	err = c.Watch(&source.Kind{Type: &examplev1alpha1.Example{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// Watch for changes to secondary resource Pods and requeue the owner Example
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &examplev1alpha1.Example{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileExample implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileExample{}

// ReconcileExample reconciles a Example object
type ReconcileExample struct {
// This client, initialized using mgr.Client() above, is a split client
// that reads objects from the cache and writes to the apiserver
client client.Client
scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Example object and makes changes based on the state read
// and what is in the Example.Spec
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileExample) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling Example")

	instance := &examplev1alpha1.Example{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			return reconcile.Result{}, nil
		}
		return reconcile.Result{}, err
	}

	if len(instance.Status.Values) == instance.Spec.Count {
		return reconcile.Result{}, nil
	}

	singleValue := "example"

	values := make([]string, instance.Spec.Count)
	for i := 0; i < instance.Spec.Count; i++ {
		values[i] = singleValue
	}
	instance.Status.Values = values

	err = r.client.Status().Update(context.TODO(), instance)
	if err != nil {
		reqLogger.Error(err, "Failed to update Example status")
		return reconcile.Result{}, err
	}


	return reconcile.Result{}, nil
}
