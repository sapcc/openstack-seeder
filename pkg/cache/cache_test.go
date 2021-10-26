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

package cache

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDefaultExpiration(t *testing.T) {
	c := New(2*time.Hour, 4*time.Hour)
	c.Add("entity", "key", "value", 0)
	c.Add("entity", "key2", "value2", 1*time.Hour)

	t1 := time.Unix(0, c.entities["entity"]["key"].Expiration)
	t2 := time.Unix(0, c.entities["entity"]["key2"].Expiration)

	assert.Equal(t, 2*time.Hour, c.defaultExpiration)
	assert.Equal(t, time.Now().Add(2*time.Hour).Hour(), t1.Hour())
	assert.Equal(t, time.Now().Add(1*time.Hour).Hour(), t2.Hour())
}

func TestAddGet(t *testing.T) {
	c := New(4*time.Second, 5*time.Second)
	c.Add("entity", "key1", "value1", 0)
	c.Add("entity", "key2", "value2", 2*time.Second)

	v, ok := c.Get("entity", "key1")
	assert.True(t, ok)
	assert.Equal(t, "value1", v)

	v, ok = c.Get("entity", "key2")
	assert.True(t, ok)
	assert.Equal(t, "value2", v)
}

func TestExpiration(t *testing.T) {
	c := New(2*time.Hour, 1*time.Second)
	c.Add("entity", "key", "value", 2*time.Second)
	c.Add("entity", "key2", "value2", 3*time.Second)

	v1, ok := c.Get("entity", "key")
	assert.True(t, ok)
	assert.Equal(t, v1, "value")
	time.Sleep(3 * time.Second)
	v2, ok := c.Get("entity", "key2")
	assert.False(t, ok)
	assert.Equal(t, v2, "")

}
