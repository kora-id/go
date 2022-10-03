// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package now

import (
	"errors"
	"time"
)

// New make new instances of now
func New(t time.Time) *Now {
	return &Now{t}
}

// NewParse  make new instances of now from parsing value
func NewParse(layout string, value string) *Now {
	if t, e := time.Parse(layout, value); e == nil {
		t = t.In(time.Local)
		return &Now{t}
	}

	return &Now{time.Now()}
}

// BeginningOfMinute get the begining of minutes
func BeginningOfMinute() time.Time {
	return New(time.Now()).BeginningOfMinute()
}

// BeginningOfHour get the begining of hour
func BeginningOfHour() time.Time {
	return New(time.Now()).BeginningOfHour()
}

// BeginningOfDay get the begining of day
func BeginningOfDay() time.Time {
	return New(time.Now()).BeginningOfDay()
}

// BeginningOfWeek get the begining of week
func BeginningOfWeek() time.Time {
	return New(time.Now()).BeginningOfWeek()
}

// BeginningOfMonth get the begining of month
func BeginningOfMonth() time.Time {
	return New(time.Now()).BeginningOfMonth()
}

// BeginningOfQuarter get the begining of quarter
func BeginningOfQuarter() time.Time {
	return New(time.Now()).BeginningOfQuarter()
}

// BeginningOfYear get the begining of year
func BeginningOfYear() time.Time {
	return New(time.Now()).BeginningOfYear()
}

// EndOfMinute get the end of minute
func EndOfMinute() time.Time {
	return New(time.Now()).EndOfMinute()
}

// EndOfHour get the end of hour
func EndOfHour() time.Time {
	return New(time.Now()).EndOfHour()
}

// EndOfDay get the end of day
func EndOfDay() time.Time {
	return New(time.Now()).EndOfDay()
}

// EndOfWeek get the end of week
func EndOfWeek() time.Time {
	return New(time.Now()).EndOfWeek()
}

// EndOfMonth get the end of month
func EndOfMonth() time.Time {
	return New(time.Now()).EndOfMonth()
}

// EndOfQuarter get the end of quarter
func EndOfQuarter() time.Time {
	return New(time.Now()).EndOfQuarter()
}

// EndOfYear get the end of year
func EndOfYear() time.Time {
	return New(time.Now()).EndOfYear()
}

// Monday get the monday
func Monday() time.Time {
	return New(time.Now()).Monday()
}

// Sunday get the sunday
func Sunday() time.Time {
	return New(time.Now()).Sunday()
}

// EndOfSunday get the end of sunday
func EndOfSunday() time.Time {
	return New(time.Now()).EndOfSunday()
}

// Parse parsing time from string
func Parse(strs ...string) (time.Time, error) {
	return New(time.Now()).Parse(strs...)
}

// MustParse required parsing time
func MustParse(strs ...string) time.Time {
	return New(time.Now()).MustParse(strs...)
}

// Between check the time is in between param date
func Between(time1, time2 string) bool {
	return New(time.Now()).Between(time1, time2)
}

// NewDateRange creating slice each day between two date
func NewDateRange(start time.Time, end time.Time) *DateRange {
	dr := new(DateRange)
	for d := start; d.Sub(end) <= (0 * time.Second); d = d.AddDate(0, 0, 1) {
		r := &EachDay{d, d.Day(), d.Month(), d.Year(), nil}
		dr.Data = append(dr.Data, r)
	}

	return dr
}

// NewTimeRange setup range date from string
func NewTimeRange(start, end string) (t *TimeRange, e error) {
	t = new(TimeRange)

	if t.Start, e = parseWithFormat(start); e == nil {
		if t.End, e = parseWithFormat(end); e == nil {
			t.End = t.End.Add(24*time.Hour - 1*time.Millisecond)

			if t.End.Before(t.Start) {
				e = errors.New("invalid date range, end date cannot be lower then start date")
			}
		}
	}

	return
}
