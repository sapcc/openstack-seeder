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

	th "github.com/gophercloud/gophercloud/testhelper"
	"github.com/gophercloud/gophercloud/testhelper/client"
)

// ListOutput provides a single page of Service results.
const ListServiceOutput = `
{
    "links": {
        "next": null,
        "previous": null
    },
    "services": [
        {
            "id": "9876",
            "links": {
                "self": "https://example.com/identity/v3/services/9876"
            },
            "type": "compute",
            "enabled": false,
            "extra": {
                "name": "service-two",
                "description": "Service Two"
            }
        }
    ]
}
`

// GetOutput provides a Get result.
const GetServiceOutput = `
{
    "service": {
        "id": "9876",
        "links": {
            "self": "https://example.com/identity/v3/services/9876"
        },
        "type": "compute",
        "enabled": false,
        "extra": {
            "name": "service-two",
            "description": "Service Two",
            "email": "service@example.com"
        }
    }
}
`

// UpdateOutput provides an update result.
const UpdateServiceOutput = `
{
    "service": {
        "type": "compute",
        "enabled": true,
        "name": "service-new",
        "description": "Service New"
    }
}
`

// UpdateRequest provides the input to as Update request.
const UpdateServiceRequest = `
{
    "service": {
        "type": "compute",
        "description": "Service New"
    }
}
`

// HandleListServicesSuccessfully creates an HTTP handler at `/services` on the
// test handler mux that responds with a list of two services.
func HandleListServicesSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "GET")
		th.TestHeader(t, r, "Accept", "application/json")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, ListServiceOutput)
	})
}

// HandleCreateServiceSuccessfully creates an HTTP handler at `/services` on the
// test handler mux that tests service creation.
func HandleCreateServiceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/services", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "POST")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, CreateRequest)

		w.WriteHeader(http.StatusCreated)
		fmt.Fprintf(w, GetServiceOutput)
	})
}

// HandleUpdateServiceSuccessfully creates an HTTP handler at `/services` on the
// test handler mux that tests service update.
func HandleUpdateServiceSuccessfully(t *testing.T) {
	th.Mux.HandleFunc("/services/9876", func(w http.ResponseWriter, r *http.Request) {
		th.TestMethod(t, r, "PATCH")
		th.TestHeader(t, r, "X-Auth-Token", client.TokenID)
		th.TestJSONRequest(t, r, UpdateServiceOutput)

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, UpdateServiceOutput)
	})
}
