// Copyright 2014 beego Author. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package orm

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"runtime"
	"strings"
	"time"

	"git.tech.kora.id/go/utility/log"
	"go.uber.org/zap"
)

// NewLog set io.Writer to create a Logs.
func NewLog() *zap.Logger {
	d := log.New(os.Getenv("APP_NAME"), os.Getenv("APP_MODE") == "DEV")

	return d
}

func ormCaller(skipFrames int) string {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	fn := frame.File
	idx := strings.LastIndexByte(frame.File, '/')
	if idx == -1 {
		return fmt.Sprintf("%s:%d", fn, frame.Line)
	}

	idx = strings.LastIndexByte(fn[:idx], '/')
	if idx == -1 {
		return fmt.Sprintf("%s:%d", fn, frame.Line)
	}

	return fmt.Sprintf("%s:%d", fn[idx+1:], frame.Line)
}

func logQuery(query string, t time.Time, err error, args ...interface{}) {
	sub := time.Now().Sub(t) / 1e5
	elsp := float64(int(sub)) / 10.0

	query = strings.Replace(query, "`", "", -1)
	query = fmt.Sprintf(strings.Replace(query, "?", "'%v'", -1), args...)
	var fields = []zap.Field{
		zap.String("query", query),
		zap.String("latecy", fmt.Sprintf("%1.1fms", elsp)),
		zap.String("caller", ormCaller(4)),
	}

	if err == nil {
		DebugLog.Debug(fmt.Sprintf("ORM/%s", "OK"), fields...)
	} else {
		DebugLog.Error(fmt.Sprintf("ORM/%s", "FAIL"), fields...)
	}
}

// statement query logger struct.
// if dev mode, use stmtQueryLog, or use stmtQuerier.
type stmtQueryLog struct {
	alias *alias
	query string
	stmt  stmtQuerier
}

var _ stmtQuerier = new(stmtQueryLog)

func (d *stmtQueryLog) Close() error {
	a := time.Now()
	err := d.stmt.Close()
	logQuery(d.query, a, err)
	return err
}

func (d *stmtQueryLog) Exec(args ...interface{}) (sql.Result, error) {
	a := time.Now()
	res, err := d.stmt.Exec(args...)
	logQuery(d.query, a, err, args...)
	return res, err
}

func (d *stmtQueryLog) Query(args ...interface{}) (*sql.Rows, error) {
	a := time.Now()
	res, err := d.stmt.Query(args...)
	logQuery(d.query, a, err, args...)
	return res, err
}

func (d *stmtQueryLog) QueryRow(args ...interface{}) *sql.Row {
	a := time.Now()
	res := d.stmt.QueryRow(args...)
	logQuery(d.query, a, nil, args...)
	return res
}

func newStmtQueryLog(alias *alias, stmt stmtQuerier, query string) stmtQuerier {
	d := new(stmtQueryLog)
	d.stmt = stmt
	d.alias = alias
	d.query = query
	return d
}

// database query logger struct.
// if dev mode, use dbQueryLog, or use dbQuerier.
type dbQueryLog struct {
	alias *alias
	db    dbQuerier
	tx    txer
	txe   txEnder
}

var _ dbQuerier = new(dbQueryLog)
var _ txer = new(dbQueryLog)
var _ txEnder = new(dbQueryLog)

func (d *dbQueryLog) Prepare(query string) (*sql.Stmt, error) {
	a := time.Now()
	stmt, err := d.db.Prepare(query)
	logQuery(query, a, err)
	return stmt, err
}

func (d *dbQueryLog) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	a := time.Now()
	stmt, err := d.db.PrepareContext(ctx, query)
	logQuery(query, a, err)
	return stmt, err
}

func (d *dbQueryLog) Exec(query string, args ...interface{}) (sql.Result, error) {
	a := time.Now()
	res, err := d.db.Exec(query, args...)
	logQuery(query, a, err, args...)
	return res, err
}

func (d *dbQueryLog) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	a := time.Now()
	res, err := d.db.ExecContext(ctx, query, args...)
	logQuery(query, a, err, args...)
	return res, err
}

func (d *dbQueryLog) Query(query string, args ...interface{}) (*sql.Rows, error) {
	a := time.Now()
	res, err := d.db.Query(query, args...)
	logQuery(query, a, err, args...)
	return res, err
}

func (d *dbQueryLog) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	a := time.Now()
	res, err := d.db.QueryContext(ctx, query, args...)
	logQuery(query, a, err, args...)
	return res, err
}

func (d *dbQueryLog) QueryRow(query string, args ...interface{}) *sql.Row {
	a := time.Now()
	res := d.db.QueryRow(query, args...)
	logQuery(query, a, nil, args...)
	return res
}

func (d *dbQueryLog) QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row {
	a := time.Now()
	res := d.db.QueryRowContext(ctx, query, args...)
	logQuery(query, a, nil, args...)
	return res
}

func (d *dbQueryLog) Begin() (*sql.Tx, error) {
	a := time.Now()
	tx, err := d.db.(txer).Begin()
	logQuery("START TRANSACTION", a, err)
	return tx, err
}

func (d *dbQueryLog) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	a := time.Now()
	tx, err := d.db.(txer).BeginTx(ctx, opts)
	logQuery("START TRANSACTION", a, err)
	return tx, err
}

func (d *dbQueryLog) Commit() error {
	a := time.Now()
	err := d.db.(txEnder).Commit()
	logQuery("COMMIT", a, err)
	return err
}

func (d *dbQueryLog) Rollback() error {
	a := time.Now()
	err := d.db.(txEnder).Rollback()
	logQuery("ROLLBACK", a, err)
	return err
}

func (d *dbQueryLog) SetDB(db dbQuerier) {
	d.db = db
}

func newDbQueryLog(alias *alias, db dbQuerier) dbQuerier {
	d := new(dbQueryLog)
	d.alias = alias
	d.db = db
	return d
}
