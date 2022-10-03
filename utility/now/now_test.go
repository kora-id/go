// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package now

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var format = "2006-01-02 15:04:05.999999999"

func TestBeginningOf(t *testing.T) {
	n := time.Date(2013, 11, 18, 17, 51, 49, 123456789, time.UTC)

	if New(n).BeginningOfMinute().Format(format) != "2013-11-18 17:51:00" {
		t.Errorf("BeginningOfMinute")
	}

	if New(n).BeginningOfHour().Format(format) != "2013-11-18 17:00:00" {
		t.Errorf("BeginningOfHour")
	}

	if New(n).BeginningOfDay().Format(format) != "2013-11-18 00:00:00" {
		t.Errorf("BeginningOfDay")
	}

	location, _ := time.LoadLocation("Japan")
	beginningOfDay := time.Date(2015, 05, 01, 0, 0, 0, 0, location)
	if New(beginningOfDay).BeginningOfDay().Format(format) != "2015-05-01 00:00:00" {
		t.Errorf("BeginningOfDay")
	}

	if New(n).BeginningOfWeek().Format(format) != "2013-11-17 00:00:00" {
		t.Errorf("BeginningOfWeek")
	}

	FirstDayMonday = true
	if New(n).BeginningOfWeek().Format(format) != "2013-11-18 00:00:00" {
		t.Errorf("BeginningOfWeek, FirstDayMonday")
	}
	FirstDayMonday = false

	if New(n).BeginningOfMonth().Format(format) != "2013-11-01 00:00:00" {
		t.Errorf("BeginningOfMonth")
	}

	if New(n).BeginningOfQuarter().Format(format) != "2013-10-01 00:00:00" {
		t.Error("BeginningOfQuarter")
	}

	if New(n.AddDate(0, -1, 0)).BeginningOfQuarter().Format(format) != "2013-10-01 00:00:00" {
		t.Error("BeginningOfQuarter")
	}

	if New(n.AddDate(0, 1, 0)).BeginningOfQuarter().Format(format) != "2013-10-01 00:00:00" {
		t.Error("BeginningOfQuarter")
	}

	if New(n).BeginningOfYear().Format(format) != "2013-01-01 00:00:00" {
		t.Errorf("BeginningOfYear")
	}
}

func TestEndOf(t *testing.T) {
	n := time.Date(2013, 11, 18, 17, 51, 49, 123456789, time.UTC)

	if New(n).EndOfMinute().Format(format) != "2013-11-18 17:51:59.999999999" {
		t.Errorf("EndOfMinute")
	}

	if New(n).EndOfHour().Format(format) != "2013-11-18 17:59:59.999999999" {
		t.Errorf("EndOfHour")
	}

	if New(n).EndOfDay().Format(format) != "2013-11-18 23:59:59.999999999" {
		t.Errorf("EndOfDay")
	}

	FirstDayMonday = true
	if New(n).EndOfWeek().Format(format) != "2013-11-24 23:59:59.999999999" {
		t.Errorf("EndOfWeek, FirstDayMonday")
	}

	FirstDayMonday = false
	if New(n).EndOfWeek().Format(format) != "2013-11-23 23:59:59.999999999" {
		t.Errorf("EndOfWeek")
	}

	if New(n).EndOfMonth().Format(format) != "2013-11-30 23:59:59.999999999" {
		t.Errorf("EndOfMonth")
	}

	if New(n).EndOfQuarter().Format(format) != "2013-12-31 23:59:59.999999999" {
		t.Errorf("EndOfQuarter")
	}

	if New(n.AddDate(0, -1, 0)).EndOfQuarter().Format(format) != "2013-12-31 23:59:59.999999999" {
		t.Errorf("EndOfQuarter")
	}

	if New(n.AddDate(0, 1, 0)).EndOfQuarter().Format(format) != "2013-12-31 23:59:59.999999999" {
		t.Errorf("EndOfQuarter")
	}

	if New(n).EndOfYear().Format(format) != "2013-12-31 23:59:59.999999999" {
		t.Errorf("EndOfYear")
	}

	n1 := time.Date(2013, 02, 18, 17, 51, 49, 123456789, time.UTC)
	if New(n1).EndOfMonth().Format(format) != "2013-02-28 23:59:59.999999999" {
		t.Errorf("EndOfMonth for 2013/02")
	}

	n2 := time.Date(1900, 02, 18, 17, 51, 49, 123456789, time.UTC)
	if New(n2).EndOfMonth().Format(format) != "1900-02-28 23:59:59.999999999" {
		t.Errorf("EndOfMonth")
	}
}

