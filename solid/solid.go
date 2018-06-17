// Ground Control: free and automated network scanner
// Copyright (C) 2018 Alessio Sanfratello
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

package solid

import (
	"encoding/json"
	"io/ioutil"
)

// Cache yet another dbms
type Cache struct {
	Name     string            `json:"name"`
	Data     map[string]string `json:"data"`
	Metadata map[string]string `json:"metadata"`
}

// MakeSolidCache constructor
func MakeSolidCache(name string) Cache {
	return Cache{
		Name:     name,
		Data:     make(map[string]string),
		Metadata: make(map[string]string),
	}
}

// Put tbd
func (c *Cache) Put(key string, value interface{}) {
	marshaledValue, _ := json.Marshal(value)
	c.Data[key] = string(marshaledValue)
	c.store()
}

// Get tbd
func (c *Cache) Get(key string, value interface{}) {
	v := c.Data[key]
	json.Unmarshal([]byte(v), &value)
}

func (c Cache) store() {
	marshaledValue, _ := json.Marshal(c)
	ioutil.WriteFile(c.Name+".solid", marshaledValue, 0644)
}