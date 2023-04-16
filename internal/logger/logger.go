package logger

import (
	"errors"
	"fmt"
	stdlog "log"
	"os"
	"syscall"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	LogLevelList = []string{"debug", "info", "warn", "error"}
	LogLevel     = zap.NewAtomicLevel()
)

//go:generate options-gen -out-filename=logger_options.gen.go -from-struct=Options
type Options struct {
	level          string `option:"mandatory" validate:"required,oneof=debug info warn error"`
	productionMode bool
}

func MustInit(opts Options) {
	if err := Init(opts); err != nil {
		panic(err)
	}
}

func Init(opts Options) error {
	if err := opts.Validate(); err != nil {
		return fmt.Errorf("validate options: %v", err)
	}

	err := LogLevel.UnmarshalText([]byte(opts.level))
	if err != nil {
		return err
	}

	zapCfg := zap.Config{
		EncoderConfig: zapcore.EncoderConfig{
			MessageKey:     "message",
			LevelKey:       "level",
			TimeKey:        "time",
			NameKey:        "logger",
			CallerKey:      "caller",
			StacktraceKey:  "component",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.FullCallerEncoder,
		},
	}
	var encoder zapcore.Encoder
	if opts.productionMode {
		encoder = zapcore.NewJSONEncoder(zapCfg.EncoderConfig)
	} else {
		zapCfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoder = zapcore.NewConsoleEncoder(zapCfg.EncoderConfig)
	}
	cores := []zapcore.Core{
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), LogLevel),
	}
	l := zap.New(zapcore.NewTee(cores...))
	zap.ReplaceGlobals(l)

	return nil
}

func Sync() {
	if err := zap.L().Sync(); err != nil && !errors.Is(err, syscall.ENOTTY) {
		stdlog.Printf("cannot sync logger: %v", err)
	}
}