func TestMondayAndSunday(t *testing.T) {
	n := time.Date(2013, 11, 19, 17, 51, 49, 123456789, time.UTC)
	n2 := time.Date(2013, 11, 24, 17, 51, 49, 123456789, time.UTC)

	if New(n).Monday().Format(format) != "2013-11-18 00:00:00" {
		t.Errorf("Monday")
	}

	if New(n2).Monday().Format(format) != "2013-11-18 00:00:00" {
		t.Errorf("Monday")
	}

	if New(n).Sunday().Format(format) != "2013-11-24 00:00:00" {
		t.Errorf("Sunday")
	}

	if New(n2).Sunday().Format(format) != "2013-11-24 00:00:00" {
		t.Errorf("Sunday")
	}

	if New(n).EndOfSunday().Format(format) != "2013-11-24 23:59:59.999999999" {
		t.Errorf("Sunday")
	}

	if New(n).BeginningOfWeek().Format(format) != "2013-11-17 00:00:00" {
		t.Errorf("BeginningOfWeek, FirstDayMonday")
	}

	FirstDayMonday = true
	if New(n).BeginningOfWeek().Format(format) != "2013-11-18 00:00:00" {
		t.Errorf("BeginningOfWeek, FirstDayMonday")
	}
}

func TestParse(t *testing.T) {
	n := time.Date(2013, 11, 18, 17, 51, 49, 123456789, time.UTC)

	assert.Equal(t, "2013-12-19 23:28:09", New(n).MustParse("2013-12-19 23:28:09.999999999 +0800 CST").Format(format))
	assert.Equal(t, "12/10/2002", New(n).MustParse("12/10/2002").Format("02/01/2006"))
	assert.Equal(t, "2002-10-12", New(n).MustParse("2002-10-12").Format("2006-01-02"))
}

func TestBetween(t *testing.T) {
	location, _ := time.LoadLocation(TimeLocation)
	tm := time.Date(2015, 06, 30, 17, 51, 49, 123456789, location)
	if !New(tm).Between("2015-05-12 12:20", "2015-06-30 17:51:50") {
		t.Errorf("Between on location: %s", time.Now().Location().String())
	}
}

func TestMakeRange(t *testing.T) {
	r, e := NewTimeRange("08/05/2018", "09/05/2018")

	assert.NoError(t, e)
	assert.Equal(t, "2018-05-08 00:00:00 +0700 WIB", r.Start.String())
	assert.Equal(t, "2018-05-09 23:59:59.999 +0700 WIB", r.End.String())

	_, e = NewTimeRange("10/05/2018", "09/05/2018")
	assert.Error(t, e)

	r, e = NewTimeRange("01/12/2018", "20/20/2018")
	assert.Error(t, e)
}

func Example() {
	time.Now() // 2013-11-18 17:51:49.123456789 Mon

	BeginningOfMinute() // 2013-11-18 17:51:00 Mon
	BeginningOfHour()   // 2013-11-18 17:00:00 Mon
	BeginningOfDay()    // 2013-11-18 00:00:00 Mon
	BeginningOfWeek()   // 2013-11-17 00:00:00 Sun

	FirstDayMonday = true // Set Monday as first day
	BeginningOfWeek()     // 2013-11-18 00:00:00 Mon
	BeginningOfMonth()    // 2013-11-01 00:00:00 Fri
	BeginningOfQuarter()  // 2013-10-01 00:00:00 Tue
	BeginningOfYear()     // 2013-01-01 00:00:00 Tue

	EndOfMinute() // 2013-11-18 17:51:59.999999999 Mon
	EndOfHour()   // 2013-11-18 17:59:59.999999999 Mon
	EndOfDay()    // 2013-11-18 23:59:59.999999999 Mon
	EndOfWeek()   // 2013-11-23 23:59:59.999999999 Sat

	FirstDayMonday = true // Set Monday as first day
	EndOfWeek()           // 2013-11-24 23:59:59.999999999 Sun
	EndOfMonth()          // 2013-11-30 23:59:59.999999999 Sat
	EndOfQuarter()        // 2013-12-31 23:59:59.999999999 Tue
	EndOfYear()           // 2013-12-31 23:59:59.999999999 Tue

	// Use another time
	t := time.Date(2013, 02, 18, 17, 51, 49, 123456789, time.UTC)
	New(t).EndOfMonth() // 2013-02-28 23:59:59.999999999 Thu

	Monday()      // 2013-11-18 00:00:00 Mon
	Sunday()      // 2013-11-24 00:00:00 Sun
	EndOfSunday() // 2013-11-24 23:59:59.999999999 Sun
}
