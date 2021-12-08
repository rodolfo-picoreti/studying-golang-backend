package services

import (
	"example/hello/config"

	"github.com/gin-contrib/cache/persistence"
)

var store *persistence.RedisStore

func GetCacheStore() *persistence.RedisStore {
	if store == nil {
		config := config.ReadConfig()

		store = persistence.NewRedisCache(config.Cache.Hostname,
			config.Cache.Password,
			config.Cache.DefaultExpiration)
	}
	return store
}
