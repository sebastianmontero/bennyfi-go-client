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

type Winner struct {
	Participant   eos.AccountName `json:"participant"`
	Prize         string          `json:"prize"`
	EntryPosition uint64          `json:"entry_position"`
}

func (m *Winner) GetPrize() eos.Asset {
	prize, err := eos.NewAssetFromString(m.Prize)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse prize: %v to asset", m.Prize))
	}
	return prize
}

func (m *Winner) IsWinner(account interface{}) bool {
	return fmt.Sprintf("%v", m.Participant) == fmt.Sprintf("%v", account)
}

type DistributionWinners []*Winner

func (m DistributionWinners) FindPos(account interface{}) int {
	for i, winner := range m {
		if winner.IsWinner(account) {
			return i
		}
	}
	return -1
}

func (m DistributionWinners) Find(account interface{}) *Winner {
	pos := m.FindPos(account)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

type DistributionWinnersEntry struct {
	Key   string              `json:"key"`
	Value DistributionWinners `json:"value"`
}

func NewWinner(participant eos.AccountName, prize string, entryPosition uint64) *Winner {
	return &Winner{
		Participant:   participant,
		Prize:         prize,
		EntryPosition: entryPosition,
	}
}

type Winners []*DistributionWinnersEntry

func (m Winners) FindPos(key string) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m Winners) Find(key string) *DistributionWinnersEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (p *Winners) Upsert(key string, winner *Winner) {
	m := *p
	pos := m.FindPos(key)
	if pos >= 0 {
		m[pos].Value = append(m[pos].Value, winner)
	} else {
		m = append(m, &DistributionWinnersEntry{
			Key: key,
			Value: DistributionWinners{
				winner,
			},
		})
	}
	*p = m
}

func (p *Winners) Remove(key string) *DistributionWinnersEntry {
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
