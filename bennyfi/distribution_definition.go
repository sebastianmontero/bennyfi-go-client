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
	"math"

	"github.com/sebastianmontero/bennyfi-go-client/util/utype"
	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/err"
	"github.com/sebastianmontero/eos-go-toolbox/util"
)

var (
	DistributionMainToken    = eos.Name("maintoken")
	DistributionMainNFT      = eos.Name("mainnft")
	DistributionProjectToken = eos.Name("projecttoken")
	DistributionProjectNFT   = eos.Name("projectnft")
	OrderedDistributionNames = []eos.Name{DistributionMainNFT, DistributionMainToken, DistributionProjectNFT, DistributionProjectToken}
)

func IsFTDistribution(distName eos.Name) bool {
	return distName == DistributionMainToken || distName == DistributionProjectToken
}

func IsNFTDistribution(distName eos.Name) bool {
	return distName == DistributionProjectNFT || distName == DistributionMainNFT
}

type IDistributionDefinition interface {
	GetNumWinners() int
	HasBeneficiaryReward() bool
	HasVesting() bool
	GetVesting() *VestingConfig
}

type BaseDistributionDefinition struct {
	*VestingConfig `json:"vesting_config"`
}

type BaseDistributionDefinitionCustomJSON struct {
	VestingConfig map[string]interface{} `json:"vesting_config"`
}

func (m BaseDistributionDefinition) ToCustomJSON() BaseDistributionDefinitionCustomJSON {
	return BaseDistributionDefinitionCustomJSON{
		VestingConfig: m.VestingConfig.ToMap(),
	}
}

func (m *BaseDistributionDefinition) HasVesting() bool {
	return m.VestingConfig.HasVesting()
}

func (m *BaseDistributionDefinition) GetVesting() *VestingConfig {
	return m.VestingConfig
}

type DistributionDefinitionFT struct {
	AllParticipantsPerc uint32   `json:"all_participants_perc_x100000"`
	BeneficiaryPerc     uint32   `json:"beneficiary_perc_x100000"`
	RoundManagerPerc    uint32   `json:"round_manager_perc_x100000"`
	WinnersPerc         []uint32 `json:"winners_perc_x100000"`
	*BaseDistributionDefinition
	Reward eos.Asset `json:"reward"`
}

type DistributionDefinitionFTCustomJSON struct {
	BaseDistributionDefinitionCustomJSON
	DistributionDefinitionFT
}

func (m DistributionDefinitionFT) ToCustomJSON() DistributionDefinitionFTCustomJSON {
	return DistributionDefinitionFTCustomJSON{
		BaseDistributionDefinitionCustomJSON: m.BaseDistributionDefinition.ToCustomJSON(),
		DistributionDefinitionFT:             m,
	}
}

func (m *DistributionDefinitionFT) GetNumWinners() int {
	return len(m.WinnersPerc)
}

func (m *DistributionDefinitionFT) HasBeneficiaryReward() bool {
	return m.BeneficiaryPerc > 0
}

func (m *DistributionDefinitionFT) CalculateDistribution(numParticipantsEntered uint32, totalReward eos.Asset) *Distribution {

	precisionAdj := math.Pow(10, float64(totalReward.Precision))
	percAdj := float64(10000000)
	reward := float64(totalReward.Amount) / precisionAdj
	rewardToAllParticipants := reward * float64((float64(m.AllParticipantsPerc) / percAdj))
	rewardToBeneficiary := reward * float64((float64(m.BeneficiaryPerc) / percAdj))
	feeToManager := reward * float64((float64(m.RoundManagerPerc) / percAdj))
	minParticipantReward := eos.Asset{Amount: eos.Int64((rewardToAllParticipants / float64(numParticipantsEntered)) * float64(precisionAdj)), Symbol: totalReward.Symbol}
	beneficiaryReward := eos.Asset{Amount: eos.Int64(rewardToBeneficiary * float64(precisionAdj)), Symbol: totalReward.Symbol}
	managerFee := eos.Asset{Amount: eos.Int64(feeToManager * float64(precisionAdj)), Symbol: totalReward.Symbol}
	remaining := totalReward.Sub(beneficiaryReward).Sub(managerFee).Sub(util.MultiplyAsset(minParticipantReward, int64(numParticipantsEntered)))
	winnerPrizes := make([]eos.Asset, 0)
	if m.GetNumWinners() > 0 {
		winnerPrize := remaining
		winnerPrizeAdj := float64(winnerPrize.Amount) / precisionAdj

		for _, winnerPerc := range m.WinnersPerc {
			prizeAmount := winnerPrizeAdj * float64((float64(winnerPerc) / percAdj))
			prize := eos.Asset{Amount: eos.Int64(prizeAmount * float64(precisionAdj)), Symbol: totalReward.Symbol}
			remaining = remaining.Sub(prize)
			winnerPrizes = append(winnerPrizes, prize)
		}
		firstPrize := winnerPrizes[0]
		winnerPrizes[0] = firstPrize.Add(remaining)
	} else {
		if m.BeneficiaryPerc > 0 {
			beneficiaryReward = beneficiaryReward.Add(remaining)
		} else {
			managerFee = managerFee.Add(remaining)
		}
	}
	return NewDistribution(&DistributionFT{
		WinnerPrizes:          winnerPrizes,
		BeneficiaryReward:     beneficiaryReward,
		BeneficiaryRewardPaid: eos.Asset{Amount: 0, Symbol: totalReward.Symbol},
		MinParticipantReward:  minParticipantReward,
		RoundManagerFee:       managerFee,
		RoundManagerFeePaid:   eos.Asset{Amount: 0, Symbol: totalReward.Symbol},
	})
}

