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

	"github.com/sebastianmontero/bennyfi-go-client/util/utype"
	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/err"
)

type IReward interface {
	Clone() interface{}
	GetFundingState() eos.Name
	SetFundingState(state eos.Name)
	GetFunder() eos.AccountName
}

type BaseReward struct {
	Funder       eos.AccountName `json:"funder"`
	FundingState eos.Name        `json:"funding_state"`
}

func (m *BaseReward) GetFundingState() eos.Name {
	return m.FundingState
}

func (m *BaseReward) SetFundingState(state eos.Name) {
	m.FundingState = state
}

func (m *BaseReward) GetFunder() eos.AccountName {
	return m.Funder
}

func (m *BaseReward) Clone() *BaseReward {
	return &BaseReward{
		Funder:       m.Funder,
		FundingState: m.FundingState,
	}
}

type RewardFT struct {
	*BaseReward
	Reward eos.Asset `json:"reward"`
}

func (m *RewardFT) Clone() interface{} {
	return &RewardFT{
		BaseReward: m.BaseReward.Clone(),
		Reward:     m.Reward,
	}
}

type RewardNFT struct {
	*BaseReward
	Reward uint32 `json:"reward"`
}

func (m *RewardNFT) Clone() interface{} {
	return &RewardNFT{
		BaseReward: m.BaseReward.Clone(),
		Reward:     m.Reward,
	}
}

var RewardVariant = eos.NewVariantDefinition([]eos.VariantType{
	{Name: "RewardFT", Type: &RewardFT{}},
	{Name: "RewardNFT", Type: &RewardNFT{}},
})

func GetRewardVariants() *eos.VariantDefinition {
	return RewardVariant
}

type Reward struct {
	eos.BaseVariant
}

func NewReward(value interface{}) *Reward {
	return &Reward{
		BaseVariant: eos.BaseVariant{
			TypeID: GetRewardVariants().TypeID(utype.TypeName(value)),
			Impl:   value,
		}}
}

func (m *Reward) Clone() *Reward {
	return NewReward(m.Impl.(IReward).Clone())
}

func (m *Reward) SetFundingState(state eos.Name) {
	m.Impl.(IReward).SetFundingState(state)
}

func (m *Reward) GetFundingState() eos.Name {
	return m.Impl.(IReward).GetFundingState()
}

func (m *Reward) GetFunder() eos.AccountName {
	return m.Impl.(IReward).GetFunder()
}

func (m *Reward) RewardNFT() *RewardNFT {
	switch v := m.Impl.(type) {
	case *RewardNFT:
		return v
	default:
		panic(&err.InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for value: %v of variant %T", v, v, m),
			ExpectedType: "RewardNFT",
			Value:        m,
		})
	}
}

func (m *Reward) RewardFT() *RewardFT {
	switch v := m.Impl.(type) {
	case *RewardFT:
		return v
	default:
		panic(&err.InvalidTypeError{
			Label:        fmt.Sprintf("received1 an unexpected type %T for value: %v of variant %T", v, v, m),
			ExpectedType: "RewardFT",
			Value:        m,
		})
	}
}

// MarshalJSON translates to []byte
func (m *Reward) MarshalJSON() ([]byte, error) {
	return m.BaseVariant.MarshalJSON(RewardVariant)
}

// UnmarshalJSON translates WinnerVariant
func (m *Reward) UnmarshalJSON(data []byte) error {
	return m.BaseVariant.UnmarshalJSON(data, RewardVariant)
}

// UnmarshalBinary ...
func (m *Reward) UnmarshalBinary(decoder *eos.Decoder) error {
	return m.BaseVariant.UnmarshalBinaryVariant(decoder, RewardVariant)
}

type RewardEntry struct {
	Key   eos.Name `json:"first"`
	Value *Reward  `json:"second"`
}

func (m *RewardEntry) Clone() *RewardEntry {
	return &RewardEntry{
		Key:   m.Key,
		Value: m.Value.Clone(),
	}
}

func (m *RewardEntry) AsRewardFT() *RewardFT {
	return m.Value.Impl.(*RewardFT)
}

