package logger

import (
	"github.com/Maxim-Ba/metriccollector/internal/models/metrics"
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func init() {
    // создаём предустановленный регистратор zap
    logger, err := zap.NewDevelopment()
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
	sugar.Infoln(
		"MType", m.MType,
		"ID", m.ID,
		"Value", m.Value,
		"Delta", m.Delta,
)}
