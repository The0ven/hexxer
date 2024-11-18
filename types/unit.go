package types

type Unit struct {
    Team int
    Coord Tile
    Movement int
    Sight int
    Strength int
    CombatRage int
    Size int
}

func (u Unit) Count() int {
    return 2 * u.Size
}

func (u Unit) MovementRange() []Tile {
    return u.Coord.Range(u.Movement)
}

func (u Unit) SightRange() []Tile {
    return u.Coord.Range(u.Sight)
}
