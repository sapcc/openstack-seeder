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
	"sync"
	"time"
)

type Cache struct {
	entities          map[string]map[string]Item
	mu                sync.RWMutex
	cleanupInterval   time.Duration
	defaultExpiration time.Duration
	Stop              chan bool
}

type Item struct {
	Expiration int64
	Value      string
}

func New(defaultExpiration, cleanupInterval time.Duration) (c *Cache) {
	e := make(map[string]map[string]Item)
	c = &Cache{
		cleanupInterval:   cleanupInterval,
		defaultExpiration: defaultExpiration,
		entities:          e,
	}
	go c.runCleanup()
	return
}

func (c *Cache) runCleanup() {
	ticker := time.NewTicker(c.cleanupInterval)
	for {
		select {
		case <-ticker.C:
			c.mu.Lock()
			now := time.Now().UnixNano()
			for e, i := range c.entities {
				for k, v := range i {
					if now > v.Expiration {
						delete(c.entities[e], k)
					}
				}

			}
			c.mu.Unlock()
		case <-c.Stop:
			ticker.Stop()
			return
		}
	}
}

func (c *Cache) Add(entity, key, value string, expire time.Duration) {
	c.mu.Lock()
	if _, ok := c.entities[entity]; !ok {
		c.entities[entity] = make(map[string]Item)
	}
	var exp int64
	if expire.Hours() != 0 {
		exp = time.Now().Add(expire).UnixNano()
	} else {
		exp = time.Now().Add(c.defaultExpiration).UnixNano()
	}
	c.entities[entity][key] = Item{
		Expiration: exp,
		Value:      value,
	}
	c.mu.Unlock()
}

func (c *Cache) Get(entity, key string) (string, bool) {
	if i, ok := c.entities[entity]; ok {
		if v, ok := i[key]; ok {
			return v.Value, true
		}
	}
	return "", false
}
