// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package orm

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// CodeType define prefered type of code
type CodeType int

// Enum the CodeGenerator
const (
	_           CodeType = iota // int enum type
	CodeRoman                   // roman code SO/2017/VII/001
	CodeNumeric                 // numeric code FZ000001
)

var (
	roman = map[time.Month]string{
		time.January:   "I",
		time.February:  "II",
		time.March:     "III",
		time.April:     "IV",
		time.May:       "V",
		time.June:      "VI",
		time.July:      "VII",
		time.August:    "VIII",
		time.September: "IX",
		time.October:   "X",
		time.November:  "XI",
		time.December:  "XII",
	}
)

// CodeGenerator help to generate code format based on data from database table
func CodeGenerator(t CodeType, prefix string, tableName string, field string) (code string) {
	var lastCode string
	NewOrm().Raw("SELECT " + field + " FROM " + tableName + " where " + field + " like '" + prefix + "%%' ORDER BY id DESC LIMIT 1").QueryRow(&lastCode)

	if t == CodeRoman {
		code = generateRomanCode(prefix, lastCode)
	} else if t == CodeNumeric {
		code = generateNumericCode(prefix, lastCode)
	}

	return
}

func generateRomanCode(prefix string, l string) (code string) {
	tm := time.Now()
	year := fmt.Sprintf("%d", tm.Year())
	month := roman[tm.Month()]

	l = strings.Replace(l, prefix, "", -1)
	splitCode := strings.Split(l, "/")

	if len(splitCode) != int(4) || (splitCode[2] != month || splitCode[1] != year) {
		code = prefix + "/" + year + "/" + month + "/" + fmt.Sprintf("%03d", 1)
	} else {
		lastCode, _ := strconv.Atoi(splitCode[3])
		code = prefix + "/" + year + "/" + month + "/" + fmt.Sprintf("%03d", lastCode+1)
	}

	return
}

func generateNumericCode(prefix string, l string) (code string) {
	lc := strings.Replace(l, prefix, "", 1)
	lastCode, _ := strconv.Atoi(lc)

	return fmt.Sprintf("%s%03d", prefix, lastCode+1)
}
