package main

import "fmt"

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main() {
    var n, k int
    if _, err := fmt.Scan(&n, &k); err != nil {
        return
    }
    // Minimum moves: 3*n + min(k-1, n-k)
    ans := 3*n + min(k-1, n-k)
    fmt.Println(ans)
}
