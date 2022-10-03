// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dbReader

import (
	"fmt"
	"git.tech.kora.id/go/utility"
	"strings"
)

type Table struct {
	Name          string
	Pk            string
	Uk            []string
	Fk            map[string]*ForeignKey
	Columns       []*Column
	ImportTimePkg bool
}

func (tb *Table) String() string {
	rv := fmt.Sprintf("type %s struct {\n", utility.ToCamelCase(tb.Name))
	for _, v := range tb.Columns {
		rv += v.String() + "\n"
	}
	rv += "}\n"
	return rv
}

func (tb *Table) MarshalColumn() string {
	var colMarshal []string

	colMarshal = append(colMarshal, fmt.Sprintf("%s %s %s", "ID", "string", "`json:\"id\"`"))
	for col, fk := range tb.Fk {
		cname := utility.ToCamelCase(fk.Name)
		if strings.HasSuffix(cname, "Id") {
			cname = utility.RightTrim(cname, "Id") + "ID"
		} else {
			cname = cname + "ID"
			col = col + "_id"
		}
		colMarshal = append(colMarshal, fmt.Sprintf("%s %s %s", cname, "string", fmt.Sprintf("`json:\"%s\"`", col)))
	}
	return strings.Join(colMarshal, "\n")
}
