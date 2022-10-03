package mw

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"git.tech.kora.id/go/rest"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

// MemorySink implements zap.Sink by writing all messages to a buffer.
type MemorySink struct {
	*bytes.Buffer
}

func (s *MemorySink) Close() error { return nil }
func (s *MemorySink) Sync() error  { return nil }

func TestRecover(t *testing.T) {
	sink := &MemorySink{new(bytes.Buffer)}
	zap.RegisterSink("memory", func(url *url.URL) (zap.Sink, error) {
		return sink, nil
	})

	e := rest.New()
	conf := zap.NewProductionConfig()
	conf.OutputPaths = []string{"memory://"}

	e.Logger, _ = conf.Build()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := Recover()(rest.HandlerFunc(func(c *rest.Context) error {
		panic("test")
	}))
	h(c)

	assert.Equal(t, http.StatusInternalServerError, rec.Code)

	output := sink.String()
	assert.Contains(t, output, "PANIC RECOVER")
}
