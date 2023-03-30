package router

import (
	"go-web/interface/http"
	"go-web/interface/resolvers"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
)

var ProviderSet = wire.NewSet(CreateInitRoutesFunc, resolvers.NewConfig, resolvers.NewGraphqlHandler)

func CreateInitRoutesFunc(gql *handler.Server) http.InitRoutersFunc {
	return func(r *gin.Engine) {
		r.POST("/query", func() gin.HandlerFunc {
			return func(context *gin.Context) {
				gql.ServeHTTP(context.Writer, context.Request)
			}
		}())
	}
}
