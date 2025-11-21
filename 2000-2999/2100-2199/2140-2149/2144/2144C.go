package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

const MOD = 998244353

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int, n)
        b := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
        }
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &b[i])
        }

        pairs := make([][2]int, n)
        for i := 0; i < n; i++ {
            pairs[i] = [2]int{min(a[i], b[i]), max(a[i], b[i])}
        }
        sort.Slice(pairs, func(i, j int) bool {
            if pairs[i][0] == pairs[j][0] {
                return pairs[i][1] < pairs[j][1]
            }
            return pairs[i][0] < pairs[j][0]
        })

        dp := make([][2]int, n)
        for i := range dp {
            dp[i] = [2]int{0, 0}
        }

        for i := 0; i < n; i++ {
            dp[i][0] = 1
            dp[i][1] = 1
            for j := 0; j < i; j++ {
                if pairs[i][0] >= pairs[j][0] {
                    dp[i][0] = (dp[i][0] + dp[j][0]) % MOD
                }
                if pairs[i][1] >= pairs[j][1] {
                    dp[i][1] = (dp[i][1] + dp[j][1]) % MOD
                }
            }
        }

        ans := 0
        for i := 0; i < n; i++ {
            ans = (ans + dp[i][0]) % MOD
            ans = (ans + dp[i][1]) % MOD
        }
        ans = (ans + 1) % MOD

        fmt.Fprintln(out, ans)
    }
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
