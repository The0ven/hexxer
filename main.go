package main

import (
	"encoding/json"
	"fmt"
	"hexxer/graphics"
	"hexxer/types"
	"io"
	"os"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func main() {
    fmt.Println("Starting Up...")

    screenWidth := int32(1600)
    screenHeight := int32(850)

    rl.InitWindow(screenWidth, screenHeight, "Hexxer")
    defer rl.CloseWindow()

    rl.SetTargetFPS(60)

    game := startGame()
    scale := float32(1)
    offsetW := screenWidth / 2
    offsetH := screenHeight / 2

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()

        rl.ClearBackground(rl.Black)
        scale = min(float32(10), max(float32(1), scale + rl.GetMouseWheelMove()))

        

        if rl.IsMouseButtonDown(rl.MouseButtonRight) {
            offsetW = int32(rl.GetMouseDelta().X) + offsetW
            offsetH = int32(rl.GetMouseDelta().Y) + offsetH
        }

        drawGame(game, offsetW, offsetH, scale)

        rl.EndDrawing()
    }
}

func drawGame(game types.Game, offsetW int32, offsetH int32, scale float32) {
    for tile := range game.Map {
        hex := graphics.Hexagon(tile, 15 * scale, game.Map[tile].Colour, offsetW, offsetH)
        rl.DrawPoly(hex.Center, hex.Sides, hex.Radius, hex.Rotation, hex.Col)
    }
}

func startGame() types.Game {
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
                results = append(results, types.NewSea(rl.DarkBlue, seaTile))
            }
        }
        results = append(results, types.NewCoast(rl.Beige, tile))
    }

    return results
}
