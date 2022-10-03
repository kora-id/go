// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation_test

import (
	"testing"
	"time"

	"git.tech.kora.id/go/validation"

	"github.com/stretchr/testify/assert"
)

func TestIsNotEmpty(t *testing.T) {
	t.Parallel()

	tt := time.Time{}
	var tests = []struct {
		param    interface{}
		expected bool
	}{
		{"", false},
		{true, true},
		{false, true},
		{0, false},
		{0.2, true},
		{nil, false},
		{time.Now(), true},
		{[]string{}, false},
		{uint(4), true},
		{int8(4), true},
		{uint8(4), true},
		{int16(4), true},
		{uint16(4), true},
		{int32(4), true},
		{uint32(4), true},
		{int64(4), true},
		{uint64(4), true},
		{tt, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsNotEmpty(test.param))
	}
}

func TestIsNumeric(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    interface{}
		expected bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", false},
		{"소주", false},
		{"ABC", false},
		{"FoObAr", false},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", false}, //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", true},  //UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", true},
		{"+00123", true},
		{"0", true},
		{"-0", true},
		{"123.123", true},
		{" ", false},
		{".", false},
		{"12𐅪3", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", false},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", true},
		{1, true},
		{0.1, true},
		{-1, true},
		{&struct{}{}, false},
		{[]string{"a"}, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsNumeric(test.param))
	}
}

func TestIsAlpha(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", true},
		{"소주", false},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", true},  //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", false}, //UTF-8(ASCII): 0
		{"123", false},
		{"0123", false},
		{"-00123", false},
		{"0", false},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", true},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsAlpha(test.param))
	}
}

func TestIsAlphaSpace(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"\n", true},
		{"\r", true},
		{"Ⅸ", false},
		{"", true},
		{"   fooo   ", true},
		{"abc!!!", false},
		{"abc1", false},
		{"abc〩", false},
		{"abc", true},
		{"소주", true},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", true},
		{"소", true},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", true},
		{"\u0026", false},
		{"\u0030", false},
		{"123", false},
		{"0123", false},
		{"-00123", false},
		{"0", false},
		{"-0", false},
		{"123.123", false},
		{" ", true},
		{".", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", true},
		{"ix", true},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsAlphaSpace(test.param))
	}
}

func TestIsAlphanumeric(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"\n", false},
		{"\r", false},
		{"Ⅸ", false},
		{"", true},
		{"   fooo   ", false},
		{"abc!!!", false},
		{"abc123", true},
		{"ABC111", true},
		{"abc1", true},
		{"abc〩", false},
		{"abc", true},
		{"소주", false},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", true},  //UTF-8(ASCII): p
		{"\u0026", false}, //UTF-8(ASCII): &
		{"\u0030", true},  //UTF-8(ASCII): 0
		{"123", true},
		{"0123", true},
		{"-00123", false},
		{"0", true},
		{"-0", false},
		{"123.123", false},
		{" ", false},
		{".", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", true},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsAlphanumeric(test.param))
	}
}

func TestIsAlphanumericSpace(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"\n", true},
		{"\r", true},
		{"Ⅸ", false},
		{"", true},
		{"   fooo   ", true},
		{"abc!!!", false},
		{"abc123", true},
		{"ABC111", true},
		{"abc1", true},
		{"abc〩", false},
		{"abc", true},
		{"소주", false},
		{"ABC", true},
		{"FoObAr", true},
		{"소aBC", false},
		{"소", false},
		{"달기&Co.", false},
		{"〩Hours", false},
		{"\ufff0", false},
		{"\u0070", true},
		{"\u0026", false},
		{"\u0030", true},
		{"123", true},
		{"0123", true},
		{"-00123", false},
		{"0", true},
		{"-0", false},
		{"123.123", false},
		{" ", true},
		{".", false},
		{"-1¾", false},
		{"1¾", false},
		{"〥〩", false},
		{"모자", false},
		{"ix", true},
		{"۳۵۶۰", false},
		{"1--", false},
		{"1-1", false},
		{"-", false},
		{"--", false},
		{"1++", false},
		{"1+1", false},
		{"+", false},
		{"++", false},
		{"+1", false},
		{"Foo Bar", true},
		{"Foo 12 Bar", true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsAlphanumericSpace(test.param))
	}
}

