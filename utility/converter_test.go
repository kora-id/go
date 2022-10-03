// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utility_test

import (
	"fmt"
	"testing"

	"git.tech.kora.id/go/utility"

	"github.com/stretchr/testify/assert"
)

func TestToInt(t *testing.T) {
	var tests = []struct {
		param    interface{}
		expected int
	}{
		{"1000", 1000},
		{"-123", -123},
		{"0.1", 0},
		{"string", 0},
		{"100000000000000000000000000000000000000000000", 0},
		{" 1", 1},
		{1, 1},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.ToInt(test.param))
	}
}

func TestToBoolean(t *testing.T) {
	var tests = []struct {
		param    interface{}
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"True", true},
		{"False", false},
		{nil, false},
		{"0", false},
		{"1", true},
		{0, false},
		{1, true},
		{"string", false},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.ToBoolean(test.param))
	}
}

func TestToString(t *testing.T) {
	var tests = []struct {
		param    interface{}
		expected string
	}{
		{"str123", "str123"},
		{"123", "123"},
		{"12.3", "12.3"},
		{true, "true"},
		{1.5 + 10i, "(1.5+10i)"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.ToString(test.param))
	}
}

func TestToFloat(t *testing.T) {
	var tests = []struct {
		param    interface{}
		expected float64
	}{
		{"", 0},
		{"123", 123},
		{"-.01", -0.01},
		{"string", 0},
		{"10.", 10.0},
		{"1.23e3", 1230},
		{".23e10", 0.23e10},
		{[]string{"test"}, 0},
		{0.1, 0.1},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.ToFloat(test.param), test.param)
	}
}

func TestToJSON(t *testing.T) {
	tests := []interface{}{"test", map[string]string{"a": "b", "b": "c"}, func() error {
		return fmt.Errorf("error")
	}}
	expected := [][]string{
		{"\"test\"", "<nil>"},
		{"{\"a\":\"b\",\"b\":\"c\"}", "<nil>"},
		{"", "json: unsupported type: func() error"},
	}
	for i, test := range tests {
		assert.Equal(t, expected[i][0], utility.ToJSON(test))
	}
}

func TestToCamelCase(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected string
	}{
		{"a_b_c", "ABC"},
		{"my_func", "MyFunc"},
		{"1ab_cd", "1abCd"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.ToCamelCase(test.param))
	}

	tests = []struct {
		param    string
		expected string
	}{
		{"abc", "abc"},
		{"a_b_c", "aBC"},
		{"my_func", "myFunc"},
		{"1ab_cd", "1abCd"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.ToLowerCamelCase(test.param))
	}
}

func TestToUnderscore(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected string
	}{
		{"MyFunc", "my_func"},
		{"ABC", "a_b_c"},
		{"1B", "1_b"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.ToUnderscore(test.param))
	}
}

func TestToUpper(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected string
	}{
		{"a_b_c", "A_B_C"},
		{"my_func", "MY_FUNC"},
		{"1ab_cd", "1AB_CD"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.ToUpper(test.param))
	}
}

func TestToLower(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected string
	}{
		{"a_b_c", "A_B_C"},
		{"my_func", "MY_FUNC"},
		{"1ab_cd", "1AB_CD"},
	}
	for _, test := range tests {
		assert.Equal(t, test.param, utility.ToLower(test.expected))
	}
}

func TestLeftTrim(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   string
		param2   string
		expected string
	}{
		{"  \r\n\tfoo  \r\n\t   ", "", "foo  \r\n\t   "},
		{"010100201000", "01", "201000"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.LeftTrim(test.param1, test.param2))
	}
}

func TestRightTrim(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   string
		param2   string
		expected string
	}{
		{"  \r\n\tfoo  \r\n\t   ", "", "  \r\n\tfoo"},
		{"010100201000", "01", "0101002"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.RightTrim(test.param1, test.param2))
	}
}

func TestTrim(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   string
		param2   string
		expected string
	}{
		{"  \r\n\tfoo  \r\n\t   ", "", "foo"},
		{"010100201000", "01", "2"},
		{"1234567890987654321", "1-8", "909"},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, utility.Trim(test.param1, test.param2))
	}
}
