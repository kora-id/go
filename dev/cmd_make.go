// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"git.tech.kora.id/go/dev/core"
	"git.tech.kora.id/go/dev/generate"
	"os"
)

var makeCommand = &core.Command{
	Name: "make",
	Info: "source code generator",
	Usage: `
dev make model [-tables=""] [-database=test] [-conn="root:@tcp(127.0.0.1:3306)"]
	generate model source code from databases.
	-tables: 	a list of table names separated by ',', default is empty, indicating all tables
	-database: 	database name
	-conn:  	the connection string used by the driver.
				default for mysql: root:@tcp(127.0.0.1:3306)

dev make request [-name=test]
	generate appcode based on an existing database
	-name: 	    request name

dev make handler [-methods=""]
	generate appcode based on an existing database
	-methods:  	[get|post|show|put|delete] a list of http method separated by ',', default is empty, indicating GET, POST, PUT, DELETE method.
`,
}

var name, database, conn, mode, tables, endpoint, methods core.DocVal

func init() {
	makeCommand.Run = actionMake
	makeCommand.Flag.Var(&database, "database", "specify the database want to use.")
	makeCommand.Flag.Var(&conn, "conn", "connection string used by the driver to connect to a database instance.")
	makeCommand.Flag.Var(&tables, "tables", "specify tables to generate datastore.")
	makeCommand.Flag.Var(&endpoint, "endpoint", "specify endpoint name, will be used as folder name as well.")
	makeCommand.Flag.Var(&methods, "methods", "specify HTTP method that will serve by the endpoint.")
	makeCommand.Flag.Var(&name, "name", "specify project name.")
}

func actionMake(cmd *core.Command, args []string) int {
	curpath, _ := os.Getwd()
	if len(args) < 1 {
		core.Log.Error("Command is missing.")
		os.Exit(2)
	}

	core.Log.Info("")
	core.Log.Info("Generating codes ...")
	core.Log.Info("--------------------------------------")

	called := args[0]
	switch called {
	case "request":
		cmd.Flag.Parse(args[1:])
		var tpl = &core.StubTemplate{
			AppPath:     curpath,
			PackageName: core.GetDirName(curpath),
		}

		if name == "" {
			core.Log.Error("name request must be specified.")
			os.Exit(2)
		}

		core.Log.Info("Making a request file ...")
		generate.FileRequest(name.String(), tpl)

	case "model":
		cmd.Flag.Parse(args[1:])
		var tpl = &core.StubTemplate{
			AppPath:     curpath,
			PackageName: core.GetDirName(curpath),
			PackagePath: core.GetPackagePath(curpath),
		}

		c := sqlConnection()

		core.Log.Info("Making a model file ...")
		generate.FileModel("mysql", c, tables.String(), tpl)

	case "handler":
		cmd.Flag.Parse(args[1:])
		var tpl = &core.StubTemplate{
			AppPath:     curpath,
			PackageName: core.GetDirName(curpath),
		}

		if methods == "" {
			methods = core.DocVal("get,post,put,delete,show")
		}

		core.Log.Info("Making a handler file ...")
		generate.FileHandler(methods.String(), tpl)

	default:
		core.Log.Error("Command is missing.")
	}

	return 0
}

func sqlConnection() string {
	if database == "" {
		core.Log.Error("database name must be specified.")
		os.Exit(2)
	}

	var c string
	if conn == "" {
		c = fmt.Sprint("root:@tcp(127.0.0.1:3306)/", database)
	} else {
		c = fmt.Sprint(conn, "/", database)
	}

	conn = core.DocVal(c)
	return conn.String()
}
