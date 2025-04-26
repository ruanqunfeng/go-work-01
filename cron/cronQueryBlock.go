package cron

import (
	"fmt"
	"go-work-01/controller"
	"go-work-01/model"
	"go-work-01/utils"
	"log"
	"time"
)

func SyncBlocks() {

	// 使用的是第三方链,可能会有限流

	for {
		fmt.Println("开始同步区块")
		// 更新三个特殊区块head，safe...

		// 从链上查询最新区块
		lb := utils.QueryLatestBlock()
		if lb != nil {
			// 查询数据库的区块
			ethBlock, err := model.GetLastBlock()
			if err != nil {
				fmt.Println("查询数据库中的最新区块失败", err)
			} else {
				// 链上最新区块大于数据库中的最新区块
				// 获取链上最新区块
				for i := ethBlock.BlockNumber + 1; i <= lb.Number().Int64(); i++ {
					block := utils.QueryBlockByNumber(i)
					// 可能会出现区块查询失败的情况
					// 理论上需要有重试机制，这里暂时先不处理
					if block != nil {
						ethBlock = controller.AssembleBlock(block)

						err = model.InsertBlock(ethBlock)
						if err != nil {
							fmt.Println("插入数据库失败", err)
							continue
						}

						txs := make([]*model.EthTx, 0, len(block.Transactions()))
						// 获取交易插入数据库
						for _, t := range block.Transactions() {
							rawTx := utils.QueryTransaction(t.Hash().Hex())
							if rawTx != nil {
								// 插入数据库
								tx := &model.EthTx{}
								tx.BlockNumber = rawTx.BlockNumber
								tx.FromAddress = rawTx.FromAddress
								tx.Gas = rawTx.Gas
								tx.GasLimit = rawTx.GasLimit
								tx.GasPrice = rawTx.GasPrice
								tx.InputData = rawTx.InputData
								tx.ToAddress = rawTx.ToAddress
								tx.TxHash = rawTx.TxHash
								tx.Value = rawTx.Value
								txs = append(txs, tx)
							}
						}
						// todo: 插入失败要重试，先不做
						model.BatchCreateTx(txs)
					}
				}
			}
		}

		// 同步区块中的交易

		// 更新3个特殊区块

		// 完成之后sleep 5秒
		time.Sleep(5 * time.Second)
	}

	log.Fatal("当前同步任务已退出")

}
