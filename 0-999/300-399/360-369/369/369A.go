package main

import (
    "fmt"
)

func main() {
    var n, m, k int
    // n: number of days, m: clean bowls, k: clean plates
    if _, err := fmt.Scan(&n, &m, &k); err != nil {
        return
    }
    washes := 0
    for i := 0; i < n; i++ {
        var dish int
        fmt.Scan(&dish)
        switch dish {
        case 1:
            if m > 0 {
                m--
            } else {
                washes++
            }
        case 2:
            if k > 0 {
                k--
            } else if m > 0 {
                m--
            } else {
                washes++
            }
        }
    }
    fmt.Println(washes)
}
