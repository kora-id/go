// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"net/http"

	"git.tech.kora.id/go/validation"
)

// ResponseFormat is standart response formater of the applicatin.
type ResponseFormat struct {
	Code    int               `json:"-"`
	Status  string            `json:"status,omitempty"`
	Message interface{}       `json:"message,omitempty"`
	Data    interface{}       `json:"data,omitempty"`
	Total   int64             `json:"total,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

// SetError set an error into response formater.
func (r *ResponseFormat) SetError(err error) {
	r.Code = http.StatusBadRequest
	r.Status = HTTPResponseFailed
	r.Data = nil
	r.Total = 0

	// Check error based on type
	if he, ok := err.(*HTTPError); ok {
		// Error cause of http failure should return status as is the errors
		// using standart http code.
		r.Code = he.Code
	} else if o, ok := err.(*validation.Response); ok {
		// Error cause of validation failure should return
		// status 422 and returning all failure messages as errors.
		r.Code = http.StatusUnprocessableEntity
		r.Errors = o.GetErrors()
	}

	r.Message = http.StatusText(r.Code)
}

// reset all data in response formater
func (r *ResponseFormat) reset() {
	r.Data = nil
	r.Errors = nil
	r.Message = nil
	r.Total = 0
}
