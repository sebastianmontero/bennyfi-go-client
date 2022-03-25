package nft

import (
	"fmt"
	"strings"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/bennyfi-go-client/nft/marble"
)

const (
	NFTAttrAccessLevel = eos.Name("accesslevel")
)

type AccessNFT struct {
	Contract eos.AccountName
	Group    eos.Name
	NFT      *marble.MarbleNFTContract
}

func (m *AccessNFT) SetupAccessNFTGroup() error {

	args := marble.NewGroupArgs{
		Title:       "Access NFT",
		Description: "Access NFT",
		GroupName:   m.Group,
		Manager:     m.Contract,
		SupplyCap:   1000,
	}
	_, err := m.NFT.NewGroup(&args, m.Contract)
	if err != nil {
		return fmt.Errorf("failed creating %v nft group: %v", m.Group, err)
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
	fmt.Printf("Checking if account: %v already has an access NFT...\n", args.Owner)
	items, err := m.NFT.GetItemsByOwnerAndGroup(args.Owner, m.Group)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		fmt.Printf("Account: %v does not have an access NFT, minting...\n", args.Owner)
		mintArgs := &marble.MintItemArgs{
			To:    args.Owner,
			Group: m.Group,
		}
		_, err = m.NFT.MintItem(mintArgs, m.Contract)

		if err != nil {
			return fmt.Errorf("failed miniting item: %v, error: %v", mintArgs, err)
		}

		fmt.Printf("Minted getting access NFT for account: %v...\n", args.Owner)
		items, err = m.NFT.GetItemsByOwnerAndGroup(args.Owner, m.Group)
		if err != nil {
			return err
		}
	}
	item := items[0]
	fmt.Printf("Adding access level attribute with value: %v to item: %v with owner: %v...\n", args.AccessLevel, item, args.Owner)
	newAttrArgs := &marble.NewAttributeArgs{
		Serial:        item.Serial,
		AttributeName: NFTAttrAccessLevel,
		InitialPoints: int64(args.AccessLevel),
	}
	_, err = m.NFT.NewAttribute(newAttrArgs, m.Contract)

	if err != nil {
		if strings.Contains(err.Error(), "shared attributes already exists") {
			fmt.Printf("Item: %v already has access level attribute\n", item)
		} else {
			return fmt.Errorf("failed adding access level attr: %v, error: %v", newAttrArgs, err)
		}
	}
	return nil

}

func (m *AccessNFT) GetAccessNFTGroup() (*marble.Group, error) {
	return m.NFT.GetGroupByName(m.Group)
}

func (m *AccessNFT) GetAccessNFTs(owner eos.AccountName) ([]*marble.Item, error) {
	return m.NFT.GetItemsByOwnerAndGroup(owner, m.Group)
}
