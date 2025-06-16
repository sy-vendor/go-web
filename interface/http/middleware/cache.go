package middleware

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go-web/pkg/cache"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CacheConfig 缓存配置
type CacheConfig struct {
	// 默认缓存时间
	DefaultTTL time.Duration
	// 是否启用缓存
	Enabled bool
	// 缓存排除的路径
	ExcludePaths []string
	// 缓存排除的方法
	ExcludeMethods []string
	// 缓存排除的状态码
	ExcludeStatusCodes []int
	// 缓存键前缀
	KeyPrefix string
	// Redis 客户端
	RedisClient *cache.RedisClient
}

// DefaultCacheConfig 返回默认的缓存配置
func DefaultCacheConfig(redisClient *cache.RedisClient) *CacheConfig {
	return &CacheConfig{
		DefaultTTL:         5 * time.Minute,
		Enabled:            true,
		ExcludePaths:       []string{"/api/v1/graphql"},
		ExcludeMethods:     []string{"POST", "PUT", "DELETE", "PATCH"},
		ExcludeStatusCodes: []int{http.StatusBadRequest, http.StatusUnauthorized, http.StatusForbidden, http.StatusInternalServerError},
		KeyPrefix:          "cache:",
		RedisClient:        redisClient,
	}
}

// Cache 缓存中间件
func Cache(logger *zap.Logger, config *CacheConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检查是否启用缓存
		if !config.Enabled {
			c.Next()
			return
		}

		// 检查是否排除当前路径
		if shouldExcludePath(c.Request.URL.Path, config.ExcludePaths) {
			c.Next()
			return
		}

		// 检查是否排除当前方法
		if shouldExcludeMethod(c.Request.Method, config.ExcludeMethods) {
			c.Next()
			return
		}

		// 生成缓存键
		cacheKey := generateCacheKey(c, config.KeyPrefix)

		// 尝试从缓存获取响应
		if cachedResponse, err := config.RedisClient.Get(c.Request.Context(), cacheKey); err == nil && cachedResponse != nil {
			logger.Debug("cache hit",
				zap.String("key", cacheKey),
				zap.String("path", c.Request.URL.Path),
			)
			c.DataFromReader(http.StatusOK, int64(len(cachedResponse)), "application/json", bytes.NewReader(cachedResponse), nil)
			c.Abort()
			return
		}

		// 创建自定义的响应写入器
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		// 处理请求
		c.Next()

		// 检查响应状态码
		if shouldExcludeStatusCode(writer.statusCode, config.ExcludeStatusCodes) {
			return
		}

		// 缓存响应
		if writer.statusCode == http.StatusOK {
			responseBody := writer.body.Bytes()
			ttl := getCacheTTL(c, config.DefaultTTL)
			if err := config.RedisClient.Set(c.Request.Context(), cacheKey, responseBody, ttl); err != nil {
				logger.Error("failed to save cache",
					zap.Error(err),
					zap.String("key", cacheKey),
					zap.String("path", c.Request.URL.Path),
				)
			} else {
				logger.Debug("cache miss, saved to cache",
					zap.String("key", cacheKey),
					zap.String("path", c.Request.URL.Path),
					zap.Duration("ttl", ttl),
				)
			}
		}
	}
}

// responseWriter 自定义响应写入器
type responseWriter struct {
	gin.ResponseWriter
	body       *bytes.Buffer
	statusCode int
}

// Write 实现 io.Writer 接口
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteHeader 实现 http.ResponseWriter 接口
func (w *responseWriter) WriteHeader(statusCode int) {
	w.statusCode = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}

// shouldExcludePath 检查是否排除当前路径
func shouldExcludePath(path string, excludePaths []string) bool {
	for _, excludePath := range excludePaths {
		if strings.HasPrefix(path, excludePath) {
			return true
		}
	}
	return false
}

// shouldExcludeMethod 检查是否排除当前方法
func shouldExcludeMethod(method string, excludeMethods []string) bool {
	for _, excludeMethod := range excludeMethods {
		if method == excludeMethod {
			return true
		}
	}
	return false
}

// shouldExcludeStatusCode 检查是否排除当前状态码
func shouldExcludeStatusCode(statusCode int, excludeStatusCodes []int) bool {
	for _, excludeStatusCode := range excludeStatusCodes {
		if statusCode == excludeStatusCode {
			return true
		}
	}
	return false
}

// generateCacheKey 生成缓存键
func generateCacheKey(c *gin.Context, prefix string) string {
	// 获取请求路径和查询参数
	path := c.Request.URL.Path
	query := c.Request.URL.RawQuery

	// 计算 MD5 哈希
	hash := md5.New()
	io.WriteString(hash, path)
	io.WriteString(hash, query)
	hashStr := hex.EncodeToString(hash.Sum(nil))

	return fmt.Sprintf("%s%s", prefix, hashStr)
}

// getCacheTTL 获取缓存时间
func getCacheTTL(c *gin.Context, defaultTTL time.Duration) time.Duration {
	// 从请求头获取缓存时间
	if ttl := c.GetHeader("X-Cache-TTL"); ttl != "" {
		if duration, err := time.ParseDuration(ttl); err == nil {
			return duration
		}
	}
	return defaultTTL
}

// getFromCache 从缓存获取数据
func getFromCache(key string) ([]byte, bool) {
	// TODO: 实现从 Redis 或其他缓存系统获取数据
	return nil, false
}

// saveToCache 保存数据到缓存
func saveToCache(key string, data []byte, ttl time.Duration) {
	// TODO: 实现保存数据到 Redis 或其他缓存系统
}
