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

	eos "github.com/eoscanada/eos-go"
)

type TotalEarlyExitRewardFeeEntry struct {
	Key   eos.Name `json:"first"`
	Value string   `json:"second"`
}

func (m *TotalEarlyExitRewardFeeEntry) GetTotalEarlyExitRewardFee() eos.Asset {
	fee, err := eos.NewAssetFromString(m.Value)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse early exit reward fee: %v to asset", m.Value))
	}
	return fee
}

type TotalEarlyExitRewardFees []*TotalEarlyExitRewardFeeEntry

func (m TotalEarlyExitRewardFees) FindPos(key eos.Name) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m TotalEarlyExitRewardFees) Find(key eos.Name) *TotalEarlyExitRewardFeeEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (p *TotalEarlyExitRewardFees) Upsert(key eos.Name, fee string) {
	m := *p
	pos := m.FindPos(key)
	entry := &TotalEarlyExitRewardFeeEntry{
		Key:   key,
		Value: fee,
	}
	if pos >= 0 {
		m[pos] = entry
	} else {
		m = append(m, entry)
	}
	*p = m
}

func (p *TotalEarlyExitRewardFees) Remove(key eos.Name) *TotalEarlyExitRewardFeeEntry {
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
