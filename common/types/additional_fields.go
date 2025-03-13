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
package types

import (
	"fmt"

	"github.com/sebastianmontero/eos-go-toolbox/dto"
)

type AdditionalField struct {
	Key   string         `json:"first"`
	Value *dto.FlexValue `json:"second"`
}

func (m AdditionalField) Clone() *AdditionalField {
	return &AdditionalField{
		Key:   m.Key,
		Value: m.Value,
	}
}

func (m *AdditionalField) String() string {
	return fmt.Sprintf("key: %v, value: %v", m.Key, m.Value)
}

type AdditionalFields []*AdditionalField

func (m AdditionalFields) FindPos(key string) int {
	for i, attr := range m {
		if attr.Key == key {
			return i
		}
	}
	return -1
}

func (m AdditionalFields) Has(key string) bool {
	return m.FindPos(key) >= 0
}

func (m AdditionalFields) Find(key string) *AdditionalField {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (m AdditionalFields) Get(key string) *AdditionalField {
	field := m.Find(key)
	if field == nil {
		panic(fmt.Sprintf("%v not found in additional fields", key))
	}
	return field
}

func (m AdditionalFields) GetValue(key string) *dto.FlexValue {
	field := m.Find(key)
	if field == nil {
		panic(fmt.Sprintf("%v not found in additional fields", key))
	}
	return field.Value
}

func (m AdditionalFields) Clone() AdditionalFields {
	clone := make(AdditionalFields, len(m))
	for i, attr := range m {
		clone[i] = attr.Clone()
	}
	return clone
}

func (m AdditionalFields) String() string {
	str := "["
	for _, field := range m {
		str += fmt.Sprintf("\n%v", field.String())
	}
	str += "]\n"
	return str
}

func (p *AdditionalFields) Set(key string, value *dto.FlexValue) {
	m := *p
	pos := m.FindPos(key)
	field := &AdditionalField{Key: key, Value: value}
	if pos >= 0 {
		m[pos] = field
	} else {
		m = append(m, field)
	}
	*p = m
}

func (p *AdditionalFields) Remove(key string) *AdditionalField {
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
