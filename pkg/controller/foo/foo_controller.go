/*

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

package foo

import (
	"context"
	"reflect"

	foogroupv1 "github.com/raker22/k8s-foo-resource/pkg/apis/foogroup/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new Foo Controller and adds it to the Manager with default RBAC. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileFoo{Client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("foo-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to Foo
	err = c.Watch(&source.Kind{Type: &foogroupv1.Foo{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create
	// Uncomment watch a Deployment created by Foo - change this for objects you create
	//err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
	//	IsController: true,
	//	OwnerType:    &foogroupv1.Foo{},
	//})
	//if err != nil {
	//	return err
	//}

	return nil
}

var _ reconcile.Reconciler = &ReconcileFoo{}

// ReconcileFoo reconciles a Foo object
type ReconcileFoo struct {
	client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a Foo object and makes changes based on the state read
// and what is in the Foo.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  The scaffolding writes
// a Deployment as an example
// Automatically generate RBAC rules to allow the Controller to read and write Foos
// +kubebuilder:rbac:groups=foogroup.raker22.com,resources=foos,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=foogroup.raker22.com,resources=foos/status,verbs=get;update;patch
func (r *ReconcileFoo) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	// Fetch the Foo instance
	instance := &foogroupv1.Foo{}
	err := r.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Object not found, return.  Created objects are automatically garbage collected.
			// For additional cleanup logic use finalizers.
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// TODO(user): Change this to be the object type created by your controller
	// Define the desired Deployment object
	foo := instance.DeepCopy()

	foo.Status = foogroupv1.FooStatus{
		Message: instance.Spec.Message,
		Value:   instance.Spec.Value,
	}

	//if err := controllerutil.SetControllerReference(instance, foo, r.scheme); err != nil {
	//	return reconcile.Result{}, err
	//}

	// TODO(user): Change this for the object type created by your controller
	// Check if the Deployment already exists
	found := &foogroupv1.Foo{}
	err = r.Get(context.TODO(), types.NamespacedName{Name: foo.Name, Namespace: foo.Namespace}, found)
	if err != nil {
		return reconcile.Result{}, err
	}

	// TODO(user): Change this for the object type created by your controller
	// Update the found object and write the result back if there are any changes
	if !reflect.DeepEqual(foo.Status, found.Status) {
		found.Status = foo.Status
		log.Info("Updating Foo", "namespace", foo.Namespace, "name", foo.Name)
		err = r.Update(context.TODO(), found)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return reconcile.Result{}, nil
}