func TestIsEmail(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    interface{}
		expected bool
	}{
		{"", true},
		{"foo@bar.com", true},
		{"x@x.x", true},
		{"foo@bar.com.au", true},
		{"foo+bar@bar.com", true},
		{"foo@bar.coffee", true},
		{"foo@bar.中文网", true},
		{"invalidemail@", false},
		{"invalid.com", false},
		{"@invalid.com", false},
		{"test|123@m端ller.com", true},
		{"hans@m端ller.com", true},
		{"hans.m端ller@test.com", true},
		{"NathAn.daVIeS@DomaIn.cOM", true},
		{"NATHAN.DAVIES@DOMAIN.CO.UK", true},
		{123, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsEmail(test.param))
	}
}

func TestIsURL(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    interface{}
		expected bool
	}{
		{"", true},
		{"http://foo.bar#com", true},
		{"http://foobar.com", true},
		{"https://foobar.com", true},
		{"foobar.com", true},
		{"http://foobar.coffee/", true},
		{"http://foobar.中文网/", true},
		{"http://foobar.org/", true},
		{"http://foobar.ORG", true},
		{"http://foobar.org:8080/", true},
		{"ftp://foobar.ru/", true},
		{"ftp.foo.bar", true},
		{"http://user:pass@www.foobar.com/", true},
		{"http://user:pass@www.foobar.com/path/file", true},
		{"http://127.0.0.1/", true},
		{"http://duckduckgo.com/?q=%2F", true},
		{"http://localhost:3000/", true},
		{"http://foobar.com/?foo=bar#baz=qux", true},
		{"http://foobar.com?foo=bar", true},
		{"http://www.xn--froschgrn-x9a.net/", true},
		{"http://foobar.com/a-", true},
		{"http://foobar.پاکستان/", true},
		{"http://foobar.c_o_m", false},
		{"xyz://foobar.com", false},
		{"invalid.", false},
		{".com", false},
		{"rtmp://foobar.com", false},
		{"http://www.foo_bar.com/", false},
		{"http://localhost:3000/", true},
		{"http://foobar.com#baz=qux", true},
		{"http://foobar.com/t$-_.+!*\\'(),", true},
		{"http://www.foobar.com/~foobar", true},
		{"http://www.-foobar.com/", false},
		{"http://www.foo---bar.com/", false},
		{"http://r6---snnvoxuioq6.googlevideo.com", true},
		{"mailto:someone@example.com", true},
		{"irc://irc.server.org/channel", false},
		{"irc://#channel@network", true},
		{"/abs/test/dir", false},
		{"./rel/test/dir", false},
		{"http://foo^bar.org", false},
		{"http://foo&*bar.org", false},
		{"http://foo&bar.org", false},
		{"http://foo bar.org", false},
		{"http://foo.bar.org", true},
		{"http://www.foo.bar.org", true},
		{"http://www.foo.co.uk", true},
		{"foo", false},
		{"http://.foo.com", false},
		{"http://,foo.com", false},
		{",foo.com", false},
		{"https://pbs.twimg.com/profile_images/560826135676588032/j8fWrmYY_normal.jpeg", true},
		{"http://prometheus-alertmanager.service.q:9093", true},
		{"https://www.logn-123-123.url.with.sigle.letter.d:12345/url/path/foo?bar=zzz#user", true},
		{"http://me.example.com", true},
		{"http://www.me.example.com", true},
		{"https://farm6.static.flickr.com", true},
		{"https://zh.wikipedia.org/wiki/Wikipedia:%E9%A6%96%E9%A1%B5", true},
		{"google", false},
		{"http://hyphenated-host-name.example.co.in", true},
		{"http://cant-end-with-hyphen-.example.com", false},
		{"http://-cant-start-with-hyphen.example.com", false},
		{"http://www.domain-can-have-dashes.com", true},
		{"http://m.abcd.com/test.html", true},
		{"http://m.abcd.com/a/b/c/d/test.html?args=a&b=c", true},
		{"http://[::1]:9093", true},
		{"http://[::1]:909388", false},
		{"1200::AB00:1234::2552:7777:1313", false},
		{"http://[2001:db8:a0b:12f0::1]/index.html", true},
		{"http://[1200:0000:AB00:1234:0000:2552:7777:1313]", true},
		{"http://user:pass@[::1]:9093/a/b/c/?a=v#abc", true},
		{"https://127.0.0.1/a/b/c?a=v&c=11d", true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsURL(test.param))
	}
}

