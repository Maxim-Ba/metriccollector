package main

import (
	"fmt"
	"net/http"

	_ "github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/server/config"
	"github.com/Maxim-Ba/metriccollector/internal/server/router"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
)


func main() {

	parameters := config.GetParameters()

	_,err:=storage.New(parameters)
	if err != nil {
		panic(err)
	}
	mux := router.New()
	fmt.Println("Running server on ", parameters.Address)
	
	err = http.ListenAndServe(parameters.Address, mux)
	if err != nil {
		panic(err)
	}
}
