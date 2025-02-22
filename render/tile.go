package render

import (
	"hexxer/graphics"
	"hexxer/types"
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)


func DrawSeletedTile(currentTile *types.Tile, offsetW int32, offsetH int32, scale float32) {   
    hex := graphics.Hexagon(*currentTile, 15*scale, rl.Red, offsetW, offsetH)
    rl.DrawPolyLines(hex.Center, hex.Sides, hex.Radius, hex.Rotation, hex.Col)
}

func DrawGame(game types.Game, offsetW int32, offsetH int32, scale float32) {
    for tile := range game.Map {
        hex := graphics.Hexagon(tile, 15*scale, game.Map[tile].Colour, offsetW, offsetH)
        rl.DrawPoly(hex.Center, hex.Sides, hex.Radius, hex.Rotation, hex.Col)
    }
}

func DrawTerrainTooltip(game types.Game, tile types.Tile) {
    options := []string{"Sea","Coast","Plains","Forest","Hills","Mountain","None"}

    x := int32(1300 - 120)
    rl.DrawRectangle(x, 850 - 120, int32(11 * len(options) + 4), 94, color.RGBA{20, 0, 30, 88})
    for i, option := range options {
        y := 850 - 40 - (int32(i) * 21)
        if _, ok := game.Map[tile]; (ok && game.Map[tile].Type == option) || (!ok && option == "None") {
            rl.DrawRectangle(x - 24, y + 1, 18, 18, rl.RayWhite)
        } else {
            rl.DrawRectangleLinesEx(rl.NewRectangle(float32(x - 24), float32(y + 1), 18, 18), 1, rl.RayWhite)
        }
        rl.DrawText(option, x, y, 20, rl.RayWhite)
    }
}
