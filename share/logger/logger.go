package logger

import (
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"gopros/share/config"
	"log"
)

var (
	gLogger *zap.Logger
	err error
	gFlags config.Flags
)

type AppLoggerConfig struct {
	Zap struct{
		Logger zap.Config
	}
}




func init()  {
	gFlags.Init()
	err = godotenv.Load(gFlags.Env)
	if err != nil {
		log.Fatal(err)
	}
	var cfg AppLoggerConfig
	err = config.Load(&cfg)
	if err != nil {
		log.Fatal(err)
	}
	gLogger, err = cfg.Zap.Logger.Build(zap.AddCallerSkip(1))
	if err != nil {
		log.Fatal(err)
	}
}


func D(fields ...zap.Field)  {
	gLogger.Debug("",fields...)
}

func Debug(msg string,fields ...zap.Field)  {
	gLogger.Debug(msg,fields...)
}

func W(fields ...zap.Field)  {
	gLogger.Warn("",fields...)
}

func Warn(msg string,fields ...zap.Field)  {
	gLogger.Warn(msg,fields...)
}

func I(fields ...zap.Field)  {
	gLogger.Info("",fields...)
}

func Info(msg string,fields ...zap.Field)  {
	gLogger.Info(msg,fields...)
}

func E(fields ...zap.Field)  {
	gLogger.Error("",fields...)
}

func Error(msg string,fields ...zap.Field)  {
	gLogger.Error(msg,fields...)
}

func F(fields ...zap.Field)  {
	gLogger.Fatal("",fields...)
}

func Fatal(msg string,fields ...zap.Field)  {
	gLogger.Fatal(msg,fields...)
}