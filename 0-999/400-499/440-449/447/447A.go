package main

import (
    "fmt"
)

func main() {
    var p, n int
    if _, err := fmt.Scan(&p, &n); err != nil {
        return
    }
    used := make([]bool, p)
    for i := 1; i <= n; i++ {
        var x int
        if _, err := fmt.Scan(&x); err != nil {
            return
        }
        h := x % p
        if used[h] {
            fmt.Println(i)
            return
        }
        used[h] = true
    }
    fmt.Println(-1)
}
