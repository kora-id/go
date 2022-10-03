// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package now

import (
	"fmt"
	"sort"
	"strconv"
	"time"
)

// TimeRange data structur for time range
type TimeRange struct {
	Start time.Time
	End   time.Time
}

// EachDay data structure for slice each date
type EachDay struct {
	Date  time.Time
	Day   int
	Month time.Month
	Year  int
	Data  interface{}
}

// EachMonth data structure for slice each month
type EachMonth struct {
	Month time.Month
	Year  int
	Days  []*EachDay
}

// DateRange data structure for slice date
type DateRange struct {
	Data []*EachDay
}

// ByDate list data each date
func (dr *DateRange) ByDate() []*EachDay {
	return dr.Data
}

// ByMonth grouping data based on month date
func (dr *DateRange) ByMonth() []*EachMonth {
	var result []*EachMonth

	days := make(map[int]*EachMonth)
	for _, day := range dr.Data {
		k, _ := strconv.Atoi(fmt.Sprintf("%d%d", day.Date.Year(), day.Date.Month()))
		if em := days[k]; em != nil {
			em.Days = append(em.Days, day)
		} else {
			em := new(EachMonth)
			em.Month = day.Date.Month()
			em.Year = day.Date.Year()

			em.Days = append(em.Days, day)
			days[k] = em
		}
	}

	// To store the keys in slice in sorted order
	var keys []int
	for k := range days {
		keys = append(keys, k)
	}

	sort.Ints(keys)

	for _, k := range keys {
		result = append(result, days[k])
	}

	return result
}
