package main

import (
    "fmt"
)

func main() {
    var n, m int
    if _, err := fmt.Scan(&n, &m); err != nil {
        return
    }
    // Minimum moves is when using as many 2-steps as possible
    kStart := (n + 1) / 2
    // Find smallest multiple of m >= kStart
    t := (kStart + m - 1) / m * m
    if t <= n {
        fmt.Println(t)
    } else {
        fmt.Println(-1)
    }
}
