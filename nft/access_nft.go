package nft

import (
	"fmt"

	"github.com/eoscanada/eos-go"
	"github.com/sebastianmontero/bennyfi-go-client/nft/marble"
)

const (
	NFTAttrAccessLevel = eos.Name("accesslevel")
)

type AccessNFT struct {
	Contract eos.AccountName
	Group    eos.Name
	Frame    eos.Name
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

func (m *AccessNFT) SetupAccessNFTFrame() error {

	args := marble.Frame{
		FrameName:   m.Frame,
		Group:       m.Group,
		DefaultTags: make(marble.Tags),
		DefaultAttributes: marble.Attributes{
			NFTAttrAccessLevel: 100,
		},
	}
	_, err := m.NFT.NewFrame(&args, m.Contract)
	if err != nil {
		return fmt.Errorf("failed creating %v frame: %v", m.Frame, err)
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
		buildArgs := marble.NewQuickBuildArgs(m.Frame, args.Owner)
		buildArgs.OverrideAttributes[NFTAttrAccessLevel] = int64(args.AccessLevel)
		_, err = m.NFT.QuickBuild(buildArgs, m.Contract)

		if err != nil {
			return fmt.Errorf("failed building item: %v, error: %v", buildArgs, err)
		}
	} else {
		fmt.Printf("Account: %v already has an access NFT\n", args.Owner)
	}
	return nil

}

func (m *AccessNFT) GetAccessNFTGroup() (*marble.Group, error) {
	return m.NFT.GetGroupByName(m.Group)
}

func (m *AccessNFT) GetAccessNFTs(owner eos.AccountName) ([]*marble.Item, error) {
	return m.NFT.GetItemsByOwnerAndGroup(owner, m.Group)
}
