package redis

import (
	"context"
	"os"
	"strconv"
	"sync"

	"github.com/google/wire"
	rd "github.com/redis/go-redis/v9"
)

var (
	r    *Redis
	once sync.Once
)

type Redis struct {
	rdb *rd.Client
}

type Service interface {
	GetRDB() *rd.Client
}

var ProviderSet = wire.NewSet(NewRedis)

// NewRedis Init redis service
func NewRedis(ctx context.Context) Service {
	once.Do(func() {
		db, err := strconv.Atoi(os.Getenv("REDIS_DB"))
		if err != nil {
			panic("redis db error")
		}
		rdb := rd.NewClient(&rd.Options{
			Addr:     os.Getenv("REDIS_URL"),
			Password: os.Getenv("REDIS_PASSWORD"),
			DB:       db,
		})

		// ping
		if err := rdb.Do(ctx, "ping").Err(); err != nil {
			panic(err)
		}

		r = &Redis{rdb: rdb}
	})

	return r
}

func (r *Redis) GetRDB() *rd.Client {
	return r.rdb
}
