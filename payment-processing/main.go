package main

import (
	"os"
	"os/signal"
	"spenmo/payment-processing/payment-processing/cmd/webservice"
	"spenmo/payment-processing/payment-processing/config"
	log "spenmo/payment-processing/payment-processing/pkg/logger"
	"syscall"
)

func main() {
	config.InitializeAppConfig()
	log.InitializeLogger(config.AppConfig.AppName)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP, syscall.SIGINT)

	log.Get().Info(log.GetEmptyContext(), "default (web service) only allowed to run standalone. default service initialization skipped. Running webservice mode")
	go webservice.Start()

	//will wait until terminate signal or interrupt happened
	for {
		<-c
		os.Exit(0)
	}
}