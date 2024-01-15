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
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"strings"
	"time"
	platformv1alpha1 "tmuntaner/platform-controller/api/v1alpha1"
)

// UpstreamRancherReconciler reconciles a UpstreamRancher object
type UpstreamRancherReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	rancherClient *RancherClient
}

const (
	finalizer = "rancher.platform.rubyrainbows.com"
)

//+kubebuilder:rbac:groups=platform.rubyrainbows.com,resources=upstreamranchers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=platform.rubyrainbows.com,resources=upstreamranchers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=platform.rubyrainbows.com,resources=upstreamranchers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *UpstreamRancherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling UpstreamRancher objects")

	instance := &platformv1alpha1.UpstreamRancher{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return r.handleError(ctx, req, logger, instance, client.IgnoreNotFound(err))
	}

	// if the instance is being deleted, remove the finalizer
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info(fmt.Sprintf("Removing finalizer for UpstreamRancher object: %s", instance.Name))
		controllerutil.RemoveFinalizer(instance, finalizer)
		if err := r.Update(ctx, instance); err != nil {
			return r.handleError(ctx, req, logger, instance, err)
		}

		return ctrl.Result{}, nil
	}

	logger.Info(fmt.Sprintf("Reconsiling UpstreamRancher account by name \"%s\" with current state \"%s\"", instance.Spec.Name, instance.Status.State))

	// the instance is without a state, so it should be initialized
	if instance.Status.State == "" {
		// add the finalizer to the object, so we can handle the object deletion within the controller
		logger.Info(fmt.Sprintf("Adding finalizer to UpstreamRancher account by id \"%s\"", instance.Name))
		controllerutil.AddFinalizer(instance, finalizer)
		if err := r.Update(ctx, instance); err != nil {
			return r.handleError(ctx, req, logger, instance, err)
		}

		// set the object to pending, so it can be picked up in the next reconcile loop
		instance.Status.State = platformv1alpha1.PendingState
		if err := r.Status().Update(ctx, instance); err != nil {
			return r.handleError(ctx, req, logger, instance, err)
		}

		return ctrl.Result{}, nil
	}

	if instance.Status.State == platformv1alpha1.PendingState {
		logger.Info(fmt.Sprintf("Working on UpstreamRancher object: %s", instance.Name))

		//logger.Info(fmt.Sprintf("Secret data: %s", string(secretData)))
		// do some work here

		instance.Status.State = platformv1alpha1.ReadyState
		if err := r.Status().Update(ctx, instance); err != nil {
			return r.handleError(ctx, req, logger, instance, err)
		}
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *UpstreamRancherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	r.rancherClient = NewRancherClient(r.Client)

	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1alpha1.UpstreamRancher{}).
		Complete(r)
}

func (r *UpstreamRancherReconciler) handleError(_ context.Context, _ ctrl.Request, logger logr.Logger, rancher *platformv1alpha1.UpstreamRancher, reconcileError error) (ctrl.Result, error) {
	if reconcileError == nil {
		return ctrl.Result{}, nil
	}

	if strings.Contains(reconcileError.Error(), optimisticLockErrorMsg) {
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	logger.Error(reconcileError, fmt.Sprintf("Reconcile failed for UpstreamRancher %s with error: %s", rancher.Spec.Name, reconcileError.Error()))
	return ctrl.Result{}, reconcileError
}
