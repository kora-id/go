package mw

import (
	"bytes"
	"compress/gzip"
	"net/http"
	"net/http/httptest"
	"testing"

	"git.tech.kora.id/go/rest"
	"github.com/stretchr/testify/assert"
)

func TestGzip(t *testing.T) {
	e := rest.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Skip if no Accept-Encoding header
	h := Gzip()(func(c *rest.Context) error {
		c.Response().Write([]byte("test")) // For Content-Type sniffing
		return nil
	})
	h(c)

	assert := assert.New(t)

	assert.Equal("test", rec.Body.String())

	// Gzip
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(rest.HeaderAcceptEncoding, gzipScheme)
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	h(c)
	assert.Equal(gzipScheme, rec.Header().Get(rest.HeaderContentEncoding))
	assert.Contains(rec.Header().Get(rest.HeaderContentType), rest.MIMETextPlain)
	r, err := gzip.NewReader(rec.Body)
	if assert.NoError(err) {
		buf := new(bytes.Buffer)
		defer r.Close()
		buf.ReadFrom(r)
		assert.Equal("test", buf.String())
	}
}

func TestGzipNoContent(t *testing.T) {
	e := rest.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(rest.HeaderAcceptEncoding, gzipScheme)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := Gzip()(func(c *rest.Context) error {
		return c.NoContent(http.StatusNoContent)
	})
	if assert.NoError(t, h(c)) {
		assert.Empty(t, rec.Header().Get(rest.HeaderContentEncoding))
		assert.Empty(t, rec.Header().Get(rest.HeaderContentType))
		assert.Equal(t, 0, len(rec.Body.Bytes()))
	}
}

func TestGzipErrorReturned(t *testing.T) {
	e := rest.New()
	e.Use(Gzip())
	e.GET("/", func(c *rest.Context) error {
		return rest.ErrNotFound
	})
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(rest.HeaderAcceptEncoding, gzipScheme)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Empty(t, rec.Header().Get(rest.HeaderContentEncoding))
}
