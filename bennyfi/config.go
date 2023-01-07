// The MIT License (MIT)

// Copyright (c) 2020, Digital Scarcity

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.
package bennyfi

import (
	"fmt"

	"github.com/sebastianmontero/eos-go-toolbox/dto"
)

type ConfigEntry struct {
	Key   string         `json:"first"`
	Value *dto.FlexValue `json:"second"`
}

type Config []*ConfigEntry

func (m Config) ToMap() map[string]interface{} {
	configMap := make(map[string]interface{})
	for _, configEntry := range m {
		configMap[configEntry.Key] = configEntry.Value.Impl
	}
	return configMap
}

func (m Config) HasConfig() bool {
	return len(m) > 0
}

func (m Config) FindPos(key string) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m Config) FindEntry(key string) *ConfigEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (m Config) Find(key string) *dto.FlexValue {
	entry := m.FindEntry(key)
	if entry != nil {
		return entry.Value
	}
	return nil
}

func (m Config) Get(key string) *dto.FlexValue {
	v := m.Find(key)
	if v == nil {
		panic(fmt.Sprintf("Config param with key: %v does not exist", key))
	}
	return v
}

func (p *Config) Upsert(key string, value *dto.FlexValue) {
	m := *p
	pos := m.FindPos(key)
	entry := &ConfigEntry{
		Key:   key,
		Value: value,
	}
	if pos >= 0 {
		m[pos] = entry
	} else {
		m = append(m, entry)
	}
	*p = m
}

func (p *Config) Remove(key string) *ConfigEntry {
	m := *p
	pos := m.FindPos(key)
	if pos >= 0 {
		def := m[pos]
		m[pos] = m[len(m)-1]
		m = m[:len(m)-1]
		*p = m
		return def
	}
	return nil
}
