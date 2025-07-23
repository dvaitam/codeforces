package main

import (
    "fmt"
)

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    var res uint64 = 1
    for i := 2; i <= n; i++ {
        res *= uint64(i)
    }
    fmt.Println(res)
}
