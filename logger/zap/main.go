package main

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"log"
	"os"
	"time"
)

func main() {
	logInit()
	log.Printf("Hello, World , %v", "zap")
}

func logInit() {
	customEncoder := zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
		TimeKey:        "T",
		LevelKey:       "L",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     customTimeEncoder, // Custom time format
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// File encoder (JSON format)
	fileEncoder := zapcore.NewJSONEncoder(zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  zapcore.OmitKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder, // JSON time format
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	})

	// Log file setup
	file, _ := os.Create("logfile.log")
	writeSyncer := zapcore.AddSync(file)

	// Cores for stdout and file logging
	core := zapcore.NewTee(
		zapcore.NewCore(customEncoder, zapcore.AddSync(os.Stdout), zap.InfoLevel), // Stdout
		zapcore.NewCore(fileEncoder, writeSyncer, zap.InfoLevel),                  // Log file
	)

	// Create logger
	logger := zap.New(core)
	defer logger.Sync()

	// 将标准库 log 的输出重定向到 Zap
	zap.ReplaceGlobals(logger) // 替换全局 Zap 日志器
	zap.RedirectStdLog(logger) // 重定向标准库日志到 Zap
	// Log an example message
	logger.Info("Exposing 2 endpoint(s) beneath base path '/actuator'",
		zap.String("method", "org.springframework.boot.actuate.endpoint.web.EndpointLinksResolver,<init>,58"))

	logger.Info("Init CodeGenerate Config [ Get Db Config From application.yml ]",
		zap.String("method", "org.jeecg.config.init.CodeGenerateDbConfig,initCodeGenerateDbConfig,46"))
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05")) // yyyy-MM-dd HH:mm:ss
}
