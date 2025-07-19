package main

import (
    "fmt"
)

func main() {
    var a, b int64
    if _, err := fmt.Scan(&a, &b); err != nil {
        return
    }
    if a < b {
        fmt.Println(-1)
        return
    }
    // Compute k = floor((a+b)/(2*b))
    k := (a + b) / (2 * b)
    // Answer x = (a + b) / (2 * k)
    ans := float64(a+b) / (2 * float64(k))
    fmt.Printf("%.12f\n", ans)
}
