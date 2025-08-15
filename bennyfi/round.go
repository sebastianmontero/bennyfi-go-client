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
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sebastianmontero/bennyfi-go-client/common/types"
	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
	"github.com/sebastianmontero/eos-go-toolbox/util"
)

var (
	PoolManagerFeeOwner      = "pool_manager_fee_owner"
	BeneficiaryRewardOwner   = "beneficiary_reward_owner"
	ReturnCycle              = "return_cycle"
	NumClaimedPartialReturns = "num_claimed_partial_returns"
	ExpectedYield            = "yield"
	RoundNotStarted          = eos.Name("notstarted")
	RoundPending             = eos.Name("pending")
	RoundAcceptingEntries    = eos.Name("open")
	RoundDrawing             = eos.Name("drawing")
	// RoundOpen                        = eos.Name("roundopen")
	RoundClosed                            = eos.Name("closed")
	RoundClosedYieldWithdrawnTestSetupOnly = eos.Name("closedyield")
	RoundUnlocked                          = eos.Name("unlocked")
	RoundTimedOut                          = eos.Name("cancelled")
	RoundStopped                           = eos.Name("stopped")
	RoundStakeStateNotStarted              = eos.Name("notstarted")
	RoundStakeStateStaked                  = eos.Name("staked")
	RoundStakeStatePayingPartialReturns    = eos.Name("payreturns")
	RoundStakeStateUnstakingTimedOut       = eos.Name("unstakingtmo")
	RoundStakeStateUnstakingUnlocked       = eos.Name("unstakingulk")
	RoundStakeStateUnstaked                = eos.Name("unstaked")
	RoundStakeStateStopped                 = eos.Name("stopped")
	VestingStateNotApplicable              = eos.Name("notaplicable")
	VestingStateNotStarted                 = eos.Name("notstarted")
	VestingStateVesting                    = eos.Name("vesting")
	VestingStateVesting1                   = eos.Name("vesting1") //used for entries to enable handling the different vesting cycles
	VestingStateVesting2                   = eos.Name("vesting2")
	VestingStateFinished                   = eos.Name("finished")
	VestingStateStopped                    = eos.Name("stopped")
	RoundTypeFunded                        = eos.Name("funded")
	RoundTypeYield                         = eos.Name("yield")
	RoundAccessPrivate                     = eos.Name("private")
	RoundAccessPublic                      = eos.Name("public")
	FundingStatePending                    = eos.Name("pending")
	FundingStateFunded                     = eos.Name("funded")
	FundingStateRefunded                   = eos.Name("refunded")
	FundingStatePartiallyCommited          = eos.Name("pcommited")
	FundingStateCommited                   = eos.Name("commited")
	FundingStateYield                      = eos.Name("yield")
)

type FundRoundArgs struct {
	RoundID uint64          `json:"pool_id"`
	Funder  eos.AccountName `json:"funder"`
}

type TstLapseTimeArgs struct {
	RoundID                uint64 `json:"pool_id"`
	LapseEnrollmentTimeEnd bool   `json:"lapse_enrollment_time_end"`
	CallCounter            uint64 `json:"call_counter"`
}

