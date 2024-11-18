package types

import(
    "math"
)

//Cubic coordinate tile (int)
type Tile struct {
    Q int
    R int
    S int
}

//Cubic coordinate tile (fractional)
type FracTile struct {
    Q float64
    R float64
    S float64
}

//Converts a cubic tile to a fractional tile
func (t Tile) ToFractional() FracTile {
    return FracTile{
        float64(t.Q),
        float64(t.R),
        float64(t.S),
    }
}

//Converts a cubic tile to an axial tile
func (t Tile) ToAxial() AxialTile {
    return AxialTile{t.Q, t.R}
}

//Adds a cubic tile (b) to the caller cubic tile (a)
func (a Tile) Add(b Tile) Tile {
    return Tile{a.Q + b.Q, a.R + b.R, a.S + b.S}
}

//Subtracts a cubic tile (b) from the caller cubic tile (a)
func (a Tile) Subtract(b Tile) Tile {
    return Tile{a.Q - b.Q, a.R - b.R, a.S - b.S}
}

//Calculates the distance from cubic tile (a) to cubic tile (b)
func (a Tile) Distance(b Tile) int {
    v := a.Subtract(b)
    return max(Abs(v.Q), Abs(v.R), Abs(v.S))
}

//Gathers a slice of cubic tiles within a step range (n) from cubic tile (a)
//One step is defined as the movement from one cubic tile to one directly adjacent
func (t Tile) Range(n int) []Tile {
    results := []Tile{}
    for q := -n; q <= n; q++ {
        for r := max(-n, -q-n); r <= min(n, -q+n); r++ {
            s := -q-r
            results = append(results, t.Add(Tile{q, r, s}))
        }
    }
    return results
}

//Rounds a fractional tile to the nearest unit tile
func (ft FracTile) Round() Tile {
    q := math.Round(ft.Q)
    r := math.Round(ft.R)
    s := math.Round(ft.S)

    q_diff := math.Abs(q - ft.Q)
    r_diff := math.Abs(r - ft.R)
    s_diff := math.Abs(s - ft.S)

    if(q_diff > r_diff && q_diff > s_diff){
        q = -r-s
    } else if(r_diff > s_diff){
        r = -q-s
    } else {
        s = -q-r
    }

    return Tile{int(q), int(r), int(s)}
}

//Linear interpolation of a fractile at a point (t) between
//one fractile (fa) and another (fb)
func (fa FracTile) Lerp(fb FracTile, t float64) FracTile {
    return FracTile{
        Lerp(fa.Q, fb.Q, t),
        Lerp(fa.R, fb.R, t),
        Lerp(fa.S, fb.S, t),
    }
}

//Returns a list of tiles along the straightest path between
//one tile (a) and another (b)
func (a Tile) LineDraw(b Tile) []Tile {
    N := a.Distance(b)
    results := []Tile{}
    for i := 0; i <= N; i++ {
        results = append(
            results, 
            a.ToFractional().Lerp(
                b.ToFractional(), 
                float64(1/N*i),
            ).Round(),
        )
    }
    return results
}

//Lighterweight cubic tile
type AxialTile struct {
    Q int
    R int
}

//Get the implied S value of the axialtile
func (at AxialTile) S() int {
    return -at.Q -at.R
}

//Convert an axialtile to a cubic tile
func (at AxialTile) ToCubic() Tile {
    return Tile{at.Q, at.R, at.S()}
}

//Lighterweight cubic tile
type FracAxialTile struct {
    Q float64
    R float64
}

//Get the implied S value of the axialtile
func (at FracAxialTile) S() float64 {
    return -at.Q -at.R
}

//Convert an axialtile to a cubic tile
func (at FracAxialTile) ToCubic() FracTile {
    return FracTile{at.Q, at.R, at.S()}
}
