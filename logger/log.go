package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
)

var logger = zap.NewNop()

func newLogger(production bool) (l *zap.Logger) {
	if production {
		l, _ = zap.NewProduction()
	} else {
		l, _ = zap.NewDevelopment()
	}
	return
}

// EnableLogging Enable logger for force block ev
func EnableLogging(production bool) {
	logger = buildLogger("./force_relay")
}

// Logger get logger
func Logger() *zap.Logger {
	return logger
}

// Sugar get Sugar logger
func Sugar() *zap.SugaredLogger {
	return logger.Sugar()
}

func buildLogger(logPath string) *zap.Logger {

	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	topicDebugging := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath + ".log",
		MaxSize:    500, // megabytes
		MaxBackups: 32,
		MaxAge:     7, // days
	})

	topicErrors := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath + ".error.log",
		MaxSize:    500, // megabytes
		MaxBackups: 64,
		MaxAge:     7, // days
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	fileEncoder := zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, topicErrors, highPriority),
		zapcore.NewCore(consoleEncoder, consoleErrors, highPriority),
		zapcore.NewCore(fileEncoder, topicDebugging, lowPriority),
		zapcore.NewCore(consoleEncoder, consoleDebugging, lowPriority),
	)

	return zap.New(core)
}
