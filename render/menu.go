package render

import (
	"hexxer/graphics"
	"hexxer/types"
	"image/color"
    "hexxer/save"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawMenu(game *types.Game, mouse rl.Vector2) {
    options := []string{"Exit", "Load", "Save"}
    for idx, option := range options {
        x := 1300/2 - 100
        y := 850 * 0.4 + (idx * 50)

        button := graphics.NewButton(float32(x), float32(y), 200, 40, option, rl.LightGray, rl.DarkGray)
        button.DrawButton(24)
        if button.IsPressed(mouse) {
            switch option {
            case "Exit":
                rl.CloseWindow()
            case "Save":
                save.SaveGame(*game)
            case "Load":
                *game = *save.LoadGame()
            }
        }
    }
}

func DrawTooltip(game types.Game, tile types.Tile, screenWidth int32, screenHeight int32) {
    options := []string{"Sea","Coast","Plains","Forest","Hills","Mountains","None"}

    x := int32(screenWidth - 120)
    rl.DrawRectangle(x, screenHeight - 120, int32(11 * len(options) + 4), 94, color.RGBA{20, 0, 30, 88})
    for i, option := range options {
        y := screenHeight - 40 - (int32(i) * 21)
        if _, ok := game.Map[tile]; (ok && game.Map[tile].Type == option) || (!ok && option == "None") {
            rl.DrawRectangle(x - 24, y + 1, 18, 18, rl.RayWhite)
        } else {
            rl.DrawRectangleLinesEx(rl.NewRectangle(float32(x - 24), float32(y + 1), 18, 18), 1, rl.RayWhite)
        }
        rl.DrawText(option, x, y, 20, rl.RayWhite)
    }
}
