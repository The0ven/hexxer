package main

import (
	"encoding/json"
	"fmt"
	"hexxer/graphics"
	"hexxer/types"
	"image/color"
	"io"
	"os"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var defaultRadius = float32(15)
var screenWidth = int32(1300)
var screenHeight = int32(850)
var game = startTestGame()

func main() {
    fmt.Println("Starting Up...")

    rl.InitWindow(screenWidth, screenHeight, "Hexxer")
    defer rl.CloseWindow()

    rl.SetTargetFPS(60)

    scale := float32(1)
    offsetW := int32(float32(screenWidth) / float32(2))
    offsetH := int32(float32(screenHeight) / float32(2))
    mode := "tile"
    var currentTile *types.Tile

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()
        rl.ClearBackground(rl.Black)
        rl.DrawText(fmt.Sprintf("Map Len: %v", len(game.Map)), 20, 20, 20, rl.RayWhite)

        mouse := rl.GetMousePosition()

        scale = min(float32(20), max(float32(0.3), scale + rl.GetMouseWheelMove()))
        //TODO: Add support for sprite changes based on scale

        if rl.IsMouseButtonDown(rl.MouseButtonRight) {
            offsetW = int32(rl.GetMouseDelta().X) + offsetW
            offsetH = int32(rl.GetMouseDelta().Y) + offsetH
        }

        drawGame(game, offsetW, offsetH, scale)
        
        if mode == "tile" {
            if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
                tile := graphics.PointToTile(mouse.X, mouse.Y, defaultRadius * scale, offsetW, offsetH)
                currentTile = &tile
            }

            if currentTile != nil {
                drawTooltip(game, *currentTile)
                changeTile(game.Map, currentTile)
                drawSeletedTile(currentTile, offsetW, offsetH, scale)
            }
        }

        if mode == "menu" {
            if rl.IsKeyPressed(rl.KeyTab) {
                mode = "tile"
            }
            drawMenu(&game, mouse, offsetW, offsetH)
        } else {
            if rl.IsKeyPressed(rl.KeyTab) {
                mode = "menu"
            }
        }

        rl.EndDrawing()
    }
}

func saveGame(game types.Game) {
    teams := []types.Team{}
    Map := []types.Terrain{}
    units := []types.Unit{}
    for _, t := range game.Teams {
        teams = append(teams, t)
    }
    for _, tr := range game.Map {
        Map = append(Map, tr)
    }
    for _, u := range game.Units {
        units = append(units, u)
    }
    unmappedGame := types.UnmappedGame{Teams: teams, Map: Map, Units: units}
    blob, err := json.Marshal(unmappedGame)
    if err != nil {
        fmt.Printf("Error marshalling!!!\n%v\n", err)
    }
    err = os.WriteFile("saved_game.json", blob, 0666)
    if err != nil {
        fmt.Printf("Error writing file!!!\n%v\n", err)
    }
}

func loadGame() *types.Game {
    file, err := os.Open("saved_game.json")
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()

    fileBytes, err := io.ReadAll(file)
    if err != nil {
        fmt.Println(err)
    }
    
    var ug types.UnmappedGame
    err = json.Unmarshal(fileBytes, &ug)
    if err != nil {
        fmt.Println(err)
    }

    game := types.NewGame(ug.Teams, ug.Map, ug.Units)
    fmt.Printf("Len loaded game map: %v", len(game.Map))
    return &game
}

func drawMenu(game *types.Game, mouse rl.Vector2, offsetW int32, offsetH int32) {
    size := int32(30)
    options := []string{"Exit", "Load", "Save"}
    for idx, option := range options {
        w := float32(rl.MeasureText(option, size)) + 40
        h := float32(size) + 10
        x := float32(offsetW) - w / float32(2)
        y := float32(offsetH) - (h + 20) * float32(idx)

        button := graphics.NewButton(x, y, w, h, option, rl.LightGray, rl.DarkGray)
        button.DrawButton(size)
        if button.IsPressed(mouse) {
            fmt.Printf("Button %v pressed!\n", option)
            switch option {
            case "Exit":
                rl.CloseWindow()
            case "Save":
                saveGame(*game)
            case "Load":
                *game = *loadGame()
            }
        }
    }
}

func drawSeletedTile(currentTile *types.Tile, offsetW int32, offsetH int32, scale float32) {   
    hex := graphics.Hexagon(*currentTile, defaultRadius * scale, rl.Red, offsetW, offsetH)
    rl.DrawPolyLines(hex.Center, hex.Sides, hex.Radius, hex.Rotation, hex.Col)
}

func changeTile(gamemap map[types.Tile]types.Terrain, currentTile *types.Tile) {
    var result = types.Terrain{}
    if rl.IsKeyPressed(rl.KeyZero) {
        delete(gamemap, *currentTile)
    }
    if rl.IsKeyPressed(rl.KeyOne) {
        result = types.NewSea(*currentTile)
    }
    if rl.IsKeyPressed(rl.KeyTwo) {
        result = types.NewCoast(*currentTile)
    }
    if rl.IsKeyPressed(rl.KeyThree) {
        result = types.NewLand(2, *currentTile)
    }
    if rl.IsKeyPressed(rl.KeyFour) {
        result = types.NewForest(2, *currentTile)
    }
    if rl.IsKeyPressed(rl.KeyFive) {
        result = types.NewHill(3, *currentTile)
    }
    if rl.IsKeyPressed(rl.KeySix) {
        result = types.NewMountain(4, *currentTile)
    }
    if result != (types.Terrain{}) {
        gamemap[*currentTile] = result
    }
}

func drawGame(game types.Game, offsetW int32, offsetH int32, scale float32) {
    for tile := range game.Map {
        hex := graphics.Hexagon(tile, defaultRadius * scale, game.Map[tile].Colour, offsetW, offsetH)
        rl.DrawPoly(hex.Center, hex.Sides, hex.Radius, hex.Rotation, hex.Col)
    }
}

func drawTooltip(game types.Game, tile types.Tile) {
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

func startTestGame() types.Game {
    file, err := os.Open("testgame.json")
    if err != nil {
        fmt.Println(err)
    }
    defer file.Close()

    fileBytes, err := io.ReadAll(file)
    if err != nil {
        fmt.Println(err)
    }
    
    var ug types.TestUnmappedGame
    err = json.Unmarshal(fileBytes, &ug)
    if err != nil {
        fmt.Println(err)
    }

    return types.NewGame([]types.Team{}, testGame(ug.Map), []types.Unit{})
}

func testGame(tiles []types.Tile) []types.Terrain {
    results := []types.Terrain{}

    for _, tile := range tiles {
        tR := tile.Range(2)
        for _, seaTile := range tR {
            if !slices.Contains(tiles, seaTile) {
                results = append(results, types.NewSea(seaTile))
            }
        }
        results = append(results, types.NewCoast(tile))
    }

    return results
}
