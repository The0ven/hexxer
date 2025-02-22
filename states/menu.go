package states

import (
	"hexxer/render"
	"hexxer/types"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type MenuMode struct {
    game         *types.Game
    previousMode GameState
}

func NewMenuMode(game *types.Game, prev GameState) *MenuMode {
    return &MenuMode{game: game, previousMode: prev}
}

func (m *MenuMode) HandleInput() {
    if rl.IsKeyPressed(rl.KeyTab) {
        return
    }
}

func (m *MenuMode) Update() GameState {
    if rl.IsKeyPressed(rl.KeyTab) {
        return m.previousMode
    }
    return m
}

func (m *MenuMode) Draw() {
    m.previousMode.Draw()
    render.DrawMenu(m.game, rl.GetMousePosition())
}
