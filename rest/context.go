// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"go.uber.org/zap"
)

// Context represents the context of the current HTTP request. It holds request and
// response objects, path, path parameters, data and registered handler.
type Context struct {
	request      *http.Request
	response     *Response
	path         string
	pnames       []string
	pvalues      []string
	query        url.Values
	handler      HandlerFunc
	store        Map
	rest         *Rest
	validator    Validator
	ResponseBody *ResponseFormat
}

func (c *Context) writeContentType(value string) {
	header := c.Response().Header()
	if header.Get(HeaderContentType) == "" {
		header.Set(HeaderContentType, value)
	}
}

// Request returns `*http.Request`.
func (c *Context) Request() *http.Request {
	return c.request
}

// SetRequest sets `*http.Request`.
func (c *Context) SetRequest(r *http.Request) {
	c.request = r
}

// Response returns `*Response`.
func (c *Context) Response() *Response {
	return c.response
}

// IsTLS returns true if HTTP connection is TLS otherwise false.
func (c *Context) IsTLS() bool {
	return c.request.TLS != nil
}

// IsWebSocket returns true if HTTP connection is WebSocket otherwise false.
func (c *Context) IsWebSocket() bool {
	upgrade := c.request.Header.Get(HeaderUpgrade)
	return upgrade == "websocket" || upgrade == "Websocket"
}

// Scheme returns the HTTP protocol scheme, `http` or `https`.
func (c *Context) Scheme() string {
	// Can't use `r.Request.URL.Scheme`
	// See: https://groups.google.com/forum/#!topic/golang-nuts/pMUkBlQBDF0
	if c.IsTLS() {
		return "https"
	}
	if scheme := c.request.Header.Get(HeaderXForwardedProto); scheme != "" {
		return scheme
	}
	if scheme := c.request.Header.Get(HeaderXForwardedProtocol); scheme != "" {
		return scheme
	}
	if ssl := c.request.Header.Get(HeaderXForwardedSsl); ssl == "on" {
		return "https"
	}
	if scheme := c.request.Header.Get(HeaderXUrlScheme); scheme != "" {
		return scheme
	}
	return "http"
}

// RealIP returns the client's network address based on `X-Forwarded-For`
// or `X-Real-IP` request header.
func (c *Context) RealIP() string {
	if ip := c.request.Header.Get(HeaderXForwardedFor); ip != "" {
		return strings.Split(ip, ", ")[0]
	}
	if ip := c.request.Header.Get(HeaderXRealIP); ip != "" {
		return ip
	}
	ra, _, _ := net.SplitHostPort(c.request.RemoteAddr)
	return ra
}

// Path returns the registered path for the handler.
func (c *Context) Path() string {
	return c.path
}

// SetPath sets the registered path for the handler.
func (c *Context) SetPath(p string) {
	c.path = p
}

// Param returns path parameter by name.
func (c *Context) Param(name string) string {
	for i, n := range c.pnames {
		if i < len(c.pvalues) {
			if n == name {
				return c.pvalues[i]
			}
		}
	}
	return ""
}

// ID return id parameters from route and request
func (c *Context) ID() int64 {
	id, _ := strconv.Atoi(c.Param("id"))

	return int64(id)
}

// ParamNames returns path parameter names.
func (c *Context) ParamNames() []string {
	return c.pnames
}

// SetParamNames sets path parameter names.
func (c *Context) SetParamNames(names ...string) {
	c.pnames = names
}

// ParamValues returns path parameter values.
func (c *Context) ParamValues() []string {
	return c.pvalues[:len(c.pnames)]
}

// SetParamValues sets path parameter values.
func (c *Context) SetParamValues(values ...string) {
	c.pvalues = values
}

// QueryParam returns the query param for the provided name.
func (c *Context) QueryParam(name string) string {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query.Get(name)
}

// QueryParams returns the query parameters as `url.Values`.
func (c *Context) QueryParams() url.Values {
	if c.query == nil {
		c.query = c.request.URL.Query()
	}
	return c.query
}

// QueryString returns the URL query string.
func (c *Context) QueryString() string {
	return c.request.URL.RawQuery
}

// Cookie returns the named cookie provided in the request.
func (c *Context) Cookie(name string) (*http.Cookie, error) {
	return c.request.Cookie(name)
}

// SetCookie adds a `Set-Cookie` header in HTTP response.
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response(), cookie)
}

// Cookies returns the HTTP cookies sent with the request.
func (c *Context) Cookies() []*http.Cookie {
	return c.request.Cookies()
}

// Get retrieves data from the context.
func (c *Context) Get(key string) interface{} {
	return c.store[key]
}

// Set saves data in the context.
func (c *Context) Set(key string, val interface{}) {
	if c.store == nil {
		c.store = make(Map)
	}
	c.store[key] = val
}

// Bind binds the request body into provided type `i`. The default binder
// does it based on Content-Type header.
func (c *Context) Bind(i interface{}) error {
	return c.rest.Binder.Bind(i, c)
}