func TestIsJSON(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param    string
		expected bool
	}{
		{"", false},
		{"145", true},
		{"asdf", false},
		{"123:f00", false},
		{"{\"Name\":\"Alice\",\"Body\":\"Hello\",\"Time\":1294706395881547000}", true},
		{"{}", true},
		{"{\"Key\":{\"Key\":{\"Key\":123}}}", true},
		{"[]", true},
		{"null", true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsJSON(test.param))
	}
}

func TestIsLowerThenEqual(t *testing.T) {
	t.Parallel()
	x := int(7)

	var tests = []struct {
		value    interface{}
		param    interface{}
		expected bool
	}{
		{nil, 7, true},
		{"", 7, true},
		{"abcdefg", 7, true},
		{"abcdefghij", 7, false},
		{"abcd", 5, true},
		{0, 7, true},
		{7, 7, true},
		{8, 7, false},
		{5, 7, true},
		{uint(0), 7, true},
		{uint(7), 7, true},
		{uint(8), 7, false},
		{uint(5), 7, true},
		{[]string{}, 1, true},
		{[]string{"a", "b"}, 1, false},
		{&struct{}{}, 1, true},
		{&x, 1, false},
		{1, &x, true},
		{"abc", "abcd", true},
		{2, "abc", true},
		{2, "a", false},
		{uint(4), uint(3), false},
		{2.4, 2.3, false},
		{3.3, 3.31, true},
		{[]string{"a", "b"}, []string{"a", "b"}, true},
		{1, &struct{}{}, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsLowerThanEqual(test.value, test.param))
	}
}

func TestIsGreaterThanEqual(t *testing.T) {
	t.Parallel()
	x := int(7)

	var tests = []struct {
		value    interface{}
		param    interface{}
		expected bool
	}{
		{nil, 7, true},
		{"", 7, false},
		{"abcdefg", 7, true},
		{"abcdefghij", 7, true},
		{"abcd", 5, false},
		{0, 7, false},
		{7, 7, true},
		{8, 7, true},
		{5, 7, false},
		{uint(0), 7, false},
		{uint(7), 7, true},
		{uint(8), 7, true},
		{uint(5), 7, false},
		{[]string{}, 1, false},
		{[]string{"a"}, 1, true},
		{&struct{}{}, 1, false},
		{&x, 1, true},
		{1, &x, false},
		{"abc", "abcd", false},
		{4, "abc", true},
		{2, "a", true},
		{uint(4), uint(3), true},
		{2.4, 2.3, true},
		{3.3, 3.31, false},
		{[]string{"a", "b"}, []string{"a", "b"}, true},
		{1, &struct{}{}, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsGreaterThanEqual(test.value, test.param))
	}
}

func TestIsGreaterThan(t *testing.T) {
	t.Parallel()
	x := int(7)
	y := "testing"
	z := []int{1, 2, 3}
	var tests = []struct {
		value    interface{}
		param    interface{}
		expected bool
	}{
		{nil, 7, true},
		{"", 7, false},
		{"abcdefg", 7, false},
		{"abcdefghij", 7, true},
		{"abcd", 5, false},
		{0, 7, false},
		{7, 7, false},
		{8, 7, true},
		{5, 7, false},
		{uint(0), 7, false},
		{uint(7), 7, false},
		{uint(8), 7, true},
		{uint(5), 7, false},
		{float64(5.5), 7, false},
		{float64(2.5), 1, true},
		{float32(5.5), 7, false},
		{float32(2.5), 1, true},
		{[]string{}, 1, false},
		{[]string{"a"}, 1, false},
		{[]string{"a", "b"}, 1, true},
		{1, &x, false},
		{"abc", "abcd", false},
		{4, "abc", true},
		{2, "a", true},
		{uint(4), uint(3), true},
		{2.4, 2.3, true},
		{3.3, 3.31, false},
		{[]string{"a", "b"}, []string{"a", "b"}, false},
		{1, &struct{}{}, true},
		{&y, 6, true},
		{&z, 8, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsGreaterThan(test.value, test.param))
	}
}

