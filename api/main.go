package main

import (
	"fmt"
	"log"
	"time"

	"github.com/jdetok/dev-jdeko.me/api/resp"
	"github.com/jdetok/dev-jdeko.me/api/store"
	"github.com/jdetok/dev-jdeko.me/pgdb"
	"github.com/jdetok/golib/envd"
	"github.com/jdetok/golib/errd"
)

func main() {
	// load environment variabels
	e := errd.InitErr()

	// err := godotenv.Load()
	err := envd.LoadDotEnv()
	if err != nil {
		fmt.Println(e.BuildErr(err).Error())
	}

	hostaddr, err := envd.GetEnvStr("SRV_IP")
	if err != nil {
		fmt.Println(e.BuildErr(err).Error())
	}

	// CONNECT TO POSTGRES
	db, err := pgdb.PostgresConn()
	if err != nil {
		fmt.Println(e.BuildErr(err).Error())
	}

	// configs go here - 8080 for testing, will derive real vals from environment
	cfg := config{
		addr: hostaddr,
	}

	// initialize the app with the configs
	app := &application{
		config:   cfg,
		database: db,
	}
	// create array of player structs
	if app.players, err = store.GetPlayers(app.database); err != nil {
		e.Msg = "failed creating players array"
		fmt.Println(e.BuildErr(err).Error())
	}

	// create array of season structs
	if app.seasons, err = store.GetSeasons(app.database); err != nil {
		e.Msg = "failed creating seasons array"
		fmt.Println(e.BuildErr(err).Error())
	}

	// create array of season structs
	if app.teams, err = store.GetTeams(app.database); err != nil {
		e.Msg = "failed creating teams array"
		fmt.Println(e.BuildErr(err).Error())
	}

	/*
		fmt.Printf("players: %d | seasons: %d | teams: %d\n", len(app.players),
			len(app.seasons), len(app.teams))
	*/
	// checks if store needs refreshed every 30 seconds, refreshes if 60 sec since last
	go store.UpdateStructs(app.database, &app.lastUpdate,
		&app.players, &app.seasons, &app.teams,
		30*time.Second, 300*time.Second)

	var rp resp.Resp
	js, err := rp.GetPlayerDash(app.database, 2544, 22024, 0)
	// js, err := rp.GetPlayerDash(app.database, 0, 22024, 1610612741)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(js))

	mux := app.mount()
	if err := app.run(mux); err != nil {
		e.Msg = "error mounting api/http server"
		log.Fatal(e.BuildErr(err).Error())
	}

}

// var tm string
// var rs []any
/*
	var ps []store.Player
	// rows, err := db.Query("select team from lg.team order by team_id desc limit 10")
	rows, err := db.Query(pgdb.PlayersSeason.Q)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var p store.Player
		rows.Scan(&p.PlayerId, &p.Name, &p.League, &p.SeasonIdMax, &p.SeasonIdMin,
			&p.PSeasonIdMax, &p.PSeasonIdMin)
		ps = append(ps, p)
	}
	fmt.Println(ps)
*/

/*
	var ss []store.Season
	// rows, err := db.Query("select team from lg.team order by team_id desc limit 10")
	rows, err := db.Query(pgdb.RSeasons.Q)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var s store.Season
		rows.Scan(&s.SeasonId, &s.Season, &s.WSeason)
		ss = append(ss, s)
	}
	fmt.Println(ss)
*/

/*
	var ts []store.Team
	rows, err := db.Query(pgdb.Teams.Q)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var t store.Team
		rows.Scan(&t.League, &t.TeamId, &t.TeamAbbr, &t.CityTeam)
		ts = append(ts, t)
	}
	fmt.Println(ts)
*/
