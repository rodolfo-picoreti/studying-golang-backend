package services

import (
	"context"
	"example/hello/config"
	"example/hello/telemetry"
	"time"

	"github.com/gin-contrib/cache/persistence"
	"go.opentelemetry.io/otel/attribute"
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

func GetCache(ctx context.Context, key string, value interface{}) error {
	_, span := telemetry.Tracer.Start(ctx, "GetCache")
	defer span.End()

	err := GetCacheStore().Get(key, value)
	hit := err == nil
	span.SetAttributes(attribute.Bool("cache-hit", hit))

	return err
}

func SetCache(ctx context.Context, key string, value interface{}) error {
	_, span := telemetry.Tracer.Start(ctx, "SetCache")
	defer span.End()

	err := GetCacheStore().Set(key, value, time.Minute)
	if err != nil {
		telemetry.GetLogger().Warn().Err(err).Msg("setCache failed")
	}
	return err
}

func ExpireCache(ctx context.Context, key string) error {
	_, span := telemetry.Tracer.Start(ctx, "ExpireCache")
	defer span.End()

	err := GetCacheStore().Delete(key)
	if err != nil {
		telemetry.GetLogger().Warn().Err(err).Msg("expireCache failed")
	}
	return err
}