type Round struct {
	RoundID                  uint64                   `json:"pool_id"`
	TermID                   uint64                   `json:"term_id"`
	ProjectID                uint64                   `json:"project_id"`
	RoundName                string                   `json:"pool_name"`
	RoundDescription         string                   `json:"pool_description"`
	RoundCategory            eos.Name                 `json:"pool_category"`
	RoundType                eos.Name                 `json:"pool_type"`
	RoundAccess              eos.Name                 `json:"pool_access"`
	StakingPeriod            *dto.Microseconds        `json:"staking_period"`
	EnrollmentTimeOut        *dto.Microseconds        `json:"enrollment_time_out"`
	NumParticipants          uint32                   `json:"num_participants"`
	ParticipantEntryFee      eos.Asset                `json:"participant_entry_fee"`
	RoundManagerEntryFee     eos.Asset                `json:"pool_manager_entry_fee"`
	BeneficiaryEntryFee      eos.Asset                `json:"beneficiary_entry_fee"`
	BeneficiaryEntryFeeState eos.Name                 `json:"beneficiary_entry_fee_state"`
	EntryStake               eos.Asset                `json:"entry_stake"`
	Rewards                  Rewards                  `json:"rewards"`
	NumParticipantsEntered   uint32                   `json:"num_participants_entered"`
	NumClaimedReturns        uint32                   `json:"num_claimed_returns"`
	NumUnstaked              uint32                   `json:"num_unstaked"`
	NumEarlyExits            uint32                   `json:"num_early_exits"`
	VestingCycle             uint16                   `json:"vesting_cycle"`
	NumVested                uint16                   `json:"num_vested"`
	CurrentState             eos.Name                 `json:"current_state"`
	StakeState               eos.Name                 `json:"stake_state"`
	VestingState             eos.Name                 `json:"vesting_state"`
	TotalDeposits            eos.Asset                `json:"total_deposits"`
	Winners                  Winners                  `json:"winners"`
	Beneficiary              eos.AccountName          `json:"beneficiary"`
	Distributions            Distributions            `json:"distributions"`
	TotalEarlyExitStake      eos.Asset                `json:"total_early_exit_stake"`
	TotalEarlyExitRewardFees TotalEarlyExitRewardFees `json:"total_early_exit_reward_fees"`
	RoundManager             eos.AccountName          `json:"pool_manager"`
	StartTime                eos.TimePoint            `json:"start_time"`
	ClosedTime               eos.TimePoint            `json:"closed_time"`
	StakedTime               eos.TimePoint            `json:"staked_time"`
	StakeEndTime             eos.TimePoint            `json:"stake_end_time"`
	EnrollmentTimeEnd        eos.TimePoint            `json:"enrollment_time_end"`
	NextVestingTime          eos.TimePoint            `json:"next_vesting_time"`
	CreatedDate              eos.TimePoint            `json:"created_date"`
	UpdatedDate              eos.TimePoint            `json:"updated_date"`
	AdditionalFields         types.AdditionalFields   `json:"additional_fields"`
}

type RoundCustomJSON struct {
	Distributions map[eos.Name]interface{} `json:"distributions"`
	Rewards       map[eos.Name]interface{} `json:"rewards"`
	Winners       map[eos.Name]interface{} `json:"winners"`
	Round
}

func (m Round) ToCustomJSON() RoundCustomJSON {
	return RoundCustomJSON{
		Distributions: m.Distributions.ToMap(),
		Rewards:       m.Rewards.ToMap(),
		Winners:       m.Winners.ToMap(),
		Round:         m,
	}
}
func (m *Round) String() string {
	result, err := json.Marshal(m)
	if err != nil {
		panic(fmt.Sprintf("Failed marshalling round: %v", err))
	}
	return string(result)
}

func (m *Round) HasBeneficiary() bool {
	return strings.Trim(m.Beneficiary.String(), " ") != ""
}

func (m *Round) RequiresRoundManagerFunding() bool {
	return m.Rewards.Has(DistributionMainNFT)
}

func (m *Round) RequiresBeneficiaryFunding() bool {
	return m.Rewards.Has(DistributionProjectToken) || m.Rewards.Has(DistributionProjectNFT) || m.BeneficiaryEntryFee.Amount > 0
}

func (m *Round) CalculateMaxTotalDeposits() eos.Asset {
	totalStake := m.EntryStake
	totalStake.Amount = totalStake.Amount * eos.Int64(m.NumParticipants)
	return totalStake
}

func (m *Round) NumEntriesToClose() uint32 {
	return m.NumParticipants - m.NumParticipantsEntered
}

func (m *Round) GetPoolManagerFeeOwner() eos.AccountName {
	if m.AdditionalFields.Has(PoolManagerFeeOwner) {
		return eos.AN(m.AdditionalFields.GetValue(PoolManagerFeeOwner).Name().String())
	}
	return m.RoundManager
}

func (m *Round) SetPoolManagerFeeOwner(accountName interface{}) {
	account, err := util.ToName(accountName)
	if err != nil {
		panic(fmt.Sprintf("could not convert %v to eos.Name, error: %v", account, err))
	}
	m.AdditionalFields.Set(PoolManagerFeeOwner, dto.FlexValueFromName(account))
}

func (m *Round) GetBeneficiaryRewardOwner() eos.AccountName {
	if m.AdditionalFields.Has(BeneficiaryRewardOwner) {
		return eos.AN(m.AdditionalFields.GetValue(BeneficiaryRewardOwner).Name().String())
	}
	return m.Beneficiary
}

