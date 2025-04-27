package main

import (
	"net/http"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/server/config"
	"github.com/Maxim-Ba/metriccollector/internal/server/router"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
)

func main() {
	defer storage.Close()
	defer logger.Sync()
	
	parameters := config.GetParameters()


	_, err := storage.New(parameters)
	if err != nil {
		panic(err)
	}
	logger.SetLogLevel(parameters.LogLevel)

	
	mux := router.New()
	logger.LogInfo("Running server on ", parameters.Address)
	err = http.ListenAndServe(parameters.Address, mux)
	if err != nil {
		panic(err)
	}
}
