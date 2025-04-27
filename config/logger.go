package config

import (
	"interviewGenius/pkg/setting"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// InitLogger 初始化日志
func InitLogger() error {
	// 配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 输出到控制台
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
	consoleWriter := zapcore.AddSync(os.Stdout)

	// 设置日志级别
	var level zapcore.Level
	if setting.ServerSetting.RunMode == "debug" {
		level = zapcore.DebugLevel
	} else {
		level = zapcore.InfoLevel
	}

	core := zapcore.NewCore(consoleEncoder, consoleWriter, level)
	logger := zap.New(core, zap.AddCaller())

	// 替换全局logger
	zap.ReplaceGlobals(logger)

	return nil
}
