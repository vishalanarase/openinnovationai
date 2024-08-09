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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// JobState represents the state of a job
type JobState string

// These are the valid statuses of jobs
const (
	// JobPending means that the job is pending
	JobPending JobState = "Pending"
	// JobRunning means that the job is running
	JobRunning JobState = "Running"
	// JobFailed means that the job is failed
	JobFailed JobState = "Failed"
	// JobCompleted means that the job is completed
	JobCompleted JobState = "Completed"
)

// ComputeJobSpec defines the desired state of ComputeJob
type ComputeJobSpec struct {
	// The command to run as a job
	Command string `json:"command,omitempty"`
	// Criteria for selecting nodes to run the job
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`
	// The number of nodes the job should run on simultaneously
	Parallelism int `json:"parallelism,omitempty"`
}

// ComputeJobStatus defines the observed state of ComputeJob
type ComputeJobStatus struct {
	// The current state of the job (e.g., Pending, Running, Completed, Failed)
	State string `json:"state,omitempty"`
	// The start time of the job
	StartTime *metav1.Time `json:"startTime,omitempty"`
	// The end time of the job
	EndTime *metav1.Time `json:"endTime,omitempty"`
	// The list of nodes where the job is currently running
	ActiveNodes []string `json:"activeNodes,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// ComputeJob is the Schema for the computejobs API
type ComputeJob struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ComputeJobSpec   `json:"spec,omitempty"`
	Status ComputeJobStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ComputeJobList contains a list of ComputeJob
type ComputeJobList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ComputeJob `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ComputeJob{}, &ComputeJobList{})
}
