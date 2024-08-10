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

package computenode

import (
	"context"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/log"

	infrav1 "github.com/vishalanarase/openinnovationai/distributed-job-scheduler-operator/api/v1"
)

// ComputeNodeReconciler reconciles a ComputeNode object
type ComputeNodeReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infra.openinnovation.ai,resources=computenodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infra.openinnovation.ai,resources=computenodes/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infra.openinnovation.ai,resources=computenodes/finalizers,verbs=update
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch
func (r *ComputeNodeReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// If kubernetes node is deleted then delete the compute node
	node := &corev1.Node{}
	err := r.Client.Get(ctx, types.NamespacedName{Name: req.Name}, node)
	if err != nil && errors.IsNotFound(err) {
		// Fetch the compute node
		computeNode := &infrav1.ComputeNode{}
		err := r.Client.Get(ctx, types.NamespacedName{Name: node.Name}, computeNode)
		if err != nil && !errors.IsNotFound(err) {
			logger.Error(err, "Failed to get compute node")
			return ctrl.Result{}, err
		}

		// Delete the compute node
		if err := r.Client.Delete(ctx, computeNode, &client.DeleteOptions{}); err != nil {
			logger.Error(err, "Failed to delete compute node")
			return ctrl.Result{}, err
		}

		return ctrl.Result{}, nil
	}

	// Fetch the list of kubernetes nodes
	nodes := &corev1.NodeList{}
	if err := r.Client.List(ctx, nodes); err != nil {
		logger.Error(err, "Failed to list kubernetes nodes")
		return ctrl.Result{}, err
	}

	for _, node := range nodes.Items {
		// Fetch the compute node
		computeNode := &infrav1.ComputeNode{}
		err := r.Client.Get(ctx, types.NamespacedName{Name: node.Name}, computeNode)
		if err != nil && !errors.IsNotFound(err) {
			logger.Error(err, "Failed to get compute node")
		}

		if errors.IsNotFound(err) {
			// Create a new compute node resource
			computeNode = &infrav1.ComputeNode{
				ObjectMeta: metav1.ObjectMeta{
					Name:        node.Name,
					Annotations: node.Annotations,
					Labels:      node.Labels,
				},
				Spec: infrav1.ComputeNodeSpec{
					Resources: node.Status.Capacity,
				},
				Status: infrav1.ComputeNodeStatus{
					State: "Pending",
				},
			}

			if err := r.Client.Create(ctx, computeNode); err != nil {
				logger.Error(err, "Failed to create compute node resource")
				continue
			}

			logger.Info("Created new compute node", "Name", computeNode.Name)
		}

		// Sync the compute node status with the actual kubernetes node status
		nodeReady := isNodeReady(node.Status.Conditions)
		desiredState := "Running"
		if !nodeReady {
			desiredState = "Failed"
		}

		if computeNode.Status.State != desiredState {
			computeNode.Status.State = desiredState
			if err := r.Status().Update(ctx, computeNode); err != nil {
				logger.Error(err, "Failed to update compute node status")
				continue
			}
			logger.Info("Updated compute node status", "Name", computeNode.Name, "State", desiredState)
		}
	}

	return ctrl.Result{}, nil
}

// isNodeReady checks if a kubernetes node is in ready state
func isNodeReady(conditions []corev1.NodeCondition) bool {
	for _, condition := range conditions {
		if condition.Type == corev1.NodeReady && condition.Status == corev1.ConditionTrue {
			return true
		}
	}

	return false
}

// SetupWithManager sets up the controller with the manager.
func (r *ComputeNodeReconciler) SetupWithManager(mgr ctrl.Manager) error {
	trueVal := true
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.ComputeNode{}).
		Watches(&corev1.Node{}, &handler.EnqueueRequestForObject{}).
		WithOptions(controller.Options{
			MaxConcurrentReconciles: 10,
			RecoverPanic:            &trueVal,
		}).
		Complete(r)
}
