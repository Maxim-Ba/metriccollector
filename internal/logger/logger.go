package logger

import (
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger
var atomicLevel = zap.NewAtomicLevel()

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"logs.log",
		"stdout",
	}
	config.Level = atomicLevel
	atomicLevel.SetLevel(zap.ErrorLevel)

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	
	// делаем регистратор SugaredLogger
	sugar = *logger.Sugar()
	LogInfo("logger is initialized")
}
func Sync(){
	if err := sugar.Sync(); err != nil {
		sugar.Error(err)
	}
}
func SetLogLevel (level string){
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
	sugar.Info(params...)
}
func LogError(params ...interface{}) {
	sugar.Error(params...)
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
	sugar.Infoln(
		"MType", m.MType,
		"ID", m.ID,
		"Value", Value,
		"Delta", Delta,
	)
}
