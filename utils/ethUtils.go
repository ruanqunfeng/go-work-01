package utils

import (
	"context"
	"fmt"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/shopspring/decimal"
)

var Client *ethclient.Client

func init() {
	var err error
	Client, err = ethclient.Dial("wss://sepolia.gateway.tenderly.co")
	if err != nil {
		log.Fatal(err)
	}
}

func QueryBlockByNumber(blockNumber int64) *types.Block {
	block, err := Client.BlockByNumber(context.Background(), big.NewInt(blockNumber))
	if err != nil {
		log.Fatal(err)
	}

	return block
}

type RawEthTx struct {
	TxHash      string
	BlockNumber int64
	FromAddress string
	ToAddress   string
	Value       decimal.Decimal
	GasPrice    int64
	GasLimit    int64
	Gas         int64
	InputData   string
}

func QueryTransaction(hash string) *RawEthTx {
	txHash := common.HexToHash(hash)
	tx, isPending, err := Client.TransactionByHash(context.Background(), txHash)
	if err != nil {
		fmt.Println("获取交易失败:", err)
	}

	signer := types.LatestSignerForChainID(tx.ChainId())
	from, err := types.Sender(signer, tx)
	if err != nil {
		fmt.Println("获取发送方失败:", err)
		return nil
	}

	var toAddress common.Address
	if tx.To() != nil {
		toAddress = *tx.To()
	} else {
		toAddress = common.HexToAddress("0x") // 合约创建交易
	}

	// 获取交易的 gasPrice 和 gasLimit
	gasPrice := tx.GasPrice()
	gasLimit := tx.Gas()

	// 获取交易的收据
	receipt := QueryReceipt(hash)

	// 实际使用的 gas
	actualGasUsed := receipt.GasUsed

	// 将 uint64 转换为 *big.Int
	gasLimitBigInt := new(big.Int).SetUint64(gasLimit)
	actualGasUsedBigInt := new(big.Int).SetUint64(actualGasUsed)

	var fee *big.Int

	// 计算手续费
	if receipt.Status == types.ReceiptStatusFailed {
		// 如果交易失败，矿工收取 gasLimit * gasPrice
		fee = new(big.Int).Mul(gasLimitBigInt, gasPrice)
		fmt.Printf("Transaction failed. Fee: %s wei\n", fee.String())
	} else {
		// 如果交易成功，矿工收取实际使用的 gas * gasPrice
		fee = new(big.Int).Mul(actualGasUsedBigInt, gasPrice)
		fmt.Printf("Transaction succeeded. Fee: %s wei\n", fee.String())
	}

	// 输出结果
	fmt.Printf("交易状态: %v\n", isPending)
	fmt.Printf("发送方地址: %s\n", from.Hex())
	fmt.Printf("接收方地址: %s\n", toAddress.Hex())
	fmt.Printf("金额 (ETH): %s\n", weiToEth(tx.Value()))
	fmt.Printf("手续费 (ETH): %s\n", weiToEth(fee))

	ethTx := &RawEthTx{}

	ethTx.TxHash = hash
	ethTx.FromAddress = from.Hex()
	ethTx.ToAddress = toAddress.Hex()
	ethTx.Value, _ = decimal.NewFromString(weiToEth(tx.Value()))
	ethTx.GasPrice = gasPrice.Int64()
	ethTx.GasLimit = int64(gasLimit)
	ethTx.Gas = fee.Int64()
	ethTx.BlockNumber = receipt.BlockNumber.Int64()
	ethTx.InputData = string(tx.Data())

	return ethTx
}

// 5. 查询收据（通过交易哈希）
func QueryReceipt(txHash string) *types.Receipt {
	receipt, err := Client.TransactionReceipt(context.Background(), common.HexToHash(txHash))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Receipt Status: %d\n", receipt.Status)
	fmt.Printf("Gas Used: %d\n", receipt.GasUsed)

	return receipt
}

// 转换wei为ETH
func weiToEth(wei *big.Int) string {
	fbalance := new(big.Float)
	fbalance.SetString(wei.String())
	ethValue := new(big.Float).Quo(fbalance, big.NewFloat(1e18))
	return ethValue.Text('f', 18) // 保留18位小数
}
