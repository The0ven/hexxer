package types

type Unit struct {
    Name string
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

func NewInfantry(tile Tile) Unit {
    return Unit{
        Name: "Infantry",
        Team: 0,
        Coord: tile,
        Movement: 2,
        Sight: 2,
        Strength: 4,
        CombatRage: 1,
        Size: 1,
    }
}

func NewHeavyInfantry(tile Tile) Unit {
    return Unit{
        Name: "Heavy Infantry",
        Team: 0,
        Coord: tile,
        Movement: 1,
        Sight: 2,
        Strength: 6,
        CombatRage: 1,
        Size: 2,
    }
}

func NewLightCavalry(tile Tile) Unit {
    return Unit{
        Name: "Light Cavalry",
        Team: 0,
        Coord: tile,
        Movement: 4,
        Sight: 2,
        Strength: 3,
        CombatRage: 1,
        Size: 2,
    }
}

func NewHeavyCavalry(tile Tile) Unit {
    return Unit{
        Name: "Heavy Cavalry",
        Team: 0,
        Coord: tile,
        Movement: 3,
        Sight: 2,
        Strength: 5,
        CombatRage: 1,
        Size: 3,
    }
}

func NewRanged(tile Tile) Unit {
    return Unit{
        Name: "Archers",
        Team: 0,
        Coord: tile,
        Movement: 2,
        Sight: 2,
        Strength: 2,
        CombatRage: 3,
        Size: 1,
    }
}

func NewScout(tile Tile) Unit {
    return Unit{
        Name: "Archers",
        Team: 0,
        Coord: tile,
        Movement: 3,
        Sight: 3,
        Strength: 1,
        CombatRage: 1,
        Size: 1,
    }
}
