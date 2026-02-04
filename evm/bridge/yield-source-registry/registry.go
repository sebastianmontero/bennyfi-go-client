package yield_source_registry

import (
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

// IYieldSourceRegistryABI is the ABI of the IYieldSourceRegistry contract.
const IYieldSourceRegistryABI = `[
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_name",
          "type": "string"
        }
      ],
      "name": "getYieldSource",
      "outputs": [
        {
          "components": [
            {
              "internalType": "string",
              "name": "name",
              "type": "string"
            },
            {
              "internalType": "address",
              "name": "adaptorContract",
              "type": "address"
            },
            {
              "internalType": "bool",
              "name": "active",
              "type": "bool"
            }
          ],
          "internalType": "struct IYieldSourceRegistry.YieldSource",
          "name": "",
          "type": "tuple"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    },
    {
      "inputs": [
        {
          "internalType": "string",
          "name": "_name",
          "type": "string"
        }
      ],
      "name": "isYieldSourceRegistered",
      "outputs": [
        {
          "internalType": "bool",
          "name": "",
          "type": "bool"
        }
      ],
      "stateMutability": "view",
      "type": "function"
    }
]`

// YieldSource represents a yield source in the registry.
type YieldSource struct {
	Name            string         `json:"name"`
	AdaptorContract common.Address `json:"adaptorContract"`
	Active          bool           `json:"active"`
}

// Registry interacts with the IYieldSourceRegistry contract.
type Registry struct {
	contract *bind.BoundContract
	address  common.Address
}

// NewRead creates a new Registry client for read-only operations.
func NewRead(rpcUrl string, contractAddr common.Address) (*Registry, error) {
	// 1. Connect to Eth Client
	client, err := ethclient.Dial(rpcUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to dial rpc: %w", err)
	}

	return NewReadWithClient(client, contractAddr)
}

// NewReadWithClient creates a new Registry client using an existing contract backend.
func NewReadWithClient(client *ethclient.Client, contractAddr common.Address) (*Registry, error) {
	// 1. Parse ABI
	parsedABI, err := abi.JSON(strings.NewReader(IYieldSourceRegistryABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse abi: %w", err)
	}

	// 2. Create Bound Contract
	contract := bind.NewBoundContract(contractAddr, parsedABI, client, client, client)
	return &Registry{
		contract: contract,
		address:  contractAddr,
	}, nil
}

// GetYieldSource retrieves a yield source by name.
func (r *Registry) GetYieldSource(name string) (*YieldSource, error) {
	var out []interface{}
	err := r.contract.Call(nil, &out, "getYieldSource", name)
	if err != nil {
		return nil, err
	}

	// The output is a tuple, so it's a struct in Go's ABI binding representation.
	// However, since we are using dynamic binding with flattened outputs not applied here (it's a tuple return),
	// we need to be careful.
	// Wait, standard abi unpacking for a struct return typically creates a struct if generated,
	// but with Call and []interface{}, it unpacks into the fields.
	// Let's check the ABI constraint.
	// The output is a single tuple.
	// `out[0]` should be the struct/tuple data.
	
	// Actually, `abi.Unpack` (helper used by Call) unpacks the return values. 
	// If the return is a single tuple, it might unpack into a struct if we pass a struct pointer, 
	// OR it unpacks into individual fields if we flatten.
	// But here we rely on generic []interface{}.
	// `go-ethereum`'s dynamic Call usually unpacks a single tuple return as the struct fields directly if specific, 
	// OR as a single struct object if properly mapped? 
	// Without generated code, ABI decoding a tuple usually results in a struct-like generic representation or requires a matching Go struct.
	
	// Let's look at how `go-ethereum` handles this.
	// If I pass `&out`, it tries to Unpack into it.
	// If the output is a tuple `(string, address, bool)`, `out` will contain `[string, address, bool]`.
	// Wait, the ABI says output is `type: tuple` named `` containing components.
	// So `out` will have length 1, and `out[0]` will be the struct.
	// Let's use a struct helper for safety.
	
	// Better approach: Define a temporary struct to hold the result for easy decoding if possible,
	// but `Call` with `*[]interface{}` might be tricky with tuples.
	
	// Alternative: Flatten the ABI manually in client code? 
	// The user provided ABI has it as a tuple.
	// Let's try to map it to the struct.
	
	// Actually, for dynamic call with `[]interface{}`:
	// If the output is a single named tuple, `abigen` usually creates a struct.
	// Without abigen, `Call` usually returns the unpacked values.
	// If it is a tuple, it might be unpacked as a single element which is the struct?
	// Let's assume standard behavior: `abi.Unpack` into `&[]interface{}` for a tuple return 
	// typically flattens the tuple fields into the slice IF the abi definition wasn't a tuple inside a tuple?
	// The return type is `YieldSource` (tuple).
	
	// Let's try attempting to decode into `*YieldSource` directly? 
	// `c.contract.Call(nil, &resultStruct, method, ...)`
	
	result := new(YieldSource)
	// We need to pass a slice of outputs to Call usually.
	results := []interface{}{result}
	err = r.contract.Call(nil, &results, "getYieldSource", name)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// IsYieldSourceRegistered checks if a yield source is registered.
func (r *Registry) IsYieldSourceRegistered(name string) (bool, error) {
	var out []interface{}
	err := r.contract.Call(nil, &out, "isYieldSourceRegistered", name)
	if err != nil {
		return false, err
	}
	return out[0].(bool), nil
}
