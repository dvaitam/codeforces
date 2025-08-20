package main

import (
    "fmt"
    "sort"
)

func abs(x int) int { if x < 0 { return -x }; return x }

const infinity = 1 << 30

func solve() {
    var n int
    fmt.Scan(&n)
    t := make([]int, n)
    for i := 0; i < n; i++ { fmt.Scan(&t[i]) }
    sort.Ints(t)

    maxTime := 2*n + 5
    if n == 0 { maxTime = 5 }
    dp := make([][]int, n+1)
    for i := range dp {
        dp[i] = make([]int, maxTime)
        for j := range dp[i] { dp[i][j] = infinity }
    }
    for j := 0; j < maxTime; j++ { dp[0][j] = 0 }
    for i := 1; i <= n; i++ {
        for j := 1; j < maxTime; j++ {
            val1 := dp[i][j-1]
            val2 := dp[i-1][j-1] + abs(t[i-1]-j)
            if val1 < val2 { dp[i][j] = val1 } else { dp[i][j] = val2 }
        }
    }
    fmt.Println(dp[n][maxTime-1])
}

func main() {
    var q int
    fmt.Scan(&q)
    for i := 0; i < q; i++ { solve() }
}
