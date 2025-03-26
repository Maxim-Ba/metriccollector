package main

import (
	"fmt"
	"net/http"

	"github.com/Maxim-Ba/metriccollector/internal/server/config"
	"github.com/Maxim-Ba/metriccollector/internal/server/router"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
)

type Parameters struct {
	Addres string
}

func main() {

	parameterts := getParameters()
	
	_,err:=storage.New()
	if err != nil {
		panic(err)
	}
	
	mux := router.New()
	fmt.Println("Running server on ", parameterts.Addres)
	
	err = http.ListenAndServe(parameterts.Addres, mux)
	if err != nil {
		panic(err)
	}
}

func getParameters() Parameters {
	flags := config.ParseFlags()
	envConfig := config.ParseEnv()
	addres := envConfig.Addres
	if addres != "" {
		return Parameters{
			Addres: addres,
		}
	}
	return Parameters{
		Addres: flags.FlagRunAddr,
	}
}
