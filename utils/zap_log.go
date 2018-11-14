package utils

import (
	"log"
	"runtime"
	"time"

	"github.com/spark-golang/spark-url/utils/env"

	"fmt"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LogLevelFatal  = "fatal"
	LogLevelError  = "error"
	LogLevelWarn   = "warning"
	LogLevelNotice = "notice"
	LogLevelInfo   = "info"
	LogLevelDebug  = "debug"
)

var (
	Logger *zap.SugaredLogger
)

func InitLog() {
	logLevel := -1
	cfg := zap.NewProductionConfig()
	if logPath := env.Getenv("LOG_PATH_FILE"); logPath != "" {
		cfg.OutputPaths = []string{logPath}
	}

	cfg.ErrorOutputPaths = []string{}
	cfg.Level.SetLevel(zapcore.Level(logLevel))

	cfg.EncoderConfig.LevelKey = ""
	cfg.EncoderConfig.TimeKey = ""
	cfg.EncoderConfig.NameKey = ""
	cfg.EncoderConfig.CallerKey = ""

	logger, err := cfg.Build()
	if err != nil {
		log.Fatalf("failed to init log: %s", err.Error())
	}

	Logger = logger.Sugar()
}

func Zlog(logLevel, logType, msg, ctx string) {
	if Logger == nil {
		return
	}

	curDate := time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")

	_, file, line, ok := runtime.Caller(1)
	caller := ""
	if ok {
		caller = fmt.Sprintf("%s, line:%d", file, line)
	}

	Logger.Infow(msg, "level", logLevel, "category", logType, "ctx", ctx, "ts", curDate, "caller", caller)
}
