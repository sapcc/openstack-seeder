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
	"fmt"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/groups"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/projects"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
	"github.com/gophercloud/gophercloud/openstack/identity/v3/users"
)

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
	k.cache.Add("project", fmt.Sprintf("%s.%s", domain, name), r[0].ID, 0)
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
	k.cache.Add("domain", name, d[0].ID, 0)
	return d[0].ID, err
}

func (k *Keystone) GetUserID(domain, name string) (id string, err error) {
	if id, ok := k.cache.Get("user", fmt.Sprintf("%s.%s", domain, name)); ok {
		return id, nil
	}
	dID, err := k.GetDomainID(domain)
	if err != nil {
		return
	}
	p, err := users.List(k.Client, users.ListOpts{
		DomainID: dID,
		Name:     name,
	}).AllPages()
	if err != nil {
		return
	}
	u, err := users.ExtractUsers(p)
	if err != nil {
		return
	}
	if len(u) != 1 {
		return id, fmt.Errorf("could not find user: %s", name)
	}
	k.cache.Add("user", fmt.Sprintf("%s.%s", domain, name), u[0].ID, 0)
	return u[0].ID, err
}

func (k *Keystone) GetGroupID(domain, name string) (id string, err error) {
	if id, ok := k.cache.Get("group", fmt.Sprintf("%s.%s", domain, name)); ok {
		return id, nil
	}
	dID, err := k.GetDomainID(domain)
	if err != nil {
		return
	}
	p, err := groups.List(k.Client, groups.ListOpts{
		DomainID: dID,
		Name:     name,
	}).AllPages()
	if err != nil {
		return
	}
	u, err := groups.ExtractGroups(p)
	if err != nil {
		return
	}
	if len(u) != 1 {
		return id, fmt.Errorf("could not find user: %s", name)
	}
	k.cache.Add("group", fmt.Sprintf("%s.%s", domain, name), u[0].ID, 0)
	return u[0].ID, err
}

func (k *Keystone) GetRoleID(name string) (id string, err error) {
	if id, ok := k.cache.Get("role", name); ok {
		return id, nil
	}
	p, err := roles.List(k.Client, roles.ListOpts{Name: name}).AllPages()
	if err != nil {
		return
	}
	r, err := roles.ExtractRoles(p)
	if err != nil {
		return
	}
	if len(r) != 1 {
		return id, fmt.Errorf("could not find user: %s", name)
	}
	k.cache.Add("role", name, r[0].ID, 0)
	return r[0].ID, err
}
