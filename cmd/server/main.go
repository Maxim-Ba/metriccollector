package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/server/config"
	"github.com/Maxim-Ba/metriccollector/internal/server/router"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
)

func main() {
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	parameters := config.New()
	logger.SetLogLevel(parameters.LogLevel)
	_, err := storage.New(parameters)
	if err != nil {

		panic(err)
	}
	mux := router.New()
	server := &http.Server{
		Addr:    parameters.Address,
		Handler: mux,
	}

	go func() {
		logger.LogInfo("Running server on ", parameters.Address)
		if err := http.ListenAndServe(parameters.Address, mux); err != nil && err != http.ErrServerClosed {
			logger.LogError("ListenAndServe: ", err)
		}
	}()

	<-exit // Ожидание сигнала завершения

	logger.LogInfo("Shutting down server...")

	if err := server.Shutdown(context.Background()); err != nil {
		logger.LogError("Server Shutdown: ", err)
	}

	logger.LogInfo("Server exiting")
	// Явное закрытие ресурсов
	storage.Close()
	logger.Sync()
}
