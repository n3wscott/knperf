/*
Copyright 2019 The Knative Authors

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

package resources

import (
	"fmt"

	perfv1alpha1 "github.com/n3wscott/knperf/pkg/apis/perf/v1alpha1"
	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func JobLabels(perfJob *perfv1alpha1.PerfJob) map[string]string {
	return map[string]string{
		"perfJob": "pj-" + perfJob.Name,
	}
}

// MakeJob creates a Job to start or stop a Feed.
func NewJob(perfJob *perfv1alpha1.PerfJob, target string) *batchv1.Job {
	podTemplate := makePodTemplate(perfJob.Spec.Image, target)
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			GenerateName: "perf-",
			Namespace:    perfJob.Namespace,
			Labels:       JobLabels(perfJob),
			OwnerReferences: []metav1.OwnerReference{
				*metav1.NewControllerRef(perfJob, perfv1alpha1.SchemeGroupVersion.WithKind("PerfJob")),
			},
		},
		Spec: batchv1.JobSpec{
			Template: *podTemplate,
		},
	}
}

func IsJobComplete(job *batchv1.Job) bool {
	for _, c := range job.Status.Conditions {
		if c.Type == batchv1.JobComplete && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func IsJobFailed(job *batchv1.Job) bool {
	for _, c := range job.Status.Conditions {
		if c.Type == batchv1.JobFailed && c.Status == corev1.ConditionTrue {
			return true
		}
	}
	return false
}

func JobFailedMessage(job *batchv1.Job) string {
	for _, c := range job.Status.Conditions {
		if c.Type == batchv1.JobFailed && c.Status == corev1.ConditionTrue {
			return fmt.Sprintf("[%s] %s", c.Reason, c.Message)
		}
	}
	return ""
}

func GetFirstTerminationMessage(pod *corev1.Pod) string {
	for _, cs := range pod.Status.ContainerStatuses {
		if cs.State.Terminated != nil && cs.State.Terminated.Message != "" {
			return cs.State.Terminated.Message
		}
	}
	return ""
}

// makePodTemplate creates a pod template for a feed stop or start Job.
func makePodTemplate(image, target string) *corev1.PodTemplateSpec {
	return &corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Annotations: map[string]string{
				"sidecar.istio.io/inject": "true",
			},
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: "default",
			RestartPolicy:      corev1.RestartPolicyNever,
			Containers: []corev1.Container{{
				Name:            "job",
				Image:           image,
				ImagePullPolicy: "Always",
				Env: []corev1.EnvVar{{
					Name:  "TARGET",
					Value: target,
				}, {
					Name: "POD_NAME",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.name",
						},
					},
				}, {
					Name: "POD_NAMESPACE",
					ValueFrom: &corev1.EnvVarSource{
						FieldRef: &corev1.ObjectFieldSelector{
							FieldPath: "metadata.namespace",
						},
					},
				}},
			}},
		},
	}
}
