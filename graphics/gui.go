package graphics

import rl "github.com/gen2brain/raylib-go/raylib"

type Button struct {
    X float32
    Y float32
    W float32
    H float32
    Rect rl.Rectangle
    Text string
    Colour rl.Color
    TextColour rl.Color
}

func NewButton(x float32, y float32, w float32, h float32, text string, colour rl.Color, textColour rl.Color) Button {
    return Button{x, y, w, h, rl.Rectangle{X: x, Y: y, Width: w, Height: h}, text, colour, textColour}
}

func (button Button) IsPressed(mouse rl.Vector2) bool {
    return rl.CheckCollisionPointRec(mouse, button.Rect) && rl.IsMouseButtonPressed(rl.MouseButtonLeft)
}

func (button Button) DrawButton(size int32) {
    rl.DrawRectangleRounded(button.Rect, 0.2, 20, button.Colour)
    rl.DrawText(button.Text, int32(button.X + float32(20)), int32(button.Y + float32(5)), size, button.TextColour)
}
