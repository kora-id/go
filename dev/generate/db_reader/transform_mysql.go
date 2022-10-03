// Copyright 2016 PT. Qasico Teknologi Indonesia. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package dbReader

import (
	"database/sql"
	"fmt"
	"git.tech.kora.id/go/dev/core"
	"git.tech.kora.id/go/utility"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

var mysqlDataType = map[string]string{
	"int":                "int", // int signed
	"integer":            "int",
	"tinyint":            "int8",
	"smallint":           "int16",
	"mediumint":          "int32",
	"bigint":             "int64",
	"int unsigned":       "uint", // int unsigned
	"integer unsigned":   "uint",
	"tinyint unsigned":   "uint8",
	"smallint unsigned":  "uint16",
	"mediumint unsigned": "uint32",
	"bigint unsigned":    "uint64",
	"bit":                "uint64",
	"bool":               "bool",   // boolean
	"enum":               "string", // enum
	"set":                "string", // set
	"varchar":            "string", // string & text
	"char":               "string",
	"tinytext":           "string",
	"mediumtext":         "string",
	"text":               "string",
	"longtext":           "string",
	"blob":               "string", // blob
	"tinyblob":           "string",
	"mediumblob":         "string",
	"longblob":           "string",
	"date":               "time.Time", // time
	"datetime":           "time.Time",
	"timestamp":          "time.Time",
	"time":               "time.Time",
	"float":              "float32", // float & decimal
	"double":             "float64",
	"decimal":            "float64",
	"binary":             "string", // binary
	"varbinary":          "string",
}

type MysqlDB struct{}

type ColumnInfo struct {
	ColumnName    string
	DataType      string
	ColumnType    string
	IsNullable    string
	ColumnDefault string
	Extra         string
}

func (*MysqlDB) GetTableNames(db *sql.DB) (tables []string) {
	rows, err := db.Query("SHOW TABLES")
	if err != nil {
		core.Log.Error("Could not show tables, Please check your connection string.")
		os.Exit(2)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			core.Log.Error("Could not show tables.")
			os.Exit(2)
		}
		if name != "migrations" {
			tables = append(tables, name)
		}
	}
	return
}

func (*MysqlDB) GetConstraints(db *sql.DB, table *Table, blackList map[string]bool) {
	rows, err := db.Query(
		`SELECT
			c.constraint_type, u.column_name, u.referenced_table_schema, u.referenced_table_name, referenced_column_name, u.ordinal_position
		FROM
			information_schema.table_constraints c
		INNER JOIN
			information_schema.key_column_usage u ON c.constraint_name = u.constraint_name
		WHERE
			c.table_schema = database() AND c.table_name = ? AND u.table_schema = database() AND u.table_name = ?`,
		table.Name, table.Name) //  u.position_in_unique_constraint,
	if err != nil {
		core.Log.Error("Could not get query INFORMATION_SCHEMA for PK/UK/FK information.")
		os.Exit(2)
	}
	for rows.Next() {
		var constraintTypeBytes, columnNameBytes, refTableSchemaBytes, refTableNameBytes, refColumnNameBytes, refOrdinalPosBytes []byte
		if err := rows.Scan(&constraintTypeBytes, &columnNameBytes, &refTableSchemaBytes, &refTableNameBytes, &refColumnNameBytes, &refOrdinalPosBytes); err != nil {
			core.Log.Error("Could not get query INFORMATION_SCHEMA for PK/UK/FK information.")
			os.Exit(2)
		}
		constraintType, columnName, refTableSchema, refTableName, refColumnName, refOrdinalPos :=
			string(constraintTypeBytes), string(columnNameBytes), string(refTableSchemaBytes),
			string(refTableNameBytes), string(refColumnNameBytes), string(refOrdinalPosBytes)
		if constraintType == "PRIMARY KEY" {
			if refOrdinalPos == "1" {
				table.Pk = columnName
			} else {
				table.Pk = ""
				// add table to blacklist so that other struct will not reference it, because we are not
				// registering blacklisted tables
				blackList[table.Name] = true
			}
		} else if constraintType == "UNIQUE" {
			table.Uk = append(table.Uk, columnName)
		} else if constraintType == "FOREIGN KEY" {
			fk := new(ForeignKey)
			fk.Name = columnName
			fk.RefSchema = refTableSchema
			fk.RefTable = refTableName
			fk.RefColumn = refColumnName
			table.Fk[columnName] = fk
		}
	}
}

