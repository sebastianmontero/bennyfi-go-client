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
package nft

import (
	"fmt"
	"strconv"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type Format struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type BaseCollection struct {
	Author             eos.AccountName   `json:"author"`
	CollectionName     eos.Name          `json:"collection_name"`
	AllowNotify        uint8             `json:"allow_notify"`
	AuthorizedAccounts []eos.AccountName `json:"authorized_accounts"`
	NotifyAccounts     []eos.AccountName `json:"notify_accounts"`
}

func (m *BaseCollection) Clone() *BaseCollection {
	return &BaseCollection{
		Author:             m.Author,
		CollectionName:     m.CollectionName,
		AllowNotify:        m.AllowNotify,
		AuthorizedAccounts: m.AuthorizedAccounts,
		NotifyAccounts:     m.NotifyAccounts,
	}
}

type Collection struct {
	*BaseCollection
	MarketFee      string  `json:"market_fee"`
	SerializedData []uint8 `json:"serialized_data"`
}

func (m *Collection) MarketFeeFloat() float64 {
	v, err := strconv.ParseFloat(m.MarketFee, 64)
	if err != nil {
		panic(fmt.Sprintf("Invalid float value: %v", m.MarketFee))
	}
	return v
}

type CreateCollectionArgs struct {
	*BaseCollection
	MarketFee float64      `json:"market_fee"`
	Data      AttributeMap `json:"data"`
}

func (m *CreateCollectionArgs) Collection() *Collection {
	return &Collection{
		BaseCollection: m.BaseCollection.Clone(),
		MarketFee:      fmt.Sprint(m.MarketFee),
	}
}

type Schema struct {
	SchemaName eos.Name  `json:"schema_name"`
	Format     []*Format `json:"format"`
}

type CreateSchemaArgs struct {
	AuthorizedCreator eos.AccountName `json:"authorized_creator"`
	CollectionName    eos.Name        `json:"collection_name"`
	SchemaName        eos.Name        `json:"schema_name"`
	Format            []*Format       `json:"schema_format"`
}

func (m *CreateSchemaArgs) Schema() *Schema {
	return &Schema{
		SchemaName: m.SchemaName,
		Format:     m.Format,
	}
}

type BaseTemplate struct {
	SchemaName   eos.Name `json:"schema_name"`
	Transferable uint8    `json:"transferable"`
	Burnable     uint8    `json:"burnable"`
	MaxSupply    uint32   `json:"max_supply"`
}

func (m *BaseTemplate) Clone() *BaseTemplate {
	return &BaseTemplate{
		SchemaName:   m.SchemaName,
		Transferable: m.Transferable,
		Burnable:     m.Burnable,
		MaxSupply:    m.MaxSupply,
	}
}

type Template struct {
	*BaseTemplate
	TemplateId              int32   `json:"template_id"`
	IssuedSupply            uint32  `json:"issued_supply"`
	ImmutableSerializedData []uint8 `json:"immutable_serialized_data"`
}

type CreateTemplateArgs struct {
	AuthorizedCreator eos.AccountName `json:"authorized_creator"`
	CollectionName    eos.Name        `json:"collection_name"`
	*BaseTemplate
	ImmutableData AttributeMap `json:"immutable_data"`
}

func (m *CreateTemplateArgs) Template() *Template {
	return &Template{
		BaseTemplate: m.BaseTemplate.Clone(),
		IssuedSupply: 0,
	}
}

type BaseAsset struct {
	CollectionName eos.Name `json:"collection_name"`
	SchemaName     eos.Name `json:"schema_name"`
	TemplateId     int32    `json:"template_id"`
}

func (m *BaseAsset) Clone() *BaseAsset {
	return &BaseAsset{
		CollectionName: m.CollectionName,
		SchemaName:     m.SchemaName,
		TemplateId:     m.TemplateId,
	}
}

type Asset struct {
	*BaseAsset
	AssetId                 string          `json:"asset_id"`
	RamPayer                eos.AccountName `json:"ram_payer"`
	BackedTokens            []eos.Asset     `json:"backed_tokens"`
	ImmutableSerializedData []uint8         `json:"immutable_serialized_data"`
	MutableSerializedData   []uint8         `json:"mutable_serialized_data"`
}

type MintAssetArgs struct {
	AuthorizedMinter eos.AccountName `json:"authorized_minter"`
	*BaseAsset
	NewAssetOwner eos.AccountName `json:"new_asset_owner"`
	ImmutableData AttributeMap    `json:"immutable_data"`
	MutableData   AttributeMap    `json:"mutable_data"`
	TokensToBack  []eos.Asset     `json:"tokens_to_back"`
}

func (m *MintAssetArgs) Asset() *Asset {
	return &Asset{
		BaseAsset:    m.BaseAsset.Clone(),
		RamPayer:     m.AuthorizedMinter,
		BackedTokens: m.TokensToBack,
	}
}

type NFTContract struct {
	*contract.Contract
}

func NewNFTContract(eos *service.EOS, contractName string) *NFTContract {
	return &NFTContract{
		&contract.Contract{
			EOS:          eos,
			ContractName: contractName,
		},
	}
}

func (m *NFTContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *NFTContract) Init() (string, error) {
	return m.ExecAction(m.ContractName, "init", nil)
}

func (m *NFTContract) CreateCollection(collection *CreateCollectionArgs) (string, error) {
	return m.ExecAction(collection.Author, "createcol", collection)
}

func (m *NFTContract) CreateSchema(schema *CreateSchemaArgs) (string, error) {
	return m.ExecAction(schema.AuthorizedCreator, "createschema", schema)
}

func (m *NFTContract) CreateTemplate(template *CreateTemplateArgs) (string, error) {
	return m.ExecAction(template.AuthorizedCreator, "createtempl", template)
}

func (m *NFTContract) MintAsset(asset *MintAssetArgs) (string, error) {
	return m.ExecAction(asset.AuthorizedMinter, "mintasset", asset)
}

func (m *NFTContract) EditCollectionFormats(formats []*Format) (string, error) {
	actionData := make(map[string]interface{})
	actionData["collection_format_extension"] = formats
	return m.ExecAction(m.ContractName, "admincoledit", formats)
}

func (m *NFTContract) InitCollectionFormats() (string, error) {
	formats := []*Format{
		{
			Name: "name",
			Type: "string",
		},
		{
			Name: "img",
			Type: "ipfs",
		},
		{
			Name: "description",
			Type: "string",
		},
		{
			Name: "url",
			Type: "string",
		},
	}
	return m.EditCollectionFormats(formats)
}

func (m *NFTContract) GetAssets(owner eos.AccountName) ([]Asset, error) {
	return m.GetAssetsReq(owner, nil)
}

func (m *NFTContract) GetLastAsset(owner eos.AccountName) (*Asset, error) {
	assets, err := m.GetAssetsReq(owner, &eos.GetTableRowsRequest{
		Reverse: true,
		Limit:   1,
	})
	if err != nil {
		return nil, err
	}
	if len(assets) > 0 {
		return &assets[0], nil
	}
	return nil, nil
}

func (m *NFTContract) GetAssetsReq(owner eos.AccountName, req *eos.GetTableRowsRequest) ([]Asset, error) {

	var assets []Asset
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "assets"
	req.Scope = string(owner)
	err := m.GetTableRows(*req, &assets)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return assets, nil
}

func (m *NFTContract) GetCollectionByName(collection eos.Name) (*Collection, error) {
	collections, err := m.GetCollectionsReq(&eos.GetTableRowsRequest{
		LowerBound: string(collection),
		UpperBound: string(collection),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(collections) > 0 {
		return &collections[0], nil
	}
	return nil, nil
}

func (m *NFTContract) GetCollectionsReq(req *eos.GetTableRowsRequest) ([]Collection, error) {

	var collections []Collection
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "collections"
	err := m.GetTableRows(*req, &collections)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return collections, nil
}

func (m *NFTContract) GetTemplate(collection eos.Name) (*Template, error) {
	templates, err := m.GetTemplatesReq(collection, &eos.GetTableRowsRequest{
		Limit: 1,
	})
	if err != nil {
		return nil, err
	}
	if len(templates) > 0 {
		return &templates[0], nil
	}
	return nil, nil
}

func (m *NFTContract) GetTemplatesReq(collectionName eos.Name, req *eos.GetTableRowsRequest) ([]Template, error) {

	var templates []Template
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "templates"
	req.Scope = string(collectionName)
	err := m.GetTableRows(*req, &templates)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return templates, nil
}

func (m *NFTContract) GetSchemaByName(collection, schemaName eos.Name) (*Schema, error) {
	schemas, err := m.GetSchemasReq(collection, &eos.GetTableRowsRequest{
		LowerBound: string(schemaName),
		UpperBound: string(schemaName),
		Limit:      1,
	})
	if err != nil {
		return nil, err
	}
	if len(schemas) > 0 {
		return &schemas[0], nil
	}
	return nil, nil
}

func (m *NFTContract) GetSchemasReq(collectionName eos.Name, req *eos.GetTableRowsRequest) ([]Schema, error) {

	var schemas []Schema
	if req == nil {
		req = &eos.GetTableRowsRequest{}
	}
	req.Table = "schemas"
	req.Scope = string(collectionName)
	err := m.GetTableRows(*req, &schemas)
	if err != nil {
		return nil, fmt.Errorf("get table rows %v", err)
	}
	return schemas, nil
}
