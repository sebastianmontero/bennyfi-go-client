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
	"github.com/sebastianmontero/bennyfi-go-client/util/utype"
	"github.com/sebastianmontero/eos-go-toolbox/err"
)

type IDistribution interface {
	PaidTotalBeneficiaryReward()
	PaidBeneficiaryReward(amount interface{})
	PaidTotalRoundManagerFee()
	PaidRoundManagerFee(amount interface{})
	Paid()
	HasBeneficiaryReward() bool
	HasRoundManagerFee() bool
}

type DistributionFT struct {
	BeneficiaryReward     string   `json:"beneficiary_reward"`
	BeneficiaryRewardPaid string   `json:"beneficiary_reward_paid"`
	RoundManagerFee       string   `json:"round_manager_fee"`
	RoundManagerFeePaid   string   `json:"round_manager_fee_paid"`
	MinParticipantReward  string   `json:"min_participant_reward"`
	WinnerPrizes          []string `json:"winner_prizes"`
}

func (m *DistributionFT) HasBeneficiaryReward() bool {
	return m.GetBeneficiaryReward().Amount > 0
}

func (m *DistributionFT) HasRoundManagerFee() bool {
	return m.GetRoundManagerFee().Amount > 0
}

func (m *DistributionFT) GetBeneficiaryReward() eos.Asset {
	beneficiaryReward, err := eos.NewAssetFromString(m.BeneficiaryReward)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse beneficiary reward: %v to asset", m.BeneficiaryReward))
	}
	return beneficiaryReward
}

func (m *DistributionFT) GetBeneficiaryRewardPaid() eos.Asset {
	beneficiaryRewardPaid, err := eos.NewAssetFromString(m.BeneficiaryRewardPaid)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse beneficiary reward paid: %v to asset", m.BeneficiaryRewardPaid))
	}
	return beneficiaryRewardPaid
}

func (m *DistributionFT) PaidTotalBeneficiaryReward() {
	m.BeneficiaryRewardPaid = m.BeneficiaryReward
}

func (m *DistributionFT) PaidBeneficiaryReward(amount interface{}) {
	amnt := amount.(eos.Asset)
	paid := m.GetBeneficiaryRewardPaid().Add(amnt)
	if paid.Amount > m.GetBeneficiaryReward().Amount {
		panic(fmt.Sprintf("Total Paid amount: %v is greater than beneficiary reward: %v, current payment: %v", paid, m.BeneficiaryReward, amount))
	}
	m.BeneficiaryRewardPaid = paid.String()
}

func (m *DistributionFT) PaidTotalRoundManagerFee() {
	m.RoundManagerFeePaid = m.RoundManagerFee
}

func (m *DistributionFT) PaidRoundManagerFee(amount interface{}) {
	amnt := amount.(eos.Asset)
	paid := m.GetRoundManagerFeePaid().Add(amnt)
	if paid.Amount > m.GetRoundManagerFee().Amount {
		panic(fmt.Sprintf("Total Paid amount: %v is greater than round manager fee: %v, current payment: %v", paid, m.RoundManagerFee, amount))
	}
	m.RoundManagerFeePaid = paid.String()
}

func (m *DistributionFT) Paid() {
	m.PaidTotalBeneficiaryReward()
	m.PaidTotalRoundManagerFee()
}

func (m *DistributionFT) GetRoundManagerFee() eos.Asset {
	roundManagerFee, err := eos.NewAssetFromString(m.RoundManagerFee)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse round manager fee: %v to asset", m.RoundManagerFee))
	}
	return roundManagerFee
}

func (m *DistributionFT) GetRoundManagerFeePaid() eos.Asset {
	roundManagerFeePaid, err := eos.NewAssetFromString(m.RoundManagerFeePaid)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse round manager fee paid: %v to asset", m.RoundManagerFeePaid))
	}
	return roundManagerFeePaid
}

