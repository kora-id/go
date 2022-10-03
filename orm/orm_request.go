// Copyright 2019 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package orm

import (
	"fmt"
	"net/url"
	"strings"

	"git.tech.kora.id/go/utility"
)

// RequestQuery struct readed from context
type RequestQuery struct {
	Conditions []map[string]string
	Fields     []string
	OrderBy    []string
	Embeds     []string
	Offset     int
	Limit      int
}

// Query make new query setter based on request query.
func (rq *RequestQuery) Query(model interface{}) (QuerySeter, Ormer) {
	o := NewOrm()

	return rq.Apply(o.QueryTable(model)), o
}

// ExcludeEmbeds will exclude RequestQuery Embeds in parameter
// example: bool:=rq.ExcludeEmbeds("table_name field")
func (rq *RequestQuery) ExcludeEmbeds(customEmbeds string) bool {
	flag := false
	for index, queryEmbeds := range rq.Embeds {
		if queryEmbeds == customEmbeds {
			rq.Embeds = append(rq.Embeds[:index], rq.Embeds[index+1:]...)
			flag = true
		}
	}
	return flag
}

// Apply set data request query into query setter.
func (rq *RequestQuery) Apply(qs QuerySeter) QuerySeter {
	// apply conditions
	qs = qs.SetCond(rq.GetCondition())

	// apply embeds
	if len(rq.Embeds) > 0 {
		j := rq.GetJoin()
		qs = qs.RelatedSel(j...)
	}

	// apply order by
	qs = qs.OrderBy(rq.OrderBy...)

	// apply limit
	qs = qs.Limit(rq.Limit, rq.Offset)

	return qs
}

// ReadFromContext reading from params url
func (rq *RequestQuery) ReadFromContext(params url.Values) *RequestQuery {
	if pl := utility.ToInt(params.Get("limit")); pl != 0 {
		rq.Limit = pl
	}

	if pp := utility.ToInt(params.Get("page")); pp != 0 {
		rq.Offset = rq.Limit * (pp - 1)
	}

	if pf := params.Get("fields"); pf != "" {
		rq.Fields = strings.Split(pf, ",")
	}

	if po := params.Get("orderby"); po != "" {
		k := strings.Replace(po, ".", "__", -1)
		rq.OrderBy = strings.Split(k, ",")
	}

	if pj := params.Get("embeds"); pj != "" {
		k := strings.Replace(pj, ".", "__", -1)
		rq.Embeds = strings.Split(k, ",")
	}

	if pc := params.Get("conditions"); pc != "" {
		for _, cond := range strings.Split(pc, "|") {
			var bc = make(map[string]string)
			for _, partcond := range strings.Split(cond, "%2C") {
				kv := strings.Split(partcond, ":")
				if len(kv) > 2 {
					bc[kv[0]] = fmt.Sprintf("%s:%s:%s", kv[1], kv[2], kv[3])
				} else if len(kv) == 2 {
					bc[kv[0]] = kv[1]
				} else {
					bc[partcond] = "true"
				}
			}
			rq.Conditions = append(rq.Conditions, bc)
		}
	}

	return rq
}

// GetCondition making condition orm from request
func (rq *RequestQuery) GetCondition() *Condition {
	c := NewCondition()
	for _, q := range rq.Conditions {
		cd := NewCondition()
		for k, v := range q {
			if strings.Contains(k, "AndNot.") {
				k = strings.Replace(k, "AndNot.", "", -1)
				k = strings.Replace(k, ".", "__", -1)

				cd = rq.condition(cd, k, v, "andnot")
			} else if strings.Contains(k, "Or.") {
				k = strings.Replace(k, "Or.", "", -1)
				k = strings.Replace(k, ".", "__", -1)
				cd = rq.condition(cd, k, v, "or")
			} else if strings.Contains(k, "OrNot.") {
				k = strings.Replace(k, "OrNot.", "", -1)
				k = strings.Replace(k, ".", "__", -1)

				cd = rq.condition(cd, k, v, "ornot")
			} else {
				k = strings.Replace(k, "And.", "", -1)
				k = strings.Replace(k, ".", "__", -1)

				cd = rq.condition(cd, k, v, "and")
			}
		}
		c = c.AndCond(cd)
	}

	return c
}

// GetJoin making join orm from request
func (rq *RequestQuery) GetJoin() []interface{} {
	new := make([]interface{}, len(rq.Embeds))
	for i, v := range rq.Embeds {
		new[i] = v
	}

	return new
}

func (rq *RequestQuery) condition(c *Condition, field string, value string, operator string) (res *Condition) {
	if strings.Contains(field, "__in") {
		v := strings.Split(value, ".")
		switch operator {
		case "or":
			res = c.Or(field, v)
		case "ornot":
			res = c.OrNot(field, v)
		case "andnot":
			res = c.AndNot(field, v)
		default:
			res = c.And(field, v)
		}
	} else if strings.Contains(field, "__between") {
		v := strings.Split(value, ".")
		switch operator {
		case "or":
			res = c.Or(field, v)
		case "ornot":
			res = c.OrNot(field, v)
		case "andnot":
			res = c.AndNot(field, v)
		default:
			res = c.And(field, v)
		}
	} else if strings.Contains(field, "__null") {
		field = strings.Replace(field, "__null", "__isnull", -1)
		switch operator {
		case "or":
			res = c.Or(field, true)
		case "ornot":
			res = c.OrNot(field, true)
		case "andnot":
			res = c.AndNot(field, true)
		default:
			res = c.And(field, true)
		}
	} else if strings.Contains(field, "__notnull") {
		field = strings.Replace(field, "__notnull", "__isnull", -1)
		switch operator {
		case "or":
			res = c.Or(field, false)
		case "ornot":
			res = c.OrNot(field, false)
		case "andnot":
			res = c.AndNot(field, false)
		default:
			res = c.And(field, false)
		}
	} else {
		switch operator {
		case "or":
			res = c.Or(field, value)
		case "ornot":
			res = c.OrNot(field, value)
		case "andnot":
			res = c.AndNot(field, value)
		default:
			res = c.And(field, value)
		}
	}
	return
}
