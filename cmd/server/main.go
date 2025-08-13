package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	_ "net/http/pprof"

	"github.com/Maxim-Ba/metriccollector/internal/logger"
	"github.com/Maxim-Ba/metriccollector/internal/server/config"
	"github.com/Maxim-Ba/metriccollector/internal/server/router"
	"github.com/Maxim-Ba/metriccollector/internal/server/services/subnet"
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
	buildinfo.PrintBuildInfo(buildVersion, buildDate, buildCommit)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)

	parameters := config.New()
	p, err := profiler.New(parameters.IsProfileOn, parameters.ProfileFileCPU, parameters.ProfileFileMem)
	if err != nil {
		logger.LogError("Profiler error ", err)
	}
	p.Start()
	signature.New(parameters.Key, parameters.CryptoKeyPath)
	subnet.New(parameters.TrustedSubnet)
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

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()

		logger.LogInfo("Running server on ", parameters.Address)
		if err = http.ListenAndServe(parameters.Address, mux); err != nil && err != http.ErrServerClosed {
			logger.LogError("ListenAndServe: ", err)
			cancel()
		}
	}()

	select {
	case <-exit:
		logger.LogInfo("Received shutdown signal...")
	case <-ctx.Done():
		logger.LogInfo("Context cancelled, shutting down...")
	}

	logger.LogInfo("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.LogError("Server Shutdown: ", err)
		if err = server.Close(); err != nil {
			logger.LogError("Server forced close error: ", err)
		}
	}
	wg.Wait()
	err = p.Close()
	logger.LogError(err)

	storage.Close()
	logger.Sync()
}
