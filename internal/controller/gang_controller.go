/*
Copyright 2024.

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

package controller

import (
	"context"
	"fmt"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	schedulingv1 "github.com/HtcOrange/gang-scheduler-api/api/v1"
	v1 "github.com/HtcOrange/gang-scheduler-api/api/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// GangReconciler reconciles a Gang object
type GangReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=scheduling.github.com,resources=gangs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=scheduling.github.com,resources=gangs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=scheduling.github.com,resources=gangs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Gang object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *GangReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := log.FromContext(ctx)

	// get gang
	var gang v1.Gang
	if err := r.Get(ctx, req.NamespacedName, &gang); err != nil {
		if errors.IsNotFound(err) {
			log.Info("Gang resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		log.Error(err, "Failed to get Gang")
		return ctrl.Result{}, err
	}

	// create or update PodGroup
	if err := r.reconcilePodGroups(ctx, &gang); err != nil {
		log.Error(err, "Failed to reconcile PodGroups")
		return ctrl.Result{}, err
	}

	// update Gang status
	if err := r.updateGangStatus(ctx, &gang); err != nil {
		log.Error(err, "Failed to update Gang status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *GangReconciler) reconcilePodGroups(ctx context.Context, gang *v1.Gang) error {
	for i := int32(0); i < gang.Spec.Replicas; i++ {
		podGroup := &v1.PodGroup{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%d", gang.Spec.Template.Spec.GroupName, i),
				Namespace: gang.Namespace,
				Labels:    gang.Spec.Labels,
			},
			Spec: gang.Spec.Template.Spec,
		}

		if err := ctrl.SetControllerReference(gang, podGroup, r.Scheme); err != nil {
			return err
		}
		if err := r.Create(ctx, podGroup); err != nil && !errors.IsAlreadyExists(err) {
			return err
		}
	}
	return nil
}

func (r *GangReconciler) updateGangStatus(ctx context.Context, gang *v1.Gang) error {
	var podGroupList v1.PodGroupList
	if err := r.List(ctx, &podGroupList, client.InNamespace(gang.Namespace), client.MatchingLabels(gang.Spec.Labels)); err != nil {
		return err
	}

	readyReplicas := int32(0)
	for _, podGroup := range podGroupList.Items {
		if podGroup.Status.Phase == "Ready" {
			readyReplicas++
		}
	}

	gang.Status.ReadyReplicas = readyReplicas
	return r.Status().Update(ctx, gang)
}

// SetupWithManager sets up the controller with the Manager.
func (r *GangReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&schedulingv1.Gang{}).
		Named("gang").
		Complete(r)
}
