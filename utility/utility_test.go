// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package utility

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type ModelTest struct {
	Image         string `orm:"column(image);null" json:"image"`
	BarcodeType   string `orm:"column(barcode_type);null;options(qr_code,ean_13,ean_8,upc_a,upc_e)" json:"barcode_type"`
	BarcodeNumber string `orm:"column(barcode_number);size(50);null" json:"barcode_number"`
	BarcodeImage  string `orm:"column(barcode_image);null" json:"barcode_image"`
	Note          string `orm:"column(note);null" json:"note"`
	Test          string `orm:"-" json:"attributes"`
}

func TestFields(t *testing.T) {
	x := Fields(ModelTest{}, "note")
	assert.Equal(t, 4, len(x))

	i := &ModelTest{
		Image: "test",
	}

	x = Fields(i, "note")
	assert.Equal(t, 4, len(x))
}

func TestPasswordHash(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		param string
	}{
		{"randompassword"},
		{"123456"},
	}
	for _, test := range tests {
		hash, e := PasswordHasher(test.param)
		assert.NoError(t, e)

		match := PasswordHash(hash, test.param)
		assert.NoError(t, match)
	}
}

func TestFloatPrecision(t *testing.T) {
	var tests = []struct {
		param    float64
		expected float64
	}{
		{2.558, 2.56},
		{2.551, 2.55},
		{2050, 2050},
	}
	for _, test := range tests {
		assert.Equal(t, test.expected, FloatPrecision(test.param, 0.5, 2))
	}
}
