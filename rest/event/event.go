// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

import (
	"errors"
	"sync"
	"time"
)

// NOTFOUND error when event not found
const NOTFOUND = "E_NOT_FOUND"

// internal mapping of event names to observing channels
var events = make(map[string][]chan interface{})

// mutex for touching the event map
var rwMutex sync.RWMutex

// Listen start observing the specified event via provided output channel
func Listen(event string, outputChan chan interface{}) {
	rwMutex.Lock()
	defer rwMutex.Unlock()

	events[event] = append(events[event], outputChan)
}

// Call a notification (arbitrary data) to the specified event
func Call(event string, data interface{}) error {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	outChans, ok := events[event]
	if !ok {
		return errors.New(NOTFOUND)
	}
	for _, outputChan := range outChans {
		outputChan <- data
	}

	return nil
}

// CallTimeout a notification to the specified event using the provided timeout for
// any output channels that are blocking
func CallTimeout(event string, data interface{}, timeout time.Duration) error {
	rwMutex.RLock()
	defer rwMutex.RUnlock()

	outChans, ok := events[event]
	if !ok {
		return errors.New(NOTFOUND)
	}
	for _, outputChan := range outChans {
		select {
		case outputChan <- data:
		case <-time.After(timeout):
		}
	}

	return nil
}
