package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    if n%2 != 0 {
        fmt.Printf("%d %d", n-9, 9)
    } else {
        fmt.Printf("%d %d", n-8, 8)
    }
}
