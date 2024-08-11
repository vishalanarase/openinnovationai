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
	"sort"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"

	infrav1 "github.com/vishalanarase/openinnovationai/distributed-job-scheduler-operator/api/v1"
)

// ComputeJobReconciler reconciles a ComputeJob object
type ComputeJobReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=infra.openinnovation.ai,resources=computejobs,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=infra.openinnovation.ai,resources=computejobs/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=infra.openinnovation.ai,resources=computejobs/finalizers,verbs=update
// +kubebuilder:rbac:groups=core,resources=pods,verbs=get;list;watch;create;update;patch;delete

func (r *ComputeJobReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	computeJob := &infrav1.ComputeJob{}
	err := r.Get(ctx, req.NamespacedName, computeJob)
	if err != nil {
		if errors.IsNotFound(err) {
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get compute job")

		return ctrl.Result{}, err
	}

	if computeJob.Status.State == "" {
		logger.Info("Update compute job status", "Status", infrav1.JobPending)
		computeJob.Status.State = string(infrav1.JobPending)
		computeJob.Status.StartTime = &metav1.Time{Time: time.Now()}
		err := r.Status().Update(ctx, computeJob)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	// Implement scheduling logic
	computeNodes := &infrav1.ComputeNodeList{}
	err = r.List(ctx, computeNodes)
	if err != nil {
		logger.Error(err, "Failed to list compute nodes")
		return ctrl.Result{}, err
	}

	// Sort the ComputeNodeList by name
	sort.Slice(computeNodes.Items, func(i, j int) bool {
		return computeNodes.Items[i].Name < computeNodes.Items[j].Name
	})

	// Select the nodes to run pod
	selectedNodes := selectNodes(computeNodes, computeJob.Spec.NodeSelector, computeJob.Spec.Parallelism)
	if len(selectedNodes) == 0 {
		logger.Info("Update compute job status", "Status", infrav1.JobFailed)
		computeJob.Status.State = string(infrav1.JobFailed)
		computeJob.Status.EndTime = &metav1.Time{Time: time.Now()}
		return ctrl.Result{}, r.Status().Update(ctx, computeJob)
	}

	activeNodes := []string{}
	// Ensure Pods are created on selected nodes
	for _, node := range selectedNodes {
		pod := &corev1.Pod{}
		err := r.Get(ctx, podName(computeJob.ObjectMeta, node.Name), pod)
		if err != nil {
			if !errors.IsNotFound(err) {
				logger.Error(err, "Failed to get pod for compute job")
				return ctrl.Result{}, err
			}

			pod = createPodForComputeJob(computeJob, &node)
			// Set the owner reference for the Pod to be the ComputeJob
			err = controllerutil.SetControllerReference(computeJob, pod, r.Scheme)
			if err != nil {
				logger.Error(err, "Failed to set controller reference on pod")
				return ctrl.Result{}, err
			}

			err = r.Client.Create(ctx, pod)
			if err != nil {
				logger.Error(err, "Failed to create pod for compute job", "Node", node.Name)
				continue
			}

			activeNodes = append(activeNodes, node.Name)
		}
	}

	pods, err := r.listPodsOwnedByComputeJob(ctx, computeJob)
	if err != nil {
		logger.Error(err, "Failed to list pods owned by compute job")
		return reconcile.Result{}, err
	}

	// Check if all Pods are completed
	if areAllPodsCompleted(pods) {
		logger.Info("Update compute job status", "Status", infrav1.JobCompleted)
		computeJob.Status.State = string(infrav1.JobCompleted)
		computeJob.Status.EndTime = &metav1.Time{Time: time.Now()}
		err := r.Status().Update(ctx, computeJob)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else {
		logger.Info("Update compute job status", "Status", infrav1.JobRunning)
		computeJob.Status.State = string(infrav1.JobRunning)
		computeJob.Status.ActiveNodes = activeNodes
		err := r.Status().Update(ctx, computeJob)
		if err != nil {
			return reconcile.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

// listPodsOwnedByComputeJob lists all Pods owned by the specified ComputeJob
func (r *ComputeJobReconciler) listPodsOwnedByComputeJob(ctx context.Context, job *infrav1.ComputeJob) ([]corev1.Pod, error) {
	podList := corev1.PodList{}
	err := r.Client.List(ctx, &podList, client.InNamespace(job.Namespace))
	if err != nil {
		return nil, err
	}

	ownedPods := []corev1.Pod{}
	for _, pod := range podList.Items {
		for _, ownerRef := range pod.OwnerReferences {
			if ownerRef.UID == job.UID {
				ownedPods = append(ownedPods, pod)
				break
			}
		}
	}

	return ownedPods, nil
}

// selectNodes selects nodes based on nodeSelector and parallelism
func selectNodes(nodes *infrav1.ComputeNodeList, selector map[string]string, parallelism int) []infrav1.ComputeNode {
	selected := []infrav1.ComputeNode{}
	for _, node := range nodes.Items {
		if matchesSelector(node, selector) && len(selected) < parallelism {
			selected = append(selected, node)
		}
	}

	return selected
}

// matchesSelector checks if a ComputeNode matches the provided nodeSelector criteria
func matchesSelector(node infrav1.ComputeNode, selector map[string]string) bool {
	// Iterate over the selector criteria
	for key, value := range selector {
		// Check if the node has the label with the specified key
		if nodeValue, exists := node.Labels[key]; !exists || nodeValue != value {
			// If the node doesn't have the label or the value doesn't match, return false
			return false
		}
	}

	return true
}

// createPodForComputeJob creates a pod that runs the job on the specified node
func createPodForComputeJob(job *infrav1.ComputeJob, node *infrav1.ComputeNode) *corev1.Pod {
	name := podName(job.ObjectMeta, node.Name)
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name.Name,
			Namespace: name.Namespace,
			Labels:    job.Spec.NodeSelector,
		},
		Spec: corev1.PodSpec{
			NodeName: node.Name,
			Containers: []corev1.Container{
				{
					Name:    "job-container",
					Image:   "ubuntu",
					Command: []string{"/bin/sh", "-c", job.Spec.Command},
				},
			},
			RestartPolicy: corev1.RestartPolicyOnFailure,
		},
	}
}

func podName(meta metav1.ObjectMeta, nodeName string) types.NamespacedName {
	return types.NamespacedName{
		Name:      meta.Name + "-" + nodeName,
		Namespace: meta.Namespace,
	}
}

// areAllPodsCompleted checks if all listed Pods are in the 'Completed' state
func areAllPodsCompleted(pods []corev1.Pod) bool {
	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodSucceeded {
			return false
		}
	}

	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *ComputeJobReconciler) SetupWithManager(mgr ctrl.Manager) error {
	trueVal := true
	return ctrl.NewControllerManagedBy(mgr).
		For(&infrav1.ComputeJob{}).
		Owns(&corev1.Pod{}).
		WithOptions(controller.Options{
			// MaxConcurrentReconciles: 10,
			RecoverPanic: &trueVal,
		}).
		Complete(r)
}