type DistributionDefinitionNFT struct {
	EachParticipantReward uint16   `json:"each_participant_reward"`
	BeneficiaryReward     uint16   `json:"beneficiary_reward"`
	RoundManagerFee       uint16   `json:"round_manager_fee"`
	WinnerPrizes          []uint16 `json:"winner_prizes"`
	*BaseDistributionDefinition
	NFTConfig *NFTConfig `json:"nft_config"`
}

type DistributionDefinitionNFTCustomJSON struct {
	BaseDistributionDefinitionCustomJSON
	DistributionDefinitionNFT
}

func (m DistributionDefinitionNFT) ToCustomJSON() DistributionDefinitionNFTCustomJSON {
	return DistributionDefinitionNFTCustomJSON{
		BaseDistributionDefinitionCustomJSON: m.BaseDistributionDefinition.ToCustomJSON(),
		DistributionDefinitionNFT:            m,
	}
}

func (m *DistributionDefinitionNFT) CalculateDistribution() *Distribution {
	return NewDistribution(
		&DistributionNFT{
			MinParticipantReward:  m.EachParticipantReward,
			BeneficiaryReward:     m.BeneficiaryReward,
			BeneficiaryRewardPaid: 0,
			RoundManagerFee:       m.RoundManagerFee,
			RoundManagerFeePaid:   0,
			WinnerPrizes:          m.WinnerPrizes,
		},
	)
}

func (m *DistributionDefinitionNFT) TotalReward(numParticipants uint32) uint32 {
	total := uint32(m.EachParticipantReward)*numParticipants + uint32(m.BeneficiaryReward) + uint32(m.RoundManagerFee)
	for _, prize := range m.WinnerPrizes {
		total += uint32(prize)
	}
	return total
}

func (m *DistributionDefinitionNFT) GetNumWinners() int {
	return len(m.WinnerPrizes)
}

func (m *DistributionDefinitionNFT) HasBeneficiaryReward() bool {
	return m.BeneficiaryReward > 0
}

var DistributionDefinitionVariant = eos.NewVariantDefinition([]eos.VariantType{
	{Name: "DistributionDefinitionFT", Type: &DistributionDefinitionFT{}},
	{Name: "DistributionDefinitionNFT", Type: &DistributionDefinitionNFT{}},
})

func GetDistributionDefinitionVariants() *eos.VariantDefinition {
	return DistributionDefinitionVariant
}

type DistributionDefinition struct {
	eos.BaseVariant
}

func (m *DistributionDefinition) GetNumWinners() int {
	return m.Impl.(IDistributionDefinition).GetNumWinners()
}

func (m *DistributionDefinition) HasBeneficiaryReward() bool {
	return m.Impl.(IDistributionDefinition).HasBeneficiaryReward()
}

func (m *DistributionDefinition) HasVesting() bool {
	return m.Impl.(IDistributionDefinition).HasVesting()
}

func (m *DistributionDefinition) GetVesting() *VestingConfig {
	return m.Impl.(IDistributionDefinition).GetVesting()
}

func NewDistributionDefinition(value interface{}) *DistributionDefinition {
	return &DistributionDefinition{
		BaseVariant: eos.BaseVariant{
			TypeID: GetDistributionDefinitionVariants().TypeID(utype.TypeName(value)),
			Impl:   value,
		}}
}

func (m *DistributionDefinition) DistributionDefinitionNFT() *DistributionDefinitionNFT {
	switch v := m.Impl.(type) {
	case *DistributionDefinitionNFT:
		return v
	default:
		panic(&err.InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for value: %v of variant %T", v, v, m),
			ExpectedType: "DistributionDefinitionNFT",
			Value:        m,
		})
	}
}

