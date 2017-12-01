package ethrpc

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/AltaApps/daikoku-server/helper"
	"github.com/GuiltyMorishita/ethutil"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/joho/godotenv"
	. "github.com/smartystreets/goconvey/convey"
)

func TestMain(m *testing.M) {
	_ = godotenv.Load(".test.env")
	exitCode := m.Run()
	defer os.Exit(exitCode)
}

func WithEthRPC(f func(rpc *EthRPC)) func() {
	return func() {
		f(NewEthRPC(os.Getenv("GETH_ENDPOINT")))
	}
}

func TestGetBalance(t *testing.T) {
	Convey("WithEthRPC", t, WithEthRPC(func(rpc *EthRPC) {
		Convey("Success", func() {
			balance, err := rpc.GetBalance(nil, "0x8FfCf7674ED27c7949Ceda9a0bD6799fe74aCf47", "latest")
			So(err, ShouldBeNil)
			So(balance, ShouldEqual, "0x56bc75e2d63100000") // 100 ETH
		})

		Convey("Empty Address", func() {
			_, err := rpc.GetBalance(nil, "", "latest")
			So(err.Error(), ShouldContainSubstring, "hex string has length 0")
		})

		Convey("Invalid Address", func() {
			_, err := rpc.GetBalance(nil, "InvalidAddress", "latest")
			So(err.Error(), ShouldContainSubstring, "cannot unmarshal hex string without 0x prefix ")
		})
	}))
}

func TestGetTransactionCount(t *testing.T) {
	Convey("WithEthRPC", t, WithEthRPC(func(rpc *EthRPC) {
		Convey("Success", func() {
			count, err := rpc.GetTransactionCount(nil, ethutil.AddressHex(1), "pending")
			So(err, ShouldBeNil)
			So(count, ShouldHaveSameTypeAs, uint64(1))
			So(count, ShouldBeGreaterThan, uint64(0))
		})

		Convey("Empty Address", func() {
			_, err := rpc.GetTransactionCount(nil, "", "pending")
			So(err.Error(), ShouldContainSubstring, "hex string has length 0")
		})

		Convey("Invalid Address", func() {
			_, err := rpc.GetTransactionCount(nil, "InvalidAddress", "pending")
			So(err.Error(), ShouldContainSubstring, "cannot unmarshal hex string without 0x prefix ")
		})
	}))
}

func TestSendRawTrancastion(t *testing.T) {
	Convey("WithEthRPC", t, WithEthRPC(func(rpc *EthRPC) {
		Convey("Success", func() {
			txHash, err := rpc.SendRawTransaction(nil, signedTransactionData(ethutil.PrivateKeyHex(1), ethutil.AddressHex(2)))
			So(err, ShouldBeNil)
			So(txHash, ShouldStartWith, "0x")
		})

		Convey("Empty Data", func() {
			_, err := rpc.SendRawTransaction(nil, "")
			So(err.Error(), ShouldContainSubstring, "EOF")
		})

		Convey("Invalid Data", func() {
			_, err := rpc.SendRawTransaction(nil, "InvalidData")
			So(err.Error(), ShouldContainSubstring, "cannot unmarshal hex string without 0x prefix ")
		})
	}))
}

func signedTransactionData(privateKeyHex string, toAddressHex string) (transactionData string) {
	from, _ := helper.PrivateKeyHexToAddressHex(privateKeyHex)
	nonce, _ := NewEthRPC(os.Getenv("GETH_ENDPOINT")).GetTransactionCount(nil, from, "pending")
	toAddress := common.HexToAddress(toAddressHex)
	value := big.NewInt(100000000000000000)
	gasLimit := big.NewInt(21000)
	gasPrice := big.NewInt(3000000000)
	data := []byte{}

	tx := types.NewTransaction(nonce, toAddress, value, gasLimit, gasPrice, data)
	signer := types.NewEIP155Signer(big.NewInt(15))
	privateKey, _ := crypto.HexToECDSA(privateKeyHex)
	signedTx, _ := types.SignTx(tx, signer, privateKey)
	ts := types.Transactions{signedTx}
	transactionData = fmt.Sprintf("0x%x", ts.GetRlp(0))
	return
}
