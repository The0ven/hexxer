package states

import (
	"fmt"
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
    rl.DrawText(fmt.Sprintf("Tile: %v", u.selectedTile), 20, 40, 20, rl.RayWhite)

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
        if u.isTileInMovementRange(newTile) {
            // Move unit in map
            u.game.Units[newTile] = *u.selectedUnit
            delete(u.game.Units, *u.selectedTile)
            u.selectedTile = &newTile
        }
    }
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
func (u *UnitMode) isTileInMovementRange(tile types.Tile) bool {
    if u.selectedUnit == nil {
        return false
    }
    for _, moveTile := range u.selectedUnit.MovementRange() {
        if moveTile == tile {
            return true
        }
    }
    return false
}

// Draw units, selected unit highlight, and movement/sight range
func (u *UnitMode) Draw() {
    render.DrawGame(*u.game, u.offsetW, u.offsetH, u.scale)
    render.DrawUnits(*u.game, u.offsetW, u.offsetH, u.scale)
    if u.selectedTile != nil {
        render.DrawSeletedTile(u.selectedTile, u.offsetW, u.offsetH, u.scale)
    }
    if u.selectedUnit != nil {
        render.DrawSeletedUnit(u.selectedUnit, u.offsetW, u.offsetH, u.scale)
    }
}
