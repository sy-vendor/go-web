package go_web

import (
	"context"
	"os/signal"
	"syscall"

	"go-web/interface/http"
	"go-web/pkg/config"
	"go-web/pkg/log"
	"go-web/pkg/mysql"

	"github.com/google/wire"
	"go.uber.org/zap"
	"golang.org/x/sys/unix"
)

var ProviderSet = wire.NewSet(NewServer, NewTopLevelCtx)

type Server struct {
	ctx    context.Context
	http   *http.Server
	logger *zap.Logger
	config *config.Config
}

func NewTopLevelCtx() context.Context {
	return context.Background()
}

func NewServer(ctx context.Context, http *http.Server, logger *zap.Logger, cfg *config.Config) *Server {
	return &Server{
		ctx:    ctx,
		http:   http,
		logger: logger,
		config: cfg,
	}
}

func (server *Server) Start() {
	if server.http != nil {
		server.http.StartServer()
	}
}

func (server *Server) AwaitSignal() {
	signal.Reset(syscall.SIGTERM, syscall.SIGINT)
	ctx, stop := signal.NotifyContext(server.ctx, unix.SIGTERM, unix.SIGINT)
	defer stop()

	server.logger.Info("waiting for shutdown signal...")
	if <-ctx.Done(); true {
		server.logger.Info("received shutdown signal")

		if err := ctx.Err(); err != nil {
			server.logger.Warn("server context error", zap.Error(err))
		}

		server.logger.Info("shutting down http server...")
		if server.http != nil {
			if err := server.http.StopServer(); err != nil {
				server.logger.Error("failed to shutdown http server",
					zap.Error(err),
					zap.String("component", "http_server"),
				)
			} else {
				server.logger.Info("http server shutdown complete")
			}
		}

		server.logger.Info("closing database connections...")
		if err := mysql.CloseMysql(); err != nil {
			server.logger.Error("failed to close database connections",
				zap.Error(err),
				zap.String("component", "database"),
			)
		} else {
			server.logger.Info("database connections closed")
		}

		// sync log
		log.Close()
		server.logger.Info("server shutdown complete")
	}
}
