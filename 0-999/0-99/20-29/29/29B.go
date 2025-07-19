package main

import (
    "fmt"
    "math"
)

func main() {
    var L, D, V, G, R float64
    if _, err := fmt.Scanf("%f %f %f %f %f", &L, &D, &V, &G, &R); err != nil {
        return
    }
    T := D / V
    cycle := G + R
    T = math.Mod(T, cycle)
    var ans float64
    if T >= G {
        ans += cycle - T
    }
    ans += L / V
    if D >= L {
        ans = L / V
    }
    fmt.Printf("%f\n", ans)
}
