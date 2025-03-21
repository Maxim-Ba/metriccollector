package main

import (
	"net/http"

	"github.com/Maxim-Ba/metriccollector/internal/server/handlers"
)



func main() {
	mux:=handlers.InitHandlers()

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}