// String sends a string response with status code.
func (c *Context) String(code int, s string) (err error) {
	return c.Blob(code, MIMETextPlainCharsetUTF8, []byte(s))
}

// Serve response json data with data that already collected
// if error is not nill will returning error responses.
func (c *Context) Serve(e error) (err error) {
	c.ResponseBody.Status = HTTPResponseSuccess
	c.ResponseBody.Code = http.StatusOK
	if e != nil {
		c.ResponseBody.SetError(e)
	}

	if c.Request().Method == http.MethodHead || c.Request().Method == http.MethodOptions {
		err = c.NoContent(http.StatusNoContent)
	} else {
		err = c.JSON(c.ResponseBody.Code, c.ResponseBody)
	}

	return
}

// CSV sends a binary file as csv as attachment
func (c *Context) CSV(fn string, v []byte) (err error) {
	c.writeContentType("text/csv")
	c.response.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", fn))

	_, err = c.response.Write(v)

	return
}

// JSON sends a JSON response with status code.
func (c *Context) JSON(code int, i interface{}) (err error) {
	_, pretty := c.QueryParams()["pretty"]
	if c.rest.Config.DevMode || pretty {
		return c.JSONPretty(code, i, "  ")
	}
	b, err := json.Marshal(i)
	if err != nil {
		return
	}
	return c.JSONBlob(code, b)
}

// JSONPretty sends a pretty-print JSON with status code.
func (c *Context) JSONPretty(code int, i interface{}, indent string) (err error) {
	b, err := json.MarshalIndent(i, "", indent)
	if err != nil {
		return
	}
	return c.JSONBlob(code, b)
}

// JSONBlob sends a JSON blob response with status code.
func (c *Context) JSONBlob(code int, b []byte) (err error) {
	return c.Blob(code, MIMEApplicationJSONCharsetUTF8, b)
}

// JSONP sends a JSONP response with status code. It uses `callback` to construct
// the JSONP payload.
func (c *Context) JSONP(code int, callback string, i interface{}) (err error) {
	b, err := json.Marshal(i)
	if err != nil {
		return
	}
	return c.JSONPBlob(code, callback, b)
}

// JSONPBlob sends a JSONP blob response with status code. It uses `callback`
// to construct the JSONP payload.
func (c *Context) JSONPBlob(code int, callback string, b []byte) (err error) {
	c.writeContentType(MIMEApplicationJavaScriptCharsetUTF8)
	c.response.WriteHeader(code)
	if _, err = c.response.Write([]byte(callback + "(")); err != nil {
		return
	}
	if _, err = c.response.Write(b); err != nil {
		return
	}
	_, err = c.response.Write([]byte(");"))
	return
}

// Blob sends a blob response with status code and content type.
func (c *Context) Blob(code int, contentType string, b []byte) (err error) {
	c.writeContentType(contentType)
	c.response.WriteHeader(code)
	_, err = c.response.Write(b)
	return
}

// Stream sends a streaming response with status code and content type.
func (c *Context) Stream(code int, contentType string, r io.Reader) (err error) {
	c.writeContentType(contentType)
	c.response.WriteHeader(code)
	_, err = io.Copy(c.response, r)
	return
}

// NoContent sends a response with no body and a status code.
func (c *Context) NoContent(code int) error {
	c.response.WriteHeader(code)
	return nil
}

// Redirect redirects the request to a provided URL with status code.
func (c *Context) Redirect(code int, url string) error {
	if code < 300 || code > 308 {
		return ErrInvalidRedirectCode
	}
	c.response.Header().Set(HeaderLocation, url)
	c.response.WriteHeader(code)
	return nil
}

// Error invokes the registered HTTP error handler. Generally used by middleware.
func (c *Context) Error(err error) {
	c.rest.HTTPErrorHandler(err, c)
}

// Rest get instances of rest
func (c *Context) Rest() *Rest {
	return c.rest
}

// Handler returns the matched handler by router.
func (c *Context) Handler() HandlerFunc {
	return c.handler
}

// SetHandler sets the matched handler by router.
func (c *Context) SetHandler(h HandlerFunc) {
	c.handler = h
}

// Logger returns the `Logger` instance.
func (c *Context) Logger() *zap.Logger {
	return c.rest.Logger
}

// Reset resets the context after request completes. It must be called along
// with `Rest#AcquireContext()` and `Rest#ReleaseContext()`.
// See `Rest#ServeHTTP()`
func (c *Context) Reset(r *http.Request, w http.ResponseWriter) {
	c.request = r
	c.response.reset(w)
	c.query = nil
	c.handler = NotFoundHandler
	c.store = nil
	c.path = ""
	c.pnames = nil
	c.ResponseBody.reset()
}

// JwtUsers get a user sessions that having jwt token in
// request header and checked again the model.
func (c *Context) JwtUsers(model jwtUser) interface{} {
	if u := c.Get("user"); u != nil {
		s := u.(*jwt.Token)
		c := s.Claims.(jwt.MapClaims)
		id := int64(c["id"].(float64))

		if users, err := model.GetUser(id); err == nil {
			return users
		}
	}

	return nil
}
