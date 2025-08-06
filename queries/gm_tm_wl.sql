/*
working on this trying to recreate 'final' field in the recent games query.
originally (8/6 1pm) returns game id, matchup from winning team, winning team id,
then losing team id
*/
select a.game_id, a.matchup, a.team_id as wtm, b.team_id as ltm
from stats.tbox a
inner join (
	select game_id, team_id
	from stats.tbox
	where wl = 'L'
) b on b.game_id = a.game_id
where a.wl = 'W';
