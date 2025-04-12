package main

import (
	"fmt"
	"net/http"

	"github.com/Maxim-Ba/metriccollector/internal/agent/config"
	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/server/router"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
)


func main() {

	parameterts := config.GetParameters()
	logger.InitLogger()
	_,err:=storage.New()
	if err != nil {
		panic(err)
	}
	logger.InitLogger()
	mux := router.New()
	fmt.Println("Running server on ", parameterts.Addres)
	
	err = http.ListenAndServe(parameterts.Addres, mux)
	if err != nil {
		panic(err)
	}
}
