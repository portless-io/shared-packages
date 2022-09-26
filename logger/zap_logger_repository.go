package logger

import (
	"github.com/mattn/go-colorable"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	log *zap.SugaredLogger
}

func NewZapLogger(lv string, logFilePath string, environment string) (Logger, error) {
	switch environment {
	case "PRODUCTION":
		return NewProductionLogger(lv, logFilePath)
	case "DEVELOPMENT":
		return NewDevelopmentLogger(lv)
	default:
		return NewDevelopmentLogger(lv)
	}
}

func NewDevelopmentLogger(lv string) (Logger, error) {
	level := zap.NewAtomicLevelAt(getLevel((lv)))
	config := zap.NewDevelopmentEncoderConfig()
	config.TimeKey = ""
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncodeTime = zapcore.RFC3339TimeEncoder
	config.CallerKey = "caller"
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(config),
			zapcore.AddSync(colorable.NewColorableStdout()), level),
	)

	logger := zap.New(core, zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(1))
	return &zapLogger{
		log: logger.Sugar(),
	}, nil
}

func NewProductionLogger(lv string, logFilePath string) (Logger, error) {
	level := zap.NewAtomicLevelAt(getLevel((lv)))

	config := zap.NewProductionConfig()

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder

	config.EncoderConfig = encoderConfig
	config.Level = level
	config.OutputPaths = []string{
		"stdout",
	}

	if logFilePath != "" {
		config.OutputPaths = append(config.OutputPaths, logFilePath)
	}

	logger, err := config.Build()
	log := zap.New(logger.Core(), zap.AddStacktrace(zapcore.ErrorLevel), zap.AddCallerSkip(1))

	return &zapLogger{
		log: log.Sugar(),
	}, err
}

func getLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

func mapStringToKeyPairs(m map[string]interface{}) []interface{} {
	args := []interface{}{}
	for k, v := range m {
		args = append(args, k)
		args = append(args, v)
	}

	return args
}

func (z *zapLogger) Debug(args ...interface{}) {
	z.log.Debug(args...)
}

func (z *zapLogger) Info(args ...interface{}) {
	z.log.Info(args...)
}

func (z *zapLogger) Warning(args ...interface{}) {
	z.log.Warn(args...)
}

func (z *zapLogger) Error(args ...interface{}) {
	z.log.Error(args...)
}

func (z *zapLogger) Fatal(args ...interface{}) {
	z.log.Fatal(args...)
}

func (z *zapLogger) Debugf(template string, args ...interface{}) {
	z.log.Debugf(template, args...)
}

func (z *zapLogger) Infof(template string, args ...interface{}) {
	z.log.Infof(template, args...)
}

func (z *zapLogger) Warningf(template string, args ...interface{}) {
	z.log.Warnf(template, args...)
}

func (z *zapLogger) Errorf(template string, args ...interface{}) {
	z.log.Errorf(template, args...)
}

func (z *zapLogger) Fatalf(template string, args ...interface{}) {
	z.log.Fatalf(template, args...)
}

func (z *zapLogger) Debugw(template string, context map[string]interface{}) {
	z.log.Debugw(template, mapStringToKeyPairs(context)...)
}

func (z *zapLogger) Infow(template string, context map[string]interface{}) {
	z.log.Infow(template, mapStringToKeyPairs(context)...)
}

func (z *zapLogger) Warningw(template string, context map[string]interface{}) {
	z.log.Warnw(template, mapStringToKeyPairs(context)...)
}

func (z *zapLogger) Errorw(template string, context map[string]interface{}) {
	z.log.Errorw(template, mapStringToKeyPairs(context)...)
}

func (z *zapLogger) Fatalw(template string, context map[string]interface{}) {
	z.log.Fatalw(template, mapStringToKeyPairs(context)...)
}
