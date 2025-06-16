package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-web/pkg/config"

	"github.com/google/wire"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// RedisConfig Redis 配置
type RedisConfig struct {
	// Redis 地址
	Addr string
	// Redis 密码
	Password string
	// Redis 数据库
	DB int
	// 连接池大小
	PoolSize int
	// 最小空闲连接数
	MinIdleConns int
	// 连接超时时间
	DialTimeout time.Duration
	// 读取超时时间
	ReadTimeout time.Duration
	// 写入超时时间
	WriteTimeout time.Duration
}

// DefaultRedisConfig 返回默认的 Redis 配置
func DefaultRedisConfig() *RedisConfig {
	return &RedisConfig{
		Addr:         "localhost:6379",
		Password:     "",
		DB:           0,
		PoolSize:     10,
		MinIdleConns: 5,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
	}
}

// RedisClient Redis 客户端
type RedisClient struct {
	client *redis.Client
	logger *zap.Logger
}

// NewRedisClient 创建 Redis 客户端
func NewRedisClient(cfg *config.Config, logger *zap.Logger) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:         cfg.Redis.Addr,
		Password:     cfg.Redis.Password,
		DB:           cfg.Redis.DB,
		PoolSize:     cfg.Redis.PoolSize,
		MinIdleConns: cfg.Redis.MinIdleConns,
		DialTimeout:  cfg.Redis.DialTimeout,
		ReadTimeout:  cfg.Redis.ReadTimeout,
		WriteTimeout: cfg.Redis.WriteTimeout,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Redis.DialTimeout)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	logger.Info("connected to Redis",
		zap.String("addr", cfg.Redis.Addr),
		zap.Int("db", cfg.Redis.DB),
		zap.Int("pool_size", cfg.Redis.PoolSize),
	)

	return &RedisClient{
		client: client,
		logger: logger,
	}, nil
}

// Get 获取缓存
func (c *RedisClient) Get(ctx context.Context, key string) ([]byte, error) {
	data, err := c.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get cache: %w", err)
	}
	return data, nil
}

// Set 设置缓存
func (c *RedisClient) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	var data []byte
	var err error

	// 根据值类型进行序列化
	switch v := value.(type) {
	case []byte:
		data = v
	case string:
		data = []byte(v)
	default:
		data, err = json.Marshal(value)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
	}

	if err := c.client.Set(ctx, key, data, ttl).Err(); err != nil {
		return fmt.Errorf("failed to set cache: %w", err)
	}

	return nil
}

// Delete 删除缓存
func (c *RedisClient) Delete(ctx context.Context, key string) error {
	if err := c.client.Del(ctx, key).Err(); err != nil {
		return fmt.Errorf("failed to delete cache: %w", err)
	}
	return nil
}

// Exists 检查缓存是否存在
func (c *RedisClient) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.client.Exists(ctx, key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check cache existence: %w", err)
	}
	return exists > 0, nil
}

// Close 关闭 Redis 连接
func (c *RedisClient) Close() error {
	if err := c.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis connection: %w", err)
	}
	return nil
}

var ProviderSet = wire.NewSet(NewRedisClient)
