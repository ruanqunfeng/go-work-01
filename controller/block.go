package controller

import (
	"encoding/json"
	"fmt"
	"go-work-01/model"
	"go-work-01/utils"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
)

func AssembleBlock(ethBlock *types.Block) *model.EthBlock {
	block := &model.EthBlock{}

	block.BlockNumber = ethBlock.Number().Int64()
	block.BlockHash = ethBlock.Hash().Hex()
	block.ParentHash = ethBlock.ParentHash().Hex()
	block.GasUsed = int64(ethBlock.GasUsed())
	block.GasLimit = int64(ethBlock.GasLimit())
	block.MinerAddress = ethBlock.Coinbase().Hex()
	block.TimeStamp = int64(ethBlock.Time())
	block.TransactionCount = len(ethBlock.Transactions())
	transactionsJSON, err := json.Marshal(ethBlock.Transactions())
	if err == nil {
		block.BlockData = datatypes.JSON(transactionsJSON) // 赋值为 JSON 数据
	}

	return block
}

func GetBlock(c *gin.Context) {
	blockNumber := c.Query("blockNumber")

	// 要支持head, finialized, safe
	num, err := model.ParseBLockNumber(blockNumber)
	if err != nil {
		c.JSON(400, gin.H{
			"error": err.Error(),
		})
		return
	}

	// head safe ...
	if num < 0 {
		ethBlock := utils.QueryBlockByNumber(num)
		if ethBlock == nil {
			c.JSON(400, gin.H{
				"error": "block not found",
			})
			return
		}

		block := AssembleBlock(ethBlock)
		c.JSON(200, gin.H{
			"block": block,
		})
		return
	}

	// 先从Redis中查

	// 再从Mysql中查, 更新缓存
	block, _ := model.GetBlockByNumber(num)

	if block != nil {
		c.JSON(200, gin.H{
			"block": block,
		})
		return
	}

	// 查不到再从链上查,并插入数据库, 更新缓存
	ethBlock := utils.QueryBlockByNumber(num)
	if ethBlock == nil {
		c.JSON(400, gin.H{
			"error": "block not found",
		})
		return
	}

	block = AssembleBlock(ethBlock)

	err = model.InsertBlock(block)
	// 插入失败暂时不考虑，只是打印输出
	if err != nil {
		fmt.Println("insert block failed", err)
	}

	c.JSON(200, gin.H{
		"block": block,
	})
}
