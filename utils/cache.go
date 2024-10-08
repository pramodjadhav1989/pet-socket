package utils

import (
	"context"
	"time"

	log "github.com/smartpet/websocket/utils/logger"

	"github.com/jellydator/ttlcache/v3"
)

var cache *ttlcache.Cache[string, interface{}]

func InitCache() {
	cache = ttlcache.New[string, interface{}](
		ttlcache.WithTTL[string, interface{}](30 * time.Minute),
	)
	go cache.Start()
}

func GetCache() *ttlcache.Cache[string, interface{}] {
	return cache
}

func CloseCache() {
	cache.Stop()
}

func LogCacheStats(ctx context.Context) {
	cacheStats := cache.Metrics()
	log.Info(ctx).Interface("buyback-Cache", cacheStats).Send()
}

func StartCacheLogging(ctx context.Context, frequency time.Duration) {
	log.ApplicationInfo(ctx).Msg("Starting Async Logging in separate Thread for cache")
	go startAsyncLogging(ctx, frequency)
}

func startAsyncLogging(ctx context.Context, frequency time.Duration) {

	for {
		select {
		case <-ctx.Done():
			return

		case <-time.After(frequency):
			LogCacheStats(ctx)
		}
	}
}
