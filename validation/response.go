// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package validation

import (
	"regexp"
	"strings"

	"git.tech.kora.id/go/utility"
)

// Response format when running validations
type Response struct {
	Valid          bool              // state of validation
	messages       map[string]string // compiled error messages
	FailMsg        map[string]string // failing error messages
	customMessages map[string]string // custom messages
	failureKeys    []string
}

// NewResponse create new instance responses
func NewResponse() *Response {
	return &Response{Valid: true}
}

// GetMessages is a map which contains all errors from validating a struct.
func (res *Response) GetMessages() map[string]string {
	return res.messages
}

// GetMessage returns failure message by key provided as parameter,
func (res *Response) GetMessage(k string) string {
	return res.messages[k]
}

// Failure set an failure message as key and value
func (res *Response) Failure(k string, e string) {
	if res.FailMsg == nil {
		res.FailMsg = make(map[string]string)
	}

	res.Valid = false
	res.FailMsg[k] = e

	res.compile()
}

// GetErrors get error message to be serve
func (res *Response) GetErrors() map[string]string {
	msg := make(map[string]string)
	for k, v := range res.GetMessages() {
		msg[trimMessage(k)] = v
	}

	return msg
}

// Error implement error type interfaces
func (res *Response) Error() string {
	return utility.ToJSON(res.GetMessages())
}

func (res *Response) applyCustomMessage() {
	for i := range res.FailMsg {
		if c := res.customMessages[i]; c != "" {
			res.Failure(i, c)
			continue
		}

		if IsMatches(i, "(\\.[0-9]+\\.[a-z]+\\.[a-z]*)$") {
			re := regexp.MustCompile("[^a-z.]")
			ix := re.ReplaceAllString(i, "*")
			if c := res.customMessages[ix]; c != "" {
				res.Failure(i, c)
			}
		}
	}
}

func (res *Response) compile() *Response {
	res.messages = make(map[string]string)
	for k, v := range res.FailMsg {
		if _, ok := res.messages[k]; !ok {
			res.messages[k] = v
		}
	}
	return res
}

func trimMessage(s string) string {
	if idx := strings.LastIndex(s, "."); idx != -1 {
		return s[:idx]
	}
	return s
}

// SetError is helper to manualy set error where ever its needed
func SetError(field string, value string) *Response {
	o := new(Response)
	o.Failure(field, value)

	return o.compile()
}
