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
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type Group struct {
	Title        string          `json:"title"`
	Description  string          `json:"description"`
	GroupName    eos.Name        `json:"group_name"`
	Manager      eos.AccountName `json:"manager"`
	Supply       uint64          `json:"supply"`
	IssuedSupply uint64          `json:"issued_supply"`
	SupplyCap    uint64          `json:"supply_cap"`
}

func (m *Group) NewGroupArgs() *NewGroupArgs {
	return &NewGroupArgs{
		Title:       m.Title,
		Description: m.Description,
		GroupName:   m.GroupName,
		Manager:     m.Manager,
		SupplyCap:   m.SupplyCap,
	}
}

type SetGroupManagerArgs struct {
	GroupName  eos.Name        `json:"group_name"`
	NewManager eos.AccountName `json:"new_manager"`
	Memo       string          `json:"memo"`
}

type NewGroupArgs struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	GroupName   eos.Name        `json:"group_name"`
	Manager     eos.AccountName `json:"manager"`
	SupplyCap   uint64          `json:"supply_cap"`
}

func (m *NewGroupArgs) String() string {
	return fmt.Sprintf(`
		Item{
			Title: %v,
			Description: %v,
			GroupName: %v,
			Manager: %v,
			SupplyCap: %v
		}
	`,
		m.Title,
		m.Description,
		m.GroupName,
		m.Manager,
		m.SupplyCap)
}

type Item struct {
	Serial uint64          `json:"serial"`
	Group  eos.Name        `json:"group"`
	Owner  eos.AccountName `json:"owner"`
}

func (m *Item) String() string {
	return fmt.Sprintf(`
		Item{
			Serial: %v,
			Group: %v,
			Owner: %v,
		}
	`,
		m.Serial,
		m.Group,
		m.Owner)
}

func (m *Item) MintItemArgs() *MintItemArgs {
	return &MintItemArgs{
		To:    m.Owner,
		Group: m.Group,
	}
}

type MintItemArgs struct {
	To    eos.AccountName `json:"to"`
	Group eos.Name        `json:"group_name"`
}

func (m *MintItemArgs) String() string {
	return fmt.Sprintf(`
		MintItemArgs{
			To: %v,
			Group: %v,
		}
	`,
		m.To,
		m.Group)
}

type Attribute struct {
	AttributeName eos.Name `json:"attribute_name"`
	Points        int64    `json:"points"`
	Locked        uint8    `json:"locked"`
}

func (m *Attribute) NewAttributeArgs(serial uint64, shared bool) *NewAttributeArgs {
	return &NewAttributeArgs{
		Serial:        serial,
		AttributeName: m.AttributeName,
		InitialPoints: m.Points,
		Shared:        shared,
	}
}

type NewAttributeArgs struct {
	Serial        uint64   `json:"serial"`
	AttributeName eos.Name `json:"attribute_name"`
	InitialPoints int64    `json:"initial_points"`
	Shared        bool     `json:"shared"`
}

func (m *NewAttributeArgs) String() string {
	return fmt.Sprintf(`
		NewAttributeArgs{
			Serial: %v,
			AttributeName: %v,
			InitialPoints: %v,
			Shared: %v,
		}
	`,
		m.Serial,
		m.AttributeName,
		m.InitialPoints,
		m.Shared)
}

type InitArgs struct {
	ContractName    string          `json:"contract_name"`
	ContractVersion string          `json:"contract_version"`
	InitialAdmin    eos.AccountName `json:"initial_admin"`
}

func NewInitArgs(initialAdmin eos.AccountName) *InitArgs {
	return &InitArgs{
		ContractName:    "Marble",
		ContractVersion: "v1.2.0",
		InitialAdmin:    initialAdmin,
	}
}

type TagEntry struct {
	Key   eos.Name `json:"first"`
	Value string   `json:"second"`
}

func (m *TagEntry) String() string {
	return fmt.Sprintf("%v=%v", m.Key, m.Value)
}

type AttributeEntry struct {
	Key   eos.Name `json:"first"`
	Value int64    `json:"second"`
}

func (m *AttributeEntry) String() string {
	return fmt.Sprintf("%v=%v", m.Key, m.Value)
}

type Tags []*TagEntry
type Attributes []*AttributeEntry

type RemoveFrameArgs struct {
	FrameName eos.Name `json:"frame_name"`
	Memo      string   `json:"memo"`
}

