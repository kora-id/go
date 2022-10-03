// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package event

import (
	"fmt"
	"os"

	"github.com/nats-io/nats.go"
)

// Nats connection instances
var Nats *nats.EncodedConn

// NatsSetup setting up nats with config
func init() {
	host := os.Getenv("NATS_HOST")

	if nc, e := nats.Connect(host); e == nil {
		Nats, _ = nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	} else {
		fmt.Println("Could not connect to nats host: ", host)
	}
}
