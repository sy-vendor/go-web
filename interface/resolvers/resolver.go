package resolvers

import (
	"context"
	"go-web/ent"
	generated "go-web/graph/generated"
	goWebErrors "go-web/pkg/errors"
	"go-web/pkg/i18n"
	"go-web/pkg/redis"
	"time"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/gin-gonic/gin"
	redisv9 "github.com/redis/go-redis/v9"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client *ent.Client
	rdb    redis.Service
	logger *zap.Logger
}

const (
	maxQueryComplexity = 300
	defaultCacheTTL    = 5 * time.Minute
)

// NewConfig
func NewConfig(client *ent.Client, rdb redis.Service, logger *zap.Logger) *generated.Config {
	return &generated.Config{
		Resolvers: &Resolver{
			client: client,
			rdb:    rdb,
			logger: logger,
		},
	}
}

// ErrorPresenter 统一 GraphQL 错误格式与国际化
func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	if e, ok := err.(*goWebErrors.Error); ok {
		lang := "en"
		if ginCtx := GinContextFromContext(ctx); ginCtx != nil {
			lang = i18n.GetLang(ginCtx)
		}
		return &gqlerror.Error{
			Message:    i18n.TByLang(lang, e.Message),
			Extensions: map[string]interface{}{"code": e.Code, "details": e.Details},
			Path:       graphql.GetPath(ctx),
		}
	}
	return graphql.DefaultErrorPresenter(ctx, err)
}

// GinContextFromContext 从 GraphQL context 获取 gin.Context
func GinContextFromContext(ctx context.Context) *gin.Context {
	if v := ctx.Value("GinContextKey"); v != nil {
		if ginCtx, ok := v.(*gin.Context); ok {
			return ginCtx
		}
	}
	return nil
}

// NewGraphqlHandler
func NewGraphqlHandler(c *generated.Config, client *ent.Client, rdb *redisv9.Client, logger *zap.Logger) *handler.Server {
	if c == nil {
		panic("graphql config is required")
	}

	if client == nil {
		panic("ent client is required")
	}

	h := handler.NewDefaultServer(generated.NewExecutableSchema(*c))

	// 添加事务支持
	h.Use(entgql.Transactioner{TxOpener: client})

	// 设置查询复杂度限制
	h.Use(extension.FixedComplexityLimit(maxQueryComplexity))

	// 添加缓存支持
	h.Use(extension.Introspection{})
	h.Use(extension.AutomaticPersistedQuery{
		Cache: redis.NewAPQCache(rdb, defaultCacheTTL),
	})

	// 配置传输层
	h.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})
	h.AddTransport(transport.Options{})
	h.AddTransport(transport.GET{})
	h.AddTransport(transport.POST{})

	// 注册自定义 ErrorPresenter
	h.SetErrorPresenter(ErrorPresenter)

	// 添加恢复中间件
	// extension.RecoverFunc 可能不存在，建议移除或用默认 panic 处理
	// h.Use(extension.RecoverFunc(func(ctx context.Context, err interface{}) error {
	// 	logger.Error("panic recovered",
	// 		zap.Any("error", err),
	// 		zap.Any("path", graphql.GetPath(ctx)),
	// 	)
	// 	return goWebErrors.New(goWebErrors.ErrSystem, "system_error", "Internal server error")
	// }))

	return h
}
