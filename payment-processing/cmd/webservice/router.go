package webservice

import (
	"net/http"
	"spenmo/payment-processing/payment-processing/config"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	"spenmo/payment-processing/payment-processing/internal/pkg/dto/response"
	"spenmo/payment-processing/payment-processing/internal/pkg/middleware"
	"spenmo/payment-processing/payment-processing/internal/pkg/store/redis"
	"spenmo/payment-processing/payment-processing/pkg/routergroup"
	"time"
)

type dependencies struct {
	db                  *component.Database
	cacher              redis.Store
}

// Init to initialize the web-service router
func routerStart(dep *dependencies) *routergroup.Router {
	router := routergroup.New()

	handlePublicRoutes(router, dep)
	handlePrivateRoutes(router, dep)
	router.GET("/ping", Ping(dep))
	return router
}

func handlePublicRoutes(router *routergroup.Router, dep *dependencies) {
	tenantMiddleware := middleware.NewTenant()

	publicGroup := router.Group("/public")
	publicGroup.Use(middleware.InjectCors, tenantMiddleware.Enforce)

	//module := handler.NewModule(
	//)
	//controller.ApplyRoutes(publicGroup, module)
}

func handlePrivateRoutes(router *routergroup.Router, dep *dependencies) {
	//privateGroup := router.Group("/private")

	//module := privateHandler.NewModule()
	//privateController.ApplyRoutes(privateGroup, module)
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
