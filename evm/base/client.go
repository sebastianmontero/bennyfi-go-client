package base

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sebastianmontero/bennyfi-go-client/evm/common/eth"
)

// ReadClient represents a base client for read-only operations.
type ReadClient struct {
	Contract *bind.BoundContract
	Address  common.Address
	Client   eth.EthClient
}

// WriteClient represents a base client for write operations.
type WriteClient struct {
	*ReadClient
	Auth *bind.TransactOpts
}

// New creates a new WriteClient.
func New(rpcUrl string, privateKeyHex string, contractAddr common.Address, additionalABIs ...string) (*WriteClient, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}
	return NewWithClient(client, privateKeyHex, contractAddr, additionalABIs...)
}

// NewWithClient creates a new WriteClient with an existing ethclient.
func NewWithClient(client eth.EthClient, privateKeyHex string, contractAddr common.Address, additionalABIs ...string) (*WriteClient, error) {
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to get chain id: %w", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("failed to create auth: %w", err)
	}

	readClient, err := NewReadWithClient(client, contractAddr, additionalABIs...)
	if err != nil {
		return nil, err
	}

	return &WriteClient{
		ReadClient: readClient,
		Auth:     auth,
	}, nil
}

// NewClientFactory creates a factory function that produces WriteClients for a given address.
// It captures the client and credentials to reuse them for creating multiple clients.
func NewClientFactory(client eth.EthClient, privateKeyHex string, additionalABIs ...string) (func(common.Address) (*WriteClient, error), error) {
	// Verify private key validity early
	_, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		return nil, fmt.Errorf("invalid private key: %w", err)
	}

	return func(contractAddr common.Address) (*WriteClient, error) {
		return NewWithClient(client, privateKeyHex, contractAddr, additionalABIs...)
	}, nil
}

// NewRead creates a new ReadClient.
func NewRead(rpcUrl string, contractAddr common.Address, additionalABIs ...string) (*ReadClient, error) {
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}
	return NewReadWithClient(client, contractAddr, additionalABIs...)
}

// NewReadWithClient creates a new ReadClient with an existing ethclient.
func NewReadWithClient(client eth.EthClient, contractAddr common.Address, additionalABIs ...string) (*ReadClient, error) {
	mergedABI, err := MergeABIs(additionalABIs...)
	if err != nil {
		return nil, fmt.Errorf("failed to merge ABIs: %w", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(mergedABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse abi: %w", err)
	}

	contract := bind.NewBoundContract(contractAddr, parsedABI, client, client, client)

	return &ReadClient{
		Contract: contract,
		Address:  contractAddr,
		Client:   client,
	}, nil
}

// MergeABIs merges multiple ABI JSON strings into a single ABI string.
// It naively concatenates the arrays, verifying they are valid JSON arrays.
func MergeABIs(abis ...string) (string, error) {
	var builder strings.Builder
	builder.WriteString("[")
	first := true
	for _, abiStr := range abis {
		trimmed := strings.TrimSpace(abiStr)
		if len(trimmed) < 2 {
			continue
		} // "[]" or empty

		if !strings.HasPrefix(trimmed, "[") || !strings.HasSuffix(trimmed, "]") {
			return "", fmt.Errorf("invalid ABI JSON: must be an array")
		}

		inner := trimmed[1 : len(trimmed)-1]
		if strings.TrimSpace(inner) == "" {
			continue
		}

		if !first {
			builder.WriteString(",")
		}
		builder.WriteString(inner)
		first = false
	}
	builder.WriteString("]")
	return builder.String(), nil
}

// NewReadFromContract creates a new ReadClient from an existing bound contract.
func NewReadFromContract(client eth.EthClient, contract *bind.BoundContract, contractAddr common.Address) *ReadClient {
	return &ReadClient{
		Contract: contract,
		Address:  contractAddr,
		Client:   client,
	}
}

// TryGetRevertReason simulates a transaction using Call to extract a more detailed revert reason.
// Many EVM nodes drop the revert data during gas estimation (which Transact uses),
// but provide it reliably during eth_call.
func (c *WriteClient) TryGetRevertReason(origErr error, method string, params ...interface{}) error {
	callOpts := &bind.CallOpts{
		Pending: false,
		From:    c.Auth.From,
		Context: c.Auth.Context,
	}
	if callOpts.Context == nil {
		callOpts.Context = context.Background()
	}
	
	var out []interface{}
	callErr := c.Contract.Call(callOpts, &out, method, params...)
	if callErr != nil {
		// Parse the simulated call error
		return ParseError(callErr)
	}
	// If the call inexplicably succeeds, or origErr cannot be parsed better, fallback to parsing origErr
	return ParseError(origErr)
}

// ParseError attempts to extract a revert reason from an error returned by go-ethereum.
// It looks for a JSON-RPC DataError and tries to decode it.
func ParseError(err error) error {
	if err == nil {
		return nil
	}
	
	// Use reflection or interface assertion to get the RPC DataError if present.
	type dataError interface {
		ErrorData() interface{}
	}

	var dErr dataError
	if errors.As(err, &dErr) {
		data := dErr.ErrorData()
		if strData, ok := data.(string); ok && strData != "" {
			// Try to unpack the standard revert reason
			if len(strData) >= 138 && strings.HasPrefix(strData, "0x08c379a0") { // 0x08c379a0 is Error(string) signature
				bytesData := common.FromHex(strData)
				revertReason, unpackErr := abi.UnpackRevert(bytesData)
				if unpackErr == nil {
					return fmt.Errorf("execution reverted: %s", revertReason)
				}
			}
			return fmt.Errorf("execution reverted (raw data: %s): %w", strData, err)
		} else if data != nil {
			return fmt.Errorf("execution reverted (complex data structure): %w", err)
		}
	}
	
	return err
}
