package main

import (
    "fmt"
)

func main() {
    var n, m int
    if _, err := fmt.Scan(&n, &m); err != nil {
        return
    }
    minR, minC := n, m
    maxR, maxC := -1, -1
    for i := 0; i < n; i++ {
        var s string
        fmt.Scan(&s)
        for j, c := range s {
            if c == 'B' {
                if i < minR {
                    minR = i
                }
                if i > maxR {
                    maxR = i
                }
                if j < minC {
                    minC = j
                }
                if j > maxC {
                    maxC = j
                }
            }
        }
    }
    // Compute center
    row := (minR + maxR) / 2
    col := (minC + maxC) / 2
    // Output 1-based indices
    fmt.Println(row+1, col+1)
}
