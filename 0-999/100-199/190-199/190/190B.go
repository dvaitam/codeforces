package main

import (
    "fmt"
    "math"
)

func main() {
    var a, b, r1, x, y, r2 float64
    if _, err := fmt.Scan(&a, &b, &r1); err != nil {
        return
    }
    if _, err := fmt.Scan(&x, &y, &r2); err != nil {
        return
    }
    dx := a - x
    dy := b - y
    d := math.Hypot(dx, dy)
    var res float64
    if d > r1+r2 {
        res = (d - r1 - r2) / 2
    } else if d+math.Min(r1, r2) < math.Max(r1, r2) {
        res = (math.Max(r1, r2) - math.Min(r1, r2) - d) / 2
    } else {
        res = 0
    }
    fmt.Printf("%.12f\n", res)
}
