/*
Copyright 2021 SAP.

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

package controllers

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/log"

	openstackstablesapccv2 "github.com/sapcc/openstack-seeder/api/v2"
	"github.com/sapcc/openstack-seeder/config/options"
	"github.com/sapcc/openstack-seeder/openstack"
)

// OpenstackSeedReconciler reconciles a OpenstackSeed object
type OpenstackSeedReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=openstack.stable.sap.cc,resources=openstackseeds,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=openstack.stable.sap.cc,resources=openstackseeds/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=openstack.stable.sap.cc,resources=openstackseeds/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the OpenstackSeed object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.8.3/pkg/reconcile
func (r *OpenstackSeedReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	l := log.FromContext(ctx, "OpenstackSeed")
	var seed openstackstablesapccv2.OpenstackSeed
	if err := r.Get(ctx, req.NamespacedName, &seed); err != nil {
		l.Error(err, "unable to fetch OpenstackSeed")
		// we'll ignore not-found errors, since they can't be fixed by an immediate
		// requeue (we'll need to wait for a new notification), and we can get them
		// on deleted requests.
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}
	r.reconcileSeeds(seed)
	if len(seed.Status.UnfinishedSeeds) > 0 {
		seedStatus.WithLabelValues(seed.Name).Set(0)
		return ctrl.Result{RequeueAfter: 10 * time.Minute}, nil
	}
	seedStatus.WithLabelValues(seed.Name).Set(1)
	return ctrl.Result{RequeueAfter: 24 * time.Hour}, nil
}

func (r *OpenstackSeedReconciler) reconcileSeeds(seed openstackstablesapccv2.OpenstackSeed) error {
	v := reflect.ValueOf(seed)
	completed := true
	if len(seed.Status.UnfinishedSeeds) > 0 {
		for i := 0; i < v.NumField(); i++ {
			if _, ok := seed.Status.UnfinishedSeeds[v.Field(i).Type().Name()]; ok {
				err := r.reconcileSeed(v.Field(i).Type().Name(), seed)
				if err == nil {
					delete(seed.Status.UnfinishedSeeds, v.Field(i).Type().Name())
				} else {
					completed = false
				}
			}
		}
	} else {
		for i := 0; i < v.NumField(); i++ {
			err := r.reconcileSeed(v.Field(i).Type().Name(), seed)
			if err != nil {
				completed = false
				seed.Status.UnfinishedSeeds[v.Field(i).Type().Name()] = err.Error()
			}
		}
	}
	if !completed {
		return fmt.Errorf("seed was not successfull")
	}
	return nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *OpenstackSeedReconciler) SetupWithManager(mgr ctrl.Manager, opts options.Options) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&openstackstablesapccv2.OpenstackSeed{}).
		WithOptions(controller.Options{MaxConcurrentReconciles: opts.MaxConcurrentReconciles}).
		Complete(r)
}

func (r *OpenstackSeedReconciler) reconcileSeed(name string, seed openstackstablesapccv2.OpenstackSeed) (err error) {
	switch name {
	case "domains":
		err = r.seedDomains(seed.Spec.Domains)
	case "regions":
		err = r.seedRegions(seed.Spec.Regions)
	case "service":
		err = r.seedServices(seed.Spec.Services)
	case "role":
		err = r.seedRoles(seed.Spec.Roles)
	}
	return err
}

func (r *OpenstackSeedReconciler) seedDomains(domains []openstackstablesapccv2.DomainSpec) (err error) {
	c, err := openstack.NewIdentityClient()
	if err != nil {
		return
	}
	k := openstack.NewKeystone(c)
	if len(domains) > 0 {
		for _, d := range domains {
			if _, err := k.SeedDomain(d); err != nil {
				return err
			}
		}
	}
	return
}

func (r *OpenstackSeedReconciler) seedRegions(regions []openstackstablesapccv2.RegionSpec) (err error) {
	c, err := openstack.NewIdentityClient()
	if err != nil {
		return
	}
	k := openstack.NewKeystone(c)
	if len(regions) > 0 {
		for _, d := range regions {
			if _, err := k.SeedRegion(d); err != nil {
				return err
			}
		}
	}
	return
}

func (r *OpenstackSeedReconciler) seedServices(services []openstackstablesapccv2.ServiceSpec) (err error) {
	c, err := openstack.NewIdentityClient()
	if err != nil {
		return
	}
	k := openstack.NewKeystone(c)
	if len(services) > 0 {
		for _, s := range services {
			if _, err := k.SeedService(s); err != nil {
				return err
			}
		}
	}
	return
}

func (r *OpenstackSeedReconciler) seedRoles(roles []openstackstablesapccv2.RoleSpec) (err error) {
	c, err := openstack.NewIdentityClient()
	if err != nil {
		return
	}
	k := openstack.NewKeystone(c)
	if len(roles) > 0 {
		for _, r := range roles {
			if _, err := k.SeedRole(r); err != nil {
				return err
			}
		}
	}
	return
}
