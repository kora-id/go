// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dbReader

import (
	"fmt"
	"git.tech.kora.id/go/utility"
	"strings"
)

type Column struct {
	Name string
	Type string
	Tag  *OrmTag
}

func (col *Column) String() string {
	if strings.HasSuffix(col.Name, "Id") {
		col.Name = utility.RightTrim(col.Name, "Id") + "ID"
	}
	return fmt.Sprintf("%s %s %s", col.Name, col.Type, col.Tag.String())
}

type ForeignKey struct {
	Name      string
	RefSchema string
	RefTable  string
	RefColumn string
	Column    *Column
}
