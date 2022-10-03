// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"git.tech.kora.id/go/rest"
	"git.tech.kora.id/go/rest/mw"

	"github.com/stretchr/testify/assert"
)

type userJwt struct {
	ID   int    `json:"id" xml:"id" form:"id"`
	Name string `json:"name" xml:"name" form:"name"`
}

type testJwtUser struct{}

func (t *testJwtUser) GetUser(id int64) (interface{}, error) {
	u := &userJwt{ID: 1, Name: "Demo"}
	return u, nil
}

func TestJwtToken(t *testing.T) {
	token := rest.JwtToken("id", 1)
	validAuth := "Bearer " + token

	assert.NotEmpty(t, token, "JWT is not working as expexted.")

	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	req.Header.Set(rest.HeaderAuthorization, validAuth)
	res := httptest.NewRecorder()

	r := rest.New()
	ctx := r.NewContext(req, res)

	handler := func(c *rest.Context) error {
		return c.String(http.StatusOK, "test")
	}

	h := mw.JWT(r.Config.JwtSecret)(handler)
	if assert.NoError(t, h(ctx), "jwt token invalid middleware handle error") {
		x := &testJwtUser{}
		i := ctx.JwtUsers(x).(*userJwt)

		assert.Equal(t, 1, i.ID)
		assert.Equal(t, "Demo", i.Name)
	}
}
