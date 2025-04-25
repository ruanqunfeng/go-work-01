package controller

import (
	"encoding/json"
	"fmt"
	"go-work-01/model"
	"go-work-01/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func GetTx(c *gin.Context) {
	txH := c.Query("txHash")
	if txH == "" || len(txH) != 66 || !strings.HasPrefix(txH, "0x") {
		c.JSON(400, gin.H{
			"error": "不正确的hash",
		})
		return
	}

	// 先从Mysql查
	tx, err := model.GetTxByHash(txH)

	if err == nil {
		c.JSON(200, gin.H{
			"tx": tx,
		})
		return
	}

	// 查不到查链上, 更新Mysql
	rawTx := utils.QueryTransaction(txH)
	if rawTx == nil {
		c.JSON(400, gin.H{
			"error": "tx not found",
		})
		return
	}

	tx = &model.EthTx{}
	tx.BlockNumber = rawTx.BlockNumber
	tx.FromAddress = rawTx.FromAddress
	tx.Gas = rawTx.Gas
	tx.GasLimit = rawTx.GasLimit
	tx.GasPrice = rawTx.GasPrice
	tx.InputData = rawTx.InputData
	tx.ToAddress = rawTx.ToAddress
	tx.TxHash = rawTx.TxHash
	tx.Value = rawTx.Value

	err = model.CreateTx(tx)
	if err != nil {
		fmt.Println("insert tx failed", err)
	}

	c.JSON(200, gin.H{
		"tx": tx,
	})
}

func GetTxReceipt(c *gin.Context) {
	txH := c.Query("txHash")
	if txH == "" || len(txH) != 66 || !strings.HasPrefix(txH, "0x") {
		c.JSON(400, gin.H{
			"error": "不正确的hash",
		})
		return
	}
	// 先从Mysql查
	txReceipt, err := model.GetTxReceipt(txH)
	if err == nil {
		c.JSON(200, gin.H{
			"txReceipt": txReceipt,
		})
		return
	}

	// 查不到查链上, 更新Mysql
	receipt := utils.QueryReceipt(txH)
	if receipt == nil {
		c.JSON(400, gin.H{
			"error": "tx not found",
		})
		return
	}
	txReceipt = &model.EthTxReceipt{}
	txReceipt.ContractAddress = receipt.ContractAddress.Hex()
	txReceipt.GasUsed = int64(receipt.GasUsed)
	txReceipt.Status = int(receipt.Status)
	txReceipt.TxHash = receipt.TxHash.Hex()
	l, err := json.Marshal(receipt.Logs)
	if err == nil {
		txReceipt.Logs = datatypes.JSON(l)
	}
	err = model.CreateTxReceipt(txReceipt)

	if err != nil {
		fmt.Println("insert txReceipt failed", err)
	}
	c.JSON(200, gin.H{
		"txReceipt": receipt,
	})
}
