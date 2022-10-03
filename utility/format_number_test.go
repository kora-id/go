// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utility_test

import (
	"testing"

	"git.tech.kora.id/go/utility"

	"github.com/stretchr/testify/assert"
)

func TestFormatNumber(t *testing.T) {
	var tests = []struct {
		param    float64
		expected string
	}{
		{1000, "1.000"},
		{10000, "10.000"},
		{100000, "100.000"},
		{1000000, "1.000.000"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.FormatNumber("#.###,", test.param))
	}
}
