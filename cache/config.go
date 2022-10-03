// Copyright 2018 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import "git.tech.kora.id/go/env"

// Config instance of configuration application
// for cache
var Config *RedisConfig

// RedisConfig the configuration of redis server
type RedisConfig struct {
	Host           string
	Password       string
	Protocol       string
	MaxIdle        int
	MaxActive      int
	IdleTimeout    int
	TimeoutConnect int
	TimeoutRead    int
	TimeoutWrite   int
	DefaultExpire  int
}

func init() {
	Config = &RedisConfig{
		Host:           env.GetString("REDIS_HOST", ":6379"),
		Password:       env.GetString("REDIS_PASSWORD", ""),
		Protocol:       env.GetString("REDIS_PROTOCOL", "tcp"),
		MaxIdle:        env.GetInt("REDIS_MAXIDLE", 5),
		MaxActive:      env.GetInt("REDIS_MAXACTIVE", 0),
		IdleTimeout:    env.GetInt("REDIS_IDLE_TIMEOUT", 240),
		TimeoutConnect: env.GetInt("REDIS_TIMEOUT_CONNECT", 10000),
		TimeoutRead:    env.GetInt("REDIS_TIMEOUT_READ", 5000),
		TimeoutWrite:   env.GetInt("REDIS_TIMEOUT_WRITE", 5000),
		DefaultExpire:  env.GetInt("REDIS_DEFAULT_EXPIRE", 10000),
	}

	Instance = NewRedisCache()
}
