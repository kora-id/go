// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generate

import (
	"fmt"
	"git.tech.kora.id/go/dev/core"
	"git.tech.kora.id/go/dev/generate/stubs"
	"os"
	"path"
	"strings"
)

func FileHandler(methods string, tpl *core.StubTemplate) {
	fileHandler := path.Join(tpl.AppPath, "handler.go")
	f, err := FileReader(fileHandler)
	if err != nil {
		os.Exit(2)
	}

	template := stubs.HandlerStruct
	me := fmt.Sprintf("")
	ms := strings.Split(methods, ",")
	if ContainsString(ms, "get") {
		template = template + stubs.HandlerGet
		me += `r.GET("", h.get, auth.Authorized(""))` + "\n"
	}

	if ContainsString(ms, "post") {
		template = template + stubs.HandlerPost
		me += `r.POST("", h.create, auth.Authorized(""))` + "\n"
	}

	if ContainsString(ms, "show") {
		template = template + stubs.HandlerShow
		me += `r.GET("/:id", h.show, auth.Authorized(""))` + "\n"
	}

	if ContainsString(ms, "put") {
		template = template + stubs.HandlerPut
		me += `r.PUT("/:id", h.update, auth.Authorized(""))` + "\n"
	}

	if ContainsString(ms, "delete") {
		template = template + stubs.HandlerDelete
		me += `r.DELETE("/:id", h.delete)`
	}

	template = strings.Replace(template, "{{ModulEndpoint}}", me, -1)
	WriteFile(f, template, tpl)
	core.FormatSourceCode(f.Name())
	core.Log.Info(fmt.Sprintf("%-20s => \t\t%s", "handler", f.Name()))
}
