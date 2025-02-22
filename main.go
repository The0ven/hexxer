package main

import (
	"encoding/json"
	"fmt"
    "hexxer/states"
    "hexxer/types"
	"io"
	"os"
	"slices"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var screenWidth = int32(1300)
var screenHeight = int32(850)


func main() {
    fmt.Println("Starting Up...")

    rl.InitWindow(screenWidth, screenHeight, "Hexxer")
    defer rl.CloseWindow()
    rl.SetTargetFPS(60)
    rl.SetExitKey(rl.KeyNumLock)

    // Initialize game and starting state
    game := startTestGame()
    var state states.GameState = states.NewTileMode(&game)


    for !rl.WindowShouldClose() {
        state.HandleInput()
        newState := state.Update()
        if newState != state {
            state = newState
        }

        rl.BeginDrawing()
        rl.ClearBackground(rl.Black)
        state.Draw()
        rl.EndDrawing()
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

    return types.NewGame(testGameTeams(), testGameTerrain(ug.Map), []types.Unit{})
}

func testGameTeams() []types.Team {
    results := []types.Team{}

    results = append(results, types.Team{Colour: rl.Blue, Name: "Goodie Guys", Id: 1})
    results = append(results, types.Team{Colour: rl.Red, Name: "Evil Lebarons", Id: 2})

    return results
}

func testGameTerrain(tiles []types.Tile) []types.Terrain {
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
