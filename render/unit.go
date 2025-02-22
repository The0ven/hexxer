package render

import (
	"fmt"
	"hexxer/graphics"
	"hexxer/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawUnits(game types.Game, offsetW int32, offsetH int32, scale float32) {
    for tile, unit := range game.Units {
        marker := graphics.Circle(tile, (15) * scale, game.Teams[unit.Team].Colour, offsetW, offsetH)
        rl.DrawCircleV(marker.Center, marker.Radius/2, marker.Col)
        rl.DrawCircleLinesV(marker.Center, marker.Radius/2, rl.DarkGray)
    }
}

func DrawSeletedUnit(selectedUnit *types.Unit, offsetW int32, offsetH int32, scale float32) {
    marker := graphics.Circle(selectedUnit.Coord, (15) * scale, rl.White, offsetW, offsetH)
    rl.DrawCircleLinesV(marker.Center, marker.Radius/2, rl.Red)

    // Draw movement range
    for _, moveTile := range selectedUnit.MovementRange() {
        movePos := graphics.PlaceTile(moveTile, 15*scale, offsetW, offsetH)
        rl.DrawCircle(int32(movePos.X), int32(movePos.Y), 5*scale, rl.Green)
    }

    // Draw sight range
    for _, sightTile := range selectedUnit.SightRange() {
        sightPos := graphics.PlaceTile(sightTile, 15*scale, offsetW, offsetH)
        rl.DrawCircleLines(int32(sightPos.X), int32(sightPos.Y), 6*scale, rl.Blue)
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
