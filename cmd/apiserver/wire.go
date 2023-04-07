//go:build wireinject
// +build wireinject

package main

import (
	go_web "go-web"
	"go-web/interface/http"
	"go-web/interface/router"
	"go-web/pkg/log"
	"go-web/pkg/mysql"
	"go-web/pkg/redis"

	"github.com/google/wire"
)

var providerSet = wire.NewSet(
	log.ProviderSet,
	go_web.ProviderSet,
	http.ProviderSet,
	mysql.ProviderSet,
	redis.ProviderSet,
	router.ProviderSet,
)

func Create() (*go_web.Server, error) {
	panic(wire.Build(providerSet))
}
