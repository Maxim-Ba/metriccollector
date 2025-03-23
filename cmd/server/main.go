package main

import (
	"fmt"
	"net/http"

	"github.com/Maxim-Ba/metriccollector/internal/server/handlers"
)

type Parameters struct {
	Addres string
}

func main() {

	parameterts := getParameters()
	mux := handlers.InitHandlers()

	fmt.Println("Running server on ", parameterts.Addres)

	err := http.ListenAndServe(parameterts.Addres, mux)
	if err != nil {
		panic(err)
	}
}

func getParameters() Parameters {
	parseFlags()
	envConfig := parseEnv()
	addres := envConfig.Addres
	if addres != "" {
		return Parameters{
			Addres: addres,
		}
	}
	return Parameters{
		Addres: flagRunAddr,
	}
}
