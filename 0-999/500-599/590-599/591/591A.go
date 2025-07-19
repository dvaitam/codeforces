package main

import (
    "fmt"
)

func main() {
    var L, p, q float64
    if _, err := fmt.Scan(&L, &p, &q); err != nil {
        return
    }
    result := p * (L / (p + q))
    // Print with sufficient precision
    fmt.Printf("%.10f\n", result)
}
