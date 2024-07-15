package bridge

import (
	"time"

	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/util"
	"github.com/sebastianmontero/eos-go/ecc"
)

type ActionProof struct {
	Action      *eos.Action             `json:"action"`
	ActReceipt  *eos.ActionTraceReceipt `json:"receipt"`
	ReturnValue []byte                  `json:"returnvalue"`
	AmpPoofPath []eos.Checksum256       `json:"ampproofpath"`
}

func NewActionProofFromActionData(account eos.AccountName, actionName eos.ActionName, globalSequence uint64, actionData interface{}) *ActionProof {
	action := &eos.Action{
		Account: account,
		Name:    actionName,
		Authorization: []eos.PermissionLevel{
			{
				Actor:      account,
				Permission: eos.PermissionName("active"),
			},
		},
		ActionData: eos.NewActionData(actionData),
	}
	actionReceipt := &eos.ActionTraceReceipt{
		Receiver:        eos.AN("rec1"),
		ActionDigest:    getDummyChecksum(),
		GlobalSequence:  eos.Uint64(globalSequence),
		ReceiveSequence: eos.Uint64(0),
		AuthSequence:    []eos.TransactionTraceAuthSequence{},
		CodeSequence:    eos.Varuint32(0),
		ABISequence:     eos.Varuint32(0),
	}

	return &ActionProof{
		Action:      action,
		ActReceipt:  actionReceipt,
		ReturnValue: []byte{},
		AmpPoofPath: []eos.Checksum256{},
	}

}

type LightProof struct {
	ChainID     eos.Checksum256   `json:"chain_id"`
	BlockHeader *eos.BlockHeader  `json:"header"`
	Root        eos.Checksum256   `json:"root"`
	BmProofPath []eos.Checksum256 `json:"bmproofpath"`
}

func NewDummyLightProof(chainId eos.Checksum256) *LightProof {
	return &LightProof{
		ChainID:     chainId,
		BlockHeader: NewDummyBlockHeader(),
		Root:        getDummyChecksum(),
		BmProofPath: []eos.Checksum256{},
	}
}

func NewDummyProducerSchedule() *eos.ProducerSchedule {
	return &eos.ProducerSchedule{
		Producers: []eos.ProducerKey{},
	}
}

func NewDummyBlockHeader() *eos.BlockHeader {
	return &eos.BlockHeader{
		Timestamp: eos.BlockTimestamp{
			Time: time.Now(),
		},
		Producer:         eos.AccountName("p1"),
		Previous:         getDummyChecksum(),
		TransactionMRoot: getDummyChecksum(),
		ActionMRoot:      getDummyChecksum(),
		NewProducersV1:   NewDummyProducerSchedule(),
		HeaderExtensions: []*eos.Extension{},
	}
}

type SignedBlockHeader struct {
	BlockHeader        *eos.BlockHeader `json:"header"`
	ProducerSignatures []ecc.Signature  `json:"producer_signatures"`
	PreviousBmroot     eos.Checksum256  `json:"previous_bmroot"`
	BmProofPath        []uint16         `json:"bmproofpath"`
}

func NewDummySignedBlockHeader() *SignedBlockHeader {
	return &SignedBlockHeader{
		BlockHeader:        NewDummyBlockHeader(),
		ProducerSignatures: []ecc.Signature{},
		PreviousBmroot:     getDummyChecksum(),
		BmProofPath:        []uint16{},
	}
}

type AnchorBlock struct {
	Block       *SignedBlockHeader `json:"block"`
	ActiveNodes []uint16           `json:"active_nodes"`
	NodeCount   uint64             `json:"node_count"`
}

func NewDummyAnchorBlock() *AnchorBlock {
	return &AnchorBlock{
		Block:       NewDummySignedBlockHeader(),
		ActiveNodes: []uint16{},
		NodeCount:   0,
	}
}

type HeavyProof struct {
	ChainID      eos.Checksum256      `json:"chain_id"`
	Hashes       []eos.Checksum256    `json:"hashes"`
	BlockToProve *AnchorBlock         `json:"blocktoprove"`
	BftProof     []*SignedBlockHeader `json:"bftproof"`
}

func NewDummyHeavyProof(chainId eos.Checksum256) *HeavyProof {
	return &HeavyProof{
		ChainID:      chainId,
		Hashes:       []eos.Checksum256{},
		BlockToProve: NewDummyAnchorBlock(),
		BftProof:     []*SignedBlockHeader{},
	}
}

func getDummyChecksum() eos.Checksum256 {
	checksum, err := util.ToChecksum256("bca376f206b8fc25a6ed44dbdc66547c36c6c33e3a119ffbeaef943642f0e906")
	if err != nil {
		panic(err)
	}
	return checksum
}
