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
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	schedulingv1 "github.com/HtcOrange/gang-scheduler-api/api/v1"
	v1 "github.com/HtcOrange/gang-scheduler-api/api/v1"
)

// PodGroupReconciler reconciles a PodGroup object
type PodGroupReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=scheduling.github.com,resources=podgroups,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=scheduling.github.com,resources=podgroups/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=scheduling.github.com,resources=podgroups/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=pods,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the PodGroup object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.19.1/pkg/reconcile
func (r *PodGroupReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// get podgroup
	var podGroup schedulingv1.PodGroup
	if err := r.Get(ctx, req.NamespacedName, &podGroup); err != nil {
		if errors.IsNotFound(err) {
			// podgroup has been deleted
			return ctrl.Result{}, err
		}
	}

	// get podlist in one podgroup
	var pods corev1.PodList
	if err := r.List(ctx, &pods, client.InNamespace(podGroup.Namespace), client.MatchingLabels{"pod-group": podGroup.Spec.GroupName}); err != nil {
		logger.Error(err, "Failed to get pod list")
		return ctrl.Result{}, err
	}
	podCount := len(pods.Items)
	logger.Info("PodGroup status", "GroupName", podGroup.Spec.GroupName, "PodCount", podCount)

	scheduledPodNum, podMetaList := getPodsScheduledNumAndMetaList(pods)
	podGroup.Status.PodMetaList = podMetaList
	podGroup.Status.Message = fmt.Sprintf("%d pods has been created, %d pods has been scheduled", podCount, scheduledPodNum)

	// sync PodGroup status
	if podCount >= podGroup.Spec.MinNum {
		if scheduledPodNum < podCount {
			// some of pods are still pending
			podGroup.Status.Phase = "PartialScheduled"
		} else {
			// all of the pods are scheduled
			podGroup.Status.Phase = "Ready"
			podGroup.Status.Message = fmt.Sprintf("PodGroup %s is Ready", podGroup.Spec.GroupName)
		}
	} else {
		// count of pods created is less than min num
		podGroup.Status.Phase = "PartialCreated"
	}

	// update PodGroup status
	if err := r.Status().Update(ctx, &podGroup); err != nil {
		logger.Error(err, "Update PodGroup status failed", podGroup.Spec.GroupName)
		return ctrl.Result{}, nil
	}

	return ctrl.Result{RequeueAfter: 3 * time.Second}, nil
}

func getPodsScheduledNumAndMetaList(podList corev1.PodList) (int, []v1.PodMeta) {
	var scheduledCount int
	var podMetaList []v1.PodMeta
	for _, pod := range podList.Items {
		if pod.Spec.NodeName != "" {
			// pod's spec nodeName has been added
			scheduledCount++
		}
		meta := v1.PodMeta{
			Name:   pod.Name,
			IP:     pod.Status.PodIP,
			HostIP: pod.Status.HostIP,
			Phase:  string(pod.Status.Phase),
		}
		podMetaList = append(podMetaList, meta)
	}
	return scheduledCount, podMetaList
}

// SetupWithManager sets up the controller with the Manager.
func (r *PodGroupReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&schedulingv1.PodGroup{}).
		Named("podgroup").
		Complete(r)
}
