package webservice

import (
	"net/http"
	"spenmo/payment-processing/payment-processing/config"
	"spenmo/payment-processing/payment-processing/internal/app/public/handler"
	controller "spenmo/payment-processing/payment-processing/internal/app/public/http"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/middleware"
	"spenmo/payment-processing/payment-processing/internal/pkg/store/mysql"
	"spenmo/payment-processing/payment-processing/internal/pkg/store/redis"
	"spenmo/payment-processing/payment-processing/pkg/routergroup"
	"time"
)

type dependencies struct {
	db                      *component.Database
	cacher                  redis.Store
	rWalletStore            redis.RWalletStore
	rCardStore              redis.RCardStore
	walletStore             mysql.WalletStore
	cardStore               mysql.CardStore
	limitStore              mysql.LimitStore
	cardTransactionLogStore mysql.CardTransactionLogStore
	walletBalanceLogStore   mysql.WalletBalanceLogStore
	rateLimiter             *middleware.MustRateLimit
}

// Init to initialize the web-service router
func routerStart(dep *dependencies) *routergroup.Router {
	router := routergroup.New()

	handlePublicRoutes(router, dep)
	router.GET("/ping", Ping(dep))
	return router
}

func handlePublicRoutes(router *routergroup.Router, dep *dependencies) {
	tenantMiddleware := middleware.NewTenant()
	authMiddleware := middleware.NewAuth()

	publicGroup := router.Group("/public")
	publicGroup.Use(middleware.InjectCors, tenantMiddleware.Enforce, authMiddleware.Enforce)

	module := handler.NewModule(
		dep.rCardStore,
		dep.rWalletStore,
		dep.walletStore,
		dep.cardStore,
		dep.limitStore,
		dep.cardTransactionLogStore,
		dep.walletBalanceLogStore)
	controller.ApplyRoutes(publicGroup, module, dep.rateLimiter)
}

func Ping(dep *dependencies) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := dep.db.Master.Ping()
		if err != nil {
			panic(err)
		}

		err = dep.cacher.Ping()
		if err != nil {
			panic(err)
		}

		response.WriteResponse(w, &response.PingResponse{
			Status:          "ok",
			ServerTimestamp: time.Now().Unix(),
			AppName:         config.AppConfig.AppName,
			Environment:     string(config.AppConfig.Env),
		}, nil)
	}
}
