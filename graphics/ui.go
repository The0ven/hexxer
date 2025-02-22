package graphics

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TextButton struct {
    BoundingBox rl.Rectangle
    Text string
    TextX float32
    TextY float32
    TextSize int32
    TextColour color.RGBA
    BGColour color.RGBA
}

func NewTextButton(x float32, y float32, size int32, text string, textColour color.RGBA, bgColour color.RGBA) TextButton {
    w := rl.MeasureText(text, size) + 20
    h := size + 10
    return TextButton{rl.NewRectangle(x - float32(10), y - float32(5), float32(w), float32(h)), text, x, y, size, textColour, bgColour}
}

func (b TextButton) DrawTextButton() {
    rl.DrawRectangleRec(b.BoundingBox, b.BGColour)
    rl.DrawText(b.Text, int32(b.TextX), int32(b.TextY), b.TextSize, b.TextColour)
}

func (b TextButton) IsPressed(mouse rl.Vector2) bool {
    return rl.CheckCollisionPointRec(mouse, b.BoundingBox) && rl.IsMouseButtonPressed(rl.MouseButtonLeft)
}
