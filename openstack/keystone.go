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
	"os"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/regions"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/services"

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
		updated, err = roles.Create(k.Client, opts).Extract()
		return updated, err
	}
	role := rs[0]
	if !isEqual(specRole, role) {
		update := roles.UpdateOpts{}
		d, _ := json.Marshal(specRole)
		json.Unmarshal(d, &update)
		update.Extra = specRole.Extra
		updated, err = roles.Update(k.Client, role.ID, update).Extract()
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
		updated, err = regions.Create(k.Client, opts).Extract()
		return updated, err
	}
	if !isEqual(specRegion, region) {
		opts := regions.UpdateOpts{}
		b, _ := json.Marshal(specRegion)
		json.Unmarshal(b, &opts)
		updated, err = regions.Update(k.Client, region.ID, opts).Extract()
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
		updated, err = services.Create(k.Client, opts).Extract()
		return
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
		updated, err = services.Update(k.Client, svc.ID, opts).Extract()
		return
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
		domains.Create(k.Client, opts)
	}
	domain := ds[0]
	d, _ := json.Marshal(spec)
	specDomain := domains.Domain{}
	json.Unmarshal(d, &specDomain)
	if !isEqual(specDomain, domain) {
		upd := domains.UpdateOpts{}
		json.Unmarshal(d, &upd)
		updated, err = domains.Update(k.Client, domain.ID, upd).Extract()
	}
	return
}

func (k *Keystone) GetProjectID(domain, name string) (id string, err error) {
	if id, ok := k.cache.Get("project", fmt.Sprintf("%s.%s", domain, name)); ok {
		return id, nil
	}
	domainID, err := k.GetDomainID(domain)
	if err != nil {
		return
	}
	p, err := projects.List(k.Client, projects.ListOpts{
		Name:     name,
		DomainID: domainID,
	}).AllPages()
	if err != nil {
		return
	}
	r, err := projects.ExtractProjects(p)
	if err != nil {
		return
	}
	if len(r) != 1 {
		return id, fmt.Errorf("could not find project: %s", name)
	}
	k.cache.Add("project", r[0].ID, name, 0)
	return r[0].ID, err
}

func (k *Keystone) GetDomainID(name string) (id string, err error) {
	if id, ok := k.cache.Get("domain", name); ok {
		return id, nil
	}
	p, err := domains.List(k.Client, domains.ListOpts{Name: name}).AllPages()
	if err != nil {
		return
	}
	d, err := domains.ExtractDomains(p)
	if err != nil {
		return
	}
	if len(d) != 1 {
		return id, fmt.Errorf("could not find domain: %s", name)
	}
	k.cache.Add("domain", d[0].ID, name, 0)
	return d[0].ID, err
}
