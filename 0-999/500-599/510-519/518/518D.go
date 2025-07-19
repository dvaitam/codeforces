package main

import (
    "fmt"
)

func main() {
    var n, t int
    var p float64
    // Read n, probability p, and number of trials t
    if _, err := fmt.Scan(&n, &p, &t); err != nil {
        return
    }
    // dp[i] is the probability of state i
    dp := make([]float64, n+1)
    dp[0] = 1.0
    for time := 1; time <= t; time++ {
        // Update dp from highest state to lowest
        dp[n] += p * dp[n-1]
        for num := n - 1; num > 0; num-- {
            dp[num] = p*dp[num-1] + (1-p)*dp[num]
        }
        dp[0] *= (1 - p)
    }
    // Compute expected value
    var ans float64
    for num := 0; num <= n; num++ {
        ans += float64(num) * dp[num]
    }
    // Print with 10 decimal places
    fmt.Printf("%.10f\n", ans)
}
