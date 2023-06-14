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
	"math/big"
	"strconv"
	"strings"
	"time"

	eos "github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/dto"
	"github.com/sebastianmontero/eos-go-toolbox/util"
)

var (
	RoundNotStarted       = eos.Name("notstarted")
	RoundPending          = eos.Name("pending")
	RoundAcceptingEntries = eos.Name("open")
	RoundDrawing          = eos.Name("drawing")
	// RoundOpen                        = eos.Name("roundopen")
	RoundClosed                      = eos.Name("closed")
	RoundUnlocked                    = eos.Name("unlocked")
	RoundTimedOut                    = eos.Name("cancelled")
	RoundStakeStateNotStarted        = eos.Name("notstarted")
	RoundStakeStateStaked            = eos.Name("staked")
	RoundStakeStateUnstakingTimedOut = eos.Name("unstakingtmo")
	RoundStakeStateUnstakingUnlocked = eos.Name("unstakingulk")
	RoundStakeStateUnstaked          = eos.Name("unstaked")
	VestingStateNotApplicable        = eos.Name("notaplicable")
	VestingStateNotStarted           = eos.Name("notstarted")
	VestingStateVesting              = eos.Name("vesting")
	VestingStateVesting1             = eos.Name("vesting1") //used for entries to enable handling the different vesting cycles
	VestingStateVesting2             = eos.Name("vesting2")
	VestingStateFinished             = eos.Name("finished")
	RoundTypeManagerFunded           = eos.Name("mgrfunded")
	RoundTypeRexPool                 = eos.Name("rexpool")
	RoundAccessPrivate               = eos.Name("private")
	RoundAccessPublic                = eos.Name("public")
	RexStateNotApplicable            = eos.Name("notaplicable")
	RexStatePreRex                   = eos.Name("prerex")
	RexStateInSavings                = eos.Name("insavings")
	RexStateInLockPeriod             = eos.Name("lockperiod")
	RexStateSold                     = eos.Name("sold")
	RexStateWithdrawn                = eos.Name("withdrawn")
	FundingStatePending              = eos.Name("pending")
	FundingStateFunded               = eos.Name("funded")
	FundingStateRefunded             = eos.Name("refunded")
	FundingStateCommited             = eos.Name("commited")
	FundingStateRex                  = eos.Name("rex")
	RexLockPeriodDays                = 5
)

type FundRoundArgs struct {
	RoundID uint64          `json:"round_id"`
	Funder  eos.AccountName `json:"funder"`
}

type TstLapseTimeArgs struct {
	RoundID     uint64 `json:"round_id"`
	CallCounter uint64 `json:"call_counter"`
}

