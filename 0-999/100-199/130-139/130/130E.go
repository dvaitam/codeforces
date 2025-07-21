package main

import "fmt"

func main() {
    var n int
    if _, err := fmt.Scan(&n); err != nil {
        return
    }
    // Base tribonacci values: t0=0, t1=0, t2=1
    if n == 0 {
        fmt.Println(0)
        return
    }
    // Use a slice to store tribonacci values up to n
    dp := make([]int, n+1)
    dp[0], dp[1] = 0, 0
    if n >= 2 {
        dp[2] = 1
    }
    const mod = 26
    for i := 3; i <= n; i++ {
        dp[i] = (dp[i-1] + dp[i-2] + dp[i-3]) % mod
    }
    fmt.Println(dp[n] % mod)
}
