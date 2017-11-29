package ethrpc

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type EthError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type EthRPC struct {
	host string
}

func NewEthRPC(host string) *GethRPC {
	return &GethRPC{host: host}
}

func (rpc *EthRPC) call(method string, target interface{}, params ...interface{}) (err error) {
	result, err := rpc.Call(method, params...)
	if err != nil {
		return
	}

	if target == nil {
		return
	}

	if err = json.Unmarshal(result, target); err != nil {
		return
	}

	return
}

type request struct {
	ID      int           `json:"id"`
	JSONRPC string        `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
}

type response struct {
	ID      int             `json:"id"`
	JSONRPC string          `json:"jsonrpc"`
	Result  json.RawMessage `json:"result"`
	Error   *EthError       `json:"error"`
}

func (rpc *EthRPC) Call(method string, params ...interface{}) (json.RawMessage, error) {
	req := request{
		ID:      1,
		JSONRPC: "2.0",
		Method:  method,
		Params:  params,
	}

	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	resp, err := http.Post(rpc.url, "application/json", bytes.NewBuffer(body))
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if rpc.Debug {
		log.Printf("%s\nRequest: %s\nResponse: %s\n", method, body, data)
	}

	resp := new(ethResponse)
	if err := json.Unmarshal(data, resp); err != nil {
		return nil, err
	}

	if resp.Error != nil {
		return nil, *resp.Error
	}

	return resp.Result, nil

}

func main() {
	fmt.Println("vim-go")
}
