package logger

import (
	"context"
	"fmt"
	golog "log"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	*zap.Logger
}

func L() *Logger {
	return logger
}

var logger *Logger = &Logger{zap.L()}

var withContextHandler func(*Logger, context.Context) *Logger = DefaultWithContextHandler

func InitialLogger() *Logger {
	zapConfig := NewZapLoggerConfig()

	zapLogger, err := zapConfig.Build(zap.AddCaller(), zap.AddCallerSkip(0))
	if err != nil {
		golog.Fatalf("can't initalize zap logger: %+v", err)
	}

	logger = &Logger{zapLogger}

	zap.ReplaceGlobals(zapLogger)

	defer logger.Sync()
	return logger
}

func NewZapLoggerConfig() zap.Config {
	lv := zapcore.DebugLevel

	return zap.Config{
		Level:       zap.NewAtomicLevelAt(lv),
		Development: true,
		// Encoding:          "json",
		Encoding:          "console",
		EncoderConfig:     NewEncoderConfig(),
		OutputPaths:       []string{"stderr"},
		ErrorOutputPaths:  []string{"stderr"},
		DisableCaller:     true,
		DisableStacktrace: true,
	}
}

func NewEncoderConfig() zapcore.EncoderConfig {
	return zapcore.EncoderConfig{
		LevelKey:       "severity",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    "func",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		TimeKey:        "time",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.MillisDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}
}

func (log *Logger) WithContext(ctx context.Context) *Logger {
	return withContextHandler(log, ctx)
}

func (log *Logger) With(fields ...zap.Field) *Logger {
	return &Logger{log.Logger.With(fields...)}
}

func (log *Logger) Named(s string) *Logger {
	return &Logger{log.Logger.Named(s)}
}

func (log *Logger) Infof(format string, param ...interface{}) {
	log.Info(fmt.Sprintf(format, param...))
}

func (log *Logger) Debugf(format string, param ...interface{}) {
	log.Debug(fmt.Sprintf(format, param...))
}

func (log *Logger) Errorf(format string, param ...interface{}) {
	log.Error(fmt.Sprintf(format, param...))
}

func (log *Logger) Warnf(format string, param ...interface{}) {
	log.Warn(fmt.Sprintf(format, param...))
}

func (log *Logger) Fatalf(format string, param ...interface{}) {
	log.Fatal(fmt.Sprintf(format, param...))
}

func DefaultWithContextHandler(log *Logger, ctx context.Context) *Logger {
	if ctx == nil {
		return log
	}
	if fs, ok := ctx.Value("logFields").([]zap.Field); ok {
		return log.With(fs...)
	} else {
		return log
	}
}
