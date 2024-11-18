package types

import(
    "image/color"
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
}

func NewSea(colour color.RGBA, tile Tile) Terrain {
    return Terrain{
        colour,
        false,
        0,
        0,
        true,
        1,
        false,
        tile,
    }
}

func NewCoast(colour color.RGBA, tile Tile) Terrain {
    return Terrain{
        colour,
        false,
        0,
        1,
        false,
        0,
        false,
        tile,
    }
}

func NewLand(colour color.RGBA, height int, tile Tile) Terrain {
    return Terrain{
        colour,
        false,
        0,
        height,
        false,
        0,
        false,
        tile,
    }
}

func NewForest(colour color.RGBA, height int, tile Tile) Terrain {
    return Terrain{
        colour,
        false,
        -2,
        height,
        false,
        -2,
        true,
        tile,
    }
}

func NewHill(colour color.RGBA, height int, tile Tile) Terrain {
    return Terrain{
        colour,
        false,
        -1,
        height,
        false,
        1,
        true,
        tile,
    }
}

func NewMountain(colour color.RGBA, height int, tile Tile) Terrain {
    return Terrain{
        colour,
        true,
        -3,
        height,
        false,
        1,
        true,
        tile,
    }
}
