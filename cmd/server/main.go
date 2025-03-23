package main

import (
	"fmt"
	"net/http"

	"github.com/Maxim-Ba/metriccollector/internal/server/handlers"
)



func main() {
	parseFlags()
	mux:=handlers.InitHandlers()
	fmt.Println("Running server on", flagRunAddr)

	err := http.ListenAndServe(flagRunAddr, mux)
	if err != nil {
		panic(err)
	}
}


