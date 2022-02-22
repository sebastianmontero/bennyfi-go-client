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

	eos "github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/util"
)

type DistributionDefinition struct {
	AllParticipantsPerc uint32   `json:"all_participants_perc_x100000"`
	BeneficiaryPerc     uint32   `json:"beneficiary_perc_x100000"`
	RoundManagerPerc    uint32   `json:"round_manager_perc_x100000"`
	WinnersPerc         []uint32 `json:"winners_perc_x100000"`
}

func (m *DistributionDefinition) GetNumWinners() int {
	return len(m.WinnersPerc)
}

func (m *DistributionDefinition) CalculateDistribution(totalReward eos.Asset, numParticipantsEntered uint32) *Distribution {

	precisionAdj := math.Pow(10, float64(totalReward.Precision))
	percAdj := float64(10000000)
	reward := float64(totalReward.Amount) / precisionAdj
	rewardToAllParticipants := reward * float64((float64(m.AllParticipantsPerc) / percAdj))
	rewardToBeneficiary := reward * float64((float64(m.BeneficiaryPerc) / percAdj))
	feeToManager := reward * float64((float64(m.RoundManagerPerc) / percAdj))
	minParticipantReward := eos.Asset{Amount: eos.Int64((rewardToAllParticipants / float64(numParticipantsEntered)) * float64(precisionAdj)), Symbol: totalReward.Symbol}
	beneficiaryReward := eos.Asset{Amount: eos.Int64(rewardToBeneficiary * float64(precisionAdj)), Symbol: totalReward.Symbol}
	managerFee := eos.Asset{Amount: eos.Int64(feeToManager * float64(precisionAdj)), Symbol: totalReward.Symbol}
	winnerPrize := totalReward.Sub(beneficiaryReward).Sub(managerFee).Sub(util.MultiplyAsset(minParticipantReward, int64(numParticipantsEntered)))
	remaining := winnerPrize
	winnerPrizeAdj := float64(winnerPrize.Amount) / precisionAdj
	winnerPrizes := make([]string, 0, m.GetNumWinners())
	for _, winnerPerc := range m.WinnersPerc {
		prizeAmount := winnerPrizeAdj * float64((float64(winnerPerc) / percAdj))
		prize := eos.Asset{Amount: eos.Int64(prizeAmount * float64(precisionAdj)), Symbol: totalReward.Symbol}
		remaining = remaining.Sub(prize)
		winnerPrizes = append(winnerPrizes, prize.String())
	}
	firstPrize, _ := eos.NewAssetFromString(winnerPrizes[0])
	winnerPrizes[0] = firstPrize.Add(remaining).String()
	return &Distribution{
		WinnerPrizes:         winnerPrizes,
		BeneficiaryReward:    beneficiaryReward.String(),
		MinParticipantReward: minParticipantReward.String(),
		RoundManagerFee:      managerFee.String(),
	}
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

func (p *DistributionDefinitions) GetNumWinners(key string) int {

	entry := p.Find(key)
	if entry == nil {
		panic(fmt.Sprintf("there is no distribution definition for key: %v", key))
	}
	return entry.Value.GetNumWinners()
}

func (p *DistributionDefinitions) GetTotalNumWinners() int {
	total := 0
	for _, entry := range *p {
		total += entry.Value.GetNumWinners()
	}
	return total
}
