package main

import (
    "fmt"
)

func main() {
    var n, k int
    if _, err := fmt.Scan(&n, &k); err != nil {
        return
    }
    for i := 1; i <= 2*n; i += 2 {
        if k > 0 {
            fmt.Printf("%d %d ", i, i+1)
            k--
        } else {
            fmt.Printf("%d %d ", i+1, i)
        }
    }
    fmt.Println()
}
