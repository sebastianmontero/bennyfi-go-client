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
package marble

import (
	"fmt"
	"strconv"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/bennyfi-go-client/bennyfi"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

var (
	SettingBennyfiContract     = "BENNYFI_CONTRACT"
	SettingRewardDescMinLength = "REWARD_DESC_MIN_LENGTH"
	SettingRewardDescMaxLength = "REWARD_DESC_MAX_LENGTH"
	RewardStatePending         = eos.Name("pending")
	RewardStateVerified        = eos.Name("verified")
	RewardStateCommited        = eos.Name("commited")
	RewardStateRefunded        = eos.Name("refunded")
)

type Reward struct {
	RoundID      uint64          `json:"round_id"`
	Description  string          `json:"description"`
	GroupName    eos.Name        `json:"group_name"`
	FrameName    eos.Name        `json:"frame_name"`
	CurrentState eos.Name        `json:"current_state"`
	Beneficiary  eos.AccountName `json:"beneficiary"`
	NFTContract  eos.AccountName `json:"nft_contract"`
	CreatedDate  string          `json:"created_date"`
	UpdatedDate  string          `json:"updated_date"`
}

func (m *Reward) NewRewardArgs() *NewRewardArgs {
	return &NewRewardArgs{
		RoundID:     m.RoundID,
		Description: m.Description,
		GroupName:   m.GroupName,
		FrameName:   m.FrameName,
		Beneficiary: m.Beneficiary,
		NFTContract: m.NFTContract,
	}
}

// Order of struct properties for action must be on the same order as the action parameters for the call to succeed
type NewRewardArgs struct {
	Beneficiary eos.AccountName `json:"beneficiary"`
	RoundID     uint64          `json:"round_id"`
	Description string          `json:"description"`
	GroupName   eos.Name        `json:"group_name"`
	FrameName   eos.Name        `json:"frame_name"`
	NFTContract eos.AccountName `json:"nft_contract"`
}

type MarbleAdaptorContract struct {
	*contract.SettingsContract
}

func NewMarbleAdaptorContract(eos *service.EOS, contractName string) *MarbleAdaptorContract {
	return &MarbleAdaptorContract{
		contract.NewSettingsContract(eos, contractName),
	}
}

func (m *MarbleAdaptorContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *MarbleAdaptorContract) SetReward(reward *Reward) (string, error) {
	return m.SetRewardFromArgs(reward.NewRewardArgs())
}

func (m *MarbleAdaptorContract) SetRewardFromArgs(args *NewRewardArgs) (string, error) {
	return m.ExecAction(args.Beneficiary, "setreward", args)
}

func (m *MarbleAdaptorContract) Verify(args *bennyfi.NFTActionParams, caller eos.AccountName) (string, error) {
	return m.ExecAction(caller, "verify", args)
}

func (m *MarbleAdaptorContract) Commit(args *bennyfi.NFTActionParams, caller eos.AccountName) (string, error) {
	return m.ExecAction(caller, "commit", args)
}

func (m *MarbleAdaptorContract) Refund(args *bennyfi.NFTActionParams, caller eos.AccountName) (string, error) {
	return m.ExecAction(caller, "refund", args)
}

func (m *MarbleAdaptorContract) Transfer(args *bennyfi.NFTActionParams, caller eos.AccountName) (string, error) {
	return m.ExecAction(caller, "transfer", args)
}

func (m *MarbleAdaptorContract) GetRewardByRound(roundId uint64) (*Reward, error) {
	request := &eos.GetTableRowsRequest{
		Limit:      1,
		LowerBound: strconv.FormatUint(roundId, 10),
		UpperBound: strconv.FormatUint(roundId, 10),
	}
	rewards, err := m.GetRewardsReq(request)
	if err != nil {
		return nil, fmt.Errorf("failed getting round by id: %v, error: %v", roundId, err)
	}
	if len(rewards) > 0 {
		return rewards[0], nil
	}
	return nil, nil
}

func (m *MarbleAdaptorContract) GetRewardsByGroup(group eos.Name) ([]*Reward, error) {
	request := &eos.GetTableRowsRequest{}
	m.FilterRewardsByGroup(request, group)
	return m.GetRewardsReq(request)
}

func (m *MarbleAdaptorContract) FilterRewardsByGroup(req *eos.GetTableRowsRequest, group eos.Name) {

	req.Index = "2"
	req.KeyType = "name"
	req.LowerBound = string(group)
	req.UpperBound = string(group)
}

func (m *MarbleAdaptorContract) GetRewardsReq(req *eos.GetTableRowsRequest) ([]*Reward, error) {

	var rewards []*Reward
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "rewards"
	err := m.GetTableRows(*req, &rewards)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return rewards, nil
}