type Frame struct {
	FrameName         eos.Name   `json:"frame_name"`
	Group             eos.Name   `json:"group"`
	DefaultTags       Tags       `json:"default_tags"`
	DefaultAttributes Attributes `json:"default_attributes"`
}

func (m *Frame) String() string {
	return fmt.Sprintf(`
		Item{
			FrameName: %v,
			Group: %v,
			DefaultTags: %v,
			DefaultAttributes: %v,
		}
	`,
		m.FrameName,
		m.Group,
		m.DefaultTags,
		m.DefaultAttributes)
}

type QuickBuildArgs struct {
	FrameName          eos.Name        `json:"frame_name"`
	To                 eos.AccountName `json:"to"`
	OverrideTags       Tags            `json:"override_tags"`
	OverrideAttributes Attributes      `json:"override_attributes"`
}

func NewQuickBuildArgs(frameName eos.Name, to eos.AccountName) *QuickBuildArgs {
	return &QuickBuildArgs{
		FrameName:          frameName,
		To:                 to,
		OverrideTags:       make(Tags, 0),
		OverrideAttributes: make(Attributes, 0),
	}
}

type MarbleNFTContract struct {
	*contract.Contract
	*contract.SettingsContract
}

func NewMarbleNFTContract(eos *service.EOS, contractName string) *MarbleNFTContract {
	return &MarbleNFTContract{
		&contract.Contract{
			EOS:          eos,
			ContractName: contractName,
		},
		contract.NewSettingsContract(eos, contractName),
	}
}

func (m *MarbleNFTContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *MarbleNFTContract) Init(initArgs *InitArgs) (string, error) {
	return m.ExecAction(m.ContractName, "init", initArgs)
}

func (m *MarbleNFTContract) NewGroup(args *NewGroupArgs, admin eos.AccountName) (string, error) {
	return m.ExecAction(admin, "newgroup", args)
}

func (m *MarbleNFTContract) SetGroupManager(groupName eos.Name, newManager eos.AccountName, memo string, manager eos.AccountName) (string, error) {
	data := &SetGroupManagerArgs{
		GroupName:  groupName,
		NewManager: newManager,
		Memo:       memo,
	}
	return m.ExecAction(manager, "setmanager", data)
}

func (m *MarbleNFTContract) MintItem(args *MintItemArgs, manager eos.AccountName) (string, error) {
	return m.ExecAction(manager, "mintitem", args)
}

func (m *MarbleNFTContract) NewAttribute(args *NewAttributeArgs, manager eos.AccountName) (string, error) {
	return m.ExecAction(manager, "newattribute", args)
}

func (m *MarbleNFTContract) NewFrame(frame *Frame, manager eos.AccountName) (string, error) {
	return m.ExecAction(manager, "newframe", frame)
}

func (m *MarbleNFTContract) QuickBuild(args *QuickBuildArgs, manager eos.AccountName) (string, error) {
	return m.ExecAction(manager, "quickbuild", args)
}

func (m *MarbleNFTContract) NewDefaultNFTReward(groupName eos.Name, frameName eos.Name, manager eos.AccountName, authorizer eos.AccountName, maxSupply uint64) error {
	title := fmt.Sprintf("Group Name: %v, Frame Name: %v", groupName, frameName)
	groupArgs := &NewGroupArgs{
		Title:       title,
		Description: title,
		GroupName:   groupName,
		Manager:     manager,
		SupplyCap:   maxSupply,
	}
	_, err := m.NewGroup(groupArgs, authorizer)
	if err != nil {
		return fmt.Errorf("failed creating %v nft group: %v", groupName, err)
	}

	frameArgs := &Frame{
		FrameName:   frameName,
		Group:       groupName,
		DefaultTags: make(Tags, 0),
		DefaultAttributes: Attributes{{
			Key:   "reward",
			Value: 100,
		}},
	}
	_, err = m.NewFrame(frameArgs, manager)
	if err != nil {
		return fmt.Errorf("failed creating %v frame: %v", frameName, err)
	}
	return nil
}

func (m *MarbleNFTContract) Reset(limit uint64) (string, error) {
	return m.ExecAction(eos.AN(m.ContractName), "reset", limit)
}

func (m *MarbleNFTContract) RemoveFrame(frameName eos.Name, memo string, manager eos.AccountName) (string, error) {
	data := &RemoveFrameArgs{
		FrameName: frameName,
		Memo:      memo,
	}
	return m.ExecAction(manager, "rmvframe", data)
}

