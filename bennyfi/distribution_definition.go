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

type DistributionDefinition struct {
	AllParticipantsPerc uint32 `json:"all_participants_perc_x100000"`
	BeneficiaryPerc     uint32 `json:"beneficiary_perc_x100000"`
	RoundManagerPerc    uint32 `json:"round_manager_perc_x100000"`
}

type DistributionDefinitionEntry struct {
	Key   string                  `json:"key"`
	Value *DistributionDefinition `json:"value"`
}

type DistributionDefinitions []*DistributionDefinitionEntry

func (m DistributionDefinitions) FindPos(key string) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m DistributionDefinitions) Find(key string) *DistributionDefinitionEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (p *DistributionDefinitions) Upsert(key string, definition *DistributionDefinition) {
	m := *p
	pos := m.FindPos(key)
	defEntry := &DistributionDefinitionEntry{
		Key:   key,
		Value: definition,
	}
	if pos >= 0 {
		m[pos] = defEntry
	} else {
		m = append(m, defEntry)
	}
	*p = m
}

func (p *DistributionDefinitions) Remove(key string) *DistributionDefinitionEntry {
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
