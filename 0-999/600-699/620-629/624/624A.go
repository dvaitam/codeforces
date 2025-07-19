package main

import (
    "fmt"
)

func main() {
    var d, l, v1, v2 float64
    if _, err := fmt.Scan(&d, &l, &v1, &v2); err != nil {
        return
    }
    res := (l - d) / (v1 + v2)
    fmt.Printf("%.10f\n", res)
}
