// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import (
	"fmt"
	"time"

	"git.tech.kora.id/go/orm"

	// mysql database connection
	_ "github.com/go-sql-driver/mysql"
)

// DatabaseSetup setting up database with config
func DatabaseSetup() error {
	orm.Debug = true
	orm.DefaultTimeLoc = time.Local
	orm.DefaultRelsDepth = 3
	orm.DebugLog = Logger

	ds := fmt.Sprintf("%s:%s@tcp(%s)/%s?%s", Config.MySQLUser, Config.MySQLPass, Config.MySQLHost, Config.MySQLDB, "charset=utf8&loc=Asia%2FJakarta")
	return orm.RegisterDataBase("default", "mysql", ds)
}
