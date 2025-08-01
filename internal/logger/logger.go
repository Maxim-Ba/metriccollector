package logger

import (
	"net/http"
	"time"

	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"go.uber.org/zap"
)

var Sugar zap.SugaredLogger
var atomicLevel = zap.NewAtomicLevel()

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"stdout",
	}
	config.Level = atomicLevel
	atomicLevel.SetLevel(zap.ErrorLevel)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	// делаем регистратор SugaredLogger
	Sugar = *logger.Sugar()
	LogInfo("logger is initialized")
}
func Sync() {
	if err := Sugar.Sync(); err != nil {
		Sugar.Error(err)
	}
}
func SetLogLevel(level string) {
	if level == "debug" {
		atomicLevel.SetLevel(zap.DebugLevel)
	} else if level == "info" {
		atomicLevel.SetLevel(zap.InfoLevel)
	} else if level == "warn" {
		atomicLevel.SetLevel(zap.WarnLevel)
	} else if level == "error" {
		atomicLevel.SetLevel(zap.ErrorLevel)
	}
}
func LogInfo(params ...interface{}) {
	Sugar.Info(params...)
}
func LogError(params ...interface{}) {
	Sugar.Error(params...)
}
func LogMetric(m metrics.Metrics) {
	var Value float64
	var Delta int64
	if m.Value != nil {
		Value = *m.Value
	}
	if m.Delta != nil {
		Delta = *m.Delta
	}
	Sugar.Infoln(
		"MType", m.MType,
		"ID", m.ID,
		"Value", Value,
		"Delta", Delta,
	)
}

func LogResponse(r *http.Request, duration time.Duration, status, size int) {

	Sugar.Infoln(
		"uri", r.RequestURI,
		"method", r.Method,
		"status", status, // получаем перехваченный код статуса ответа
		"duration", duration,
		"size", size, // получаем перехваченный размер ответа
	)
}
