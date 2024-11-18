package types

type TestUnmappedGame struct {
    Map []Tile
}

type UnmappedGame struct {
    Teams []Team
    Map []Terrain
    Units []Unit
}

type Game struct {
    Teams map[int]Team
    Map map[Tile]Terrain
    Units map[Tile]Unit
}

func NewGame(teams []Team, gamemap []Terrain, units []Unit) Game {
    ts := make(map[int]Team)
    ms := make(map[Tile]Terrain)
    us := make(map[Tile]Unit)
    for i, t := range teams {
        ts[i] = t
    }
    for _, m := range gamemap {
        ms[m.Coord] = m
    }
    for _, u := range units {
        us[u.Coord] = u
    }
    return Game{
        ts,
        ms,
        us,
    }
}
