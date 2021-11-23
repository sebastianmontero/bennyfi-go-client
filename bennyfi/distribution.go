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

type Distribution struct {
	AmountToWinner       string `json:"amount_to_winner"`
	BeneficiaryReward    string `json:"beneficiary_reward"`
	RoundManagerFee      string `json:"round_manager_fee"`
	MinParticipantReward string `json:"min_participant_reward"`
}

func (m *Distribution) GetAmountToWinner() eos.Asset {
	amountToWinner, err := eos.NewAssetFromString(m.AmountToWinner)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse amount to winner: %v to asset", m.AmountToWinner))
	}
	return amountToWinner
}

func (m *Distribution) GetBeneficiaryReward() eos.Asset {
	beneficiaryReward, err := eos.NewAssetFromString(m.BeneficiaryReward)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse beneficiary reward: %v to asset", m.BeneficiaryReward))
	}
	return beneficiaryReward
}

func (m *Distribution) GetRoundManagerFee() eos.Asset {
	roundManagerFee, err := eos.NewAssetFromString(m.RoundManagerFee)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse round manager fee: %v to asset", m.RoundManagerFee))
	}
	return roundManagerFee
}

func (m *Distribution) GetMinParticipantReward() eos.Asset {
	minParticipantReward, err := eos.NewAssetFromString(m.MinParticipantReward)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse min participant reward: %v to asset", m.MinParticipantReward))
	}
	return minParticipantReward
}

type DistributionEntry struct {
	Key   string        `json:"key"`
	Value *Distribution `json:"value"`
}

type Distributions []*DistributionEntry

func (m Distributions) FindPos(key string) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m Distributions) Find(key string) *DistributionEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (p *Distributions) Upsert(key string, distribution *Distribution) {
	m := *p
	pos := m.FindPos(key)
	entry := &DistributionEntry{
		Key:   key,
		Value: distribution,
	}
	if pos >= 0 {
		m[pos] = entry
	} else {
		m = append(m, entry)
	}
	*p = m
}

func (p *Distributions) Remove(key string) *DistributionEntry {
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
