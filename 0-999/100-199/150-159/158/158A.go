package main

import (
    "fmt"
)

func main() {
    var n, k int
    if _, err := fmt.Scan(&n, &k); err != nil {
        return
    }
    scores := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Scan(&scores[i])
    }
    // Determine the k-th place score (1-indexed)
    threshold := scores[k-1]
    // Count participants with score >= threshold and positive
    count := 0
    for _, s := range scores {
        if s >= threshold && s > 0 {
            count++
        }
    }
    fmt.Println(count)
}
