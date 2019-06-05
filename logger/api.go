package logger

import "go.uber.org/zap"

// LogError log error
func LogError(msg string, err error) {
	logger.Error(msg, zap.Error(err))
}

// Debugf formats message according to format specifier
// and writes to default logger with log level = Debug.
func Debugf(format string, params ...interface{}) {
	logger.Sugar().Debugf(format, params...)
}

// Infof formats message according to format specifier
// and writes to default logger with log level = Info.
func Infof(format string, params ...interface{}) {
	logger.Sugar().Infof(format, params...)
}

// Warnf formats message according to format specifier and writes to default logger with log level = Warn
func Warnf(format string, params ...interface{}) {
	logger.Sugar().Warnf(format, params...)
}

// Errorf formats message according to format specifier and writes to default logger with log level = Error
func Errorf(format string, params ...interface{}) {
	logger.Sugar().Errorf(format, params...)
}
