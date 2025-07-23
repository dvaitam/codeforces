package main

import (
    "fmt"
)

func main() {
    var n, x int
    if _, err := fmt.Scan(&n, &x); err != nil {
        return
    }
    count := 0
    for i := 1; i <= n; i++ {
        if x%i == 0 && x/i <= n {
            count++
        }
    }
    fmt.Println(count)
}

