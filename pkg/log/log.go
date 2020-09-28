package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"go.uber.org/zap/zaptest/observer"
)

var (
	logger *zap.Logger
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
	return nil
}

func Logger() *zap.Logger {
	return logger
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

// With creates a child logger and adds structured context to it. Fields added
// to the child don't affect the parent, and vice versa.
func With(f ...zap.Field) *zap.Logger {
	if len(f) == 0 {
		return logger
	}
	return logger.With(f...)
}

// NewForTest returns a new logger and the corresponding observed logs which can be used in unit tests to verify log entries.
func NewForTest() (*zap.Logger, *observer.ObservedLogs) {
	core, recorded := observer.New(zapcore.InfoLevel)
	return zap.New(core), recorded
}
