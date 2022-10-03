// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"sort"
	"time"

	"github.com/dgrijalva/jwt-go"
)

// jwtUser model user jwt token interface
// to check is the id given valid as users.
type jwtUser interface {
	GetUser(int64) (interface{}, error)
}

// JwtKey byte of jwt secret keys.
func JwtKey() []byte {
	return []byte(Config.JwtSecret)
}

// JwtToken make an JWT token keys and values
// the return will become a valid token with
// a life time 72 hours from the time generated.
func JwtToken(k string, v interface{}, neverExpire ...bool) (token string) {
	// new instances jwt
	jwts := jwt.New(jwt.SigningMethodHS256)

	// Set claims
	claims := jwts.Claims.(jwt.MapClaims)
	claims[k] = v

	if len(neverExpire) > 0 && neverExpire[0] {
		claims["exp"] = time.Now().Add(time.Hour * 8766).Unix()
	} else {
		claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	}

	// Generate encoded token
	var e error
	if token, e = jwts.SignedString(JwtKey()); e != nil {
		panic(e)
	}

	return token
}

// DebugRoutes print all route available, only show on debug mode.
func DebugRoutes(e *Rest) {
	if e.Config.DevMode {
		e.StdLogger.Printf("%0120v", "")
		e.StdLogger.Printf("%-10s | %-50s | %-54s", "METHOD", "URL PATH", "REQ. HANDLER")
		e.StdLogger.Printf("%0120v", "")

		routes := e.Routes()
		sort.Sort(sortByPath(routes))
		for _, v := range routes {
			if v.Path[len(v.Path)-1:] != "*" {
				e.StdLogger.Printf("%-10s | %-50s | %-54s", v.Method, v.Path, v.Name)
			}
		}
		e.StdLogger.Printf("%0120v", "")
	}
}

// sortByPath Sorting echo.Routes by path
// so it make more pretty when printed on console.
type sortByPath []*Route

func (a sortByPath) Len() int {
	return len(a)
}

func (a sortByPath) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a sortByPath) Less(i, j int) bool {
	if a[i].Path < a[j].Path {
		return true
	}
	if a[i].Path > a[j].Path {
		return false
	}
	return a[i].Path < a[j].Path
}
