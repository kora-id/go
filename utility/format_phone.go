// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utility

import (
	"regexp"
	"strings"
)

// FormatPhone formating string to phone number
func FormatPhone(text string) (result string) {
	if len(text) < 10 {
		return ""
	}

	// remove non numeric string
	reg, _ := regexp.Compile("[^0-9]+")
	result = reg.ReplaceAllString(text, "")

	if result == "" || len(result) < 10 {
		return ""
	}

	prefix := string(result[0:2])
	if prefix == "08" {
		result = strings.Replace(result, "08", "628", 1)
	} else {
		prefix2 := string(result[0:1])
		if prefix2 == "8" {
			result = strings.Replace(result, "8", "628", 1)
		}
	}

	fp := string(result[0:2])
	if fp != "62" {
		result = ""
	}

	return
}
