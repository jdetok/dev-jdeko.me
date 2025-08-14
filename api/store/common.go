package store

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/jdetok/dev-jdeko.me/pgdb"
	"github.com/jdetok/golib/errd"
	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Player struct {
	PlayerId     uint64
	Name         string
	League       string
	SeasonIdMax  uint64
	SeasonIdMin  uint64
	PSeasonIdMax uint64
	PSeasonIdMin uint64
}

type Season struct {
	SeasonId string `json:"season_id"`
	Season   string `json:"season"`
	WSeason  string `json:"wseason"`
}

type Team struct {
	League   string `json:"league"`
	TeamId   string `json:"team_id"`
	TeamAbbr string `json:"team"`
	CityTeam string `json:"team_long"`
	LogoUrl  string `json:"-"`
}

type SeasonLeague struct {
	Szn  string
	WSzn string
}

/*
returns slice of season strings for date (generally pass time.Now())
calling in 2025 will return 2024-25 and 2025-26 and so on
*/
func CurrentSzns(dt time.Time) []string {
	var cyyy string = dt.Format("2006")
	var cy string = dt.AddDate(1, 0, 0).Format("06")

	var pyyy string = dt.AddDate(-1, 0, 0).Format("2006")
	var py string = dt.Format("06")

	return []string{
		fmt.Sprint(pyyy, "-", py),
		fmt.Sprint(cyyy, "-", cy),
	}
}

func LgSeasons() SeasonLeague {
	e := errd.InitErr()
	var sl SeasonLeague
	var crnt []string = CurrentSzns(time.Now())

	m, err := strconv.Atoi(time.Now().Format("1"))
	if err != nil {
		e.Msg = "error converting year to int"
		fmt.Println(e.BuildErr(err))
	}

	// beginning of year through april
	sl.Szn = crnt[0]
	sl.WSzn = crnt[0]

	// may through september
	if m > 5 && m < 10 {
		sl.WSzn = crnt[1]
	}

	// october through end of year
	if m > 10 {
		sl.Szn = crnt[1]
		sl.WSzn = crnt[1]
	}

	// fmt.Printf("NBA Season: %s | WNBA Season: %s\n", sl.Szn, sl.WSzn)
	return sl
}

// REMOVE NON SPACING CHARACTERS -- e.g. Dončić becomes doncic
func Unaccent(input string) string {
	t := transform.Chain(
		norm.NFD,
		runes.Remove(runes.In(unicode.Mn)),
		norm.NFC,
	)
	output, _, _ := transform.String(t, input)
	return output
}

// makes src url for team img
func (t Team) MakeLogoUrl() string {
	lg := strings.ToLower(t.League)
	return ("https://cdn." + lg + ".com/logos/" +
		lg + "/" + t.TeamId + "/primary/L/logo.svg")
}

// QUERY FOR PLAYER ID, PLAYER AND SAVE TO A LIST OF PLAYER STRUCTS
func GetPlayers(db *sql.DB) ([]Player, error) {
	e := errd.InitErr()
	// rows, err := db.Query(mdb.PlayersSeason.Q)
	rows, err := db.Query(pgdb.PlayersSeason.Q)
	if err != nil {
		e.Msg = "query failed"
		return nil, e.BuildErr(err)
	}
	var players []Player
	for rows.Next() {
		var p Player
		rows.Scan(&p.PlayerId, &p.Name, &p.League, &p.SeasonIdMax, &p.SeasonIdMin, &p.PSeasonIdMax, &p.PSeasonIdMin)
		p.Name = Unaccent(p.Name) // REMOVE ACCENTS FROM NAMES
		players = append(players, p)
	}
	return players, nil
}

// seasons
func GetSeasons(db *sql.DB) ([]Season, error) {
	// fmt.Println("querying seasons & saving to struct")
	e := errd.InitErr()
	// rows, err := db.Query(mdb.RSeasons.Q)
	rows, err := db.Query(pgdb.AllSeasons.Q)
	if err != nil {
		e.Msg = "error querying db"
		e.BuildErr(err)
	}

	var seasons []Season
	for rows.Next() {
		var szn Season
		rows.Scan(&szn.SeasonId, &szn.Season, &szn.WSeason)
		seasons = append(seasons, szn)
	}

	return seasons, nil
}

// teams
func GetTeams(db *sql.DB) ([]Team, error) {
	e := errd.InitErr()
	// rows, err := db.Query(mdb.Teams.Q)
	rows, err := db.Query(pgdb.Teams.Q)
	if err != nil {
		e.Msg = "error querying db"
		e.BuildErr(err)
	}

	var teams []Team
	for rows.Next() {
		var tm Team
		rows.Scan(&tm.League, &tm.TeamId, &tm.TeamAbbr, &tm.CityTeam)
		teams = append(teams, tm)
	}
	return teams, nil
}
