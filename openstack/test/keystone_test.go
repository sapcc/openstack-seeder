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

package test

import (
	"testing"

	openstackstablesapccv2 "github.com/sapcc/openstack-seeder/api/v2"
	"github.com/sapcc/openstack-seeder/openstack"
	"github.com/stretchr/testify/assert"

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

func TestSeedService(t *testing.T) {
	specEqual := openstackstablesapccv2.ServiceSpec{
		Type:        "compute",
		Enabled:     false,
		Name:        "service-two",
		Description: "Service Two",
	}
	specNotEqual := openstackstablesapccv2.ServiceSpec{
		Type:        "compute",
		Enabled:     true,
		Name:        "service-new",
		Description: "Service New",
	}
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListServicesSuccessfully(t)
	HandleUpdateServiceSuccessfully(t)

	kc := openstack.Keystone{Client: client.ServiceClient()}
	upd, err := kc.SeedService(specEqual)
	assert.NoError(t, err)
	assert.Nil(t, upd)

	upd, err = kc.SeedService(specNotEqual)
	assert.NoError(t, err)
	assert.NotNil(t, upd)
}

func TestSeedDomain(t *testing.T) {
	specEqual := openstackstablesapccv2.DomainSpec{
		Name:        "domain one",
		Description: "some description",
		Enabled:     true,
	}
	specNotEqual := openstackstablesapccv2.DomainSpec{
		Name:        "domain new",
		Description: "some other description",
		Enabled:     true,
	}
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListDomainsSuccessfully(t)
	HandleUpdateDomainSuccessfully(t)

	kc := openstack.Keystone{Client: client.ServiceClient()}
	upd, err := kc.SeedDomain(specEqual)
	assert.NoError(t, err, "patch should not be called")
	assert.Nil(t, upd, "patch should not be called")

	updt, err := kc.SeedDomain(specNotEqual)
	assert.NoError(t, err, "patch should be called")
	assert.Equal(t, updt.Description, specNotEqual.Description, "patch should be called")
	assert.Equal(t, updt.Name, specNotEqual.Name, "patch should be called")
}

func TestSeedRole(t *testing.T) {
	specEqual := openstackstablesapccv2.RoleSpec{
		Name:        "some_role",
		Description: "some description",
		DomainID:    "default",
	}
	specNotEqual := openstackstablesapccv2.RoleSpec{
		Name:        "some_role",
		Description: "new description",
	}
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListAndCreateRolesSuccessfully(t)
	HandleUpdateRoleSuccessfully(t)

	kc := openstack.Keystone{Client: client.ServiceClient()}

	upd, err := kc.SeedRole(specEqual)
	assert.NoError(t, err, "error should be nil")
	assert.Nil(t, upd, "patch should not be called")

	upd, err = kc.SeedRole(specNotEqual)
	assert.NoError(t, err, "error should be nil")
	assert.NotNil(t, upd, "patch should not be called")

}

func TestSeedRoleNotExist(t *testing.T) {
	specNotExist := openstackstablesapccv2.RoleSpec{
		Name:        "does_not_exist",
		Description: "some description",
	}
	th.SetupHTTP()
	defer th.TeardownHTTP()
	HandleListAndCreateRolesSuccessfully(t)

	kc := openstack.Keystone{Client: client.ServiceClient()}
	upd, err := kc.SeedRole(specNotExist)
	assert.NoError(t, err, "error should be nil")
	assert.Equal(t, upd.Name, specNotExist.Name)
	assert.Equal(t, upd.Extra["description"], specNotExist.Description)
}
