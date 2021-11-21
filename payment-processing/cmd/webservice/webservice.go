package webservice

import (
	"net/http"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
)

func Start() {
	//database, err := component.InitializeDatabase()
	//if err != nil {
	//	panic(err)
	//}
	//cache, err := component.InitializeCache()
	//if err != nil {
	//	panic(err)
	//}

	log.Get().Error(log.GetEmptyContext(), http.ListenAndServe(":8080", routerStart(&dependencies{
		//db:                  database,
		//cacher:              redis.NewStore(cache),
	})).Error())
}