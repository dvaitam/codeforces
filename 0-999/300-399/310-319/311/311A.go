package main

import (
    "fmt"
)

func main() {
    var n int
    var k int64
    if _, err := fmt.Scan(&n, &k); err != nil {
        return
    }
    // total number of point pairs: n*(n-1)/2
    r := int64(n) * int64(n-1) / 2
    if r <= k {
        fmt.Println("no solution")
    } else {
        // output n distinct points on a line
        for i := 0; i < n; i++ {
            fmt.Println(0, i)
        }
    }
}
