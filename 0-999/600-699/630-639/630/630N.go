package main

import (
    "fmt"
    "math"
)

func main() {
    var a, b, c int
    if _, err := fmt.Scan(&a, &b, &c); err != nil {
        return
    }
    // compute discriminant
    da := float64(b*b - 4*a*c)
    d := math.Sqrt(da)
    af := float64(a)
    bf := float64(b)
    // roots
    x1 := (-bf + d) / (2 * af)
    x2 := (-bf - d) / (2 * af)
    // output larger first, then smaller
    if x1 < x2 {
        x1, x2 = x2, x1
    }
    fmt.Printf("%.6f\n%.6f", x1, x2)
}
