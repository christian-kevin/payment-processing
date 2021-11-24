package http

import (
	"spenmo/payment-processing/payment-processing/internal/app/public/handler"
	"spenmo/payment-processing/payment-processing/internal/app/public/http/api"
	"spenmo/payment-processing/payment-processing/internal/pkg/middleware"
	"spenmo/payment-processing/payment-processing/pkg/routergroup"
)

func ApplyRoutes(r *routergroup.Router, m *handler.Module, rr *middleware.MustRateLimit) {
	applyWalletRoutes(r, m, rr)
	applyCardRoutes(r, m, rr)
	applyPublicCardRoutes(r, m, rr)
}

func applyWalletRoutes(r *routergroup.Router, m *handler.Module, rr *middleware.MustRateLimit) {
	r = r.Group("/v1/wallet")
	tenantMiddleware := middleware.NewTenant()
	authMiddleware := middleware.NewAuth()
	r.Use(tenantMiddleware.Enforce, authMiddleware.Enforce)
	r.POST("", api.CreateWallet(m.CreateWallet))
	r.GET("", api.GetWallet(m.GetWallet))
}

func applyCardRoutes(r *routergroup.Router, m *handler.Module, rr *middleware.MustRateLimit) {
	r = r.Group("/v1/card")
	tenantMiddleware := middleware.NewTenant()
	authMiddleware := middleware.NewAuth()
	r.Use(tenantMiddleware.Enforce, authMiddleware.Enforce)
	r.POST("", rr.Enforce(api.CreateCard(m.CreateCard), "create-card") )
	r.GET("/multiple", rr.Enforce(api.GetCards(m.GetCards), "get-card-multiple"))
	r.DELETE("", api.DeleteCard(m.DeleteCard))
}

func applyPublicCardRoutes(r *routergroup.Router, m *handler.Module, rr *middleware.MustRateLimit) {
	r.POST("/v1/card/public/transaction", api.CreateTransaction(m.CreateTransaction))
}
