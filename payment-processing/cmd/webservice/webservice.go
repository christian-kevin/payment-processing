package webservice

import (
	"net/http"
	"spenmo/payment-processing/payment-processing/internal/pkg/component"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

func Start() {
	_, err := component.InitializeDatabase()
	if err != nil {
		panic(err)
	}
	_, err = component.InitializeCache()
	if err != nil {
		panic(err)
	}

	log.Get().Error(log.GetEmptyContext(), http.ListenAndServe(":8080", routerStart(&dependencies{
		//db:                  database,
		//cacher:              redis.NewStore(cache),
	})).Error())
}