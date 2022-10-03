// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package mw

import (
	"fmt"
	"time"

	"git.tech.kora.id/go/rest"
	"go.uber.org/zap"
)

// HTTPLogger returns a middleware that logs HTTP requests.
func HTTPLogger() rest.MiddlewareFunc {
	return func(n rest.HandlerFunc) rest.HandlerFunc {
		return func(c *rest.Context) error {
			return logRequest(n, c)
		}
	}
}

// logRequest print all http request on consoles.
func logRequest(hand rest.HandlerFunc, c *rest.Context) (err error) {
	start := time.Now()
	req := c.Request()
	res := c.Response()
	if err = hand(c); err != nil {
		c.Error(err)
	}
	end := time.Now()
	latency := end.Sub(start) / 1e5

	var fields = []zap.Field{
		zap.String("path", req.URL.Path),
		zap.String("query", req.URL.RawQuery),
		zap.String("ip", req.RemoteAddr),
		zap.String("user-agent", req.UserAgent()),
		zap.String("latecy", fmt.Sprintf("%1.1fms", float64(latency))),
	}

	if err == nil {
		c.Logger().Info(fmt.Sprintf("%s/%d", req.Method, res.Status), fields...)
	} else {
		c.Logger().Warn(fmt.Sprintf("%s/%d", req.Method, res.Status), fields...)
	}

	return
}
