package redis

import (
	"context"
	"fmt"
	"naive-admin/pkg/config"
	"naive-admin/pkg/log"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

var once sync.Once

func NewRedis(cfg *config.Redis) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			cfg.Host,
			cfg.Port,
		),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Panic(err)
	}

	return rdb
}

func Init() {
	singleton := func() {
		if config.Conf == nil {
			panic("config is nil!")
		}
		Rdb = NewRedis(config.Conf.Data.Redis)
	}
	once.Do(singleton)
}

func Close() {
	_ = Rdb.Close()
	log.Info("redis disconnect")
}
