// Copyright (c) 2012-2016 The Revel Framework Authors, All rights reserved.
// Revel Framework source code and usage is governed by a MIT style
// license that can be found in the LICENSE file.

package cache

import (
	"testing"
	"time"
)

// Tests against a generic Cache interface.
// They should pass for all implementations.
type cacheFactory func(*testing.T, time.Duration) Cache

// Test typical cache interactions
func typicalGetSet(t *testing.T, newCache cacheFactory) {
	var err error
	cache := newCache(t, time.Hour)

	value := "foo"
	if err = cache.Set("value", value, DefaultExpiryTime); err != nil {
		t.Errorf("Error setting a value: %s", err)
	}

	value = ""
	err = cache.Get("value", &value)
	if err != nil {
		t.Errorf("Error getting a value: %s", err)
	}
	if value != "foo" {
		t.Errorf("Expected to get foo back, got %s", value)
	}
}

func expiration(t *testing.T, newCache cacheFactory) {
	// memcached does not support expiration times less than 1 second.
	var err error
	cache := newCache(t, time.Second)
	// Test Set w/ DefaultExpiryTime
	value := 10

	if err = cache.Set("int", value, time.Duration(1)*time.Second); err != nil {
		t.Errorf("Set failed: %s", err)
	}

	time.Sleep(3 * time.Second)
	if err = cache.Get("int", &value); err != ErrCacheMiss {
		t.Errorf("Expected CacheMiss, but got: %s, %d", err, value)
	}

	// Test Set w/ short time
	if err = cache.Set("int", value, time.Second); err != nil {
		t.Errorf("Set failed: %s", err)
	}
	time.Sleep(2 * time.Second)
	if err = cache.Get("int", &value); err != ErrCacheMiss {
		t.Errorf("Expected CacheMiss, but got: %s", err)
	}

	// Test Set w/ longer time.
	if err = cache.Set("int", value, time.Hour); err != nil {
		t.Errorf("Set failed: %s", err)
	}
	time.Sleep(2 * time.Second)
	if err = cache.Get("int", &value); err != nil {
		t.Errorf("Expected to get the value, but got: %s", err)
	}

	// Test Set w/ forever.
	if err = cache.Set("int", value, ForEverNeverExpiry); err != nil {
		t.Errorf("Set failed: %s", err)
	}
	time.Sleep(2 * time.Second)
	if err = cache.Get("int", &value); err != nil {
		t.Errorf("Expected to get the value, but got: %s", err)
	}
}

func emptyCache(t *testing.T, newCache cacheFactory) {
	var err error
	cache := newCache(t, time.Hour)

	err = cache.Get("notexist", 0)
	if err == nil {
		t.Errorf("Error expected for non-existent key")
	}
	if err != ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss for non-existent key: %s", err)
	}

	err = cache.Delete("notexist")
	if err != ErrCacheMiss {
		t.Errorf("Expected ErrCacheMiss for non-existent key: %s", err)
	}
}

func testReplace(t *testing.T, newCache cacheFactory) {
	var err error
	cache := newCache(t, time.Hour)

	// Replace in an empty cache.
	if err = cache.Replace("notexist", 1, ForEverNeverExpiry); err != ErrNotStored {
		t.Errorf("Replace in empty cache: expected ErrNotStored, got: %s", err)
	}

	// Set a value of 1, and replace it with 2
	if err = cache.Set("int", 1, time.Second); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	if err = cache.Replace("int", 2, time.Second); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	var i int
	if err = cache.Get("int", &i); err != nil {
		t.Errorf("Unexpected error getting a replaced item: %s", err)
	}
	if i != 2 {
		t.Errorf("Expected 2, got %d", i)
	}

	// Wait for it to expire and replace with 3 (unsuccessfully).
	time.Sleep(2 * time.Second)
	if err = cache.Replace("int", 3, time.Second); err != ErrNotStored {
		t.Errorf("Expected ErrNotStored, got: %s", err)
	}
	if err = cache.Get("int", &i); err != ErrCacheMiss {
		t.Errorf("Expected cache miss, got: %s", err)
	}
}

func testAdd(t *testing.T, newCache cacheFactory) {
	var err error
	cache := newCache(t, time.Hour)
	// Add to an empty cache.
	if err = cache.Add("int", 1, time.Second*3); err != nil {
		t.Errorf("Unexpected error adding to empty cache: %s", err)
	}

	// Try to add again. (fail)
	if err = cache.Add("int", 2, time.Second*3); err != nil {
		if err != ErrNotStored {
			t.Errorf("Expected ErrNotStored adding dupe to cache: %s", err)
		}
	}

	// Wait for it to expire, and add again.
	time.Sleep(8 * time.Second)
	if err = cache.Add("int", 3, time.Second*5); err != nil {
		t.Errorf("Unexpected error adding to cache: %s", err)
	}

	// Get and verify the value.
	var i int
	if err = cache.Get("int", &i); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if i != 3 {
		t.Errorf("Expected 3, got: %d", i)
	}
}

func testGetMulti(t *testing.T, newCache cacheFactory) {
	cache := newCache(t, time.Hour)

	m := map[string]interface{}{
		"str": "foo",
		"num": 42,
		"foo": struct{ Bar string }{"baz"},
	}

	var keys []string
	for key, value := range m {
		keys = append(keys, key)
		if err := cache.Set(key, value, time.Second*30); err != nil {
			t.Errorf("Error setting a value: %s", err)
		}
	}

	g, err := cache.GetMulti(keys...)
	if err != nil {
		t.Errorf("Error in get-multi: %s", err)
	}

	var str string
	if err = g.Get("str", &str); err != nil || str != "foo" {
		t.Errorf("Error getting str: %s / %s", err, str)
	}

	var num int
	if err = g.Get("num", &num); err != nil || num != 42 {
		t.Errorf("Error getting num: %s / %v", err, num)
	}

	var foo struct{ Bar string }
	if err = g.Get("foo", &foo); err != nil || foo.Bar != "baz" {
		t.Errorf("Error getting foo: %s / %v", err, foo)
	}
}
