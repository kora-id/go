// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utility

import (
	"math"
	"reflect"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// PasswordHash compares hashed password with its possible
// plaintext equivalent using bcrypt algorithm.
func PasswordHash(hashed string, plain string) error {
	h := []byte(hashed)
	p := []byte(plain)

	return bcrypt.CompareHashAndPassword(h, p)
}

// PasswordHasher returns the bcrypt hash of the password
// using DefaultCost
func PasswordHasher(p string) (h string, err error) {
	if hx, err := bcrypt.GenerateFromPassword([]byte(p), bcrypt.DefaultCost); err == nil {
		h = string(hx)
	}
	return
}

// Contains cek is slice contains a strings.
func Contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// Fields get only fields from struct model
func Fields(model interface{}, exclude ...string) (result []string) {
	sv := reflect.ValueOf(model)
	st := reflect.TypeOf(model)

	if sv.Kind() == reflect.Ptr {
		return Fields(sv.Elem().Interface(), exclude...)
	}

	nf := sv.NumField()
	for i := 0; i < nf; i++ {
		ftag := st.Field(i).Tag.Get("orm")
		if val := parseStructTag(ftag); val != "" {
			if !Contains(exclude, val) {
				result = append(result, val)
			}
		}
	}

	return
}

func parseStructTag(data string) string {
	tags := make(map[string]string)
	for _, v := range strings.Split(data, ";") {
		if v == "" {
			continue
		}
		v = strings.TrimSpace(v)
		t := strings.ToLower(v)
		if i := strings.Index(v, "("); i > 0 && strings.Index(v, ")") == len(v)-1 {
			name := t[:i]
			if name == "column" {
				v = v[i+1 : len(v)-1]
				tags[name] = v
			}
		}
	}

	return tags["column"]
}

// FloatPrecision to round float number with some intermed and precision decimal
func FloatPrecision(num float64, intermed float64, precision int) float64 {
	pow := math.Pow(10, float64(precision))
	digit := pow * num
	_, div := math.Modf(digit)

	var round float64
	if num > 0 {
		if div >= intermed {
			round = math.Ceil(digit)
		} else {
			round = math.Floor(digit)
		}
	} else {
		if div >= intermed {
			round = math.Floor(digit)
		} else {
			round = math.Ceil(digit)
		}
	}
	return round / pow
}

// Encrypt perform simple encryption and decription values.
func Encrypt(n interface{}) string {
	num := ToInt(n)
	return ToString(((0x0000FFFF & num) << 16) + ((0xFFFF0000 & num) >> 16))
}

// Decrypt return real values of encripted values.
func Decrypt(v interface{}) int64 {
	num := ToInt(Encrypt(v))

	return int64(num)
}
