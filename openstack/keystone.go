/**
 * Copyright 2021 SAP SE
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package openstack

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/endpoints"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"

	openstackstablesapccv2 "github.com/sapcc/openstack-seeder/api/v2"
	"github.com/sapcc/openstack-seeder/pkg/cache"
)

type Keystone struct {
	Client *gophercloud.ServiceClient
	cache  *cache.Cache
}

func NewKeystone(client *gophercloud.ServiceClient) (k *Keystone) {
	return &Keystone{
		cache:  cache.New(time.Until(time.Now().AddDate(0, 0, 7)), 30*time.Minute),
		Client: client,
	}
}

func NewIdentityClient() (client *gophercloud.ServiceClient, err error) {
	opts, err := openstack.AuthOptionsFromEnv()
	if err != nil {
		return
	}
	provider, err := openstack.AuthenticatedClient(opts)
	if err != nil {
		return
	}
	client, err = openstack.NewIdentityV3(provider, gophercloud.EndpointOpts{
		Region: os.Getenv("OS_REGION_NAME"),
	})
	return
}

func (k *Keystone) SeedRole(spec openstackstablesapccv2.RoleSpec) (updated *roles.Role, err error) {
	p, err := roles.List(k.Client, roles.ListOpts{
		Name:     spec.Name,
		DomainID: spec.DomainID,
	}).AllPages()
	if err != nil {
		return
	}
	rs, err := roles.ExtractRoles(p)
	if err != nil {
		return
	}
	d, _ := json.Marshal(spec)
	specRole := roles.Role{}
	specRole.UnmarshalJSON(d)
	if len(rs) == 0 {
		opts := roles.CreateOpts{}
		b, err := json.Marshal(specRole)
		if err != nil {
			return updated, err
		}
		json.Unmarshal(b, &opts)
		opts.Extra = specRole.Extra
		return roles.Create(k.Client, opts).Extract()
	}
	role := rs[0]
	if !isEqual(specRole, role) {
		update := roles.UpdateOpts{}
		d, _ := json.Marshal(specRole)
		json.Unmarshal(d, &update)
		update.Extra = specRole.Extra
		return roles.Update(k.Client, role.ID, update).Extract()
	}
	return
}

func (k *Keystone) SeedRegion(spec openstackstablesapccv2.RegionSpec) (updated *regions.Region, err error) {
	if spec.ID == "" {
		//TODO: warning log
		return
	}
	r := regions.Get(k.Client, spec.ID)
	region, err := r.Extract()
	if err != nil {
		return
	}
	d, _ := json.Marshal(spec)
	specRegion := regions.Region{}
	specRegion.UnmarshalJSON(d)
	if region == nil {
		opts := regions.CreateOpts{}
		b, err := json.Marshal(specRegion)
		if err != nil {
			return updated, err
		}
		json.Unmarshal(b, &opts)
		opts.Extra = specRegion.Extra
		return regions.Create(k.Client, opts).Extract()
	}
	if !isEqual(specRegion, region) {
		opts := regions.UpdateOpts{}
		b, _ := json.Marshal(specRegion)
		json.Unmarshal(b, &opts)
		return regions.Update(k.Client, region.ID, opts).Extract()
	}
	return
}

func (k *Keystone) SeedService(spec openstackstablesapccv2.ServiceSpec) (updated *services.Service, err error) {
	p, err := services.List(k.Client, services.ListOpts{Name: spec.Name, ServiceType: spec.Type}).AllPages()
	svcs, err := services.ExtractServices(p)
	if err != nil {
		return
	}
	if len(svcs) == 0 {
		opts := services.CreateOpts{}
		b, _ := json.Marshal(spec)
		json.Unmarshal(b, &opts)
		return services.Create(k.Client, opts).Extract()
	}
	svc := svcs[0]
	d, _ := json.Marshal(spec)
	specService := services.Service{}
	specService.UnmarshalJSON(d)
	if !isEqual(specService, svc) {
		opts := services.UpdateOpts{}
		b, _ := json.Marshal(specService)
		json.Unmarshal(b, &opts)
		opts.Extra = specService.Extra
		return services.Update(k.Client, svc.ID, opts).Extract()
	}
	return
}

func (k *Keystone) SeedDomain(spec openstackstablesapccv2.DomainSpec) (updated *domains.Domain, err error) {
	p, err := domains.List(k.Client, domains.ListOpts{Name: spec.Name}).AllPages()
	if err != nil {
		return
	}
	ds, err := domains.ExtractDomains(p)
	if err != nil {
		return
	}
	if len(ds) == 0 {
		opts := domains.CreateOpts{}
		b, err := json.Marshal(spec)
		if err != nil {
			return updated, err
		}
		json.Unmarshal(b, &opts)
		return domains.Create(k.Client, opts).Extract()
	}
	domain := ds[0]
	d, _ := json.Marshal(spec)
	specDomain := domains.Domain{}
	json.Unmarshal(d, &specDomain)
	if !isEqual(specDomain, domain) {
		upd := domains.UpdateOpts{}
		json.Unmarshal(d, &upd)
		return domains.Update(k.Client, domain.ID, upd).Extract()
	}
	return
}

func (k *Keystone) SeedEndpoints(spec openstackstablesapccv2.EndpointSpec, serviceID string) (updated *endpoints.Endpoint, err error) {
	_, err = url.Parse(spec.URL)
	if err != nil {
		return
	}

	p, err := endpoints.List(k.Client, endpoints.ListOpts{
		ServiceID: serviceID,
		RegionID:  spec.Region,
	}).AllPages()
	if err != nil {
		return
	}
	e, err := endpoints.ExtractEndpoints(p)
	if err != nil {
		return
	}
	if len(e) == 0 {
		opts := endpoints.CreateOpts{}
		b, err := json.Marshal(spec)
		if err != nil {
			return updated, err
		}
		json.Unmarshal(b, &opts)
		return endpoints.Create(k.Client, opts).Extract()
	}
	endpoint := e[0]
	d, _ := json.Marshal(spec)
	specEndpoint := endpoints.Endpoint{}
	json.Unmarshal(d, &specEndpoint)
	if !isEqual(specEndpoint, endpoint) {
		upd := endpoints.UpdateOpts{}
		json.Unmarshal(d, &upd)
		return endpoints.Update(k.Client, endpoint.ID, upd).Extract()
	}
	return
}

func (k *Keystone) SeedUser(spec openstackstablesapccv2.UserSpec) (updated *users.User, err error) {
	p, err := users.List(k.Client, users.ListOpts{
		DomainID: "",
		Name:     spec.Name,
	}).AllPages()
	if err != nil {
		return
	}
	u, err := users.ExtractUsers(p)
	if err != nil {
		return
	}
	if len(u) == 0 {
		opts := users.CreateOpts{}
		b, err := json.Marshal(spec)
		if err != nil {
			return updated, err
		}
		json.Unmarshal(b, &opts)
		return users.Create(k.Client, opts).Extract()
	}
	user := u[0]
	d, _ := json.Marshal(spec)
	specUser := users.User{}
	specUser.UnmarshalJSON(d)
	if !isEqual(specUser, user) {
		opts := users.UpdateOpts{}
		b, _ := json.Marshal(specUser)
		json.Unmarshal(b, &opts)
		opts.Extra = specUser.Extra
		return users.Update(k.Client, user.ID, opts).Extract()
	}
	return
}

func (k *Keystone) SeedRoleAssignment(spec openstackstablesapccv2.RoleAssignmentSpec) (err error) {
	assignOpts := roles.AssignOpts{}
	if spec.User != "" {
		r := strings.Split(spec.User, "@")
		if len(r) != 2 {
			return fmt.Errorf("user name wrong format: domain@name")
		}
		assignOpts.UserID, err = k.GetUserID(r[0], r[1])
		if err != nil {
			return
		}
	}
	roleID, err := k.GetRoleID(spec.Role)
	if err != nil {
		return
	}
	if spec.ProjectID != "" {
		assignOpts.ProjectID = spec.ProjectID
	} else {
		r := strings.Split(spec.Project, "@")
		assignOpts.ProjectID, err = k.GetProjectID(r[0], r[1])
		if err != nil {
			return
		}
	}
	if spec.Domain != "" {
		assignOpts.DomainID, err = k.GetDomainID(spec.Domain)
		if err != nil {
			return
		}
	}
	return roles.Assign(k.Client, roleID, assignOpts).ExtractErr()
}
