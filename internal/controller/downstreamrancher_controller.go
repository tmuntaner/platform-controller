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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	platformv1alpha1 "tmuntaner/platform-controller/api/v1alpha1"
)

// DownstreamRancherReconciler reconciles a DownstreamRancher object
type DownstreamRancherReconciler struct {
	client.Client
	Scheme        *runtime.Scheme
	rancherClient *RancherClient
}

//+kubebuilder:rbac:groups=platform.rubyrainbows.com,resources=downstreamranchers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=platform.rubyrainbows.com,resources=downstreamranchers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=platform.rubyrainbows.com,resources=downstreamranchers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the DownstreamRancher object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *DownstreamRancherReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Reconciling DownstreamRancher objects")

	instance := &platformv1alpha1.DownstreamRancher{}
	if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
		return r.handleError(ctx, req, logger, instance, client.IgnoreNotFound(err))
	}

	// if the instance is being deleted, remove the finalizer
	if !instance.ObjectMeta.DeletionTimestamp.IsZero() {
		logger.Info(fmt.Sprintf("Removing finalizer for DownstreamRancher object: %s", instance.Name))
		controllerutil.RemoveFinalizer(instance, finalizer)
		if err := r.Update(ctx, instance); err != nil {
			return r.handleError(ctx, req, logger, instance, err)
		}

		return ctrl.Result{}, nil
	}

	logger.Info(fmt.Sprintf("Reconsiling DownstreamRancher account by name \"%s\" with current state \"%s\"", instance.Spec.Name, instance.Status.State))

	// the instance is without a state, so it should be initialized
	if instance.Status.State == "" {
		// add the finalizer to the object, so we can handle the object deletion within the controller
		logger.Info(fmt.Sprintf("Adding finalizer to DownstreamRancher account by id \"%s\"", instance.Name))
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
		logger.Info(fmt.Sprintf("Working on DownstreamRancher object: %s", instance.Name))

		upstream := &platformv1alpha1.UpstreamRancher{}
		upstreamNamespacedName := types.NamespacedName{
			Name:      instance.Spec.UpstreamRancher,
			Namespace: req.Namespace,
		}
		if err := r.Get(ctx, upstreamNamespacedName, upstream); err != nil {
			return r.handleError(ctx, req, logger, instance, err)
		}

		clientset, err := r.rancherClient.GetKubeConfigSecret(ctx, req, upstream)
		if err != nil {
			return r.handleError(ctx, req, logger, instance, err)
		}

		pods, err := clientset.CoreV1().Pods("kube-system").List(ctx, metav1.ListOptions{})
		if err != nil {
			return r.handleError(ctx, req, logger, instance, err)
		}
		for _, pod := range pods.Items {
			fmt.Println(pod.Name)
		}
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
func (r *DownstreamRancherReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.rancherClient == nil {
		r.rancherClient = NewRancherClient(r.Client)
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&platformv1alpha1.DownstreamRancher{}).
		Complete(r)
}

func (r *DownstreamRancherReconciler) handleError(_ context.Context, _ ctrl.Request, logger logr.Logger, rancher *platformv1alpha1.DownstreamRancher, reconcileError error) (ctrl.Result, error) {
	if reconcileError == nil {
		return ctrl.Result{}, nil
	}

	if strings.Contains(reconcileError.Error(), optimisticLockErrorMsg) {
		return ctrl.Result{RequeueAfter: time.Second}, nil
	}

	logger.Error(reconcileError, fmt.Sprintf("Reconcile failed for DownstreamRancher %s with error: %s", rancher.Spec.Name, reconcileError.Error()))
	return ctrl.Result{}, reconcileError
}
