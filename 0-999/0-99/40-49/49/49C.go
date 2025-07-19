package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    // Print the original number first
    fmt.Printf("%d ", n)
    // Print numbers from 1 to n-1
    for i := 1; i < n; i++ {
        fmt.Printf("%d ", i)
    }
}
