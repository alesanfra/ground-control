package cache

import (
	"encoding/json"
	"io/ioutil"
)

// SolidCache yet another dbms
type SolidCache struct {
	name     string
	data     map[string]string
	metadata map[string]string
}

// MakeSolidCache constructor
func MakeSolidCache(name string) SolidCache {
	c := SolidCache{name: name}
	c.data = make(map[string]string)
	c.metadata = make(map[string]string)
	return c
}

// Put tbd
func (c *SolidCache) Put(key string, value interface{}) {
	marshaledValue, _ := json.Marshal(value)
	c.data[key] = string(marshaledValue)
	c.store()
}

// Get tbd
func (c *SolidCache) Get(key string, value interface{}) {
	v := c.data[key]
	json.Unmarshal([]byte(v), &value)
}

func (c SolidCache) store() {
	marshaledValue, _ := json.Marshal(c)
	ioutil.WriteFile(c.name+".solid", marshaledValue, 0644)
}
