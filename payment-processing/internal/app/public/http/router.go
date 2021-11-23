package http

import (
	"spenmo/payment-processing/payment-processing/internal/app/public/handler"
	"spenmo/payment-processing/payment-processing/internal/app/public/http/api"
	"spenmo/payment-processing/payment-processing/pkg/routergroup"
)

func ApplyRoutes(r *routergroup.Router, m *handler.Module) {
	r.POST("/v1/wallet", api.CreateWallet(m.CreateWallet))
	r.GET("/v1/wallet", api.GetWallet(m.GetWallet))
	r.POST("/v1/card", api.CreateCard(m.CreateCard))
	r.GET("/v1/card/multiple", api.GetCards(m.GetCards))
}
