package cron

import (
	"log"
	"time"
)

func SyncBlocks() {

	// 最小从这里开始查询
	// var small int64 = 500000

	// 使用的是第三方链,可能会有限流

	for {
		// 从数据库中查询最近的区块

		// 同步区块中的交易

		// 更新3个特殊区块

		// 完成之后sleep30秒
		time.Sleep(30 * time.Second)
	}

	log.Fatal("当前同步任务已退出")

}
