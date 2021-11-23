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

type Reward struct {
	Funder       eos.AccountName `json:"funder"`
	Reward       string          `json:"reward"`
	FundingState eos.Name        `json:"funding_state"`
}

func (m *Reward) GetReward() eos.Asset {
	reward, err := eos.NewAssetFromString(m.Reward)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse reward: %v to asset", m.Reward))
	}
	return reward
}

func (m *Reward) Clone() *Reward {
	return &Reward{
		Funder:       m.Funder,
		FundingState: m.FundingState,
		Reward:       m.Reward,
	}
}

type RewardEntry struct {
	Key   string  `json:"key"`
	Value *Reward `json:"value"`
}

func (m *RewardEntry) Clone() *RewardEntry {
	return &RewardEntry{
		Key:   m.Key,
		Value: m.Value.Clone(),
	}
}

type RewardArgEntry struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type RewardsArg []*RewardArgEntry

func (m *RewardEntry) ToRewardArgEntry() *RewardArgEntry {
	return &RewardArgEntry{
		Key:   m.Key,
		Value: m.Value.Reward,
	}
}

type Rewards []*RewardEntry

func (m Rewards) ToRewardsArg() RewardsArg {
	rewardsArg := make(RewardsArg, 0, len(m))
	for _, rewardEntry := range m {
		rewardsArg = append(rewardsArg, rewardEntry.ToRewardArgEntry())
	}
	return rewardsArg
}

func (m Rewards) FindPos(key string) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m Rewards) Find(key string) *RewardEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (p *Rewards) Upsert(key string, reward *Reward) {
	m := *p
	pos := m.FindPos(key)
	entry := &RewardEntry{
		Key:   key,
		Value: reward,
	}
	if pos >= 0 {
		m[pos] = entry
	} else {
		m = append(m, entry)
	}
	*p = m
}

func (p *Rewards) Remove(key string) *RewardEntry {
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

func (m Rewards) UpdateFundingStateAll(state eos.Name) {
	for _, def := range m {
		r := def.Value
		if r.FundingState != RewardFundingStateRex {
			r.FundingState = state
		}
	}
}

func (m Rewards) UpdateFundingState(dist string, state eos.Name) {
	m.Find(dist).Value.FundingState = state
}

func (m Rewards) Clone() Rewards {
	rewards := make(Rewards, 0, len(m))
	for _, re := range m {
		rewards = append(rewards, re.Clone())
	}
	return rewards
}
