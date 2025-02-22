package states

type GameState interface {
    HandleInput()
    Update() GameState
    Draw()
}
