package types

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type TestTerrain struct {
    Coord Tile
}

type Terrain struct {
    Colour color.RGBA
    Impassable bool
    SpeedModifier int
    Height int
    Naval bool
    VisionModifier int
    BlocksVision bool
    Coord Tile
    Type string
}

func NewSea(tile Tile) Terrain {
    return Terrain{
        rl.DarkBlue,
        false,
        0,
        0,
        true,
        1,
        false,
        tile,
        "Sea",
    }
}

func NewCoast(tile Tile) Terrain {
    return Terrain{
        rl.Beige,
        false,
        0,
        1,
        false,
        0,
        false,
        tile,
        "Coast",
    }
}

func NewLand(height int, tile Tile) Terrain {
    return Terrain{
        rl.Brown,
        false,
        0,
        height,
        false,
        0,
        false,
        tile,
        "Plains",
    }
}

func NewForest(height int, tile Tile) Terrain {
    return Terrain{
        rl.DarkGreen,
        false,
        -1,
        height,
        false,
        -2,
        true,
        tile,
        "Forest",
    }
}

func NewHill(height int, tile Tile) Terrain {
    return Terrain{
        rl.DarkGray,
        false,
        -1,
        height,
        false,
        1,
        true,
        tile,
        "Hills",
    }
}

func NewMountain(height int, tile Tile) Terrain {
    return Terrain{
        rl.White,
        true,
        -3,
        height,
        false,
        1,
        true,
        tile,
        "Mountains",
    }
}
