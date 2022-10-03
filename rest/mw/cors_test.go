package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.tech.kora.id/go/rest"
	"github.com/stretchr/testify/assert"
)

func TestCORS(t *testing.T) {
	e := rest.New()

	// Wildcard origin
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := CORS()(rest.NotFoundHandler)
	h(c)
	assert.Equal(t, "*", rec.Header().Get(rest.HeaderAccessControlAllowOrigin))

	// Allow origins
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h = CORSWithConfig(CORSConfig{
		AllowOrigins: []string{"localhost"},
	})(rest.NotFoundHandler)
	req.Header.Set(rest.HeaderOrigin, "localhost")
	h(c)
	assert.Equal(t, "localhost", rec.Header().Get(rest.HeaderAccessControlAllowOrigin))

	// Preflight request
	req = httptest.NewRequest(http.MethodOptions, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	req.Header.Set(rest.HeaderOrigin, "localhost")
	req.Header.Set(rest.HeaderContentType, rest.MIMEApplicationJSON)
	cors := CORSWithConfig(CORSConfig{
		AllowOrigins:     []string{"localhost"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
	h = cors(rest.NotFoundHandler)
	h(c)
	assert.Equal(t, "localhost", rec.Header().Get(rest.HeaderAccessControlAllowOrigin))
	assert.NotEmpty(t, rec.Header().Get(rest.HeaderAccessControlAllowMethods))
	assert.Equal(t, "true", rec.Header().Get(rest.HeaderAccessControlAllowCredentials))
	assert.Equal(t, "3600", rec.Header().Get(rest.HeaderAccessControlMaxAge))

	// Preflight request with `AllowOrigins` *
	req = httptest.NewRequest(http.MethodOptions, "/", nil)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	req.Header.Set(rest.HeaderOrigin, "localhost")
	req.Header.Set(rest.HeaderContentType, rest.MIMEApplicationJSON)
	cors = CORSWithConfig(CORSConfig{
		AllowOrigins:     []string{"*"},
		AllowCredentials: true,
		MaxAge:           3600,
	})
	h = cors(rest.NotFoundHandler)
	h(c)
	assert.Equal(t, "localhost", rec.Header().Get(rest.HeaderAccessControlAllowOrigin))
	assert.NotEmpty(t, rec.Header().Get(rest.HeaderAccessControlAllowMethods))
	assert.Equal(t, "true", rec.Header().Get(rest.HeaderAccessControlAllowCredentials))
	assert.Equal(t, "3600", rec.Header().Get(rest.HeaderAccessControlMaxAge))
}
