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

package perfjob

import (
	"context"
	"fmt"
	"github.com/n3wscott/knperf/pkg/reconciler/v1alpha1/perfcmd/resources"
	"k8s.io/api/batch/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime/schema"

	"github.com/knative/pkg/logging"
	perfv1alpha1 "github.com/n3wscott/knperf/pkg/apis/perf/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	controllers "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func Add(manager controllers.Manager) error {

	if err := perfv1alpha1.AddToScheme(manager.GetScheme()); err != nil {
		fmt.Print("failed to add scheme to manager")
	}

	return controllers.
		NewControllerManagedBy(manager).
		For(&perfv1alpha1.PerfJob{}).
		Owns(&v1.Job{}).
		Complete(&Reconciler{Client: manager.GetClient()})
}

type Reconciler struct {
	client.Client
}

func (r *Reconciler) Reconcile(req controllers.Request) (controllers.Result, error) {
	ctx := context.TODO()

	// Read the PerfJob
	pj := &perfv1alpha1.PerfJob{}
	err := r.Get(ctx, req.NamespacedName, pj)
	if errors.IsNotFound(err) {
		logging.FromContext(ctx).Info("Could not find PerfJob", req.Name)
		return controllers.Result{}, nil
	} else if err != nil {
		return controllers.Result{}, err
	}

	return r.reconcilePerfJob(ctx, pj)
}

func (r *Reconciler) reconcilePerfJob(ctx context.Context, pj *perfv1alpha1.PerfJob) (controllers.Result, error) {
	//	logger := logging.FromContext(ctx)

	if pj.DeletionTimestamp != nil {
		// Everything is cleaned up by the garbage collector.
		return controllers.Result{}, nil
	}

	pj.Status.InitializeConditions()

	newJob := resources.NewJob(pj)

	job, err := r.getJob(ctx, pj, labels.SelectorFromSet(resources.JobLabels(pj)))
	// If the resource doesn't exist, we'll create it
	if k8serrors.IsNotFound(err) {
		job = newJob
		err = r.Create(ctx, job)
		if err != nil {
			return controllers.Result{}, err
		}
	} else if err != nil {
		return controllers.Result{}, err
	}

	return controllers.Result{}, nil
}

// getChannel returns the Channel object for Broker 'b' if exists, otherwise it returns an error.
func (r *Reconciler) getJob(ctx context.Context, perfJob *perfv1alpha1.PerfJob, ls labels.Selector) (*v1.Job, error) {
	list := &v1.JobList{}
	opts := &client.ListOptions{
		Namespace:     perfJob.Namespace,
		LabelSelector: ls,
		// Set Raw because if we need to get more than one page, then we will put the continue token
		// into opts.Raw.Continue.
		Raw: &metav1.ListOptions{},
	}

	err := r.List(ctx, opts, list)
	if err != nil {
		return nil, err
	}
	for _, i := range list.Items {
		if metav1.IsControlledBy(&i, perfJob) {
			return &i, nil
		}
	}

	return nil, k8serrors.NewNotFound(schema.GroupResource{}, "")
}
