package http

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"go-web/interface/http/middleware"
	"go-web/pkg/cache"

	"github.com/gin-contrib/gzip"
	gin_zap "github.com/gin-contrib/zap"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"go.uber.org/zap"
)

var ProviderSet = wire.NewSet(NewServer, NewRouter)

type Server struct {
	logger     *zap.Logger
	router     *gin.Engine
	httpServer *http.Server
}

const closeWaitTime = 30 * time.Second

type InitRoutersFunc func(r *gin.Engine)

func NewRouter(logger *zap.Logger, redisClient *cache.RedisClient, initRoutersFunc InitRoutersFunc) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)

	r := gin.New()

	r.Use(middleware.RequestID(logger))

	r.Use(middleware.RateLimit(logger, middleware.DefaultRateLimitConfig()))

	r.Use(middleware.Logger(logger))

	r.Use(middleware.Security(middleware.DefaultSecurityConfig()))

	r.Use(middleware.CSRF(logger, middleware.DefaultCSRFConfig()))

	r.Use(middleware.Validator(logger, middleware.DefaultValidatorConfig()))

	r.Use(middleware.Cache(logger, middleware.DefaultCacheConfig(redisClient)))

	r.Use(gzip.Gzip(gzip.DefaultCompression, gzip.WithExcludedExtensions([]string{})))
	r.Use(gin_zap.Ginzap(logger, time.RFC3339, true))
	r.Use(middleware.Recovery(logger))
	r.Use(middleware.Cors(middleware.DefaultCORSConfig()))

	// 添加 Prometheus 性能监控中间件
	r.Use(middleware.Metrics())

	// 添加 /metrics 路由
	r.GET("/metrics", middleware.MetricsHandler())

	initRoutersFunc(r)

	return r
}

func NewServer(logger *zap.Logger, router *gin.Engine) *Server {
	srv := &http.Server{
		Addr:              os.Getenv("HTTP_SERVER_PORT"),
		Handler:           router,
		ReadHeaderTimeout: 20 * time.Second,
	}

	return &Server{
		logger:     logger.With(zap.String("component", "http_server")),
		router:     router,
		httpServer: srv,
	}
}

func (server *Server) StartServer() {
	go func() {
		server.logger.Info("starting http server", zap.String("addr", server.httpServer.Addr))
		if err := server.httpServer.ListenAndServe(); err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				server.logger.Info("http server closed")
				return
			}
			server.logger.Error("http server error",
				zap.Error(err),
				zap.String("addr", server.httpServer.Addr),
			)
		}
	}()
}

func (server *Server) StopServer() error {
	ctx, cancel := context.WithTimeout(context.Background(), closeWaitTime)
	defer cancel()
	if err := server.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("http server shutdown error: %w", err)
	}

	return nil
}