func (*MysqlDB) GetGoDataType(dataType string) (goType string) {
	if v, ok := mysqlDataType[dataType]; ok {
		return v
	} else {
		core.Log.Error(fmt.Sprintf("data type (%s) not found.", dataType))
		os.Exit(2)
	}
	return goType
}

func (m *MysqlDB) GetColumns(db *sql.DB, table *Table, blackList map[string]bool) {
	cols, _ := db.Query(`SELECT column_name, data_type, column_type, is_nullable, column_default, extra
		FROM information_schema.columns WHERE table_schema = database() AND table_name = ?`,
		table.Name)
	defer cols.Close()
	for cols.Next() {
		var colNameBytes, dataTypeBytes, columnTypeBytes, isNullableBytes, columnDefaultBytes, extraBytes []byte
		if err := cols.Scan(&colNameBytes, &dataTypeBytes, &columnTypeBytes, &isNullableBytes, &columnDefaultBytes, &extraBytes); err != nil {
			core.Log.Error("Could not get query INFORMATION_SCHEMA for column information.")
			os.Exit(2)
		}

		ci := &ColumnInfo{
			ColumnName:    string(colNameBytes),
			DataType:      string(dataTypeBytes),
			ColumnType:    string(columnTypeBytes),
			IsNullable:    string(isNullableBytes),
			ColumnDefault: string(columnDefaultBytes),
			Extra:         string(extraBytes),
		}

		col := new(Column)
		col.Name = utility.ToCamelCase(ci.ColumnName)
		col.Type = m.GetGoDataType(ci.DataType)

		tag := new(OrmTag)
		tag.Column = ci.ColumnName

		if table.Pk == ci.ColumnName {
			col.Name = "ID"
			col.Type = "int64"
			if ci.Extra == "auto_increment" {
				tag.Auto = true
			} else {
				tag.Pk = true
			}
		} else {
			fkCol, isFk := table.Fk[ci.ColumnName]
			isBl := false
			if isFk {
				_, isBl = blackList[fkCol.RefTable]
			}

			if ci.IsNullable == "YES" {
				tag.Null = true
			}

			// get enum options
			if ci.DataType == "enum" {
				tag.Options = extractEnumOptions(ci.ColumnType)
			}

			// check if the current column is a foreign key
			if isFk && !isBl {
				tag.RelFk = true
				col.Name = utility.ToCamelCase(ci.ColumnName)
				col.Type = "*" + utility.ToCamelCase(fkCol.RefTable)

				if strings.HasSuffix(ci.ColumnName, "_id") {
					col.Name = utility.RightTrim(col.Name, "Id")
				}

				fkCol.Column = col
			} else {
				// if the name of column is Id, and it's not primary key
				if ci.ColumnName == "id" {
					col.Name = "Id_RENAME"
				}

				if isSQLSignedIntType(ci.DataType) {
					sign := extractIntSignness(ci.ColumnType)
					if sign == "unsigned" && ci.Extra != "auto_increment" {
						col.Type = m.GetGoDataType(ci.DataType + " " + sign)
					}
				}
				if isSQLStringType(ci.DataType) {
					tag.Size = extractColSize(ci.ColumnType)
				}
				if isSQLTemporalType(ci.DataType) {
					// need to import time package
					table.ImportTimePkg = true

					tag.Type = ci.DataType
					//check auto_now, auto_now_add
					if ci.ColumnDefault == "CURRENT_TIMESTAMP" && ci.Extra == "on update CURRENT_TIMESTAMP" {
						tag.AutoNow = true
					}
				}
				if isSQLDecimal(ci.DataType) {
					tag.Digits, tag.Decimals = extractDecimal(ci.ColumnType)
				}
				if isSQLBinaryType(ci.DataType) {
					tag.Size = extractColSize(ci.ColumnType)
				}
				if isSQLBitType(ci.DataType) {
					tag.Size = extractColSize(ci.ColumnType)
				}
			}
		}

		col.Tag = tag
		table.Columns = append(table.Columns, col)
	}
}
