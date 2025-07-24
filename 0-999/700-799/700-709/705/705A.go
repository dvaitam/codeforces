package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    for i := 1; i <= n; i++ {
        if i%2 == 1 {
            fmt.Print("I hate")
        } else {
            fmt.Print("I love")
        }
        if i < n {
            fmt.Print(" that ")
        } else {
            fmt.Print(" it")
        }
    }
    fmt.Println()
}
