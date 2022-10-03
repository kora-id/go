// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utility

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
)

// ToString convert the input to a string.
func ToString(value interface{}) string {
	res := fmt.Sprintf("%v", value)
	return string(res)
}

// ToJSON convert the input to a valid JSON string
func ToJSON(value interface{}) string {
	res, err := json.Marshal(value)
	if err != nil {
		res = []byte("")
	}
	return string(res)
}

// ToFloat convert the input string to a float, or 0.0 if the input is not a float.
func ToFloat(value interface{}) float64 {
	floatType := reflect.TypeOf(float64(0))

	v := reflect.ValueOf(value)
	if v.Kind() == reflect.String {
		res, err := strconv.ParseFloat(value.(string), 64)
		if err != nil {
			res = 0.0
		}

		return res
	}

	v = reflect.Indirect(v)
	if !v.Type().ConvertibleTo(floatType) {
		return 0
	}

	return v.Convert(floatType).Float()
}

// ToInt convert the input string to an integer, or 0 if the input is not an integer.
func ToInt(value interface{}) int {
	res, err := strconv.Atoi(Trim(ToString(value), ""))
	if err != nil {
		res = 0
	}
	return res
}

// ToBoolean convert the input string to a boolean.
func ToBoolean(value interface{}) bool {
	res, err := strconv.ParseBool(ToString(value))
	if err != nil {
		res = false
	}
	return res
}
