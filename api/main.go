package main

import (
	"fmt"
	"time"

	"github.com/jdetok/dev-jdeko.me/api/cache"
	"github.com/jdetok/dev-jdeko.me/applog"
	"github.com/jdetok/dev-jdeko.me/getenv"
	"github.com/jdetok/dev-jdeko.me/pgdb"
)

func main() {
	// load environment variabels
	e := applog.AppErr{Process: "Main function"}

	// err := godotenv.Load()
	err := getenv.LoadDotEnv()
	if err != nil {
		fmt.Println(e.BuildError(err).Error())
	}

	hostaddr, err := getenv.GetEnvStr("SRV_IP")
	if err != nil {
		fmt.Println(e.BuildError(err).Error())
	}

	db, err := pgdb.PostgresConn()
	if err != nil {
		fmt.Println(e.BuildError(err).Error())
	}

	// configs go here - 8080 for testing, will derive real vals from environment
	cfg := config{
		addr: hostaddr,
		// cachePath: cacheP,
	}

	// initialize the app with the configs
	app := &application{
		config:   cfg,
		database: db,
	}
	// create array of player structs
	if app.players, err = cache.GetPlayers(app.database); err != nil {
		e.Msg = "failed creating players array"
		fmt.Println(e.BuildError(err).Error())
	}

	// create array of season structs
	if app.seasons, err = cache.GetSeasons(app.database); err != nil {
		e.Msg = "failed creating seasons array"
		fmt.Println(e.BuildError(err).Error())
	}

	// create array of season structs
	if app.teams, err = cache.GetTeams(app.database); err != nil {
		e.Msg = "failed creating teams array"
		fmt.Println(e.BuildError(err).Error())
	}

	fmt.Printf("players: %d | seasons: %d | teams: %d\n", len(app.players),
		len(app.seasons), len(app.teams))

	// checks if cache needs refreshed every 30 seconds, refreshes if 60 sec since last
	go cache.UpdateStructs(app.database, &app.lastUpdate,
		&app.players, &app.seasons, &app.teams,
		30*time.Second, 300*time.Second)

	/*
		mux := app.mount()
		if err := app.run(mux); err != nil {
			e.Msg = "error mounting api/http server"
			log.Fatal(e.BuildError(err).Error())
		}
	*/
}

// var tm string
// var rs []any
/*
	var ps []cache.Player
	// rows, err := db.Query("select team from lg.team order by team_id desc limit 10")
	rows, err := db.Query(pgdb.PlayersSeason.Q)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var p cache.Player
		rows.Scan(&p.PlayerId, &p.Name, &p.League, &p.SeasonIdMax, &p.SeasonIdMin,
			&p.PSeasonIdMax, &p.PSeasonIdMin)
		ps = append(ps, p)
	}
	fmt.Println(ps)
*/

/*
	var ss []cache.Season
	// rows, err := db.Query("select team from lg.team order by team_id desc limit 10")
	rows, err := db.Query(pgdb.RSeasons.Q)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var s cache.Season
		rows.Scan(&s.SeasonId, &s.Season, &s.WSeason)
		ss = append(ss, s)
	}
	fmt.Println(ss)
*/

/*
	var ts []cache.Team
	rows, err := db.Query(pgdb.Teams.Q)
	if err != nil {
		fmt.Println(err)
	}

	for rows.Next() {
		var t cache.Team
		rows.Scan(&t.League, &t.TeamId, &t.TeamAbbr, &t.CityTeam)
		ts = append(ts, t)
	}
	fmt.Println(ts)
*/
