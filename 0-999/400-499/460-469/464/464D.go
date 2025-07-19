package main

import (
    "fmt"
)

func main() {
    var n, k int
    if _, err := fmt.Scan(&n, &k); err != nil {
        return
    }
    // Extend slice to length 711 to safely access d[t+1] when t=709
    d := make([]float64, 711)
    for i := 0; i < n; i++ {
        for t := 1; t < 710; t++ {
            term1 := float64(k-1) / float64(k) * d[t]
            coef := 1.0 / (float64(k) * float64(t+1))
            term2 := float64(t)*d[t] + d[t+1] + float64(t*(t+3))/2.0
            d[t] = term1 + coef*term2
        }
    }
    res := d[1] * float64(k)
    fmt.Printf("%.10f\n", res)
}
