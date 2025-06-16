package main

import (
	go_web "go-web"
	"go-web/pkg/config"
	"go-web/pkg/log"
	"os"
	"runtime/debug"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
}

func main() {
	logger := log.NewLogger()
	defer func() {
		if r := recover(); r != nil {
			logger.Sugar().Errorf("panic: %v\n%s", r, debug.Stack())
			os.Exit(1)
		}
	}()

	// 加载配置
	cfg, err := config.Load(logger)
	if err != nil {
		logger.Sugar().Errorf("config load error: %v", err)
		os.Exit(1)
	}

	// 执行数据库迁移
	if err := go_web.Migrate(cfg); err != nil {
		logger.Sugar().Errorf("db migrate error: %v", err)
		os.Exit(1)
	}

	// 创建并启动服务器
	server, err := Create(cfg)
	if err != nil {
		logger.Sugar().Errorf("server create error: %v", err)
		os.Exit(1)
	}

	server.Start()
	server.AwaitSignal()
}