func (m *Round) SetBeneficiaryRewardOwner(accountName interface{}) {
	account, err := util.ToName(accountName)
	if err != nil {
		panic(fmt.Sprintf("could not convert %v to eos.Name, error: %v", account, err))
	}
	m.AdditionalFields.Set(BeneficiaryRewardOwner, dto.FlexValueFromName(account))
}

func (m *Round) GetExpectedYield() uint32 {
	if m.AdditionalFields.Has(ExpectedYield) {
		return m.AdditionalFields.GetValue(ExpectedYield).Uint32()
	}
	return 0
}

func (m *Round) SetExpectedYield(expectedYield uint32) {
	m.AdditionalFields.Set(ExpectedYield, dto.FlexValueFromUint32(expectedYield))
}

func (m *Round) GetReturnCycle() uint32 {
	if m.AdditionalFields.Has(ReturnCycle) {
		return m.AdditionalFields.GetValue(ReturnCycle).Uint32()
	}
	return 0
}

func (m *Round) SetReturnCycle(cycle uint32) {
	m.AdditionalFields.Set(ReturnCycle, dto.FlexValueFromUint32(cycle))
}

func (m *Round) IncReturnCycle() {
	m.AdditionalFields.Set(ReturnCycle, dto.FlexValueFromUint32(m.GetReturnCycle()+1))
}

func (m *Round) GetNumClaimedPartialReturns() uint32 {
	if m.AdditionalFields.Has(NumClaimedPartialReturns) {
		return m.AdditionalFields.GetValue(NumClaimedPartialReturns).Uint32()
	}
	return 0
}

func (m *Round) SetNumClaimedPartialReturns(numClaimedPartialReturns uint32) {
	m.AdditionalFields.Set(NumClaimedPartialReturns, dto.FlexValueFromUint32(numClaimedPartialReturns))
}

func (m *Round) IncNumClaimedPartialReturns() {
	numClaimedPartialReturns := m.GetNumClaimedPartialReturns() + 1
	if numClaimedPartialReturns == m.NumParticipantsEntered {
		numClaimedPartialReturns = 0
	}
	m.AdditionalFields.Set(NumClaimedPartialReturns, dto.FlexValueFromUint32(numClaimedPartialReturns))
}

func (m *Round) UpsertDistribution(name eos.Name, distribution interface{}) {
	if m.Distributions == nil {
		m.Distributions = make(Distributions, 0, 1)
	}
	m.Distributions.Upsert(name, distribution)
}

func (m *Round) AssignWinnerPrizes(distName eos.Name, dist *Distribution) error {
	winnersEntry := m.Winners.Find(distName)
	if winnersEntry == nil {
		return fmt.Errorf("failed assigning winner prizes, there is no winners array for distribution name: %v", distName)
	}
	err := winnersEntry.Value.AssignPrizes(dist)
	if err != nil {
		return fmt.Errorf("failed assigning winner prizes for dist: %v, error: %v", distName, err)
	}
	return nil
}

func (m *Round) RemoveDistribution(name eos.Name) {
	m.Distributions.Remove(name)
}

func (m *Round) UpsertReward(name eos.Name, reward interface{}) {
	if m.Rewards == nil {
		m.Rewards = make(Rewards, 0, 1)
	}
	m.Rewards.Upsert(name, reward)
}

func (m *Round) RemoveReward(name eos.Name) {
	m.Rewards.Remove(name)
}

func (m *Round) UpsertEarlyExitRewardFee(name eos.Name, earlyExitRewardFee eos.Asset) {
	if m.TotalEarlyExitRewardFees == nil {
		m.TotalEarlyExitRewardFees = make(TotalEarlyExitRewardFees, 0, 1)
	}
	m.TotalEarlyExitRewardFees.Upsert(name, earlyExitRewardFee)
}

func (m *Round) RemoveEarlyExitRewardFee(name eos.Name) {
	m.TotalEarlyExitRewardFees.Remove(name)
}

func (m *Round) UpsertWinner(name eos.Name, winner interface{}) {
	if m.Winners == nil {
		m.Winners = make(Winners, 0, 1)
	}
	m.Winners.Upsert(name, winner)
}

func (m *Round) RemoveWinner(name eos.Name) {
	m.Winners.Remove(name)
}

