package ethrpc

import (
	"github.com/GuiltyMorishita/jsonrpc"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

type EthRPC struct {
	RPCClient *jsonrpc.RPCClient
}

func NewEthRPC(endpoint string) *EthRPC {
	return &EthRPC{
		RPCClient: jsonrpc.NewRPCClient(endpoint),
	}
}

func (rpc *EthRPC) GetBalance(address, block string) (balance string, err error) {
	response, err := rpc.RPCClient.Call("eth_getBalance", address, block)
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
	response, err := rpc.RPCClient.Call("eth_getTransactionCount", address, block)
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

func (rpc *EthRPC) SendRawTransaction(data string) (txHash string, err error) {
	response, err := rpc.RPCClient.Call("eth_sendRawTransaction", data)
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
