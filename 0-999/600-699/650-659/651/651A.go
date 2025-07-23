package main

import (
    "fmt"
)

func main() {
    var a, b int
    if _, err := fmt.Scan(&a, &b); err != nil {
        return
    }
    minutes := 0
    for a > 0 && b > 0 {
        if a == 1 && b == 1 {
            break
        }
        if a < b {
            a, b = b, a
        }
        a -= 2
        b += 1
        minutes++
    }
    fmt.Println(minutes)
}
