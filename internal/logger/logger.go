package logger

import (
	"go.uber.org/zap"
)

var sugar zap.SugaredLogger

func InitLogger() {
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

func LogInfo( message string)  {
	sugar.Infoln(
		"message", message,
)
}
