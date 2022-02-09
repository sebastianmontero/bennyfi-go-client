package nft

import (
	"fmt"

	"github.com/eoscanada/eos-go"
)

const (
	NFTAttrName                  = "name"
	NFTAttrAccessLevel           = "access level"
	NFTAttrAllowedParallelRounds = "allowed parallel rounds"
	NFTTemplateSchema            = eos.Name("accessnfttmp")
)

type AccessNFT struct {
	Contract           eos.AccountName
	Collection         eos.Name
	CollectionNameAttr string
	NFT                *NFTContract
}

func (m *AccessNFT) SetupAccessNFTCollection() error {

	createCollectionArgs := &CreateCollectionArgs{
		BaseCollection: &BaseCollection{
			CollectionName:     m.Collection,
			Author:             m.Contract,
			AllowNotify:        0,
			AuthorizedAccounts: []eos.AccountName{m.Contract},
			NotifyAccounts:     []eos.AccountName{},
		},
		MarketFee: 0.10,
		Data:      AttributeMap{NFTAttrName: ToAtomicAttribute(m.CollectionNameAttr)},
	}
	_, err := m.NFT.CreateCollection(createCollectionArgs)
	if err != nil {
		return fmt.Errorf("failed creating %v collection: %v", m.Collection, err)
	}

	createSchemaArgs := &CreateSchemaArgs{
		AuthorizedCreator: m.Contract,
		CollectionName:    m.Collection,
		SchemaName:        NFTTemplateSchema,
		Format: []*Format{
			{
				Name: NFTAttrName,
				Type: "string",
			},
			{
				Name: NFTAttrAccessLevel,
				Type: "uint64",
			},
			{
				Name: NFTAttrAllowedParallelRounds,
				Type: "uint16",
			},
		},
	}
	_, err = m.NFT.CreateSchema(createSchemaArgs)
	if err != nil {
		return fmt.Errorf("failed creating access nft template schema: %v", err)
	}

	createTemplateArgs := &CreateTemplateArgs{
		AuthorizedCreator: m.Contract,
		CollectionName:    m.Collection,
		BaseTemplate: &BaseTemplate{
			SchemaName:   NFTTemplateSchema,
			Transferable: 1,
			Burnable:     0,
			MaxSupply:    0,
		},
		ImmutableData: AttributeMap{
			NFTAttrName: ToAtomicAttribute(m.CollectionNameAttr),
		},
	}
	_, err = m.NFT.CreateTemplate(createTemplateArgs)
	if err != nil {
		return fmt.Errorf("failed creating access nft template: %v", err)
	}

	return nil
}

type AccessNFTMintArgs struct {
	Owner                 eos.AccountName
	AccessLevel           uint64
	AllowedParallelRounds uint16
}

func (m *AccessNFTMintArgs) String() string {
	return fmt.Sprintf(
		`
		AccessNFTMintArgs {
				Owner: %v
				AccessLevel: %v
				AllowedParallelRounds: %v
			}
		`,
		m.Owner,
		m.AccessLevel,
		m.AllowedParallelRounds,
	)
}

func (m *AccessNFT) Mint(args *AccessNFTMintArgs) error {
	fmt.Println("Getting access NFT template...")
	template, err := m.NFT.GetTemplate(m.Collection)
	if err != nil {
		return fmt.Errorf("failed getting template for collection: %v, error: %v", m.Collection, err)
	}
	fmt.Println("template: ", template.TemplateId)
	mintAssetArgs := &MintAssetArgs{
		AuthorizedMinter: m.Contract,
		BaseAsset: &BaseAsset{
			CollectionName: m.Collection,
			SchemaName:     NFTTemplateSchema,
			TemplateId:     template.TemplateId,
		},
		NewAssetOwner: args.Owner,
		ImmutableData: AttributeMap{
			NFTAttrAccessLevel:           ToAtomicAttribute(args.AccessLevel),
			NFTAttrAllowedParallelRounds: ToAtomicAttribute(args.AllowedParallelRounds),
		},
	}
	_, err = m.NFT.MintAsset(mintAssetArgs)
	if err != nil {
		return fmt.Errorf("failed minting NFT asset: %v, error: %v", args, err)
	}
	return nil

}
