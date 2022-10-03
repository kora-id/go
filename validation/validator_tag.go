// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation

import (
	"errors"
	"fmt"
	"strings"
)

type validatorTag struct {
	Name  string
	Param string
	Fn    validatorFn
}

func fetchTag(tag string, tFn map[string]validatorFn) (vt []validatorTag, e error) {
	if tag == "-" {
		e = errors.New("tag skipped")
		return
	}

	tl := strings.Split(tag, "|")
	vt = make([]validatorTag, 0, len(tl))

	for _, i := range tl {
		t := validatorTag{}
		p := strings.SplitN(i, ":", 2)

		if t.Name = strings.Trim(p[0], " "); t.Name == "" {
			e = errors.New("tag validation cannot be empty")
			break
		}

		if len(p) > 1 {
			t.Param = strings.Trim(p[1], " ")
		}

		var found bool
		if t.Fn, found = tFn[t.Name]; !found {
			e = fmt.Errorf("cannot find any tag function with name %s", t.Name)
			break
		}

		vt = append(vt, t)
	}

	return
}
