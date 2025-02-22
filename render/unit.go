package render

import (
	"fmt"
	"hexxer/graphics"
	"hexxer/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawUnits(game types.Game, offsetW int32, offsetH int32, scale float32) {
    for tile, unit := range game.Units {
        marker := graphics.Circle(tile, float32(15) * scale, game.Teams[unit.Team].Colour, offsetW, offsetH)
        rl.DrawCircleV(marker.Center, marker.Radius/2, marker.Col)
        rl.DrawCircleLinesV(marker.Center, marker.Radius/2, rl.DarkGray)
    }
}

func DrawSeletedUnit(selectedUnit *types.Unit, movementRange []types.Tile, offsetW int32, offsetH int32, scale float32) {
    marker := graphics.Circle(selectedUnit.Coord, float32(15) * scale, rl.White, offsetW, offsetH)
    rl.DrawCircleLinesV(marker.Center, marker.Radius/2, rl.Red)

    // Draw movement range
    for _, moveTile := range movementRange {
        if moveTile == selectedUnit.Coord {
            continue
        }
        movePos := graphics.Circle(moveTile, float32(15) * scale, rl.White, offsetW, offsetH)
        rl.DrawCircleV(movePos.Center, movePos.Radius/2, rl.Green)
    }

    // Draw sight range
    for _, sightTile := range selectedUnit.SightRange() {
        if sightTile == selectedUnit.Coord {
            continue
        }
        sightPos := graphics.Circle(sightTile, float32(15) * scale, rl.White, offsetW, offsetH)
        rl.DrawCircleLinesV(sightPos.Center, 6*scale, rl.Blue)
    }
}

func drawUnitTooltip(game types.Game, tile types.Tile) {
    unit, ok := game.Units[tile]

    x := int32(1300 - 170)
    y := int32(850 - 40)
    if ok {
        team := game.Teams[unit.Team]
        
        rl.DrawRectangle(x-6, y-26, 170, 50, team.Colour)

        rl.DrawText(unit.Name, x, y, 20, rl.RayWhite)
        rl.DrawText(fmt.Sprintf("Strength: %v", unit.Strength), x, y - 21, 20, rl.RayWhite)
    }
}