func (m *DistributionDefinition) DistributionDefinitionFT() *DistributionDefinitionFT {
	switch v := m.Impl.(type) {
	case *DistributionDefinitionFT:
		return v
	default:
		panic(&err.InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for value: %v of variant %T", v, v, m),
			ExpectedType: "DistributionDefinitionFT",
			Value:        m,
		})
	}
}

// MarshalJSON translates to []byte
func (m *DistributionDefinition) MarshalJSON() ([]byte, error) {
	return m.BaseVariant.MarshalJSON(DistributionDefinitionVariant)
}

// UnmarshalJSON translates WinnerVariant
func (m *DistributionDefinition) UnmarshalJSON(data []byte) error {
	return m.BaseVariant.UnmarshalJSON(data, DistributionDefinitionVariant)
}

// UnmarshalBinary ...
func (m *DistributionDefinition) UnmarshalBinary(decoder *eos.Decoder) error {
	return m.BaseVariant.UnmarshalBinaryVariant(decoder, DistributionDefinitionVariant)
}

type DistributionDefinitionEntry struct {
	Key   eos.Name                `json:"first"`
	Value *DistributionDefinition `json:"second"`
}

type DistributionDefinitions []*DistributionDefinitionEntry

func (m DistributionDefinitions) ToMap() map[eos.Name]interface{} {
	distDefMap := make(map[eos.Name]interface{})
	for _, distDefEntry := range m {
		var distDef interface{}
		if IsFTDistribution(distDefEntry.Key) {
			distDef = distDefEntry.Value.DistributionDefinitionFT().ToCustomJSON()
		} else {
			distDef = distDefEntry.Value.DistributionDefinitionNFT().ToCustomJSON()
		}
		distDefMap[distDefEntry.Key] = distDef
	}
	return distDefMap
}

func (m DistributionDefinitions) FindPos(key eos.Name) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m DistributionDefinitions) Has(key eos.Name) bool {
	return m.FindPos(key) >= 0
}

func (m DistributionDefinitions) Find(key eos.Name) *DistributionDefinitionEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (m DistributionDefinitions) FindFT(key eos.Name) *DistributionDefinitionFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.DistributionDefinitionFT()
	}
	return nil
}

func (m DistributionDefinitions) FindNFT(key eos.Name) *DistributionDefinitionNFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.DistributionDefinitionNFT()
	}
	return nil
}

// Used during testing to make sure dist definitions are ordered as they are oredered in the c++ map
func (m DistributionDefinitions) GetOrderedDistDefs() DistributionDefinitions {
	dists := make(DistributionDefinitions, 0)
	for _, name := range OrderedDistributionNames {
		dist := m.Find(name)
		if dist != nil {
			dists = append(dists, dist)
		}
	}
	return dists
}

func (p *DistributionDefinitions) Upsert(key eos.Name, definition interface{}) {
	m := *p
	pos := m.FindPos(key)
	defEntry := &DistributionDefinitionEntry{
		Key:   key,
		Value: NewDistributionDefinition(definition),
	}
	if pos >= 0 {
		m[pos] = defEntry
	} else {
		m = append(m, defEntry)
	}
	*p = m
}

func (p *DistributionDefinitions) Remove(key eos.Name) *DistributionDefinitionEntry {
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

func (m DistributionDefinitions) GetNumWinners(key eos.Name) int {

	entry := m.Find(key)
	if entry == nil {
		panic(fmt.Sprintf("there is no distribution definition for key: %v", key))
	}
	return entry.Value.GetNumWinners()
}

func (m DistributionDefinitions) GetTotalNumWinners() int {
	total := 0
	for _, entry := range m {
		total += entry.Value.GetNumWinners()
	}
	return total
}

func (m DistributionDefinitions) HasVesting() bool {
	for _, entry := range m {
		if entry.Value.HasVesting() {
			return true
		}
	}
	return false
}

func (m DistributionDefinitions) GetDistributionDefinitionsFT() []*DistributionDefinitionFT {
	distDefsFT := make([]*DistributionDefinitionFT, 0)
	for _, distDefEntry := range m {
		switch v := distDefEntry.Value.Impl.(type) {
		case *DistributionDefinitionFT:
			distDefsFT = append(distDefsFT, v)
		}
	}
	return distDefsFT
}

func (m DistributionDefinitions) GetVestingTrackers() VestingTrackers {
	trackers := make(VestingTrackers, 0)
	for _, entry := range m {
		if entry.Value.HasVesting() {
			trackers = append(trackers, &VestingTracker{
				VestingConfig: entry.Value.GetVesting(),
				DistName:      entry.Key,
			})
		}
	}
	return trackers
}

func (m DistributionDefinitions) GetVestingContext(cycle uint16, startTime eos.TimePoint) *VestingContext {
	return m.GetVestingTrackers().GetContextForCycle(cycle, startTime)
}
