// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generate

import (
	"fmt"
	"os"
	"path"
	"strings"

	"git.tech.kora.id/go/dev/core"
	"git.tech.kora.id/go/dev/generate/stubs"
	"git.tech.kora.id/go/utility"
)

func FileRequest(name string, tpl *core.StubTemplate) {
	fileHandler := path.Join(tpl.AppPath, fmt.Sprintf("request_%s.go", name))
	f, err := FileReader(fileHandler)
	if err != nil {
		os.Exit(2)
	}

	rtemplate := stubs.RequestHeader
	rtemplate += strings.Replace(stubs.RequestStruct, "{{RequestName}}", utility.ToLower(name), -1)
	WriteFile(f, rtemplate, tpl)
	core.FormatSourceCode(f.Name())
	core.Log.Info(fmt.Sprintf("%-20s => \t\t%s", "request", f.Name()))
}