func (m *Round) UpdateFundingStateAll(state eos.Name) {
	m.Rewards.UpdateFundingStateAll(state)
}

func (m *Round) UpdateFundingState(dist eos.Name, state eos.Name) {
	m.Rewards.UpdateFundingState(dist, state)
}

func (m *Round) UpdateEntryFees(entryFees *EntryFees) {
	m.RoundManagerEntryFee = entryFees.PoolManagerEntryFee
	m.BeneficiaryEntryFee = entryFees.BeneficiaryEntryFee
	m.ParticipantEntryFee = entryFees.ParticipantEntryFee
}

func (m *Round) GetTotalEntryFee() eos.Asset {
	// fmt.Printf("Calculating total entry fee, Round Manager Entry Fee: %v, Beneficiary Entry Fee: %v, Participant Entry Fee: %v \n", m.RoundManagerEntryFee, m.BeneficiaryEntryFee, m.ParticipantEntryFee)
	return m.BeneficiaryEntryFee.Add(m.RoundManagerEntryFee).Add(util.MultiplyAsset(m.ParticipantEntryFee, int64(m.NumParticipantsEntered)))
}

func (m *Round) CalculateReturns(entryPos uint64, distName eos.Name, isEarlyExit bool, earlyExitFeePerc uint32, currentReturns *Returns) interface{} {

	if IsFTDistribution(distName) {
		dist := m.Distributions.FindFT(distName)
		minParticipantReward := dist.MinParticipantReward
		fmt.Printf("Winners: %v\n", m.Winners)
		winner := m.Winners.FindWinnerFT(distName, entryPos)
		winnerPrize := eos.Asset{Amount: 0, Symbol: minParticipantReward.Symbol}
		if winner != nil {
			winnerPrize = winner.Prize
		}
		earlyExitRewardFee := util.CalculateAssetPercentage(winnerPrize, earlyExitFeePerc)
		if isEarlyExit {
			winnerPrize = winnerPrize.Sub(earlyExitRewardFee)
			earlyExitRewardFee = earlyExitRewardFee.Add(minParticipantReward)
			minParticipantReward = eos.Asset{Amount: 0, Symbol: minParticipantReward.Symbol}
		}

		amountPaidOut := eos.Asset{Amount: 0, Symbol: minParticipantReward.Symbol}
		if currentReturns != nil {
			amountPaidOut = currentReturns.ReturnsFT().AmountPaidOut
		}
		return &ReturnsFT{
			Prize:              winnerPrize,
			MinimumPayout:      minParticipantReward,
			EarlyExitReturnFee: earlyExitRewardFee,
			AmountPaidOut:      amountPaidOut,
		}
	} else {
		dist := m.Distributions.FindNFT(distName)
		winner := m.Winners.FindWinnerNFT(distName, entryPos)
		winnerPrize := uint16(0)
		if winner != nil {
			winnerPrize = winner.Prize
		}
		return &ReturnsNFT{
			Prize:         winnerPrize,
			MinimumPayout: dist.MinParticipantReward,
		}
	}
}

func (m *Round) CalculateUnlockTime() eos.TimePoint {
	return eos.TimePoint(m.StakedTime.Time().Add(time.Hour * time.Duration(m.StakingPeriod.Hrs())).UnixMicro())
}

func (m *Round) SetYieldReward(totalReturn eos.Asset, partial bool) (reward eos.Asset, totalReward eos.Asset) {
	r := m.Rewards.FindFT(DistributionMainToken)
	reward = totalReturn
	reward = reward.Sub(m.TotalDeposits)
	if partial {
		r.FundingState = FundingStatePartiallyCommited
	} else {
		r.FundingState = FundingStateCommited
	}
	if reward.Amount < 0 {
		reward = eos.Asset{Amount: 0, Symbol: totalReturn.Symbol}
	}
	totalReward = r.Reward.Add(reward)
	r.Reward = totalReward
	return
}

type NewRoundArgs struct {
	RoundManager     eos.AccountName `json:"pool_manager"`
	TermID           uint64          `json:"term_id"`
	ProjectID        uint64          `json:"project_id"`
	RoundName        string          `json:"pool_name"`
	RoundDescription string          `json:"pool_description"`
	RoundCategory    eos.Name        `json:"pool_category"`
	StartTime        eos.TimePoint   `json:"start_time"`
}

