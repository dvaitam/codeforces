package main

import (
    "fmt"
)

func main() {
    var n, m int
    if _, err := fmt.Scan(&n, &m); err != nil {
        return
    }
    min := n
    if m < n {
        min = m
    }
    // Number of pairs is min+1
    fmt.Println(min + 1)
    for i := 0; i <= min; i++ {
        // Each pair (i, min-i)
        fmt.Printf("%d %d", i, min-i)
        if i < min {
            fmt.Println()
        }
    }
}
