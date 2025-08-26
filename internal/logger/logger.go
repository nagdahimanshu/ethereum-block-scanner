// Package logger implements logger
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	LevelTrace = "trace"
	LevelDebug = "debug"
	LevelInfo  = "info"
	LevelWarn  = "warn"
	LevelError = "error"
	LevelFatal = "fatal"
)

// DefaultLogger is a default setup logger.
var DefaultLogger Logger

// Logger interface used in the sdk.
type Logger interface {
	Debug(...interface{})
	Debugf(msg string, others ...interface{})
	Debugw(msg string, keysAndValues ...interface{})
	Info(...interface{})
	Infof(msg string, others ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warn(...interface{})
	Warnf(msg string, others ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Error(...interface{})
	Errorf(msg string, others ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Fatal(...interface{})
	Fatalf(msg string, others ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})
	With(kv ...interface{}) Logger
}

func init() {
	zlog, err := zap.NewDevelopment(
		zap.AddCallerSkip(1),
	)
	if err != nil {
		panic(err)
	}
	defer func() {
		zlog.Sync() //nolint:errcheck // https://github.com/uber-go/zap/issues/880
	}()
	zlog.WithOptions()
	DefaultLogger = &logger{SugaredLogger: zlog.Sugar()}
}

func getZapcoreLevel(logLevel string) zapcore.Level {
	level, _ := zapcore.ParseLevel(logLevel)
	return level
}

// NewDefaultProductionLogger returns zap logger with default production setting.
func NewDefaultProductionLogger(logLevel string) (Logger, error) {
	zapLogLevel := getZapcoreLevel(logLevel)
	var zLog *zap.Logger
	var err error
	if zapLogLevel > zapcore.DebugLevel {
		config := zap.Config{
			Level:       zap.NewAtomicLevelAt(zap.InfoLevel),
			Development: false,
			Sampling: &zap.SamplingConfig{
				Initial:    100,
				Thereafter: 100,
			},
			Encoding: "json",
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:        "@timestamp",
				LevelKey:       "log.level",
				NameKey:        "logger",
				CallerKey:      "caller",
				FunctionKey:    zapcore.OmitKey,
				MessageKey:     "message",
				StacktraceKey:  "stacktrace",
				LineEnding:     zapcore.DefaultLineEnding,
				EncodeLevel:    zapcore.LowercaseLevelEncoder,
				EncodeTime:     zapcore.ISO8601TimeEncoder,
				EncodeDuration: zapcore.SecondsDurationEncoder,
				EncodeCaller:   zapcore.ShortCallerEncoder,
			},
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
		}
		zLog, err = config.Build(zap.AddCallerSkip(1))
		if err != nil {
			return nil, err
		}
	} else {
		zLog, err = zap.NewDevelopment(zap.AddCallerSkip(1))
		if err != nil {
			return nil, err
		}
	}
	logger := &logger{
		SugaredLogger: zLog.Sugar(),
	}
	return logger, nil
}

func NewSilentLogger() (Logger, error) {
	config := &zap.Config{
		Level:       zap.NewAtomicLevelAt(zap.ErrorLevel),
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}
	zlog, err := config.Build()
	if err != nil {
		return nil, err
	}
	logger := &logger{
		SugaredLogger: zlog.Sugar(),
	}
	return logger, nil
}

type logger struct {
	*zap.SugaredLogger
}

// Debug methods
func (l *logger) Debug(args ...interface{}) {
	l.SugaredLogger.Debug(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.SugaredLogger.Debugf(format, args...)
}

func (l *logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Debugw(msg, keysAndValues...)
}

// Info methods
func (l *logger) Info(args ...interface{}) {
	l.SugaredLogger.Info(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.SugaredLogger.Infof(format, args...)
}

func (l *logger) Infow(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Infow(msg, keysAndValues...)
}

// Warn methods
func (l *logger) Warn(args ...interface{}) {
	l.SugaredLogger.Warn(args...)
}

func (l *logger) Warnf(format string, args ...interface{}) {
	l.SugaredLogger.Warnf(format, args...)
}

func (l *logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Warnw(msg, keysAndValues...)
}

// Error methods
func (l *logger) Error(args ...interface{}) {
	l.SugaredLogger.Error(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.SugaredLogger.Errorf(format, args...)
}

func (l *logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Errorw(msg, keysAndValues...)
}

// Fatal methods
func (l *logger) Fatal(args ...interface{}) {
	l.SugaredLogger.Fatal(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.SugaredLogger.Fatalf(format, args...)
}

func (l *logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Fatalw(msg, keysAndValues...)
}

// Panic methods
func (l *logger) Panic(args ...interface{}) {
	l.SugaredLogger.Panic(args...)
}

func (l *logger) Panicf(format string, args ...interface{}) {
	l.SugaredLogger.Panicf(format, args...)
}

func (l *logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.SugaredLogger.Panicw(msg, keysAndValues...)
}

func (l *logger) With(kv ...interface{}) Logger {
	return &logger{
		SugaredLogger: l.SugaredLogger.With(kv...),
	}
}

// Ensure logger conforms to the Logger interface.
var _ Logger = (*logger)(nil)