func RoundToNewRoundArgs(round *Round) *NewRoundArgs {
	return &NewRoundArgs{
		TermID:           round.TermID,
		ProjectID:        round.ProjectID,
		RoundName:        round.RoundName,
		RoundDescription: round.RoundDescription,
		RoundCategory:    round.RoundCategory,
		RoundManager:     round.RoundManager,
		StartTime:        round.StartTime,
	}
}

func (m *Round) Clone() *Round {
	return &Round{
		RoundID:                  m.RoundID,
		TermID:                   m.TermID,
		ProjectID:                m.ProjectID,
		RoundName:                m.RoundName,
		RoundDescription:         m.RoundDescription,
		RoundCategory:            m.RoundCategory,
		RoundType:                m.RoundType,
		RoundAccess:              m.RoundAccess,
		StakingPeriod:            m.StakingPeriod,
		EnrollmentTimeOut:        m.EnrollmentTimeOut,
		NumParticipants:          m.NumParticipants,
		ParticipantEntryFee:      m.ParticipantEntryFee,
		RoundManagerEntryFee:     m.RoundManagerEntryFee,
		BeneficiaryEntryFee:      m.BeneficiaryEntryFee,
		BeneficiaryEntryFeeState: m.BeneficiaryEntryFeeState,
		EntryStake:               m.EntryStake,
		Rewards:                  m.Rewards.Clone(),
		NumParticipantsEntered:   m.NumParticipantsEntered,
		NumClaimedReturns:        m.NumClaimedReturns,
		NumUnstaked:              m.NumUnstaked,
		NumEarlyExits:            m.NumEarlyExits,
		VestingCycle:             m.VestingCycle,
		NumVested:                m.NumVested,
		CurrentState:             m.CurrentState,
		StakeState:               m.StakeState,
		VestingState:             m.VestingState,
		TotalDeposits:            m.TotalDeposits,
		Winners:                  m.Winners,
		Beneficiary:              m.Beneficiary,
		Distributions:            m.Distributions,
		TotalEarlyExitStake:      m.TotalEarlyExitStake,
		TotalEarlyExitRewardFees: m.TotalEarlyExitRewardFees,
		RoundManager:             m.RoundManager,
		StartTime:                m.StartTime,
		ClosedTime:               m.ClosedTime,
		StakedTime:               m.StakedTime,
		StakeEndTime:             m.StakeEndTime,
		EnrollmentTimeEnd:        m.EnrollmentTimeEnd,
		NextVestingTime:          m.NextVestingTime,
		CreatedDate:              m.CreatedDate,
		UpdatedDate:              m.UpdatedDate,
	}
}

func (m *BennyfiContract) NewRound(round *Round) (string, error) {
	return m.NewRoundFromRoundArgs(RoundToNewRoundArgs(round))
}

func (m *BennyfiContract) NewRoundFromRoundArgs(roundArgs *NewRoundArgs) (string, error) {
	return m.ExecAction(roundArgs.RoundManager, "newpool", roundArgs)
}

func (m *BennyfiContract) StartRound(roundID uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "startpool", roundID)
}

func (m *BennyfiContract) TimeoutRound(roundID uint64, sudo bool, authorizer interface{}) (string, error) {
	var permissionLevel interface{}
	var author eos.AccountName
	var err error
	if authorizer == nil {
		permissionLevel = fmt.Sprintf("%v@open", m.ContractName)
		author = eos.AccountName(m.ContractName)
	} else {
		permissionLevel = authorizer
		author, err = util.ToAccountName(authorizer)
		if err != nil {
			return "", fmt.Errorf("failed parsing authorizer account: %v, error: %v", authorizer, err)
		}
	}
	actionData := struct {
		Authorizer eos.AccountName
		RoundId    uint64
		Sudo       bool
	}{author, roundID, sudo}
	// fmt.Printf("Permission level: %v\n", permissionLevel)
	return m.ExecAction(permissionLevel, "timeoutpool", actionData)
}

func (m *BennyfiContract) DeleteStoppedRound(roundID uint64, authorizer eos.AccountName) (string, error) {
	actionData := struct {
		Authorizer eos.AccountName
		RoundId    uint64
	}{authorizer, roundID}
	// fmt.Printf("Permission level: %v\n", permissionLevel)
	return m.ExecAction(authorizer, "delstppdpool", actionData)
}

