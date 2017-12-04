package ethrpc

import (
	"net/http"

	"github.com/GuiltyMorishita/jsonrpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"google.golang.org/appengine"
	"google.golang.org/appengine/urlfetch"
)

// EthRPCer ...
type EthRPCer interface {
	GetBalance(address, block string) (balance string, err error)
	GetTransactionCount(address, block string) (count uint64, err error)
	SendRawTransaction(txData string) (txHash string, err error)
	UseAppEngineContext(r *http.Request)
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

func (rpc *EthRPC) GetBalance(address, block string) (balance string, err error) {
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

func (rpc *EthRPC) GetTransactionCount(address, block string) (count uint64, err error) {
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

func (rpc *EthRPC) SendRawTransaction(txData string) (txHash string, err error) {
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

func (rpc *EthRPC) UseAppEngineContext(r *http.Request) {
	aeCtx := appengine.NewContext(r)
	rpc.rpcClient.SetHTTPClient(urlfetch.Client(aeCtx))
}
