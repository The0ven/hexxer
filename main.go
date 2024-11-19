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
var screenWidth = int32(1600)
var screenHeight = int32(850)

func main() {
    fmt.Println("Starting Up...")

    rl.InitWindow(screenWidth, screenHeight, "Hexxer")
    defer rl.CloseWindow()

    rl.SetTargetFPS(60)

    game := startGame()
    scale := float32(1)
    offsetW := screenWidth / 2
    offsetH := screenHeight / 2
    mode := "tile"
    var currentTile *types.Tile

    for !rl.WindowShouldClose() {
        rl.BeginDrawing()
        rl.ClearBackground(rl.Black)

        scale = min(float32(10), max(float32(1), scale + rl.GetMouseWheelMove()))
        //TODO: Add support for sprite changes based on scale

        if rl.IsMouseButtonDown(rl.MouseButtonRight) {
            offsetW = int32(rl.GetMouseDelta().X) + offsetW
            offsetH = int32(rl.GetMouseDelta().Y) + offsetH
        }

        drawGame(game, offsetW, offsetH, scale, currentTile)
        
        if mode == "tile" {
            if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
                mouse := rl.GetMousePosition()
                tile := graphics.PointToTile(mouse.X, mouse.Y, defaultRadius * scale, offsetH, offsetW)
                currentTile = &tile
            }

            if currentTile != nil {
                drawTooltip(game, *currentTile, offsetW, offsetH, scale)
                changeTile(game.Map, currentTile)
            }

        }

        rl.EndDrawing()
    }
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

func drawGame(game types.Game, offsetW int32, offsetH int32, scale float32, currentTile *types.Tile) {
    for tile := range game.Map {
        hex := graphics.Hexagon(tile, defaultRadius * scale, game.Map[tile].Colour, offsetW, offsetH)
        if currentTile != nil && *currentTile == tile {
            rl.DrawText(fmt.Sprintf("Tile: %v\n %v", tile, hex.Center), 20, 61, 20, rl.White)
        }
        rl.DrawPoly(hex.Center, hex.Sides, hex.Radius, hex.Rotation, hex.Col)
    }
}

func drawTooltip(game types.Game, tile types.Tile, offsetW int32, offsetH int32, scale float32) {
    hex := graphics.Hexagon(tile, defaultRadius * scale, rl.Black, offsetW, offsetH)
    rl.DrawText(fmt.Sprintf("Tile: %v\n %v", tile, hex.Center), 20, 20, 20, rl.White)
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
                results = append(results, types.NewSea(seaTile))
            }
        }
        results = append(results, types.NewCoast(tile))
    }

    return results
}