func (m *BennyfiContract) StopRound(roundID uint64, sudo bool, authorizer interface{}) (string, error) {
	var permissionLevel interface{}
	var author eos.AccountName
	var err error
	if authorizer == nil {
		permissionLevel = fmt.Sprintf("%v@open", m.ContractName)
		author = eos.AccountName(m.ContractName)
	} else {
		permissionLevel = authorizer
		author, err = util.ToAccountName(authorizer)
		if err != nil {
			return "", fmt.Errorf("failed parsing authorizer account: %v, error: %v", authorizer, err)
		}
	}
	actionData := struct {
		Authorizer eos.AccountName
		RoundId    uint64
		Sudo       bool
	}{author, roundID, sudo}
	// fmt.Printf("Permission level: %v\n", permissionLevel)
	return m.ExecAction(permissionLevel, "stoppool", actionData)
}

func (m *BennyfiContract) FundRound(roundID uint64, funder interface{}) (string, error) {
	f, err := util.ToAccountName(funder)
	if err != nil {
		return "", fmt.Errorf("failed parsing funder account: %v, error: %v", funder, err)
	}
	actionData := &FundRoundArgs{
		RoundID: roundID,
		Funder:  f,
	}

	return m.ExecAction(funder, "fundpool", actionData)
}

func (m *BennyfiContract) StartRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "startpools", callCounter)
}

func (m *BennyfiContract) EndEnrollment(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "endenrollmnt", callCounter)
}

func (m *BennyfiContract) ClaimPartialReturns(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "clmprtrtrnpl", callCounter)
}

func (m *BennyfiContract) UnlockRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "unlockpools", callCounter)
}

func (m *BennyfiContract) UnlockRound(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "unlockpool", roundId)
}

func (m *BennyfiContract) UnstakeUnlockedRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "ustkulkpools", callCounter)
}

func (m *BennyfiContract) UnstakeTimedoutRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "ustktmdpools", callCounter)
}

func (m *BennyfiContract) DeleteTimedoutRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "deltmdpools", callCounter)
}

func (m *BennyfiContract) Redraw(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "redraw", callCounter)
}

func (m *BennyfiContract) VestingRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "vestingpools", callCounter)
}

func (m *BennyfiContract) TstLapseTime(roundId uint64) (string, error) {
	return m.TstLapseTimeDetailed(roundId, false)
}

func (m *BennyfiContract) TstLapseTimeDetailed(roundId uint64, lapseEnrollmentTimeEnd bool) (string, error) {
	actionData := &TstLapseTimeArgs{
		RoundID:                roundId,
		LapseEnrollmentTimeEnd: lapseEnrollmentTimeEnd,
		CallCounter:            m.NextCallCounter(),
	}
	return m.ExecAction(eos.AN(m.ContractName), "tstlapsetime", actionData)
}

func (m *BennyfiContract) ReceiveRand(actor eos.AccountName, roundId uint64, randomNumber string) (string, error) {
	rand, err := hex.DecodeString(randomNumber)
	if err != nil {
		return "", fmt.Errorf("failed decoding random number: %v, error: %v", randomNumber, err)
	}
	actionData := struct {
		AssocID uint64
		Random  eos.Checksum256
	}{roundId, eos.Checksum256(rand)}
	return m.ExecAction(actor, "receiverand", actionData)
}

func (m *BennyfiContract) GetRounds() ([]Round, error) {

	return m.GetRoundsReq(nil)
}

func (m *BennyfiContract) GetAllRoundsAsMap() ([]map[string]interface{}, error) {
	return m.GetAllRoundsFromAsMap(0)
}

func (m *BennyfiContract) GetAllRoundsFromAsMap(roundID uint64) ([]map[string]interface{}, error) {
	req := eos.GetTableRowsRequest{
		Table: "pools",
	}
	return m.GetAllTableRowsFromAsMap(req, "pool_id", strconv.FormatUint(roundID, 10), nil)
}

