package go_web

import (
	"context"
	"os/signal"
	"syscall"

	"go-web/interface/http"
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
}

func NewTopLevelCtx() context.Context {
	return context.Background()
}

func NewServer(ctx context.Context, http *http.Server, logger *zap.Logger) *Server {
	return &Server{
		ctx:    ctx,
		http:   http,
		logger: logger,
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

	if <-ctx.Done(); true {
		if err := ctx.Err(); err != nil {
			server.logger.Sugar().Warnf("server top-level context error: %v", err)
		}
		if server.http != nil {
			err := server.http.StopServer()
			if err != nil {
				server.logger.Sugar().Errorf("close http server failed: %v", err)
			}
		}
	}

	if err := mysql.CloseMysql(); err != nil {
		server.logger.Sugar().Errorf("close db failed: %v", err)
	}

	// sync log
	log.Close()
}
