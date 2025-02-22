package states

import (
	"hexxer/graphics"
	"hexxer/render"
	"hexxer/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TileMode struct {
    game        *types.Game
    selectedTile *types.Tile
    offsetW, offsetH int32
    scale       float32
}

func NewTileMode(game *types.Game) *TileMode {
    return &TileMode{
        game:   game,
        scale:  1,
        offsetW: 1300 / 2,
        offsetH: 850 / 2,
    }
}


func (t *TileMode) HandleInput() {
    if rl.IsKeyPressed(rl.KeyEscape) {
        return
    }

    t.scale = min(float32(20), max(float32(0.3), t.scale + rl.GetMouseWheelMove()))
    // rl.DrawText(fmt.Sprintf("Tile: %v", t.selectedTile), 20, 40, 20, rl.RayWhite)

    if rl.IsMouseButtonDown(rl.MouseButtonRight) {
        t.offsetW += int32(rl.GetMouseDelta().X)
        t.offsetH += int32(rl.GetMouseDelta().Y)
    }

    if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
        mouse := rl.GetMousePosition()
        tile := graphics.PointToTile(mouse.X, mouse.Y, 15*t.scale, t.offsetW, t.offsetH)
        t.selectedTile = &tile
    }

    t.changeTile()
}

func (t *TileMode) Update() GameState {
    if rl.IsKeyPressed(rl.KeyEscape) {
        return NewMenuMode(t.game, t)
    }
    if rl.IsKeyPressed(rl.KeyU) {
        return NewUnitMode(t.game)
    }
    return t
}

func (t *TileMode) Draw() {
    render.DrawGame(*t.game, t.offsetW, t.offsetH, t.scale)

    if t.selectedTile != nil {
        render.DrawSeletedTile(t.selectedTile, t.offsetW, t.offsetH, t.scale)
        render.DrawTerrainTooltip(*t.game, *t.selectedTile)
    }
}

var tileKeyMappings = map[int32]func(t types.Tile) types.Terrain{
        rl.KeyOne:   func(t types.Tile) types.Terrain { return types.NewSea(t) },
        rl.KeyTwo:   func(t types.Tile) types.Terrain { return types.NewCoast(t) },
        rl.KeyThree: func(t types.Tile) types.Terrain { return types.NewLand(2, t) },
        rl.KeyFour:  func(t types.Tile) types.Terrain { return types.NewForest(2, t) },
        rl.KeyFive:  func(t types.Tile) types.Terrain { return types.NewHill(3, t) },
        rl.KeySix:   func(t types.Tile) types.Terrain { return types.NewMountain(4, t) },
    }

func (t *TileMode) changeTile() {
    if t.selectedTile == nil {
        return
    }

    if rl.IsKeyPressed(rl.KeyZero) {
        delete(t.game.Map, *t.selectedTile)
        return
    }

    if action, exists := tileKeyMappings[rl.GetKeyPressed()]; exists {
        t.game.Map[*t.selectedTile] = action(*t.selectedTile)
    }
}
