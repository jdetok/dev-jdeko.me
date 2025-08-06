package pgdb

import (
	"database/sql"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jdetok/dev-jdeko.me/applog"
	"github.com/jdetok/golib/pgresd"
)

func PostgresConn() (*sql.DB, error) {
	e := applog.AppErr{Process: "postgres connection"}
	pg := pgresd.GetEnvPG()
	pg.MakeConnStr()
	db, err := pg.Conn()
	if err != nil {
		e.Msg = "error connecting to postgres"
		return nil, e.BuildError(err)
	}
	return db, nil
}

/*
func CreateDBConn(connStr string) (*sql.DB, error) {
	e := applog.AppErr{Process: "InitDB(): initialize database connection"}

	db, err := sql.Open("mysql", connStr)
	if err != nil {
		e.Msg = fmt.Sprintf("sql.Open() failed with connStr = %s", connStr)
		return nil, e.BuildError(err)
	}

	if err := db.Ping(); err != nil {
		e.Msg = "db.Ping() failed with returned db connection"
		return nil, e.BuildError(err)
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(200)
	return db, nil
}
*/
