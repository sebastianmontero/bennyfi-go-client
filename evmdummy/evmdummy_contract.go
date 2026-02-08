package evmdummy

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/sebastianmontero/eos-go"
	"github.com/sebastianmontero/eos-go-toolbox/contract"
	"github.com/sebastianmontero/eos-go-toolbox/service"
)

type EVMDummyContract struct {
	*contract.Contract
	callCounter uint64
}

func NewEVMDummyContract(eos *service.EOS, contractName string) *EVMDummyContract {
	return &EVMDummyContract{
		contract.NewContract(eos, contractName),
		0,
	}
}

func (m *EVMDummyContract) ExecAction(permissionLevel interface{}, action string, actionData interface{}) (string, error) {
	resp, err := m.Contract.ExecAction(permissionLevel, action, actionData)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("Tx ID: %v", resp.TransactionID), nil
}

type onBridgeMsg struct {
	Version   eos.Varuint32 `json:"version"`
	Receiver  eos.Name      `json:"receiver"`
	Sender    eos.HexBytes  `json:"sender"`
	Timestamp eos.TimePoint `json:"timestamp"`
	Value     eos.HexBytes  `json:"value"`
	Data      eos.HexBytes  `json:"data"`
}

func (m *EVMDummyContract) CallOnBridge(target eos.AccountName, message onBridgeMsg) (string, error) {
	actionData := struct {
		Target  eos.AccountName `json:"target"`
		Message onBridgeMsg     `json:"message"`
	}{
		Target:  target,
		Message: message,
	}
	return m.ExecAction(eos.AN(m.ContractName), "callonbridge", actionData)
}

func (m *EVMDummyContract) Unstake(target eos.AccountName, poolId uint64, stakeAmount uint64, rewardAmount uint64, bennyEvmBridgeAddress string) (string, error) {
	// app_type (4 bytes) + pool_id (32 bytes) + stake_amount (32 bytes) + reward_amount (32 bytes)
	data := make([]byte, 4+32+32+32)

	// UNSTAKE_SIGNATURE = 0xbc66aee6 (Little Endian)
	binary.LittleEndian.PutUint32(data[0:4], 0xbc66aee6)

	// pool_id (32 bytes) - Big Endian, last 8 bytes of slot
	binary.BigEndian.PutUint64(data[28:36], poolId)

	// stake_amount_evm (32 bytes) - Big Endian, last 8 bytes of slot
	binary.BigEndian.PutUint64(data[36+24:68], stakeAmount)

	// reward_amount_evm (32 bytes) - Big Endian, last 8 bytes of slot
	binary.BigEndian.PutUint64(data[68+24:100], rewardAmount)

	senderBytes, err := hex.DecodeString(strings.TrimPrefix(bennyEvmBridgeAddress, "0x"))
	if err != nil {
		return "", fmt.Errorf("invalid bridge address hex: %v", err)
	}

	msg := onBridgeMsg{
		Version:   0,
		Receiver:  eos.Name(target),
		Sender:    senderBytes,
		Timestamp: eos.TimePoint(time.Now().UnixMicro()),
		Value:     make([]byte, 0),
		Data:      data,
	}

	return m.CallOnBridge(target, msg)
}
