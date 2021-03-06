/*
 * Copyright 2019 The Knative Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package v1alpha1

import (
	duckv1alpha1 "github.com/knative/pkg/apis/duck/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

type PerfJob struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec defines the desired state of the PerfJob.
	Spec PerfJobSpec `json:"spec,omitempty"`

	// Status represents the current state of the PerfJob.
	// +optional
	Status PerfJobStatus `json:"status,omitempty"`
}

// Check that PerfJob can be validated, can be defaulted, and has immutable fields.
var _ runtime.Object = (*PerfJob)(nil)

type PerfJobSpec struct {
	TestImage    string `json:"testImage,omitempty"`
	ControlImage string `json:"ctrlImage,omitempty"`
	Exclusive    *bool  `json:"exclusive,omitempty"`
	Target       string `json:"target,omitempty"`
}

var pjCondSet = duckv1alpha1.NewBatchConditionSet()

// BrokerStatus represents the current state of a Broker.
type PerfJobStatus struct {
	// inherits duck/v1alpha1 Status, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last processed by the controller.
	// * Conditions - the latest available observations of a resource's current state.
	duckv1alpha1.Status `json:",inline"`

	Phase string
}

const (
	PerfJobConditionSucceeded = duckv1alpha1.ConditionSucceeded
	BatchJobConditionFinished = "BatchJobFinished"
)

// GetCondition returns the condition currently associated with the given type, or nil.
func (ps *PerfJobStatus) GetCondition(t duckv1alpha1.ConditionType) *duckv1alpha1.Condition {
	return pjCondSet.Manage(ps).GetCondition(t)
}

// IsReady returns true if the resource is ready overall.
func (ps *PerfJobStatus) IsReady() bool {
	return pjCondSet.Manage(ps).IsHappy()
}

// MarkJobSucceeded sets the BatchJobConditionSucceeded condition to True state.
func (ps *PerfJobStatus) MarkJobRunning() {
	pjCondSet.Manage(ps).MarkUnknown(BatchJobConditionFinished, "running", "todo")
}

// MarkJobSucceeded sets the BatchJobConditionSucceeded condition to True state.
func (ps *PerfJobStatus) MarkJobSucceeded() {
	pjCondSet.Manage(ps).MarkTrue(BatchJobConditionFinished)
}

// MarkJobFailed sets the BatchJobConditionSucceeded condition to True state.
func (ps *PerfJobStatus) MarkJobFailed() {
	pjCondSet.Manage(ps).MarkFalse(BatchJobConditionFinished, "failed", "todo")
}

// InitializeConditions sets relevant unset conditions to Unknown state.
func (ps *PerfJobStatus) InitializeConditions() {
	pjCondSet.Manage(ps).InitializeConditions()
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PerfJobList is a collection of PerfJobs.
type PerfJobList struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PerfJob `json:"items"`
}
