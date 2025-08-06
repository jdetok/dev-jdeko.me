package pgdb

type Query struct {
	Args []string // arguments to accept
	Q    string   // query
}

// GetPlayerDash
var Player = Query{
	Q: `
		select a.*, b.season_desc, b.wseason_desc
		from api_player_stats a
		join season b on b.season_id = a.season_id
		where a.player_id = ? and a.season_id = ?
	`,
}

// GetPlayerDash
var TeamSeasonTopP = Query{
	Q: `
		select a.*, b.season_desc, b.wseason_desc
		from api_player_stats a
		join season b on b.season_id = a.season_id
		where team_id = ? and a.season_id = ?
		order by a.points desc
		limit 2;
	`,
}

var RecentGamePlayers = Query{
	Q: `
	select a.game_id, a.team_id, e.player_id, f.player, b.lg, c.team, c.team_name,
	b.game_date, b.matchup, b.final, b.ot, a.pts, e.pts
	from t_box a
	inner join game b on b.game_id = a.game_id
	inner join team c on c.team_id = a.team_id
	inner join p_box d on d.game_id = a.game_id
		and d.team_id = a.team_id
	inner join (
		select game_id, player_id, team_id, pts
		from p_box
		group by game_id, team_id, player_id
		order by pts desc
		limit 1
	) e on e.game_id = a.game_id and e.team_id = a.team_id
	inner join player f on f.player_id = e.player_id and f.team_id = a.team_id
	where b.game_date = (
		select max(game_date) from game 
		where left(season_id, 1) in ('2', '4')
		and lg in ('NBA', 'WNBA')
	)
	and left(a.season_id, 1) in ('2', '4')
	and b.lg in ('NBA', 'WNBA')
	group by a.game_id, a.team_id
	order by e.pts desc
	`,
}

var PlayersSeason = Query{
	Q: `
	select 
		a.player_id,
		lower(a.player) as plr,
		case 
			when a.lg_id = 0 then 'nba'
			when a.lg_id = 1 then 'wnba'
		end,
		b.rs_max, 
		b.rs_min,
		coalesce(c.po_max, b.rs_max),
		coalesce(c.po_min, b.rs_min)
	from lg.plr a
	inner join (
		select player_id, min(season_id) as rs_min, max(season_id) as rs_max
		from api.plr_agg
		where left(cast(season_id as varchar(5)), 1) = '2'
		and right(cast(season_id as varchar(5)), 4) != '9999'
		group by player_id
	) b on b.player_id = a.player_id
	left join (
		select player_id, min(season_id) as po_min, max(season_id) as po_max
		from api.plr_agg
		where left(cast(season_id as varchar(5)), 1) = '4'
		and right(cast(season_id as varchar(5)), 4) != '9999'
		group by player_id
	) c on c.player_id = a.player_id
	`,
}

// GetSeasons
var RSeasons = Query{
	Args: []string{},
	Q: `
	select szn_id, szn_desc, wszn_desc
	from lg.szn	
	where left(cast(szn_id as varchar(5)), 1) in ('2', '4')
	and right(cast(szn_id as varchar(5)), 4) != '9999'
	order by right(cast(szn_id as varchar(5)), 4) desc, 
	left(cast(szn_id as varchar(5)), 1);
	`,
}

var Teams = Query{
	Args: []string{},
	Q: `
	select a.lg, a.team_id, a.team, a.team_name
	from team a
	inner join ( 
		select season_id, team_id
		from t_box
		where left(season_id, 1) = '2'
		and right(season_id, 4) >= '2000'
		group by season_id, team_id
		) b on b.team_id = a.team_id
	where a.lg in ('NBA', 'WNBA')
	group by a.lg, a.team_id, a.team, a.team_name
	`,
}
