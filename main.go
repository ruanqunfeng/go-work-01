package main

import (
	"go-work-01/controller"
	"go-work-01/cron"
	"go-work-01/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 创建限流器
	limiter := handler.NewIPRateLimiter()

	// 使用限流中间件
	r.Use(handler.RateLimitMiddleware(limiter))

	// 公共路由
	r.GET("/health", func(c *gin.Context) {
		c.String(200, "OK")
	})

	r.GET("/block", controller.GetBlock)

	r.GET("/tx", controller.GetTx)

	r.GET("/txReceipt", controller.GetTxReceipt)

	go cron.SyncBlocks()

	r.Run(":18181")
}
