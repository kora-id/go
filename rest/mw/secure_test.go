package mw

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.tech.kora.id/go/rest"
	"github.com/stretchr/testify/assert"
)

func TestSecure(t *testing.T) {
	e := rest.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	h := func(c *rest.Context) error {
		return c.String(http.StatusOK, "test")
	}

	// Default
	Secure()(h)(c)
	assert.Equal(t, "1; mode=block", rec.Header().Get(rest.HeaderXXSSProtection))
	assert.Equal(t, "nosniff", rec.Header().Get(rest.HeaderXContentTypeOptions))
	assert.Equal(t, "SAMEORIGIN", rec.Header().Get(rest.HeaderXFrameOptions))
	assert.Equal(t, "", rec.Header().Get(rest.HeaderStrictTransportSecurity))
	assert.Equal(t, "", rec.Header().Get(rest.HeaderContentSecurityPolicy))

	// Custom
	req.Header.Set(rest.HeaderXForwardedProto, "https")
	rec = httptest.NewRecorder()
	c = e.NewContext(req, rec)
	SecureWithConfig(SecureConfig{
		XSSProtection:         "",
		ContentTypeNosniff:    "",
		XFrameOptions:         "",
		HSTSMaxAge:            3600,
		ContentSecurityPolicy: "default-src 'self'",
	})(h)(c)
	assert.Equal(t, "", rec.Header().Get(rest.HeaderXXSSProtection))
	assert.Equal(t, "", rec.Header().Get(rest.HeaderXContentTypeOptions))
	assert.Equal(t, "", rec.Header().Get(rest.HeaderXFrameOptions))
	assert.Equal(t, "max-age=3600; includeSubdomains", rec.Header().Get(rest.HeaderStrictTransportSecurity))
	assert.Equal(t, "default-src 'self'", rec.Header().Get(rest.HeaderContentSecurityPolicy))
}
