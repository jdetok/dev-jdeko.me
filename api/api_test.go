package main

import (
	"database/sql"
	"fmt"
	"testing"

	"github.com/jdetok/dev-jdeko.me/api/resp"
	"github.com/jdetok/dev-jdeko.me/pgdb"
	"github.com/jdetok/golib/envd"
	"github.com/jdetok/golib/errd"
	"github.com/jdetok/golib/logd"
)

func StartupTest(t *testing.T) *sql.DB {
	e := errd.InitErr()
	err := envd.LoadDotEnvFile("../.env")
	if err != nil {
		e.Msg = "failed loading .env file"
		t.Fatal(e.BuildErr(err))
	}

	db, err := pgdb.PostgresConn()
	if err != nil {
		e.Msg = "failed connecting to postgres"
		t.Fatal(e.BuildErr(err))
	}
	return db
}

func TestGetPlayerDash(t *testing.T) {
	e := errd.InitErr()
	db := StartupTest(t)
	logd.Logc("testing GetPlayerDash with player query")
	var pIds = []uint64{2544, 2544}    // lebron
	var sIds = []uint64{22024, 22024}  // 2425 reg season
	var tIds = []uint64{0, 1610612743} // first should be plr query, second tm

	for i := range pIds {
		var rp resp.Resp
		plr := pIds[i]
		szn := sIds[i]
		tm := tIds[i]
		msg := fmt.Sprintf("pId: %d | sId: %d | tId: %d", plr, szn, tm)
		js, err := rp.GetPlayerDash(db, plr, szn, tm)
		if err != nil {
			e.Msg = fmt.Sprintf("failed getting player dash\n%s", msg)
			t.Error(e.BuildErr(err))
		}
		fmt.Println(string(js))
	}
}
