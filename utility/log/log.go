// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New making new instances
func New(name string, dev bool) (l *zap.Logger) {
	l = mode(os.Getenv("APP_MODE") == "DEV" || dev)

	return l.Named(name)
}

func mode(isDev bool) (l *zap.Logger) {
	cfg := zap.Config{
		Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:      false,
		Encoding:         "json",
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig: zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "lvl",
			NameKey:        "eng",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		},
	}

	if isDev {
		cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
		cfg.Development = true
		cfg.Encoding = "console"
		cfg.EncoderConfig = zapcore.EncoderConfig{
			NameKey:      "log",
			MessageKey:   "message",
			TimeKey:      "time",
			EncodeTime:   customTimeEncoder,
			LevelKey:     "level",
			EncodeLevel:  customLevelEncoder,
			CallerKey:    "file",
			EncodeCaller: customCallerEncoder,
		}
	}

	l, _ = cfg.Build()

	return
}

func customTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(grey(t.Format("02/01 15:04:05")))
}

func customLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(colorized(level))
}

func customCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(grey(fmt.Sprintf("@%s", caller.TrimmedPath())))
}

func colorized(level zapcore.Level) string {
	if level < zapcore.DebugLevel {
		return dim(level.CapitalString())
	}
	switch level {
	case zapcore.DebugLevel:
		return green(level.CapitalString())
	case zapcore.InfoLevel:
		return magenta(level.CapitalString())
	case zapcore.WarnLevel:
		return yellow(level.CapitalString())
	default:
		return red(level.CapitalString())
	}
}
