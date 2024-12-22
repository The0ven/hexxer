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
var scale = float32(1)
var offsetW = int32(float32(screenWidth) / float32(2))
var offsetH = int32(float32(screenHeight) / float32(2))
var modes = []string{"terrain", "units", "teams"}
var mode = 0
var menuIsActive = false
var currentTile *types.Tile
var mouse rl.Vector2
var fogOfWarTeam = 0

func main() {
    fmt.Println("Starting Up...")

    rl.InitWindow(screenWidth, screenHeight, "Hexxer")
    var game = startTestGame()
    defer rl.CloseWindow()

    rl.SetTargetFPS(60)
    rl.SetExitKey(rl.KeyNumLock)


    for !rl.WindowShouldClose() {
        rl.BeginDrawing()
        rl.ClearBackground(rl.Black)
        rl.DrawText(fmt.Sprintf("Mode: %v", modes[mode]), 20, 20, 20, rl.RayWhite)

        mouse = rl.GetMousePosition()

        scale = min(float32(20), max(float32(0.3), scale + rl.GetMouseWheelMove()))

        //Set current tile
        if rl.IsMouseButtonDown(rl.MouseButtonLeft) && !menuIsActive {
            tile := graphics.PointToTile(mouse.X, mouse.Y, defaultRadius * scale, offsetW, offsetH)
            currentTile = &tile
        }

        //Zoom in-out
        if rl.IsMouseButtonDown(rl.MouseButtonRight) {
            offsetW = int32(rl.GetMouseDelta().X) + offsetW
            offsetH = int32(rl.GetMouseDelta().Y) + offsetH
        }

        //Switch mode
        if rl.IsKeyPressed(rl.KeyTab) {
            if mode < len(modes)-1 {
                mode += 1
            } else {
                mode = 0
            }
        }

        //TODO: Add support for sprite changes based on scale
        gameLoop(&game)
        drawGame(game, offsetW, offsetH, scale)
        drawUI(&game, mouse) 

        rl.EndDrawing()
    }
}

func gameLoop(game *types.Game) {
    if rl.IsKeyPressed(rl.KeyEscape) {
        menuIsActive = !menuIsActive
    }
    switch modes[mode] {
    case "terrain":
        if currentTile != nil {
            changeTileTerrain(game.Map, currentTile)
        }
    case "units":
        if currentTile != nil {
            changeTileUnit(game.Units, currentTile)
            // TEAM CHANGE BROKEN TODO: FIX IT
            changeTileTeam(game, currentTile)
        }
    }
}

