package resolvers

import (
	"context"
	"go-web/ent"
	generated "go-web/graph/generated"
	"go-web/pkg/errors"
	"go-web/pkg/i18n"
	"go-web/pkg/redis"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/gqlerror"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/gin-gonic/gin"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	client *ent.Client
	rdb    redis.Service
}

const maxQueryComplexity = 300

// NewConfig
func NewConfig(client *ent.Client, rdb redis.Service) *generated.Config {
	return &generated.Config{
		Resolvers: &Resolver{
			client: client,
			rdb:    rdb,
		},
	}
}

// ErrorPresenter 统一 GraphQL 错误格式与国际化
func ErrorPresenter(ctx context.Context, err error) *gqlerror.Error {
	gqlErr := graphql.DefaultErrorPresenter(ctx, err)
	if e, ok := err.(*errors.Error); ok {
		lang := "en"
		if ginCtx := GinContextFromContext(ctx); ginCtx != nil {
			lang = i18n.GetLang(ginCtx)
		}
		gqlErr.Message = i18n.TByLang(lang, e.Message)
		if gqlErr.Extensions == nil {
			gqlErr.Extensions = map[string]interface{}{}
		}
		gqlErr.Extensions["code"] = e.Code
		gqlErr.Extensions["details"] = e.Details
	}
	return gqlErr
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
func NewGraphqlHandler(c *generated.Config, client *ent.Client) *handler.Server {
	if c == nil {
		panic("graphql config is required")
	}

	if client == nil {
		panic("ent client is required")
	}

	h := handler.NewDefaultServer(generated.NewExecutableSchema(*c))
	h.Use(entgql.Transactioner{TxOpener: client})
	h.Use(extension.FixedComplexityLimit(maxQueryComplexity))

	// 注册自定义 ErrorPresenter
	h.SetErrorPresenter(ErrorPresenter)

	return h
}
