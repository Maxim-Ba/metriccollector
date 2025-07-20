package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	_ "net/http/pprof"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/server/config"
	"github.com/Maxim-Ba/metriccollector/internal/server/router"
	"github.com/Maxim-Ba/metriccollector/internal/server/storage"
	"github.com/Maxim-Ba/metriccollector/internal/signature"
	"github.com/Maxim-Ba/metriccollector/pkg/buildinfo"
	"github.com/Maxim-Ba/metriccollector/pkg/profiler"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

func main() {
	buildinfo.PrintBuildInfo(buildVersion,buildDate , buildCommit)
	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	parameters := config.New()
	p, err := profiler.New(parameters.IsProfileOn, parameters.ProfileFileCPU, parameters.ProfileFileMem)
	if err != nil {
		logger.LogError("Profiler error ", err)
	}
	p.Start()
	signature.New(parameters.Key)
	logger.SetLogLevel(parameters.LogLevel)

	_, err = storage.New(parameters)
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
		if err = http.ListenAndServe(parameters.Address, mux); err != nil && err != http.ErrServerClosed {
			logger.LogError("ListenAndServe: ", err)
		}
	}()
	<-exit // Ожидание сигнала завершения

	logger.LogInfo("Shutting down server...")

	if err = server.Shutdown(context.Background()); err != nil {
		logger.LogError("Server Shutdown: ", err)
	}
	logger.LogInfo("Server exiting")
	// Явное закрытие ресурсов
	err = p.Close()
	logger.LogError(err)

	storage.Close()
	logger.Sync()
}
