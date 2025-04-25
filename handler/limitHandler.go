package handler

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// IPRateLimiter 用于存储每个IP的请求时间
type IPRateLimiter struct {
	ips map[string]time.Time
	mu  sync.Mutex
}

// NewIPRateLimiter 创建一个新的IP限流器
func NewIPRateLimiter() *IPRateLimiter {
	return &IPRateLimiter{
		ips: make(map[string]time.Time),
	}
}

// Allow 检查是否允许该IP的请求
func (rl *IPRateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	lastRequestTime, exists := rl.ips[ip]
	if !exists || time.Since(lastRequestTime) >= time.Second {
		rl.ips[ip] = time.Now()
		return true
	}
	return false
}

// RateLimitMiddleware 限流中间件
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := c.ClientIP()
		if !limiter.Allow(ip) {
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests",
			})
			return
		}
		c.Next()
	}
}