type Round struct {
	RoundID                  uint64                   `json:"round_id"`
	TermID                   uint64                   `json:"term_id"`
	ProjectID                uint64                   `json:"project_id"`
	RoundName                string                   `json:"round_name"`
	RoundDescription         string                   `json:"round_description"`
	RoundCategory            eos.Name                 `json:"round_category"`
	RoundType                eos.Name                 `json:"round_type"`
	RoundAccess              eos.Name                 `json:"round_access"`
	StakingPeriod            *dto.Microseconds        `json:"staking_period"`
	EnrollmentTimeOut        *dto.Microseconds        `json:"enrollment_time_out"`
	NumParticipants          uint32                   `json:"num_participants"`
	ParticipantEntryFee      eos.Asset                `json:"participant_entry_fee"`
	RoundManagerEntryFee     eos.Asset                `json:"round_manager_entry_fee"`
	BeneficiaryEntryFee      eos.Asset                `json:"beneficiary_entry_fee"`
	BeneficiaryEntryFeeState eos.Name                 `json:"beneficiary_entry_fee_state"`
	EntryStake               eos.Asset                `json:"entry_stake"`
	Rewards                  Rewards                  `json:"rewards"`
	RexBalance               eos.Asset                `json:"rex_balance"`
	NumParticipantsEntered   uint32                   `json:"num_participants_entered"`
	NumClaimedReturns        uint32                   `json:"num_claimed_returns"`
	NumUnstaked              uint32                   `json:"num_unstaked"`
	NumEarlyExits            uint32                   `json:"num_early_exits"`
	VestingCycle             uint16                   `json:"vesting_cycle"`
	NumVested                uint16                   `json:"num_vested"`
	CurrentState             eos.Name                 `json:"current_state"`
	RexState                 eos.Name                 `json:"rex_state"`
	StakeState               eos.Name                 `json:"stake_state"`
	VestingState             eos.Name                 `json:"vesting_state"`
	TotalDeposits            eos.Asset                `json:"total_deposits"`
	Winners                  Winners                  `json:"winners"`
	Beneficiary              eos.AccountName          `json:"beneficiary"`
	Distributions            Distributions            `json:"distributions"`
	TotalEarlyExitStake      eos.Asset                `json:"total_early_exit_stake"`
	TotalEarlyExitRewardFees TotalEarlyExitRewardFees `json:"total_early_exit_reward_fees"`
	RoundManager             eos.AccountName          `json:"round_manager"`
	StartTime                eos.TimePoint            `json:"start_time"`
	ClosedTime               eos.TimePoint            `json:"closed_time"`
	StakedTime               eos.TimePoint            `json:"staked_time"`
	MovedFromSavingsTime     eos.TimePoint            `json:"moved_from_savings_time"`
	StakeEndTime             eos.TimePoint            `json:"stake_end_time"`
	EnrollmentTimeEnd        eos.TimePoint            `json:"enrollment_time_end"`
	NextVestingTime          eos.TimePoint            `json:"next_vesting_time"`
	CreatedDate              eos.TimePoint            `json:"created_date"`
	UpdatedDate              eos.TimePoint            `json:"updated_date"`
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
func (m *Round) GetTotalEntryFee() eos.Asset {
	return m.BeneficiaryEntryFee.Add(m.RoundManagerEntryFee).Add(util.MultiplyAsset(m.ParticipantEntryFee, int64(m.NumParticipants)))
}

func (m *Round) NumEntriesToClose() uint32 {
	return m.NumParticipants - m.NumParticipantsEntered
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

func (m *Round) CalculateEntryFee(settings *EntryFeeSettings) eos.Asset {
	if m.RoundType == RoundTypeManagerFunded {
		// fmt.Println("Manager funded entry fee: ", util.MultiplyAsset(settings.SelfFundedPerUser, int64(m.NumParticipants)))
		return util.MultiplyAsset(settings.SelfFundedPerUser, int64(m.NumParticipants))
	} else {

		totalStake := util.MultiplyAsset(m.EntryStake, int64(m.NumParticipants))
		yield := util.CalculateAssetPercentage(util.MultiplyAsset(totalStake, int64(m.StakingPeriod.Hrs())), settings.HourlyYield())
		yieldUSD := util.DivideAssets(yield, settings.ValueTLOS)
		yieldPerc := util.CalculateAssetPercentage(yieldUSD, settings.PercOfYield)
		entryFee := util.DivideAssets(yieldPerc, settings.ValueBENY)
		adjustedEntryFee := util.AdjustPrecision(big.NewInt(int64(entryFee.Amount)), entryFee.Precision, settings.BENYToken.Precision)
		// fmt.Printf("Entry fee values, total stake: %v, yield: %v, yieldUSD: %v, yieldPerc: %v, entryFee: %v, adjustedEntryFee: %v \n", totalStake, yield, yieldUSD, yieldPerc, entryFee, adjustedEntryFee)
		return eos.Asset{Amount: eos.Int64(adjustedEntryFee.Int64()), Symbol: settings.BENYToken.Symbol}
	}
}

func (m *Round) CalculateEntryFees(settings *EntryFeeSettings, term *Term) {
	entryFee := m.CalculateEntryFee(settings)
	roundManagerEntryFee := util.CalculateAssetPercentage(entryFee, term.RoundManagerEntryFeePerc)
	beneficiaryEntryFee := util.CalculateAssetPercentage(entryFee, term.BeneficiaryEntryFeePerc)
	participantEntryFee := entryFee.Sub(roundManagerEntryFee).Sub(beneficiaryEntryFee)
	// fmt.Printf("Round manager percent fee: %v, beneficiary percent fee: %v\n", term.RoundManagerEntryFeePerc, term.BeneficiaryEntryFeePerc)
	// fmt.Printf("Entryfee: %v, Beneficiary Entry fee: %v, Round Manager Entry fee: %v, Participant Entry Fee total: %v \n", entryFee, beneficiaryEntryFee, roundManagerEntryFee, participantEntryFee)
	participantEntryFee = util.DivideAsset(participantEntryFee, uint64(m.NumParticipants))
	m.RoundManagerEntryFee = roundManagerEntryFee
	m.BeneficiaryEntryFee = beneficiaryEntryFee
	m.ParticipantEntryFee = participantEntryFee
}

func (m *Round) CalculateReturns(entryOwner eos.AccountName, distName eos.Name, isEarlyExit bool, earlyExitFeePerc uint32) interface{} {

	if IsFTDistribution(distName) {
		dist := m.Distributions.FindFT(distName)
		minParticipantReward := dist.MinParticipantReward
		winner := m.Winners.FindWinnerFT(distName, entryOwner)
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
		return &ReturnsFT{
			Prize:              winnerPrize,
			MinimumPayout:      minParticipantReward,
			EarlyExitReturnFee: earlyExitRewardFee,
			AmountPaidOut:      eos.Asset{Amount: 0, Symbol: winnerPrize.Symbol},
		}
	} else {
		dist := m.Distributions.FindNFT(distName)
		winner := m.Winners.FindWinnerNFT(distName, entryOwner)
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

func (m *Round) CalculateRexLockPeriodTime() eos.TimePoint {

	return eos.TimePoint(m.StakedTime.Time().Add(time.Hour * time.Duration(m.StakingPeriod.Hrs()-int64(24*RexLockPeriodDays))).UnixMicro())
}

func (m *Round) CalculateSellRexTime() eos.TimePoint {

	return eos.TimePoint(m.MovedFromSavingsTime.Time().Add(time.Hour * time.Duration(24*RexLockPeriodDays)).UnixMicro())
}

func (m *Round) CalculateUnlockTime() eos.TimePoint {
	return eos.TimePoint(m.StakedTime.Time().Add(time.Hour * time.Duration(m.StakingPeriod.Hrs())).UnixMicro())
}

type NewRoundArgs struct {
	RoundManager     eos.AccountName `json:"round_manager"`
	TermID           uint64          `json:"term_id"`
	ProjectID        uint64          `json:"project_id"`
	RoundName        string          `json:"round_name"`
	RoundDescription string          `json:"round_description"`
	RoundCategory    eos.Name        `json:"round_category"`
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
		RexBalance:               m.RexBalance,
		NumParticipantsEntered:   m.NumParticipantsEntered,
		NumClaimedReturns:        m.NumClaimedReturns,
		NumUnstaked:              m.NumUnstaked,
		NumEarlyExits:            m.NumEarlyExits,
		VestingCycle:             m.VestingCycle,
		NumVested:                m.NumVested,
		CurrentState:             m.CurrentState,
		RexState:                 m.RexState,
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
		MovedFromSavingsTime:     m.MovedFromSavingsTime,
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
	return m.ExecAction(roundArgs.RoundManager, "newround", roundArgs)
}

func (m *BennyfiContract) StartRound(roundID uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "startround", roundID)
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

	return m.ExecAction(funder, "fundround", actionData)
}

func (m *BennyfiContract) TimedEvents() (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "timedevents", nil)
}

func (m *BennyfiContract) StartRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "startrounds", callCounter)
}

