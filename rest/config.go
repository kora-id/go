// Copyright 2018 Kora ID. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package rest

import "os"

type config struct {
	Name         string // Service name
	DisableHTTP2 bool   // Force disable http/2
	DevMode      bool   // Switch dev mode for production or development
	JwtSecret    []byte // Secret key for Json web token algorithm
	RestHost     string // IP Application will run, default is 0.0.0.0:8080
	MySQLHost    string // IP Database server, default is 0.0.0.0:3306
	MySQLDB      string // Database name will be used
	MySQLUser    string // Database username
	MySQLPass    string // Database password
	FileCert     string
	FilePem      string
}

// loadConfig set config value from environment variable.
// If not exists, it will have a default values.
func loadConfig() *config {
	c := new(config)

	c.Name = os.Getenv("APP_NAME")
	c.DevMode = os.Getenv("APP_MODE") == "DEV"
	c.DisableHTTP2 = os.Getenv("APP_HTTP2") == "DISABLE"
	c.JwtSecret = []byte(os.Getenv("APP_JWT_SECRET"))
	c.RestHost = os.Getenv("APP_HOST")

	c.MySQLHost = os.Getenv("MYSQL_HOST")
	c.MySQLDB = os.Getenv("MYSQL_DB")
	c.MySQLUser = os.Getenv("MYSQL_USER")
	c.MySQLPass = os.Getenv("MYSQL_PASS")

	c.FileCert = os.Getenv("FILE_CERT")
	c.FilePem = os.Getenv("FILE_PEM")

	return c
}
