package types

func Abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}

func Lerp(a float64, b float64, t float64) float64 {
    return a + (b - a) * t
}
