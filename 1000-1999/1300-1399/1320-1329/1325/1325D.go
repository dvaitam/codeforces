package main

import (
    "fmt"
)

func main() {
    var u, v int64
    if _, err := fmt.Scan(&u, &v); err != nil {
        return
    }
    c := v - u
    // If v < u or difference is odd, no solution
    if c < 0 || c&1 == 1 {
        fmt.Println(-1)
        return
    }
    // No difference
    if c == 0 {
        if u == 0 {
            fmt.Println(0)
        } else {
            fmt.Println(1)
            fmt.Println(u)
        }
        return
    }
    // Split difference
    c >>= 1
    // If c and u have no common bits, can use two numbers
    if (c & u) == 0 {
        fmt.Println(2)
        fmt.Printf("%d %d\n", c, c^u)
    } else {
        // Otherwise, three numbers
        fmt.Println(3)
        fmt.Printf("%d %d %d\n", c, c, u)
    }
}
