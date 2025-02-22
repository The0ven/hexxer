package states

import (
	"hexxer/graphics"
	"hexxer/render"
	"hexxer/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type UnitMode struct {
    game         *types.Game
    selectedUnit *types.Unit
    selectedTile *types.Tile
    offsetW, offsetH int32
    scale       float32
}

func NewUnitMode(game *types.Game) *UnitMode {
    return &UnitMode{
        game:   game,
        scale:  1,
        offsetW: 1300 / 2,
        offsetH: 850 / 2,
    }
}

// Handle input: selecting & moving units
func (u *UnitMode) HandleInput() {
    mouse := rl.GetMousePosition()

    u.scale = min(float32(20), max(float32(0.3), u.scale + rl.GetMouseWheelMove()))
    // rl.DrawText(fmt.Sprintf("Tile: %v", u.selectedTile), 20, 40, 20, rl.RayWhite)

    // Select a unit with left-click
    if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
        tile := graphics.PointToTile(mouse.X, mouse.Y, 15*u.scale, u.offsetW, u.offsetH)
        u.selectedTile = &tile
        if unit, exists := u.game.Units[tile]; exists {
            u.selectedUnit = &unit
        } else {
            u.selectedUnit = nil
        }
    }

    // Move selected unit with right-click, respecting movement range
    if u.selectedUnit != nil && rl.IsMouseButtonPressed(rl.MouseButtonRight) {
        newTile := graphics.PointToTile(mouse.X, mouse.Y, 15*u.scale, u.offsetW, u.offsetH)
        if u.canMovetoTile(newTile) {
            // Move unit in map
            u.game.Units[newTile] = *u.selectedUnit
            u.selectedUnit.Coord = newTile
            delete(u.game.Units, *u.selectedTile)
            u.selectedTile = &newTile
        }
    }

    u.changeUnit()
}

func (u *UnitMode) Update() GameState {
    if rl.IsKeyPressed(rl.KeyEscape) {
        return NewMenuMode(u.game, u)
    }
    if rl.IsKeyPressed(rl.KeyT) {
        return NewTileMode(u.game)
    }
    return u
}

// Check if the new tile is within movement range
func (u *UnitMode) canMovetoTile(tile types.Tile) bool {
    if u.selectedUnit == nil {
        return false
    }
    for _, moveTile := range u.MovementRange(*u.selectedUnit) {
        if moveTile == tile {
            return true
        }
    }
    return false
}

func (u *UnitMode) MovementRange(unit types.Unit) []types.Tile {
    visited := make(map[types.Tile]bool)
    queue := []struct {
        tile     types.Tile
        movement int
    }{{unit.Coord, unit.Movement}}
    
    var moveRange []types.Tile
    
    for len(queue) > 0 {
        current := queue[0]
        queue = queue[1:]

        if visited[current.tile] {
            continue
        }
        visited[current.tile] = true
        moveRange = append(moveRange, current.tile)

        if current.movement <= 0 {
            continue
        }
        
        for _, neighbor := range current.tile.Range(1) {
            terrain, terrainExists := u.game.Map[neighbor]
            unit, unitExists := u.game.Units[neighbor]
            
            if !terrainExists || terrain.Impassable {
                continue
            }
            
            if unitExists && unit.Team != unit.Team {
                continue
            }
            
            moveCost := 1 - terrain.SpeedModifier
            if moveCost <= 0 {
                moveCost = 1
            }
            
            if current.movement >= moveCost {
                queue = append(queue, struct {
                    tile     types.Tile
                    movement int
                }{neighbor, current.movement - moveCost})
            }
        }
    }
    
    return moveRange
}

var unitKeyMappings = map[int32]func(t types.Tile) types.Unit{
        rl.KeyOne:   func(t types.Tile) types.Unit { return types.NewInfantry(t) },
        rl.KeyTwo:   func(t types.Tile) types.Unit { return types.NewHeavyInfantry(t) },
        rl.KeyThree: func(t types.Tile) types.Unit { return types.NewLightCavalry(t) },
        rl.KeyFour:  func(t types.Tile) types.Unit { return types.NewHeavyCavalry(t) },
        rl.KeyFive:  func(t types.Tile) types.Unit { return types.NewRanged(t) },
        rl.KeySix:   func(t types.Tile) types.Unit { return types.NewScout(t) },
    }

func (u *UnitMode) changeUnit() {
    if u.selectedTile == nil {
        return
    }

    if rl.IsKeyPressed(rl.KeyZero) {
        delete(u.game.Map, *u.selectedTile)
        return
    }

    if action, exists := unitKeyMappings[rl.GetKeyPressed()]; exists {
        unit := action(*u.selectedTile)
        u.game.Units[*u.selectedTile] = unit
        u.selectedUnit = &unit
    }
}

// Draw units, selected unit highlight, and movement/sight range
func (u *UnitMode) Draw() {
    render.DrawGame(*u.game, u.offsetW, u.offsetH, u.scale)
    if u.selectedTile != nil {
        render.DrawSeletedTile(u.selectedTile, u.offsetW, u.offsetH, u.scale)
    }
    render.DrawUnits(*u.game, u.offsetW, u.offsetH, u.scale)
    if u.selectedUnit != nil {
        render.DrawSeletedUnit(u.selectedUnit, u.MovementRange(*u.selectedUnit), u.offsetW, u.offsetH, u.scale)
    }
}
