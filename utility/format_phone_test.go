// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utility_test

import (
	"testing"

	"git.tech.kora.id/go/utility"

	"github.com/stretchr/testify/assert"
)

func TestFormatPhone(t *testing.T) {
	var tests = []struct {
		param    string
		expected string
	}{
		{"08121166567", "628121166567"},
		{"081-211-66567", "628121166567"},
		{"628121166567", "628121166567"},
		{"8121166567", "628121166567"},
		{"08121166508", "628121166508"},
		{"28121166567", ""},
		{"0812", ""},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.FormatPhone(test.param))
	}
}
