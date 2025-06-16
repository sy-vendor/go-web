package config

import (
	"fmt"
	"time"

	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

// Config holds all configuration for the application
type Config struct {
	// 服务器配置
	Server struct {
		// 监听地址
		Addr string `mapstructure:"addr"`
		// 读取超时时间
		ReadTimeout time.Duration `mapstructure:"read_timeout"`
		// 写入超时时间
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
		// 空闲超时时间
		IdleTimeout time.Duration `mapstructure:"idle_timeout"`
	} `mapstructure:"server"`

	// 日志配置
	Log struct {
		// 日志级别
		Level string `mapstructure:"level"`
		// 日志文件路径
		Filename string `mapstructure:"filename"`
		// 是否输出到控制台
		Console bool `mapstructure:"console"`
		// 是否输出到文件
		File bool `mapstructure:"file"`
		// 是否输出调用者信息
		Caller bool `mapstructure:"caller"`
		// 是否输出堆栈信息
		Stacktrace bool `mapstructure:"stacktrace"`
	} `mapstructure:"log"`

	// MySQL 配置
	MySQL struct {
		// 数据库地址
		Addr string `mapstructure:"addr"`
		// 数据库用户名
		Username string `mapstructure:"username"`
		// 数据库密码
		Password string `mapstructure:"password"`
		// 数据库名称
		Database string `mapstructure:"database"`
		// 最大空闲连接数
		MaxIdleConns int `mapstructure:"max_idle_conns"`
		// 最大打开连接数
		MaxOpenConns int `mapstructure:"max_open_conns"`
		// 连接最大生命周期
		ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime"`
	} `mapstructure:"mysql"`

	// Redis 配置
	Redis struct {
		// Redis 地址
		Addr string `mapstructure:"addr"`
		// Redis 密码
		Password string `mapstructure:"password"`
		// Redis 数据库
		DB int `mapstructure:"db"`
		// 连接池大小
		PoolSize int `mapstructure:"pool_size"`
		// 最小空闲连接数
		MinIdleConns int `mapstructure:"min_idle_conns"`
		// 连接超时时间
		DialTimeout time.Duration `mapstructure:"dial_timeout"`
		// 读取超时时间
		ReadTimeout time.Duration `mapstructure:"read_timeout"`
		// 写入超时时间
		WriteTimeout time.Duration `mapstructure:"write_timeout"`
	} `mapstructure:"redis"`
}

var (
	config *Config
)

var ProviderSet = wire.NewSet(Load)

// Load loads the configuration from the environment and config file
func Load(logger *zap.Logger) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")
	viper.AddConfigPath("/etc/go-web/")
	viper.AutomaticEnv()

	// Set default values
	setDefaults()

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("failed to read config file: %w", err)
		}
	}

	// Unmarshal config
	config = &Config{}
	if err := viper.Unmarshal(config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// Validate config
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	logger.Info("loaded config",
		zap.String("file", viper.ConfigFileUsed()),
	)

	return config, nil
}

// Get returns the current configuration
func Get() *Config {
	return config
}

// setDefaults sets default values for configuration
func setDefaults() {
	// Server defaults
	viper.SetDefault("server.addr", ":8080")
	viper.SetDefault("server.read_timeout", 5*time.Second)
	viper.SetDefault("server.write_timeout", 10*time.Second)
	viper.SetDefault("server.idle_timeout", 120*time.Second)

	// Log defaults
	viper.SetDefault("log.level", "info")
	viper.SetDefault("log.console", true)
	viper.SetDefault("log.file", false)
	viper.SetDefault("log.caller", true)
	viper.SetDefault("log.stacktrace", false)

	// MySQL defaults
	viper.SetDefault("mysql.addr", "localhost:3306")
	viper.SetDefault("mysql.max_idle_conns", 10)
	viper.SetDefault("mysql.max_open_conns", 100)
	viper.SetDefault("mysql.conn_max_lifetime", time.Hour)

	// Redis defaults
	viper.SetDefault("redis.addr", "localhost:6379")
	viper.SetDefault("redis.db", 0)
	viper.SetDefault("redis.pool_size", 10)
	viper.SetDefault("redis.min_idle_conns", 5)
	viper.SetDefault("redis.dial_timeout", 5*time.Second)
	viper.SetDefault("redis.read_timeout", 3*time.Second)
	viper.SetDefault("redis.write_timeout", 3*time.Second)
}

// validateConfig validates the configuration
func validateConfig(cfg *Config) error {
	if cfg.Server.Addr == "" {
		return fmt.Errorf("server.addr is required")
	}

	if cfg.Log.Level == "" {
		return fmt.Errorf("log.level is required")
	}

	if cfg.MySQL.Addr == "" {
		return fmt.Errorf("mysql.addr is required")
	}

	if cfg.MySQL.Username == "" {
		return fmt.Errorf("mysql.username is required")
	}

	if cfg.MySQL.Database == "" {
		return fmt.Errorf("mysql.database is required")
	}

	if cfg.Redis.Addr == "" {
		return fmt.Errorf("redis.addr is required")
	}

	return nil
}

// GetDSN returns the database connection string
func (c *Config) GetDSN() string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8mb4&parseTime=True&loc=Local&multiStatements=true",
		c.MySQL.Username,
		c.MySQL.Password,
		c.MySQL.Addr,
		c.MySQL.Database,
	)
}
