package main

import "fmt"

func main() {
    const LMT = 55
    var n, p int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    a := make([]int, n+1)
    for i := 1; i <= n; i++ {
        fmt.Scan(&a[i])
    }
    fmt.Scan(&p)
    // dp[j][t]: number of ways to choose j elements summing to t
    var dp [LMT][LMT]float64
    dp[0][0] = 1.0
    for i := 1; i <= n; i++ {
        ai := a[i]
        for j := n; j >= 1; j-- {
            for t := p; t >= ai; t-- {
                dp[j][t] += dp[j-1][t-ai]
            }
        }
    }
    // factorials
    fi := make([]float64, LMT)
    fi[0] = 1.0
    for i := 1; i < LMT; i++ {
        fi[i] = float64(i) * fi[i-1]
    }
    // accumulate answer
    var ans float64
    for i := 1; i <= n; i++ {
        for t := 1; t <= p; t++ {
            ans += dp[i][t] * fi[i] * fi[n-i]
        }
    }
    // normalize by n!
    fmt.Printf("%.6f\n", ans/fi[n])
}
