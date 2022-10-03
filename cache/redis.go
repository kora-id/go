// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"fmt"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
)

// RedisCache wraps the Redis client to meet the Cache interface.
type RedisCache struct {
	pool              *redis.Pool
	defaultExpiration time.Duration
}

// NewRedisCache returns a new RedisCache with given parameters
// until redigo supports sharding/clustering, only one host will be in hostList
func NewRedisCache() RedisCache {
	var pool = &redis.Pool{
		MaxIdle:     Config.MaxIdle,
		MaxActive:   Config.MaxActive,
		IdleTimeout: time.Duration(Config.IdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			toc := time.Millisecond * time.Duration(Config.TimeoutConnect)
			tor := time.Millisecond * time.Duration(Config.TimeoutRead)
			tow := time.Millisecond * time.Duration(Config.TimeoutWrite)
			c, err := redis.DialURL(fmt.Sprintf("redis://%s", os.Getenv("REDIS_HOST")),
				redis.DialConnectTimeout(toc),
				redis.DialReadTimeout(tor),
				redis.DialWriteTimeout(tow))
			if err != nil {
				return nil, err
			}
			if len(Config.Password) > 0 {
				if _, err = c.Do("AUTH", Config.Password); err != nil {
					_ = c.Close()
					return nil, err
				}
			} else {
				// check with PING
				if _, err = c.Do("PING"); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, err
		},
		// custom connection test method
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	defaultExpiration := time.Hour * time.Duration(Config.DefaultExpire)

	return RedisCache{pool, defaultExpiration}
}

func generalizeStringSlice(strs []string) []interface{} {
	ret := make([]interface{}, len(strs))
	for i, str := range strs {
		ret[i] = str
	}
	return ret
}

func exists(conn redis.Conn, key string) (bool, error) {
	return redis.Bool(conn.Do("EXISTS", key))
}

// Set add new cache data based on the key
func (c RedisCache) Set(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	return c.invoke(conn.Do, key, value, expires)
}

// Add stored cache data but it will see if the key already exist
func (c RedisCache) Add(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	existed, err := exists(conn, key)
	if err != nil {
		return err
	} else if existed {
		return ErrNotStored
	}
	return c.invoke(conn.Do, key, value, expires)
}

// Replace stored new cache data to existing one
func (c RedisCache) Replace(key string, value interface{}, expires time.Duration) error {
	conn := c.pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	existed, err := exists(conn, key)
	if err != nil {
		return err
	} else if !existed {
		return ErrNotStored
	}

	err = c.invoke(conn.Do, key, value, expires)
	if value == nil {
		return ErrNotStored
	}
	return err
}

// Get retrive cache data based on the key
func (c RedisCache) Get(key string, ptrValue interface{}) error {
	conn := c.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	raw, err := conn.Do("GET", key)
	if err != nil {
		return err
	} else if raw == nil {
		return ErrCacheMiss
	}
	item, err := redis.Bytes(raw, err)
	if err != nil {
		return err
	}
	return Deserialize(item, ptrValue)
}

// GetMulti retrive cache data from multiple keys
func (c RedisCache) GetMulti(keys ...string) (Getter, error) {
	conn := c.pool.Get()
	defer func() {
		_ = conn.Close()
	}()

	items, err := redis.Values(conn.Do("MGET", generalizeStringSlice(keys)...))
	if err != nil {
		return nil, err
	} else if items == nil {
		return nil, ErrCacheMiss
	}

	m := make(map[string][]byte)
	for i, key := range keys {
		m[key] = nil
		if i < len(items) && items[i] != nil {
			s, ok := items[i].([]byte)
			if ok {
				m[key] = s
			}
		}
	}
	return RedisItemMapGetter(m), nil
}

// Delete all cache data based on the key
func (c RedisCache) Delete(key string) error {
	conn := c.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	existed, err := redis.Bool(conn.Do("DEL", key))
	if err == nil && !existed {
		err = ErrCacheMiss
	}
	return err
}

// Flush clear all cache data
func (c RedisCache) Flush() error {
	conn := c.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	_, err := conn.Do("FLUSHALL")
	return err
}

func (c RedisCache) invoke(f func(string, ...interface{}) (interface{}, error), key string, value interface{}, expires time.Duration) error {

	switch expires {
	case DefaultExpiryTime:
		expires = c.defaultExpiration
	case ForEverNeverExpiry:
		expires = time.Duration(0)
	}

	b, err := Serialize(value)
	if err != nil {
		return err
	}
	conn := c.pool.Get()
	defer func() {
		_ = conn.Close()
	}()
	if expires > 0 {
		_, err = f("SETEX", key, int32(expires/time.Second), b)
		return err
	}
	_, err = f("SET", key, b)
	return err
}

// RedisItemMapGetter implements a Getter on top of the returned item map.
type RedisItemMapGetter map[string][]byte

// Get desirialization the value into the pointer provided
func (g RedisItemMapGetter) Get(key string, ptrValue interface{}) error {
	item, ok := g[key]
	if !ok {
		return ErrCacheMiss
	}
	return Deserialize(item, ptrValue)
}
