package bennyevmbridge

import (
	"fmt"

	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
	"github.com/sebastianmontero/eos-go/ecc"
)

type InitArgs struct {
	BridgeAddress        eos.Checksum160
	TokenRegistryAddress eos.Checksum160
	StakeLocalContract   eos.Name
	Version              string
	Admin                eos.Name
}

type BennyEVMBridgeContract struct {
	*contract.Contract
	callCounter uint64
}

func NewBennyEVMBridgeContract(eos *service.EOS, contractName string) *BennyEVMBridgeContract {
	return &BennyEVMBridgeContract{
		&contract.Contract{
			EOS:          eos,
			ContractName: contractName,
		},
		0,
	}
}

func (m *BennyEVMBridgeContract) NextCallCounter() uint64 {
	m.callCounter++
	return m.callCounter
}

func (m *BennyEVMBridgeContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

func (m *BennyEVMBridgeContract) ConfigureOpenPermission(publicKey *ecc.PublicKey) error {
	openActions := []string{
		"delexpirdoff",
	}
	err := m.EOS.CreateSimplePermission(m.ContractName, "open", publicKey)
	if err != nil {
		return fmt.Errorf("failed to create open permission, error: %v", err)
	}
	for _, action := range openActions {
		err = m.EOS.LinkPermission(m.ContractName, action, "open", false)
		if err != nil {
			return fmt.Errorf("failed to link open permission to the %v action, error: %v", action, err)
		}
	}
	return nil
}

func (m *BennyEVMBridgeContract) Init(args *InitArgs, authorizer interface{}) (string, error) {
	if authorizer == nil {
		authorizer = m.ContractName
	}
	return m.ExecAction(authorizer, "init", args)
}

func (m *BennyEVMBridgeContract) EVMNotify(authorizer interface{}, sender eos.Checksum160, msg []byte) (string, error) {
	actionData := struct {
		Sender eos.Checksum160
		Msg    []byte
	}{sender, msg}
	return m.ExecAction(authorizer, "evmnotify", actionData)
}
