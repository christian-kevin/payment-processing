package component

import (
	"fmt"
	"spenmo/payment-processing/payment-processing/config"
	"spenmo/payment-processing/payment-processing/pkg/cache"
	"time"
)

func InitializeCache() (redis cache.Cache, err error) {
	rConfig := cache.Config{
		ServerAddr: fmt.Sprintf("%s:%s",
			config.AppConfig.RedisAddress,
			config.AppConfig.RedisPort),
		ReadTimeout:  time.Duration(config.AppConfig.RedisReadWriteTimeoutMillis) * time.Millisecond,
		WriteTimeout: time.Duration(config.AppConfig.RedisReadWriteTimeoutMillis) * time.Millisecond,
		MaxIdle:      config.AppConfig.RedisMaxIdle,
		PoolSize:     config.AppConfig.RedisConnPool,
		IdleTimeout:  time.Duration(config.AppConfig.RedisMaxIdleTimeoutMinute) * time.Minute,
	}
	cacher, err := cache.New(cache.Redis, &rConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to open redis connection. %+v", err)
	}
	return cacher, nil
}
