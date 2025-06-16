//go:build wireinject
// +build wireinject

package main

import (
	go_web "go-web"
	"go-web/interface/http"
	"go-web/interface/router"
	"go-web/pkg/cache"
	"go-web/pkg/config"
	"go-web/pkg/log"
	"go-web/pkg/mysql"
	"go-web/pkg/redis"

	"github.com/google/wire"
)

func Create(cfg *config.Config) (*go_web.Server, error) {
	wire.Build(
		go_web.ProviderSet,
		http.ProviderSet,
		log.ProviderSet,
		mysql.ProviderSet,
		redis.ProviderSet,
		cache.ProviderSet,
		router.ProviderSet,
	)
	return nil, nil
}
