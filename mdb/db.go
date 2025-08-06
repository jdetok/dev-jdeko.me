package mdb

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jdetok/golib/errd"
)

func CreateDBConn(connStr string) (*sql.DB, error) {
	e := errd.InitErr()

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		e.Msg = fmt.Sprintf("sql.Open() failed with connStr = %s", connStr)
		return nil, e.BuildErr(err)
	}

	if err := db.Ping(); err != nil {
		e.Msg = "db.Ping() failed with returned db connection"
		return nil, e.BuildErr(err)
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(200)
	return db, nil
}
