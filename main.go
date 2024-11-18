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

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()

        rl.ClearBackground(rl.RayWhite)
        drawGame(game, screenWidth, screenHeight)

        rl.EndDrawing()
    }
}

func drawGame(game types.Game, width int32, height int32) {
    offsetW := width / 2
    offsetH := height / 2
    for tile := range game.Map {
        hex := graphics.Hexagon(tile, 30, game.Map[tile].Colour, offsetW, offsetH)
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
