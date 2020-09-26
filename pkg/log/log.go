package log

import (
	"context"
	"go.uber.org/zap"
	"moul.io/zapgorm2"
)

var (
	logger   *zap.Logger
	dbLogger zapgorm2.Logger
)

func Init(ctx context.Context, debug bool) error {
	var err error
	if debug {
		logger, err = zap.NewDevelopment()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		_ = logger.Sync()
	}()
	dbLogger = zapgorm2.New(logger)
	dbLogger.SetAsDefault()
	return nil
}

func Logger() *zap.Logger {
	return logger
}

func DbLogger() zapgorm2.Logger {
	return dbLogger
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Debug(msg string, f ...zap.Field) {
	logger.Debug(msg, f...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Info(msg string, f ...zap.Field) {
	logger.Info(msg, f...)
}

// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
func Error(msg string, f ...zap.Field) {
	logger.Error(msg, f...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func Fatal(msg string, f ...zap.Field) {
	logger.Fatal(msg, f...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields accumulated on the logger.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func Panic(msg string, f ...zap.Field) {
	logger.Panic(msg, f...)
}
