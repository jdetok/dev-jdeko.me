/* returns the top scorer in each game (between both team) with the highest
scorer of the day first */
select * from (
    select distinct on (a.game_id)
    a.game_id,  
    a.team_id, 
    d.player_id, 
    e.player, 
    case 
        when c.lg_id = 0 then 'NBA'
        when c.lg_id = 1 then 'WNBA'
        end as lg, 
    c.team,
    c.team_long,
    a.gdate,
    a.matchup, 
    a.wl, 
    a.pts as tm_pts, 
    d.pts as plr_pts
    from stats.tbox a
    inner join (
        select max(gdate) as md
        from stats.tbox
    ) b on a.gdate = b.md
    inner join lg.team c on c.team_id = a.team_id
    inner join stats.pbox d on d.game_id = a.game_id and d.team_id = a.team_id
    inner join lg.plr e on e.player_id = d.player_id
    order by a.game_id, d.pts desc, (d.ast + d.reb + d.stl + d.blk) desc)
order by plr_pts desc;
