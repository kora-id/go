// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package generate

import (
	"fmt"
	"os"
	"path"
	"strings"

	"database/sql"

	"git.tech.kora.id/go/dev/core"
	dbReader "git.tech.kora.id/go/dev/generate/db_reader"
	"git.tech.kora.id/go/dev/generate/stubs"
	"git.tech.kora.id/go/utility"
)

func FileModel(driver string, conn string, selectedTables string, tpl *core.StubTemplate) {

	var tables map[string]bool
	if selectedTables != "" {
		tables = make(map[string]bool)
		for _, v := range strings.Split(selectedTables, ",") {
			tables[v] = true
		}
	}

	db, err := sql.Open(driver, conn)
	if err != nil {
		core.Log.Error("Could not connect to database ")
		core.Log.Error(fmt.Sprintf("using: %s, %s, %s", driver, conn, err.Error()))
		os.Exit(2)
	}
	defer db.Close()

	if trans, ok := dbReader.DBDriver[driver]; ok {
		core.Log.Info("")
		core.Log.Info("Generating model file ...")
		core.Log.Info("--------------------------------------")

		tableNames := trans.GetTableNames(db)
		dbTables := dbReader.GetTableObjects(tableNames, db, trans)

		makeModels(dbTables, tables, tpl.AppPath, tpl)
	} else {
		core.Log.Error(fmt.Sprintf("%s database is not supported yet.", driver))
		os.Exit(2)
	}
}

func makeModels(tables []*dbReader.Table, selectedTables map[string]bool, modelPath string, tpl *core.StubTemplate) {
	for _, tb := range tables {
		if selectedTables != nil {
			if _, selected := selectedTables[tb.Name]; !selected {
				continue
			}
		}

		filename := dbReader.GetFileName(tb.Name)
		file := path.Join(modelPath, filename+".go")
		f, err := FileReader(file)
		if err != nil {
			continue
		}

		template := stubs.Model
		template = strings.Replace(template, "{{modelStruct}}", tb.String(), 1)

		var tPkg string
		if tb.ImportTimePkg {
			tPkg = "\"time\"\n"
		}

		template = strings.Replace(template, "{{timePkg}}", tPkg, -1)

		tpl.ModelName = utility.ToCamelCase(tb.Name)
		tpl.ModelNameSingular = strings.TrimSuffix(tpl.ModelName, "s")
		tpl.TableName = tb.Name
		WriteFile(f, template, tpl)
		core.FormatSourceCode(f.Name())
		core.Log.Info(fmt.Sprintf("%-20s => \t\t%s", "model", file))
	}
}