func (m *MarbleNFTContract) GetItems() ([]*Item, error) {
	return m.GetItemsReq(nil)
}

func (m *MarbleNFTContract) AreTablesEmpty() (bool, error) {
	items, err := m.GetItems()
	if err != nil {
		return false, err
	}
	groups, err := m.GetGroupsReq(nil)
	if err != nil {
		return false, err
	}
	frames, err := m.GetFramesReq(nil)
	if err != nil {
		return false, err
	}
	return len(items) == 0 && len(groups) == 0 && len(frames) == 0, err
}

func (m *MarbleNFTContract) GetLastItem() (*Item, error) {
	items, err := m.GetItemsReq(&eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(items) > 0 {
		return items[0], nil
	}
	return nil, nil
}

func (m *MarbleNFTContract) GetItemsByOwner(owner eos.AccountName) ([]*Item, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterItemsByOwner(request, owner)
	if err != nil {
		return nil, err
	}
	return m.GetItemsReq(request)
}

func (m *MarbleNFTContract) GetItemsByOwnerAndGroup(owner eos.AccountName, group eos.Name) ([]*Item, error) {
	items, err := m.GetItemsByOwner(owner)
	if err != nil {
		return nil, fmt.Errorf("failed getting items by owner: %v, group: %v, err: %v", owner, group, err)
	}
	filtered := make([]*Item, 0)
	for _, item := range items {
		if item.Group == group {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (m *MarbleNFTContract) FilterItemsByOwner(req *eos.GetTableRowsRequest, owner eos.AccountName) error {

	req.Index = "3"
	req.KeyType = "name"
	req.LowerBound = string(owner)
	req.UpperBound = string(owner)
	return nil
}

func (m *MarbleNFTContract) GetItemsByGroup(group eos.Name) ([]*Item, error) {
	request := &eos.GetTableRowsRequest{}
	err := m.FilterItemsByGroup(request, group)
	if err != nil {
		return nil, err
	}
	return m.GetItemsReq(request)
}

func (m *MarbleNFTContract) FilterItemsByGroup(req *eos.GetTableRowsRequest, group eos.Name) error {

	req.Index = "2"
	req.KeyType = "name"
	req.LowerBound = string(group)
	req.UpperBound = string(group)
	return nil
}

func (m *MarbleNFTContract) GetItemsReq(req *eos.GetTableRowsRequest) ([]*Item, error) {

	var items []*Item
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "items"
	err := m.GetTableRows(*req, &items)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return items, nil
}

func (m *MarbleNFTContract) GetGroupByName(group eos.Name) (*Group, error) {
	groups, err := m.GetGroupsReq(&eos.GetTableRowsRequest{
		LowerBound: string(group),
		UpperBound: string(group),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(groups) > 0 {
		return groups[0], nil
	}
	return nil, nil
}

func (m *MarbleNFTContract) GetGroupsReq(req *eos.GetTableRowsRequest) ([]*Group, error) {

	var groups []*Group
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "groups"
	err := m.GetTableRows(*req, &groups)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return groups, nil
}

func (m *MarbleNFTContract) GetFrameByName(frame eos.Name) (*Frame, error) {
	frames, err := m.GetFramesReq(&eos.GetTableRowsRequest{
		LowerBound: string(frame),
		UpperBound: string(frame),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(frames) > 0 {
		return frames[0], nil
	}
	return nil, nil
}

func (m *MarbleNFTContract) GetFramesReq(req *eos.GetTableRowsRequest) ([]*Frame, error) {

	var frames []*Frame
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "frames"
	err := m.GetTableRows(*req, &frames)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return frames, nil
}

func (m *MarbleNFTContract) GetAttribute(serial uint64, attr eos.Name) (*Attribute, error) {
	attributes, err := m.GetAttributesReq(serial, &eos.GetTableRowsRequest{
		Limit:      1,
		LowerBound: string(attr),
		UpperBound: string(attr),
	})
	if err != nil {
		return nil, err
	}
	if len(attributes) > 0 {
		return attributes[0], nil
	}
	return nil, nil
}

func (m *MarbleNFTContract) GetAttributesReq(serial uint64, req *eos.GetTableRowsRequest) ([]*Attribute, error) {

	var attributes []*Attribute
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "attributes"
	req.Scope = strconv.FormatUint(serial, 10)
	err := m.GetTableRows(*req, &attributes)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return attributes, nil
}
