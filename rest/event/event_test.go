// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"git.tech.kora.id/go/rest/event"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	// declaring listener
	listener("event:one")
	listener("event:two")

	res := m.Run()

	os.Exit(res)
}

func listener(eventName string) {
	eventChan := make(chan interface{})
	event.Listen(eventName, eventChan)

	go func() {
		for {
			data := <-eventChan
			fmt.Println(data)
		}
	}()
}

func TestCall(t *testing.T) {
	e := event.Call("event:one", "hello")
	assert.NoError(t, e)
	e = event.Call("event:two", "ladies")
	assert.NoError(t, e)
}

func TestCallTimeout(t *testing.T) {
	e := event.CallTimeout("event:one", "hello", time.Duration(10)*time.Minute)
	assert.NoError(t, e)
	e = event.CallTimeout("event:two", "ladies", time.Duration(10)*time.Minute)
	assert.NoError(t, e)
}
