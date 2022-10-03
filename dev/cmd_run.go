// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package main

import (
	"io/ioutil"
	"os"
	"path"
	"runtime"
	"strings"

	"git.tech.kora.id/go/dev/core"
)

var (
	runCommand = &core.Command{
		Name: "run",
		Info: "start watching any changes in directory and rebuild it.",
		Usage: `
dev run
	Run command will watching any changes in directory of go project,
	it will recompile and restart the application binary.
`,
	}

	mainFiles core.ListOpts
)

func init() {
	runCommand.Run = actionRun
	runCommand.Flag.Var(&mainFiles, "ext", "specify file extension to watch")
}

// actionRun, perform to scan directory and get ready to watching.
func actionRun(_ *core.Command, _ []string) int {
	gps := getGoPath()
	if len(gps) == 0 {
		core.Log.Error("$GOPATH not found, Please set $GOPATH in your environment variables.")
		os.Exit(2)
	}

	exit := make(chan bool)
	cwd, _ := os.Getwd()
	appName := path.Base(cwd)

	core.Log.Info("")
	core.Log.Info("Run applications ...")
	core.Log.Info("--------------------------------------")

	var paths []string
	readDirectory(cwd, &paths)

	var files []string
	for _, arg := range mainFiles {
		if len(arg) > 0 {
			files = append(files, arg)
		}
	}

	core.Watch(appName, paths, files)
	core.Build()

	for {
		select {
		case <-exit:
			runtime.Goexit()
		}
	}
}

// readDirectory binds paths with list of existing directory.
func readDirectory(directory string, paths *[]string) {
	fileInfos, err := ioutil.ReadDir(directory)
	if err != nil {
		return
	}

	useDirectory := false
	for _, fileInfo := range fileInfos {
		if strings.HasSuffix(fileInfo.Name(), "docs") {
			continue
		}

		if fileInfo.IsDir() == true && fileInfo.Name()[0] != '.' {
			readDirectory(directory+"/"+fileInfo.Name(), paths)
			continue
		}

		if useDirectory == true {
			continue
		}

		if path.Ext(fileInfo.Name()) == ".go" {
			*paths = append(*paths, directory)
			useDirectory = true
		}
	}

	return
}

// getGoPath returns list of go path on system.
func getGoPath() (p []string) {
	gopath := os.Getenv("GOPATH")
	p = strings.Split(gopath, ":")

	return
}