func (m *DistributionFT) GetMinParticipantReward() eos.Asset {
	minParticipantReward, err := eos.NewAssetFromString(m.MinParticipantReward)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse min participant reward: %v to asset", m.MinParticipantReward))
	}
	return minParticipantReward
}

func (m *DistributionFT) GetNumWinners() uint32 {
	return uint32(len(m.WinnerPrizes))
}

func (m *DistributionFT) GetWinnerPrize(pos uint32) eos.Asset {
	if pos >= m.GetNumWinners() {
		panic(fmt.Sprintf("There is no winner for pos: %v number of winners are: %v", pos, m.GetNumWinners()))
	}
	winnerPrize, err := eos.NewAssetFromString(m.WinnerPrizes[pos])
	if err != nil {
		panic(fmt.Sprintf("Unable to parse winner prize: %v to asset", m.WinnerPrizes[pos]))
	}
	return winnerPrize
}

type DistributionNFT struct {
	BeneficiaryReward     uint16   `json:"beneficiary_reward"`
	BeneficiaryRewardPaid uint16   `json:"beneficiary_reward_paid"`
	RoundManagerFee       uint16   `json:"round_manager_fee"`
	RoundManagerFeePaid   uint16   `json:"round_manager_fee_paid"`
	MinParticipantReward  uint16   `json:"each_participant_reward"`
	WinnerPrizes          []uint16 `json:"winner_prizes"`
}

func (m *DistributionNFT) HasBeneficiaryReward() bool {
	return m.BeneficiaryReward > 0
}

func (m *DistributionNFT) HasRoundManagerFee() bool {
	return m.RoundManagerFee > 0
}

func (m *DistributionNFT) PaidTotalBeneficiaryReward() {
	m.BeneficiaryRewardPaid = m.BeneficiaryReward
}

func (m *DistributionNFT) PaidBeneficiaryReward(amount interface{}) {
	amnt := amount.(uint16)
	paid := m.BeneficiaryRewardPaid + amnt
	if paid > m.BeneficiaryReward {
		panic(fmt.Sprintf("Total Paid amount: %v is greater than beneficiary reward: %v, current payment: %v", paid, m.BeneficiaryReward, amount))
	}
	m.BeneficiaryRewardPaid = paid
}

func (m *DistributionNFT) PaidTotalRoundManagerFee() {
	m.RoundManagerFeePaid = m.RoundManagerFee
}

func (m *DistributionNFT) PaidRoundManagerFee(amount interface{}) {
	amnt := amount.(uint16)
	paid := m.RoundManagerFeePaid + amnt
	if paid > m.RoundManagerFee {
		panic(fmt.Sprintf("Total Paid amount: %v is greater than round manager fee: %v, current payment: %v", paid, m.RoundManagerFee, amount))
	}
	m.RoundManagerFeePaid = paid
}

func (m *DistributionNFT) Paid() {
	m.PaidTotalBeneficiaryReward()
	m.PaidTotalRoundManagerFee()
}

var DistributionVariant = eos.NewVariantDefinition([]eos.VariantType{
	{Name: "DistributionFT", Type: &DistributionFT{}},
	{Name: "DistributionNFT", Type: &DistributionNFT{}},
})

func GetDistributionVariants() *eos.VariantDefinition {
	return DistributionVariant
}

type Distribution struct {
	eos.BaseVariant
}

func NewDistribution(value interface{}) *Distribution {
	return &Distribution{
		BaseVariant: eos.BaseVariant{
			TypeID: GetDistributionVariants().TypeID(utype.TypeName(value)),
			Impl:   value,
		}}
}

func (m *Distribution) PaidTotalBeneficiaryReward() {
	m.Impl.(IDistribution).PaidTotalBeneficiaryReward()
}

func (m *Distribution) PaidBeneficiaryReward(amount interface{}) {
	m.Impl.(IDistribution).PaidBeneficiaryReward(amount)
}