func (m *BennyfiContract) GetAllRoundsFrom(roundID uint64) ([]Round, error) {

	roundsAsMap, err := m.GetAllRoundsFromAsMap(roundID)
	if err != nil {
		return nil, err
	}
	roundsAsStr, err := json.Marshal(roundsAsMap)
	if err != nil {
		return nil, fmt.Errorf("failed marshalling rounds map: %v", err)
	}
	var rounds []Round
	err = json.Unmarshal(roundsAsStr, &rounds)
	if err != nil {
		return nil, fmt.Errorf("failed marshalling rounds map: %v", err)
	}
	return rounds, nil
}

func (m *BennyfiContract) GetRoundsbyTermAndId(termId uint64) ([]Round, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterRoundsbyTermAndId(request, termId)
	if err != nil {
		return nil, err
	}
	return m.GetRoundsReq(request)
}

func (m *BennyfiContract) FilterRoundsbyTermAndId(req *eos.GetTableRowsRequest, term uint64) error {

	req.Index = "15"
	req.KeyType = "i128"
	req.Reverse = true
	termAndRndLB, err := m.EOS.GetComposedIndexValue(term, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	termAndRndUB, err := m.EOS.GetComposedIndexValue(term, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	fmt.Println("LB: ", termAndRndLB, "UB: ", termAndRndUB)
	req.LowerBound = termAndRndLB
	req.UpperBound = termAndRndUB
	return err
}

func (m *BennyfiContract) GetRound(roundID uint64) (*Round, error) {
	rounds, err := m.GetRoundsReq(&eos.GetTableRowsRequest{
		LowerBound: strconv.FormatUint(roundID, 10),
		UpperBound: strconv.FormatUint(roundID, 10),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(rounds) > 0 {
		return &rounds[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) GetLastRound() (*Round, error) {
	rounds, err := m.GetRoundsReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(rounds) > 0 {
		return &rounds[0], nil
	}
	return nil, nil
}

func (m *BennyfiContract) GetRoundsReq(req *eos.GetTableRowsRequest) ([]Round, error) {

	var rounds []Round
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "pools"
	err := m.GetTableRows(*req, &rounds)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return rounds, nil
}

func (m *BennyfiContract) GetRoundsbyStateAndId(state eos.Name) ([]Round, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterRoundsbyStateAndId(request, state)
	if err != nil {
		return nil, err
	}
	return m.GetRoundsReq(request)
}

func (m *BennyfiContract) FilterRoundsbyStateAndId(req *eos.GetTableRowsRequest, state eos.Name) error {

	req.Index = "2"
	req.KeyType = "i128"
	req.Reverse = true
	stateAndRndLB, err := m.EOS.GetComposedIndexValue(state, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	stateAndRndUB, err := m.EOS.GetComposedIndexValue(state, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	fmt.Println("LB: ", stateAndRndLB, "UB: ", stateAndRndUB)
	req.LowerBound = stateAndRndLB
	req.UpperBound = stateAndRndUB
	return err
}

func (m *BennyfiContract) GetRoundsbyManagerAndId(manager interface{}) ([]Round, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterRoundsbyManagerAndId(request, manager)
	if err != nil {
		return nil, err
	}
	return m.GetRoundsReq(request)
}

func (m *BennyfiContract) FilterRoundsbyManagerAndId(req *eos.GetTableRowsRequest, manager interface{}) error {

	req.Index = "3"
	req.KeyType = "i128"
	req.Reverse = true
	mgrAndRndLB, err := m.EOS.GetComposedIndexValue(manager, 0)
	if err != nil {
		return fmt.Errorf("failed to generate lower bound composed index, err: %v", err)
	}
	mgrAndRndUB, err := m.EOS.GetComposedIndexValue(manager, uint64(18446744073709551615))
	if err != nil {
		return fmt.Errorf("failed to generate upper bound composed index, err: %v", err)
	}
	// fmt.Println("LB: ", mgrAndRndLB, "UB: ", mgrAndRndUB)
	req.LowerBound = mgrAndRndLB
	req.UpperBound = mgrAndRndUB
	return err
}

func createEOSProof(randomNumber uint64) map[string]interface{} {
	proof := make(map[string]interface{})
	proof["block_num"] = 1
	proof["block_id"] = "blockid"
	proof["seed"] = 1
	proof["final_seed"] = "finalseed"
	proof["public_key"] = "publickey"
	proof["gamma"] = "gamma"
	proof["c"] = "c"
	proof["s"] = "s"
	proof["output_u256"] = "u256"
	proof["output_u64"] = randomNumber
	return proof
}
