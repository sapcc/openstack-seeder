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
	"fmt"
	"net/http"
	"testing"

	"github.com/gophercloud/gophercloud/openstack/identity/v3/roles"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// ListOutput provides a single page of Role results.
const ListRoleOutput = `
{
    "links": {
        "next": null,
        "previous": null,
        "self": "http://example.com/identity/v3/roles"
    },
    "roles": [
        {
            "domain_id": "default",
            "id": "2844b2a08be147a08ef58317d6471f1f",
            "links": {
                "self": "http://example.com/identity/v3/roles/2844b2a08be147a08ef58317d6471f1f"
            },
            "name": "some
			_role",
			"extra": {
                "description": "some description"
            }
        },
        {
            "domain_id": "1789d1",
            "id": "9fe1d3",
            "links": {
                "self": "https://example.com/identity/v3/roles/9fe1d3"
            },
            "name": "support",
            "extra": {
                "description": "read-only support role"
            }
        }
    ]
}
`

const ListEmptyRoleOutput = `
{
    "roles": []
}
`

// GetOutput provides a Get result.
const GetRoleOutput = `
{
    "role": {
        "domain_id": "1789d1",
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/roles/9fe1d3"
        },
        "name": "does_not_exist",
        "extra": {
            "description": "some description"
        }
    }
}
`

// CreateRequest provides the input to a Create request.
const CreateRoleRequest = `
{
    "role": {
        "name": "does_not_exist",
        "description": "some description"
    }
}
`

// UpdateRequest provides the input to as Update request.
const UpdateRoleRequest = `
{
    "role": {
		"name": "some_role",
        "description": "new description"
    }
}
`

// UpdateOutput provides an update result.
const UpdateRoleOutput = `
{
    "role": {
        "domain_id": "1789d1",
        "id": "2844b2a08be147a08ef58317d6471f1f",
        "name": "some_role",
        "extra": {
            "description": "new description"
        }
    }
}
`

// ListAssignmentOutput provides a result of ListAssignment request.
const ListAssignmentOutput = `
{
    "role_assignments": [
        {
            "links": {
                "assignment": "http://identity:35357/v3/domains/161718/users/313233/roles/123456"
            },
            "role": {
                "id": "123456"
            },
            "scope": {
                "domain": {
                    "id": "161718"
                }
            },
            "user": {
                "id": "313233"
            }
        },
        {
            "links": {
                "assignment": "http://identity:35357/v3/projects/456789/groups/101112/roles/123456",
                "membership": "http://identity:35357/v3/groups/101112/users/313233"
            },
            "role": {
                "id": "123456"
            },
            "scope": {
                "project": {
                    "id": "456789"
                }
            },
            "user": {
                "id": "313233"
            }
        }
    ],
    "links": {
        "self": "http://identity:35357/v3/role_assignments?effective",
        "previous": null,
        "next": null
    }
}
`

// FirstRole is the first role in the List request.
var FirstRole = roles.Role{
	DomainID: "default",
	ID:       "2844b2a08be147a08ef58317d6471f1f",
	Links: map[string]interface{}{
		"self": "http://example.com/identity/v3/roles/2844b2a08be147a08ef58317d6471f1f",
	},
	Name:  "admin-read-only",
	Extra: map[string]interface{}{},
}

// HandleListRolesSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that responds with a list of two roles.
func HandleListAndCreateRolesSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if r.URL.Query().Get("name") == "does_not_exist" {
				fmt.Println("eeeeeeeeeeeeeempty")
				fmt.Fprintf(w, ListEmptyRoleOutput)
			} else {
				fmt.Fprintf(w, ListRoleOutput)
			}
		case http.MethodPost:
			th.TestMethod(t, r, "POST")
			th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
			th.TestJSONRequest(t, r, CreateRoleRequest)

			w.WriteHeader(http.StatusCreated)
			fmt.Fprintf(w, GetRoleOutput)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
}

// HandleUpdateRoleSuccessfully creates an HTTP handler at `/roles` on the
// test handler mux that tests role update.
func HandleUpdateRoleSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/roles/2844b2a08be147a08ef58317d6471f1f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateRoleRequest)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateRoleOutput)
	})
}
