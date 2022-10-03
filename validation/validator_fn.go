// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation

import (
	"fmt"
	"strconv"
	"strings"
)

type validatorFn func(value interface{}, param string) (valid bool, m string)

func validRequired(value interface{}, _ string) (v bool, m string) {
	if v = IsNotEmpty(value); !v {
		m = "The %s field is required"
	}
	return
}

func validNumeric(value interface{}, _ string) (v bool, m string) {
	if v = IsNumeric(value); !v {
		m = "The %s must be a number"
	}
	return
}

func validAlpha(value interface{}, _ string) (v bool, m string) {
	if v = IsAlpha(value); !v {
		m = "The %s may only contain letters"
	}
	return
}

func validAlphaNum(value interface{}, _ string) (v bool, m string) {
	if v = IsAlphanumeric(value); !v {
		m = "The %s may only contain letters and numbers"
	}
	return
}

func validAlphaNumSpace(value interface{}, _ string) (v bool, m string) {
	if v = IsAlphanumericSpace(value); !v {
		m = "The %s may only contain letters, numbers and spaces"
	}
	return
}

func validAlphaSpace(value interface{}, _ string) (v bool, m string) {
	if v = IsAlphaSpace(value); !v {
		m = "The %s may only contain letters and spaces"
	}
	return
}

func validEmail(value interface{}, _ string) (v bool, m string) {
	if v = IsEmail(value); !v {
		m = "The %s must be a valid email address"
	}
	return
}

func validLatitude(value interface{}, _ string) (v bool, m string) {
	if v = IsLatitude(value); !v {
		m = "The %s must be a valid latitude."
	}
	return
}

func validLongitude(value interface{}, _ string) (v bool, m string) {
	if v = IsLongitude(value); !v {
		m = "The %s must be a valid longitude."
	}
	return
}

func validURL(value interface{}, _ string) (v bool, m string) {
	if v = IsURL(value); !v {
		m = "The %s format is invalid"
	}
	return
}

func validJSON(value interface{}, _ string) (v bool, m string) {
	if v = IsJSON(value); !v {
		m = "The %s must be a valid JSON string"
	}
	return
}

func validLte(value interface{}, param string) (v bool, m string) {
	p := convert(param)
	if v = IsLowerThanEqual(value, p); !v {
		m = fmt.Sprintf("The %s may not be greater than %v", "%s", p)
	}
	return
}

func validGte(value interface{}, param string) (v bool, m string) {
	p := convert(param)
	if v = IsGreaterThanEqual(value, p); !v {
		m = fmt.Sprintf("The %s should be greater than %v", "%s", p)
	}
	return
}

func validLt(value interface{}, param string) (v bool, m string) {
	p := convert(param)
	if v = IsLowerThan(value, p); !v {
		m = fmt.Sprintf("The %s may not be greater than %v", "%s", p)
	}
	return
}

func validGt(value interface{}, param string) (v bool, m string) {
	p := convert(param)
	if v = IsGreaterThan(value, p); !v {
		m = fmt.Sprintf("The %s should be greater than %v", "%s", p)
	}
	return
}

func validRange(value interface{}, param string) (v bool, m string) {
	p := strings.Split(param, ",")

	if len(p) == 2 {
		min := convert(p[0])
		max := convert(p[1])

		if v = IsOnRange(value, min, max); !v {
			m = fmt.Sprintf("The %s must be between %v and %v", "%s", min, max)
		}
	}
	return
}

func validContains(value interface{}, param string) (v bool, m string) {
	if v = IsContains(value, param); !v {
		m = "The %s format is invalid"
	}
	return
}

func validMatch(value interface{}, param string) (v bool, m string) {
	if v = IsMatches(value, param); !v {
		m = "The %s format is invalid"
	}
	return
}

func validSame(value interface{}, param string) (v bool, m string) {
	if v = IsSame(value, param); !v {
		m = "The %s format is invalid"
	}
	return
}

func validIn(value interface{}, param string) (v bool, m string) {
	p := strings.Split(param, ",")
	if v = IsIn(value, p...); !v {
		m = "The selected %s is invalid"
	}
	return
}

func validNotIn(value interface{}, param string) (v bool, m string) {
	p := strings.Split(param, ",")
	if v = IsNotIn(value, p...); !v {
		m = "The selected %s is invalid"
	}
	return
}

func convert(param string) (p interface{}) {
	var errInt, errFlt error
	p, errInt = strconv.Atoi(param)
	if errInt != nil {
		p, errFlt = strconv.ParseFloat(param, 64)
		if errFlt != nil {
			p = param
		}
	}
	return p
}
