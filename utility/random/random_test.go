package random

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {
	assert.Len(t, String(32), 32)
	assert.Regexp(t, regexp.MustCompile("[0-9]+$"), String(8, Numeric))
}
