// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package log

import (
	"fmt"
	"os"

	"go.uber.org/zap"
)

// Logger global instance of zap logger
var Logger = New(os.Getenv("APP_NAME"), false)

// Debug an helper to call logger debug
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Info an helper to call logger info
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Warn an helper to call logger warn
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error an helper to call logger error
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Debugf an helper to call logger debug with fmt string
func Debugf(f string, v ...interface{}) {
	Logger.Debug(fmt.Sprintf(f, v...))
}

// Infof an helper to call logger info with fmt string
func Infof(f string, v ...interface{}) {
	Logger.Info(fmt.Sprintf(f, v...))
}

// Warnf an helper to call logger warn with fmt string
func Warnf(f string, v ...interface{}) {
	Logger.Warn(fmt.Sprintf(f, v...))
}

// Errorf an helper to call logger error with fmt string
func Errorf(f string, v ...interface{}) {
	Logger.Error(fmt.Sprintf(f, v...))
}

// Debugw an print logger with suggared log
func Debugw(msg string, kv ...interface{}) {
	Logger.Sugar().Debugw(msg, kv...)
}

// Infow an print logger with suggared log
func Infow(msg string, kv ...interface{}) {
	Logger.Sugar().Infow(msg, kv...)
}

// Warnw an print logger with suggared log
func Warnw(msg string, kv ...interface{}) {
	Logger.Sugar().Warnw(msg, kv...)
}

// Errorw an print logger with suggared log
func Errorw(msg string, kv ...interface{}) {
	Logger.Sugar().Errorw(msg, kv...)
}