func (m *Distribution) PaidTotalRoundManagerFee() {
	m.Impl.(IDistribution).PaidTotalRoundManagerFee()
}

func (m *Distribution) PaidRoundManagerFee(amount interface{}) {
	m.Impl.(IDistribution).PaidRoundManagerFee(amount)
}

func (m *Distribution) Paid() {
	m.Impl.(IDistribution).Paid()
}

func (m *Distribution) HasBeneficiaryReward() bool {
	return m.Impl.(IDistribution).HasBeneficiaryReward()
}

func (m *Distribution) HasRoundManagerFee() bool {
	return m.Impl.(IDistribution).HasRoundManagerFee()
}

func (m *Distribution) DistributionNFT() *DistributionNFT {
	switch v := m.Impl.(type) {
	case *DistributionNFT:
		return v
	default:
		panic(&err.InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for value: %v of variant %T", v, v, m),
			ExpectedType: "DistributionNFT",
			Value:        m,
		})
	}

}

func (m *Distribution) DistributionFT() *DistributionFT {
	switch v := m.Impl.(type) {
	case *DistributionFT:
		return v
	default:
		panic(&err.InvalidTypeError{
			Label:        fmt.Sprintf("received an unexpected type %T for value: %v of variant %T", v, v, m),
			ExpectedType: "DistributionFT",
			Value:        m,
		})
	}

}

// MarshalJSON translates to []byte
func (m *Distribution) MarshalJSON() ([]byte, error) {
	return m.BaseVariant.MarshalJSON(DistributionVariant)
}

// UnmarshalJSON translates WinnerVariant
func (m *Distribution) UnmarshalJSON(data []byte) error {
	return m.BaseVariant.UnmarshalJSON(data, DistributionVariant)
}

// UnmarshalBinary ...
func (m *Distribution) UnmarshalBinary(decoder *eos.Decoder) error {
	return m.BaseVariant.UnmarshalBinaryVariant(decoder, DistributionVariant)
}

type DistributionEntry struct {
	Key   eos.Name      `json:"first"`
	Value *Distribution `json:"second"`
}

type Distributions []*DistributionEntry

func (m Distributions) ToMap() map[eos.Name]interface{} {
	distMap := make(map[eos.Name]interface{})
	for _, distEntry := range m {
		distMap[distEntry.Key] = distEntry.Value.Impl
	}
	return distMap
}

func (m Distributions) FindPos(key eos.Name) int {
	for i, def := range m {
		if def.Key == key {
			return i
		}
	}
	return -1
}

func (m Distributions) Find(key eos.Name) *DistributionEntry {
	pos := m.FindPos(key)
	if pos >= 0 {
		return m[pos]
	}
	return nil
}

func (m Distributions) FindFT(key eos.Name) *DistributionFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.DistributionFT()
	}
	return nil
}

func (m Distributions) FindNFT(key eos.Name) *DistributionNFT {
	v := m.Find(key)
	if v != nil {
		return v.Value.DistributionNFT()
	}
	return nil
}

func (p *Distributions) Upsert(key eos.Name, distribution interface{}) {
	m := *p
	pos := m.FindPos(key)
	var dist *Distribution
	switch v := distribution.(type) {
	case *Distribution:
		dist = v
	default:
		dist = NewDistribution(distribution)
	}
	entry := &DistributionEntry{
		Key:   key,
		Value: dist,
	}
	if pos >= 0 {
		m[pos] = entry
	} else {
		m = append(m, entry)
	}
	*p = m
}

func (p *Distributions) Remove(key eos.Name) *DistributionEntry {
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

func (m Distributions) Paid() {
	for _, d := range m {
		d.Value.Paid()
	}
}

func (m Distributions) PaidNonVesting(distDefs DistributionDefinitions) {
	for _, d := range m {
		if !distDefs.Find(d.Key).Value.Impl.(IDistributionDefinition).HasVesting() {
			d.Value.Paid()
		}

	}
}