func (m *RewardEntry) AsRewardNFT() *RewardNFT {
	return m.Value.Impl.(*RewardNFT)
}

type FTRewardArgEntry struct {
	Key   eos.Name  `json:"first"`
	Value eos.Asset `json:"second"`
}

func (m *FTRewardArgEntry) String() string {
	return fmt.Sprintf(`
		FTRewardArgEntry {
			Key: %v,
			Value: %v,
		}
	`,
		m.Key,
		m.Value)
}

type FTRewardsArg []*FTRewardArgEntry

type Rewards []*RewardEntry

func (m Rewards) ToMap() map[eos.Name]interface{} {
	rewardMap := make(map[eos.Name]interface{})
	for _, rewardEntry := range m {
		rewardMap[rewardEntry.Key] = rewardEntry.Value.Impl
	}
	return rewardMap
}

func (m Rewards) ToFTRewardsArg() FTRewardsArg {
	rewardsArg := make(FTRewardsArg, 0)
	for _, rewardEntry := range m {
		switch v := rewardEntry.Value.Impl.(type) {
		case *RewardFT:
			rewardsArg = append(rewardsArg, &FTRewardArgEntry{
				Key:   rewardEntry.Key,
				Value: v.Reward,
			})
		}
	}
	return rewardsArg
}

func (m Rewards) GetRewardsFT() []*RewardFT {
	rewardsFT := make([]*RewardFT, 0)
	for _, rewardEntry := range m {
		switch v := rewardEntry.Value.Impl.(type) {
		case *RewardFT:
			rewardsFT = append(rewardsFT, v)
		}
	}
	return rewardsFT
}

func (m Rewards) GetFundedRewardsFT() []*RewardEntry {
	rewardsFT := make([]*RewardEntry, 0)
	for _, rewardEntry := range m {
		switch v := rewardEntry.Value.Impl.(type) {
		case *RewardFT:
			if v.FundingState == FundingStateFunded || v.FundingState == FundingStateCommited {
				rewardsFT = append(rewardsFT, rewardEntry)
			}
		}
	}
	return rewardsFT
}

func (m Rewards) GetFundedRewardsNFT() []*RewardEntry {
	rewardsNFT := make([]*RewardEntry, 0)
	for _, rewardEntry := range m {
		switch v := rewardEntry.Value.Impl.(type) {
		case *RewardNFT:
			if v.FundingState == FundingStateFunded || v.FundingState == FundingStateCommited {
				rewardsNFT = append(rewardsNFT, rewardEntry)
			}
		}
	}
	return rewardsNFT
}

func (m Rewards) FindPos(key eos.Name) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m Rewards) Find(key eos.Name) *RewardEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (m Rewards) Has(key eos.Name) bool {
	return m.Find(key) != nil
}

func (m Rewards) FindFT(key eos.Name) *RewardFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.RewardFT()
	}
	return nil
}

func (m Rewards) FindNFT(key eos.Name) *RewardNFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.RewardNFT()
	}
	return nil
}

func (p *Rewards) Upsert(key eos.Name, reward interface{}) {
	m := *p
	pos := m.FindPos(key)
	entry := &RewardEntry{
		Key:   key,
		Value: NewReward(reward),
	}
	if pos >= 0 {
		m[pos] = entry
	} else {
		m = append(m, entry)
	}
	*p = m
}

func (p *Rewards) Remove(key eos.Name) *RewardEntry {
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
		if r.GetFundingState() != FundingStateYield {
			r.SetFundingState(state)
		}
	}
}

func (m Rewards) IsFundingPending() bool {
	for _, def := range m {
		r := def.Value
		if r.GetFundingState() == FundingStatePending {
			return true
		}
	}
	return false
}

func (m Rewards) UpdateFundingState(dist eos.Name, state eos.Name) {
	m.Find(dist).Value.SetFundingState(state)
}

func (m Rewards) Clone() Rewards {
	rewards := make(Rewards, 0, len(m))
	for _, re := range m {
		rewards = append(rewards, re.Clone())
	}
	return rewards
}
