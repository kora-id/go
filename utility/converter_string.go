// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utility

import (
	"regexp"
	"strings"
	"unicode"
)

// ToLower convert the value string into lowercase format.
func ToLower(value interface{}) string {
	return strings.ToLower(ToString(value))
}

// ToUpper convert the value string into uppercase format.
func ToUpper(value interface{}) string {
	return strings.ToUpper(ToString(value))
}

// ToCamelCase converts from underscore separated form to camel case form.
// Ex.: my_func => MyFunc
func ToCamelCase(value interface{}) string {
	s := ToString(value)
	return strings.Replace(strings.Title(strings.Replace(strings.ToLower(s), "_", " ", -1)), " ", "", -1)
}

// ToLowerCamelCase converts from underscore separated form to lower camel case form.
// Ex.: my_func => myFunc
func ToLowerCamelCase(value interface{}) string {
	a := []rune(ToCamelCase(value))
	if len(a) > 0 {
		a[0] = unicode.ToLower(a[0])
	}
	return string(a)
}

// ToUnderscore converts from camel case form to underscore separated form.
// Ex.: MyFunc => my_func
func ToUnderscore(value interface{}) string {
	s := ToString(value)
	var output []rune
	var segment []rune
	for _, r := range s {
		if !unicode.IsLower(r) {
			output = addSegment(output, segment)
			segment = nil
		}
		segment = append(segment, unicode.ToLower(r))
	}
	output = addSegment(output, segment)
	return string(output)
}

// LeftTrim trim characters from the left-side of the input.
// If second argument is empty, it's will be remove leading spaces.
func LeftTrim(str, chars string) string {
	pattern := "^\\s+"
	if chars != "" {
		pattern = "^[" + chars + "]+"
	}
	r, _ := regexp.Compile(pattern)
	return string(r.ReplaceAll([]byte(str), []byte("")))
}

// RightTrim trim characters from the right-side of the input.
// If second argument is empty, it's will be remove spaces.
func RightTrim(str, chars string) string {
	pattern := "\\s+$"
	if chars != "" {
		pattern = "[" + chars + "]+$"
	}
	r, _ := regexp.Compile(pattern)
	return string(r.ReplaceAll([]byte(str), []byte("")))
}

// Trim trim characters from both sides of the input.
// If second argument is empty, it's will be remove spaces.
func Trim(str, chars string) string {
	return LeftTrim(RightTrim(str, chars), chars)
}

func addSegment(inrune, segment []rune) []rune {
	if len(segment) == 0 {
		return inrune
	}
	if len(inrune) != 0 {
		inrune = append(inrune, '_')
	}
	inrune = append(inrune, segment...)
	return inrune
}
