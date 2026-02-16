package eth

import (
	"context"
	"math/big"
	"sync/atomic"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// EthClient defines the interface for an ethereum client.
// It includes all methods from bind.ContractBackend and verify methods like ChainID.
type EthClient interface {
	bind.ContractBackend
	ChainID(ctx context.Context) (*big.Int, error)
	Close()
}

// MultiClient wraps multiple EthClients to provide failover support.
type MultiClient struct {
	clients []EthClient
	index   uint64 // atomic counter for round-robin or simple selection
}

// NewMultiClient creates a new MultiClient from a list of RPC URLs.
func NewMultiClient(urls []string) (*MultiClient, error) {
	clients := make([]EthClient, 0, len(urls))
	for _, url := range urls {
		client, err := ethclient.Dial(url)
		if err != nil {
			// If one fails, we close already opened ones and return error?
			// Or just skip/log? Strictly following "all or nothing" for now usually safer.
			for _, c := range clients {
				c.Close()
			}
			return nil, err
		}
		clients = append(clients, client)
	}
	return &MultiClient{clients: clients}, nil
}

// NewMultiClientFromClients creates a new MultiClient from existing clients.
func NewMultiClientFromClients(clients []EthClient) *MultiClient {
	return &MultiClient{clients: clients}
}

func (mc *MultiClient) Close() {
	for _, c := range mc.clients {
		c.Close()
	}
}

// executeTry tries to execute a function on clients in order, starting from current index/priority.
func (mc *MultiClient) executeTry(fn func(EthClient) error) error {
	// Simple failover: start at 0 (or current index if we implemented load balancing)
	// For now, let's try round robin start or just sequential failover.
	// Sequential failover is robust for "primary + backups".
	// Round robin is good for load balancing.
	// Let's do simple sequential failover for now.

	var firstErr error
	for _, client := range mc.clients {
		err := fn(client)
		if err == nil {
			return nil
		}
		if firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// CodeAt implementation.
func (mc *MultiClient) CodeAt(ctx context.Context, contract common.Address, blockNumber *big.Int) ([]byte, error) {
	var val []byte
	err := mc.executeTry(func(c EthClient) error {
		var err error
		val, err = c.CodeAt(ctx, contract, blockNumber)
		return err
	})
	return val, err
}

// CallContract implementation.
func (mc *MultiClient) CallContract(ctx context.Context, call ethereum.CallMsg, blockNumber *big.Int) ([]byte, error) {
	var val []byte
	err := mc.executeTry(func(c EthClient) error {
		var err error
		val, err = c.CallContract(ctx, call, blockNumber)
		return err
	})
	return val, err
}

// HeaderByNumber implementation.
func (mc *MultiClient) HeaderByNumber(ctx context.Context, number *big.Int) (*types.Header, error) {
	var val *types.Header
	err := mc.executeTry(func(c EthClient) error {
		var err error
		val, err = c.HeaderByNumber(ctx, number)
		return err
	})
	return val, err
}

// PendingCodeAt implementation.
func (mc *MultiClient) PendingCodeAt(ctx context.Context, account common.Address) ([]byte, error) {
	var val []byte
	err := mc.executeTry(func(c EthClient) error {
		var err error
		val, err = c.PendingCodeAt(ctx, account)
		return err
	})
	return val, err
}

// PendingNonceAt implementation.
func (mc *MultiClient) PendingNonceAt(ctx context.Context, account common.Address) (uint64, error) {
	var val uint64
	err := mc.executeTry(func(c EthClient) error {
		var err error
		val, err = c.PendingNonceAt(ctx, account)
		return err
	})
	return val, err
}

// SuggestGasPrice implementation.
func (mc *MultiClient) SuggestGasPrice(ctx context.Context) (*big.Int, error) {
	var val *big.Int
	err := mc.executeTry(func(c EthClient) error {
		var err error
		val, err = c.SuggestGasPrice(ctx)
		return err
	})
	return val, err
}

// SuggestGasTipCap implementation.
func (mc *MultiClient) SuggestGasTipCap(ctx context.Context) (*big.Int, error) {
	var val *big.Int
	err := mc.executeTry(func(c EthClient) error {
		var err error
		val, err = c.SuggestGasTipCap(ctx)
		return err
	})
	return val, err
}

// EstimateGas implementation.
func (mc *MultiClient) EstimateGas(ctx context.Context, call ethereum.CallMsg) (uint64, error) {
	var val uint64
	err := mc.executeTry(func(c EthClient) error {
		var err error
		val, err = c.EstimateGas(ctx, call)
		return err
	})
	return val, err
}

// SendTransaction implementation.
func (mc *MultiClient) SendTransaction(ctx context.Context, tx *types.Transaction) error {
	return mc.executeTry(func(c EthClient) error {
		return c.SendTransaction(ctx, tx)
	})
}

// FilterLogs implementation.
func (mc *MultiClient) FilterLogs(ctx context.Context, query ethereum.FilterQuery) ([]types.Log, error) {
	var val []types.Log
	err := mc.executeTry(func(c EthClient) error {
		var err error
		val, err = c.FilterLogs(ctx, query)
		return err
	})
	return val, err
}

// SubscribeFilterLogs implementation.
func (mc *MultiClient) SubscribeFilterLogs(ctx context.Context, query ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error) {
	// Subscription is tricky for failover.
	// For now, let's just subscribe to the first one that works?
	// Or maybe creating a wrapper subscription that handles reconnection is too complex for this task.
	// Let's just try to subscribe to the first working one.
	var sub ethereum.Subscription
	err := mc.executeTry(func(c EthClient) error {
		var err error
		sub, err = c.SubscribeFilterLogs(ctx, query, ch)
		return err
	})
	return sub, err
}

// ChainID implementation.
func (mc *MultiClient) ChainID(ctx context.Context) (*big.Int, error) {
	var val *big.Int
	err := mc.executeTry(func(c EthClient) error {
		var err error
		val, err = c.ChainID(ctx)
		return err
	})
	return val, err
}

// CurrentClient returns the next client in the list (round-robin) - useful if atomic access needed
func (mc *MultiClient) nextClient() EthClient {
	if len(mc.clients) == 0 {
		return nil
	}
	idx := atomic.AddUint64(&mc.index, 1) % uint64(len(mc.clients))
	return mc.clients[idx]
}

func (mc *MultiClient) GetFirstClient() (*ethclient.Client, error) {
	if len(mc.clients) == 0 {
		return nil, nil // Or error "no clients"
	}
	// We know NewMultiClient populates with *ethclient.Client, but interface is EthClient.
	// Type assertion needed.
	client, ok := mc.clients[0].(*ethclient.Client)
	if !ok {
		// Should not happen if created via NewMultiClient, but safe check.
		return nil, context.DeadlineExceeded // Just some error
	}
	return client, nil
}
