package ethrpc

import (
	"context"

	"github.com/GuiltyMorishita/jsonrpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"google.golang.org/appengine/urlfetch"
)

// EthRPCer ...
type EthRPCer interface {
	GetBalance(aeCtx context.Context, address, block string) (balance string, err error)
	GetTransactionCount(aeCtx context.Context, address, block string) (count uint64, err error)
	SendRawTransaction(aeCtx context.Context, txData string) (txHash string, err error)
}

// EthRPC ...
type EthRPC struct {
	rpcClient *jsonrpc.RPCClient
}

// NewEthRPC ...
func NewEthRPC(endpoint string) *EthRPC {
	return &EthRPC{
		rpcClient: jsonrpc.NewRPCClient(endpoint),
	}
}

func (rpc *EthRPC) GetBalance(aeCtx context.Context, address, block string) (balance string, err error) {
	rpc.setAppEngineContext(aeCtx)
	response, err := rpc.rpcClient.Call("eth_getBalance", address, block)
	if err != nil {
		return
	}

	if response.Error != nil {
		err = response.Error
		return
	}

	response.GetObject(&balance)
	return
}

func (rpc *EthRPC) GetTransactionCount(aeCtx context.Context, address, block string) (count uint64, err error) {
	rpc.setAppEngineContext(aeCtx)
	response, err := rpc.rpcClient.Call("eth_getTransactionCount", address, block)
	if err != nil {
		return
	}

	if response.Error != nil {
		err = response.Error
		return
	}

	var countHex string
	response.GetObject(&countHex)
	count, _ = hexutil.DecodeUint64(countHex)
	return
}

func (rpc *EthRPC) SendRawTransaction(aeCtx context.Context, txData string) (txHash string, err error) {
	rpc.setAppEngineContext(aeCtx)
	response, err := rpc.rpcClient.Call("eth_sendRawTransaction", txData)
	if err != nil {
		return
	}

	if response.Error != nil {
		err = response.Error
		return
	}

	response.GetObject(&txHash)
	return
}

func (rpc *EthRPC) setAppEngineContext(aeCtx context.Context) {
	if aeCtx != nil {
		rpc.rpcClient.SetHTTPClient(urlfetch.Client(aeCtx))
	}
}
