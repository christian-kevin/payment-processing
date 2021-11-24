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
}

func applyWalletRoutes(r *routergroup.Router, m *handler.Module, rr *middleware.MustRateLimit) {
	r.POST("/v1/wallet", api.CreateWallet(m.CreateWallet))
	r.GET("/v1/wallet", api.GetWallet(m.GetWallet))
}

func applyCardRoutes(r *routergroup.Router, m *handler.Module, rr *middleware.MustRateLimit) {
	r.Use(rr.SetPage("create_card").Enforce).POST("/v1/card", api.CreateCard(m.CreateCard))
	r.Use(rr.SetPage("get_card_multiple").Enforce).GET("/v1/card/multiple", api.GetCards(m.GetCards))
	r.DELETE("/v1/card", api.DeleteCard(m.DeleteCard))
}
