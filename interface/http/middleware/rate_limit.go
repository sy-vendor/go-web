package middleware

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/time/rate"
)

// RateLimiter 是限流器的接口
type RateLimiter interface {
	Allow() bool
}

// TokenBucketLimiter 是基于令牌桶算法的限流器
type TokenBucketLimiter struct {
	limiter *rate.Limiter
}

// NewTokenBucketLimiter 创建一个新的令牌桶限流器
func NewTokenBucketLimiter(rps float64, burst int) *TokenBucketLimiter {
	return &TokenBucketLimiter{
		limiter: rate.NewLimiter(rate.Limit(rps), burst),
	}
}

// Allow 检查是否允许请求通过
func (l *TokenBucketLimiter) Allow() bool {
	return l.limiter.Allow()
}

// RateLimitConfig 是限流中间件的配置选项
type RateLimitConfig struct {
	// RPS 每秒允许的请求数
	RPS float64
	// Burst 允许的突发请求数
	Burst int
	// IPBased 是否基于 IP 进行限流
	IPBased bool
	// SkipPaths 不需要限流的路径列表
	SkipPaths []string
}

// DefaultRateLimitConfig 返回默认的限流配置
func DefaultRateLimitConfig() *RateLimitConfig {
	return &RateLimitConfig{
		RPS:     100,  // 每秒 100 个请求
		Burst:   200,  // 允许突发 200 个请求
		IPBased: true, // 基于 IP 限流
		SkipPaths: []string{
			"/health",
			"/metrics",
		},
	}
}

// RateLimit 返回一个限流中间件
func RateLimit(logger *zap.Logger, config *RateLimitConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultRateLimitConfig()
	}

	// 创建限流器映射
	limiters := make(map[string]RateLimiter)
	var mu sync.Mutex

	// 检查路径是否需要跳过限流
	shouldSkip := func(path string) bool {
		for _, skipPath := range config.SkipPaths {
			if path == skipPath {
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		// 检查是否需要跳过限流
		if shouldSkip(c.Request.URL.Path) {
			c.Next()
			return
		}

		// 获取限流器的键
		key := "global"
		if config.IPBased {
			key = c.ClientIP()
		}

		// 获取或创建限流器
		mu.Lock()
		limiter, exists := limiters[key]
		if !exists {
			limiter = NewTokenBucketLimiter(config.RPS, config.Burst)
			limiters[key] = limiter
		}
		mu.Unlock()

		// 检查是否允许请求通过
		if !limiter.Allow() {
			// 获取请求 ID
			requestID := c.GetString(RequestIDKey)

			// 记录限流日志
			logger.Warn("rate limit exceeded",
				zap.String("ip", c.ClientIP()),
				zap.String("path", c.Request.URL.Path),
				zap.String("request_id", requestID),
			)

			// 返回 429 状态码
			c.AbortWithStatusJSON(http.StatusTooManyRequests, gin.H{
				"code":       http.StatusTooManyRequests,
				"message":    "Too Many Requests",
				"request_id": requestID,
			})
			return
		}

		c.Next()
	}
}
