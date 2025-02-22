package save

import (
	"encoding/json"
	"fmt"
	"hexxer/types"
	"io"
	"os"
)

func SaveGame(game types.Game) {
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

func LoadGame() *types.Game {
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


