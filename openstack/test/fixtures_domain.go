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

	"github.com/gophercloud/gophercloud/openstack/identity/v3/domains"
	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// ListOutput provides a single page of Domain results.
const ListOutput = `
{
    "domains": [
        {
            "enabled": true,
            "id": "2844b2a08be147a08ef58317d6471f1f",
            "name": "domain one",
            "description": "some description"
        },
        {
            "enabled": true,
            "id": "9fe1d3",
            "name": "domain two"
        }
    ]
}
`

// UpdateRequest provides the input to as Update request.
const UpdateRequest = `
{
    "domain": {
        "name": "domain new",
        "description": "some other description",
        "enabled": true
    }
}
`

// UpdateOutput provides an update result.
const UpdateOutput = `
{
    "domain": {
		"enabled": true,
        "id": "2844b2a08be147a08ef58317d6471f1f",
        "name": "domain new",
        "description": "some other description"
    }
}
`

// CreateRequest provides the input to a Create request.
const CreateRequest = `
{
    "domain": {
        "name": "domain two"
    }
}
`

// GetOutput provides a Get result.
const GetOutput = `
{
    "domain": {
        "enabled": true,
        "id": "9fe1d3",
        "links": {
            "self": "https://example.com/identity/v3/domains/9fe1d3"
        },
        "name": "domain two"
    }
}
`

// SecondDomainUpdated is how SecondDomain should look after an Update.
var DomainUpdated = domains.Domain{
	Enabled:     true,
	ID:          "9fe1d3",
	Name:        "domain two",
	Description: "Staging Domain",
}

// HandleListDomainsSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that responds with a list of two domains.
func HandleListDomainsSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListOutput)
	})
}

// HandleCreateDomainSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that tests domain creation.
func HandleCreateDomainSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, GetOutput)
	})
}

// HandleUpdateDomainSuccessfully creates an HTTP handler at `/domains` on the
// test handler mux that tests domain update.
func HandleUpdateDomainSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/domains/2844b2a08be147a08ef58317d6471f1f", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateRequest)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateOutput)
	})
}
