package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// APQCache 实现 GraphQL 持久化查询缓存
type APQCache struct {
	Client *redis.Client
	TTL    time.Duration
}

// Get 从缓存获取查询
func (c *APQCache) Get(ctx context.Context, key string) (interface{}, bool) {
	val, err := c.Client.Get(ctx, key).Result()
	if err != nil {
		return nil, false
	}
	return val, true
}

// Add 添加查询到缓存
func (c *APQCache) Add(ctx context.Context, key string, value interface{}) {
	c.Client.Set(ctx, key, fmt.Sprintf("%v", value), c.TTL)
}

// NewAPQCache 创建新的 APQ 缓存实例
func NewAPQCache(client *redis.Client, ttl time.Duration) *APQCache {
	return &APQCache{
		Client: client,
		TTL:    ttl,
	}
}
