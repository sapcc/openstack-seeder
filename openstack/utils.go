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
)

func isEqual(spec, os interface{}) (equal bool) {
	var specInterface map[string]interface{}
	var osInterface map[string]interface{}
	specB, _ := json.Marshal(spec)
	osB, _ := json.Marshal(os)
	json.Unmarshal(specB, &specInterface)
	json.Unmarshal(osB, &osInterface)

	for field, val := range specInterface {
		if v, ok := osInterface[field]; ok {
			if field == "id" || field == "links" {
				continue
			}
			fmt.Println("KV Pair: ", field, val)
			if val == nil {
				continue
			}
			if v != val {
				return false
			}
		}
	}
	return true
}

func mapSpecToOptions(spec, options interface{}) interface{} {
	d, _ := json.Marshal(spec)
	json.Unmarshal(d, &options)
	return options
}
