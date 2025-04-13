package logger

import (
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func init() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{
		"logs.log", 
		"stdout",   
	}

		logger, err := config.Build()
		if err != nil {
			panic(err)
		}
		defer logger.Sync()
    // делаем регистратор SugaredLogger
    sugar = *logger.Sugar()
		LogInfo("logger is initialized")
}

func LogInfo( message interface{})  {
	sugar.Infoln(
		"message", message,
)

}
func LogMetric( m metrics.Metrics)  {
	var Value float64 
	var Delta int64
	if m.Value !=nil {
		Value =  *m.Value
	}
	if m.Delta !=nil {
		Delta =  *m.Delta
	}
	sugar.Infoln(
		"MType", m.MType,
		"ID", m.ID,
		"Value", Value,
		"Delta", Delta,
)}