func changeTileTeam(game *types.Game, currentTile *types.Tile) {
    if unit, ok := game.Units[*currentTile]; ok && rl.IsKeyPressed(rl.KeyT) {
        if unit.Team < len(game.Teams){
            unit.Team += 1
            fmt.Printf("Unit: %v\n", game.Teams[unit.Team].Name)
        } else {
            unit.Team = 1
        }
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
    fmt.Printf("Len loaded game map: %v\n", len(game.Map))
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

        button := graphics.NewTextButton(x, y, size, option, rl.LightGray, rl.DarkGray)
        button.DrawTextButton()
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
    if currentTile == nil {
        return
    }
    hex := graphics.Hexagon(*currentTile, defaultRadius * scale, rl.Red, offsetW, offsetH)
    rl.DrawPolyLines(hex.Center, hex.Sides, hex.Radius, hex.Rotation, hex.Col)
}

func changeTileTerrain(gamemap map[types.Tile]types.Terrain, currentTile *types.Tile) {
    var result = types.Terrain{}

    switch rl.GetKeyPressed() {
    case rl.KeyZero:
        delete(gamemap, *currentTile)
    case rl.KeyOne:
        result = types.NewSea(*currentTile)
    case rl.KeyTwo:
        result = types.NewCoast(*currentTile)
    case rl.KeyThree:
        result = types.NewLand(2, *currentTile)
    case rl.KeyFour:
        result = types.NewForest(2, *currentTile)
    case rl.KeyFive:
        result = types.NewHill(3, *currentTile)
    case rl.KeySix:
        result = types.NewMountain(4, *currentTile)
    }

    if result != (types.Terrain{}) {
        gamemap[*currentTile] = result
    }
}

func changeTileUnit(gameUnits map[types.Tile]types.Unit, currentTile *types.Tile) {
    var result = types.Unit{}

    switch rl.GetKeyPressed() {
    case rl.KeyZero:
        delete(gameUnits, *currentTile)
    case rl.KeyOne:
        result = types.NewInfantry(*currentTile)
    case rl.KeyTwo:
        result = types.NewHeavyInfantry(*currentTile)
    case rl.KeyThree:
        result = types.NewLightCavalry(*currentTile)
    case rl.KeyFour:
        result = types.NewHeavyCavalry(*currentTile)
    case rl.KeyFive:
        result = types.NewRanged(*currentTile)
    case rl.KeySix:
        result = types.NewScout(*currentTile)
    }

    if result != (types.Unit{}) {
        gameUnits[*currentTile] = result
    }
}

func drawGame(game types.Game, offsetW int32, offsetH int32, scale float32) {
    drawTiles(game, offsetW, offsetH, scale)
    switch modes[mode] {
    case "units":
        drawUnits(game, offsetW, offsetH, scale)
    case "teams":
        drawUnits(game, offsetW, offsetH, scale)
        drawFogOfWar(game, offsetW, offsetH, scale)
    }
}

func drawFogOfWar(game types.Game, offsetW int32, offsetH int32, scale float32) {
    return
}

func drawTiles(game types.Game, offsetW int32, offsetH int32, scale float32) {
    for tile := range game.Map {
        hex := graphics.Hexagon(tile, defaultRadius * scale, game.Map[tile].Colour, offsetW, offsetH)
        rl.DrawPoly(hex.Center, hex.Sides, hex.Radius, hex.Rotation, hex.Col)
    }
}

func drawUnits(game types.Game, offsetW int32, offsetH int32, scale float32) {
    for unitTile := range game.Units {
        unit := game.Units[unitTile]
        marker := graphics.Circle(unitTile, (defaultRadius) * scale, game.Teams[unit.Team].Colour, offsetW, offsetH)
        rl.DrawCircleV(marker.Center, marker.Radius/2, marker.Col)
        rl.DrawCircleLinesV(marker.Center, marker.Radius/2, rl.DarkGray)
    }
}

func drawUI(game *types.Game, mouse rl.Vector2) {
    switch modes[mode] {
    case "terrain":
        if currentTile != nil {
            drawTerrainTooltip(*game, *currentTile)
        }
    case "units":
        if currentTile != nil {
            drawUnitTooltip(*game, *currentTile)
        }
    case "teams":
        if currentTile != nil {
            drawTeamTooltip(*game, *currentTile)
        }
    }

    drawSeletedTile(currentTile, offsetW, offsetH, scale)

    if menuIsActive {
        drawMenu(game, mouse, offsetW, offsetH)
    }
}

func drawUnitTooltip(game types.Game, tile types.Tile) {
    unit, ok := game.Units[tile]

    x := int32(screenWidth - 170)
    y := int32(screenHeight - 40)
    if ok {
        team := game.Teams[unit.Team]
        
        rl.DrawRectangle(x-6, y-26, 170, 50, team.Colour)

        rl.DrawText(unit.Name, x, y, 20, rl.RayWhite)
        rl.DrawText(fmt.Sprintf("Strength: %v", unit.Strength), x, y - 21, 20, rl.RayWhite)
    }
}

func drawTeamTooltip(game types.Game, tile types.Tile) {
    return
}

func drawTerrainTooltip(game types.Game, tile types.Tile) {
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
