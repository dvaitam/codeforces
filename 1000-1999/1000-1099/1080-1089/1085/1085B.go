package main

import (
    "fmt"
)

func main() {
    var n, k int64
    if _, err := fmt.Scan(&n, &k); err != nil {
        return
    }
    const inf = int64(4e18)
    res := inf
    // Looking for positive integers d, r such that d * r = n and r < k
    // x = k*d + r, minimize x
    for r := int64(1); r < k; r++ {
        if n%r == 0 {
            d := n / r
            x := k*d + r
            if x < res {
                res = x
            }
        }
    }
    fmt.Println(res)
}
