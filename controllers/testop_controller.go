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

package controllers

import (
	"context"
	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hellov1beta1 "example.com/api/v1beta1"
)

// TestOpReconciler reconciles a TestOp object
type TestOpReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=hello.example.com,resources=testops,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hello.example.com,resources=testops/status,verbs=get;update;patch

func (r *TestOpReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("testop", req.NamespacedName)

	resource := &hellov1beta1.TestOp{}
	err := r.Get(ctx, req.NamespacedName, resource)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		return ctrl.Result{Requeue: false}, nil
	}
	log.Info("resource found. Creating deployment.")
	d := r.deploymentForImage(resource)
	log.Info("Creating new deployment", "deployment namespace", d.Namespace, "deployment name", d.Name)
	err = r.Create(ctx, d)
	if err != nil {
		log.Info("failed to create deployment", "deployment namespace", d.Namespace, "deployment name", d.Name)
		return ctrl.Result{}, err
	}
	//resource.status.DeploymentName = d.Name
	// your logic here

	return ctrl.Result{}, nil
}

func (r *TestOpReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hellov1beta1.TestOp{}).
		Complete(r)
}

func (r *TestOpReconciler) deploymentForImage(s *hellov1beta1.TestOp) *appsv1.Deployment {
	ls := map[string]string{"app": "testop", "name": s.Name}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: s.Spec.Image,
						Name:  "containersTest",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
						}},
					}},
				},
			},
		},
	}
	return dep
}