func (m *BennyfiContract) TimeoutRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "timeoutrnds", callCounter)
}

func (m *BennyfiContract) MoveFromSavings(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "mvfrmsavings", callCounter)
}

func (m *BennyfiContract) SellRex(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "sellrex", callCounter)
}

func (m *BennyfiContract) WithdrawRex(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "withdrawrex", callCounter)
}

func (m *BennyfiContract) UnlockRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "unlockrnds", callCounter)
}

func (m *BennyfiContract) UnlockRound(roundId uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "unlockrnd", roundId)
}

func (m *BennyfiContract) UnstakeUnlockedRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "ustkulckrnds", callCounter)
}

func (m *BennyfiContract) UnstakeTimedoutRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "ustktmdrnds", callCounter)
}

func (m *BennyfiContract) Redraw() (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "redraw", nil)
}

func (m *BennyfiContract) VestingRounds(callCounter uint64) (string, error) {
	return m.ExecAction(fmt.Sprintf("%v@open", m.ContractName), "vestingrnds", callCounter)
}

func (m *BennyfiContract) TstLapseTime(roundId uint64) (string, error) {
	actionData := &TstLapseTimeArgs{
		RoundID:     roundId,
		CallCounter: m.NextCallCounter(),
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
		Table: "rounds",
	}
	return m.GetAllTableRowsFromAsMap(req, "round_id", strconv.FormatUint(roundID, 10))
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
	req.Table = "rounds"
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