func TestIsLowerThan(t *testing.T) {
	t.Parallel()
	x := int(7)
	var tests = []struct {
		value    interface{}
		param    interface{}
		expected bool
	}{
		{nil, 7, true},
		{"", 7, true},
		{"abcdefg", 7, false},
		{"abcdefghij", 7, false},
		{"abcd", 5, true},
		{0, 7, true},
		{7, 7, false},
		{8, 7, false},
		{5, 7, true},
		{uint(0), 7, true},
		{uint(7), 7, false},
		{uint(8), 7, false},
		{uint(5), 7, true},
		{float64(5.5), 7, true},
		{float64(5.5), 5, false},
		{float32(5.5), 7, true},
		{float32(5.5), 5, false},
		{[]string{}, 1, true},
		{[]string{"a"}, 1, false},
		{[]string{"a"}, 2, true},
		{&x, 1, false},
		{1, &x, true},
		{"abc", "abcd", true},
		{2, "abc", true},
		{2, "a", false},
		{uint(4), uint(3), false},
		{2.4, 2.3, false},
		{3.3, 3.31, true},
		{[]string{"a", "b"}, []string{"a", "b"}, false},
		{1, &struct{}{}, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsLowerThan(test.value, test.param))
	}
}

func TestIsOnRange(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		value    interface{}
		min      interface{}
		max      interface{}
		expected bool
	}{
		{nil, 1, 3, true},
		{"", 1, 3, false},
		{"abcdefg", 5, 10, true},
		{"abcdefghij", 5, 7, false},
		{0, 1, 7, false},
		{7, 1, 7, true},
		{8, 1, 7, false},
		{5, 1, 7, true},
		{float64(7.9), float64(1.2), float64(7.8), false},
		{5, float64(1.9), 7, true},
		{uint(0), 1, 7, false},
		{uint(7), 1, 7, true},
		{uint(8), 1, 7, false},
		{uint(5), 1, 7, true},
		{float64(5.5), 1, 7, true},
		{float64(5.5), 1, 3, false},
		{[]string{}, 1, 3, false},
		{[]string{"a", "b"}, 1, 3, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsOnRange(test.value, test.min, test.max))
	}
}

func TestIsContains(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   string
		param2   string
		expected bool
	}{
		{"abacada", "", true},
		{"abacada", "ritir", false},
		{"abacada", "a", true},
		{"", "a", true},
		{"abacada", "aca", true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsContains(test.param1, test.param2))
	}
}

func TestIsMatches(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   string
		param2   string
		expected bool
	}{
		{"123456789", "[0-9]+", true},
		{"abacada", "cab$", false},
		{"", "cab$", true},
		{"111222333", "((111|222|333)+)+", true},
		{"abacaba", "((123+]", false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsMatches(test.param1, test.param2))
	}
}

func TestIsSame(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   interface{}
		param2   interface{}
		expected bool
	}{
		{"123456789", "123546789", false},
		{"abacada", "abacada", true},
		{"", "abacada", true},
		{[]string{"a", "b"}, []string{"a", "b"}, true},
		{[]string{"a", "b"}, []string{"b", "c"}, false},
		{&struct{ name string }{name: "Wow"}, &struct{ name string }{name: "Wow"}, true},
		{&struct{ name string }{name: "Wow"}, &struct{ name string }{name: "wow"}, false},
		{1, 1, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsSame(test.param1, test.param2))
	}
}

func TestIsIn(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   interface{}
		param2   []string
		expected bool
	}{
		{"", []string{"abcd", "cdba"}, true},
		{"abcd", []string{"abcd", "cdba"}, true},
		{"abcd", []string{"abcde", "cdba"}, false},
		{"abcd", []string{}, false},
		{1, []string{"1", "2"}, true},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsIn(test.param1, test.param2...))
	}
}

func TestIsNotIn(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param1   interface{}
		param2   []string
		expected bool
	}{
		{"", []string{"abcd", "cdba"}, true},
		{"abcd", []string{"abcd", "cdba"}, false},
		{"abcd", []string{"abcde", "cdba"}, true},
		{"abcd", []string{}, true},
		{1, []string{"1", "2"}, false},
	}

	for _, test := range tests {
		assert.Equal(t, test.expected, validation.IsNotIn(test.param1, test.param2...))
	}
}
