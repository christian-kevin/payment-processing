package webservice

import (
	"net/http"
	"spenmo/payment-processing/payment-processing/config"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	"spenmo/payment-processing/payment-processing/internal/pkg/middleware"
	"spenmo/payment-processing/payment-processing/internal/pkg/ratelimiter"
	"spenmo/payment-processing/payment-processing/internal/pkg/store/mysql"
	"spenmo/payment-processing/payment-processing/internal/pkg/store/redis"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

func Start() {
	database, err := component.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	cache, err := component.InitializeCache()
	if err != nil {
		panic(err)
	}

	rCardStore := redis.NewCardStore(cache)
	rWalletStore := redis.NewWalletStore(cache)
	rRateLimitStore := redis.NewRateLimitStore(cache)
	rateLimiter := ratelimiter.NewRateLimiter(rRateLimitStore)
	r := middleware.NewMustRateLimit(rateLimiter, config.AppConfig.RateLimitEnabled)

	walletStore := mysql.NewWalletStore(database)
	cardStore := mysql.NewCardStore(database)
	limitStore := mysql.NewLimitStore(database)
	cardTransactionLogStore := mysql.NewCardTransactionLogStore(database)
	walletBalanceLogStore := mysql.NewWalletBalanceLogStore(database)

	log.Get().Error(log.GetEmptyContext(), http.ListenAndServe(":8080", routerStart(&dependencies{
		db:                      database,
		cacher:                  redis.NewStore(cache),
		rCardStore:              rCardStore,
		rWalletStore:            rWalletStore,
		walletStore:             walletStore,
		cardStore:               cardStore,
		limitStore:              limitStore,
		cardTransactionLogStore: cardTransactionLogStore,
		walletBalanceLogStore:   walletBalanceLogStore,
		rateLimiter:             r,
	})).Error())
}
