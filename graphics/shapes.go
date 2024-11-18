package graphics

import (
	"hexxer/types"
	"image/color"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Shape struct {
    Center rl.Vector2
    Sides int32
    Radius float32
    Rotation float32
    Col color.RGBA
}

func Hexagon(tile types.Tile, radius float32, colour color.RGBA, offsetH int32, offsetW int32) Shape {
    center := PlaceTile(tile, radius, offsetW, offsetH)
    return Shape{
        center,
        6,
        radius,
        30,
        colour,
    }
}

func PlaceTile(tile types.Tile, radius float32, offsetW int32, offsetH int32) rl.Vector2 {
    x := radius * (float32(math.Sqrt(3) * tile.ToFractional().Q + math.Sqrt(3)/2 * tile.ToFractional().R))
    y := radius * 3/2 * float32(tile.R)
    return rl.Vector2{X: float32(offsetH) + x, Y: float32(offsetW) + y}
}

func PointToTile(x float32, y float32, radius float32) types.Tile {
    q := (math.Sqrt(3)/3 * float64(x - 1/3 * y)) / float64(radius)
    r := float64((2/3 * y) / radius)
    return types.FracAxialTile{Q: q, R: r}.ToCubic().Round()
}
