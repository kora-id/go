package stubs

var MarshalModel = `
	return json.Marshal(&struct {
		{{MarshalModelColumn}}
		*Alias
	}{
		Alias: (*Alias)(m),
	})`

var Model = `
package model

import (
	{{timePkg}}

	"git.tech.kora.id/go/orm"
)

{{modelStruct}}

func (m *{{ModelName}}) Save(fields ...string) (err error) {
	o := orm.NewOrm()
	if m.ID > 0 {
		_, err = o.Update(m, fields...)
	} else {
		m.ID, err = o.Insert(m)
	}
	return
}

func (m *{{ModelName}}) Delete() (err error) {
	_, err = orm.NewOrm().Delete(m)
	return
}

func (m *{{ModelName}}) Read(fields ...string) error {
	o := orm.NewOrm()
	return o.Read(m, fields...)
}
`
