package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.tech.kora.id/go/rest"
	"github.com/stretchr/testify/assert"
)

func TestRequestID(t *testing.T) {
	e := rest.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	handler := func(c *rest.Context) error {
		return c.String(http.StatusOK, "test")
	}

	rid := RequestIDWithConfig(RequestIDConfig{})
	h := rid(handler)
	h(c)
	assert.Len(t, rec.Header().Get(rest.HeaderXRequestID), 32)

	// Custom generator
	rid = RequestIDWithConfig(RequestIDConfig{
		Generator: func() string { return "customGenerator" },
	})
	h = rid(handler)
	h(c)
	assert.Equal(t, rec.Header().Get(rest.HeaderXRequestID), "customGenerator")
}
