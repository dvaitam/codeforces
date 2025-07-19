package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var x, k, p int
    if _, err := fmt.Fscan(reader, &x, &k, &p); err != nil {
        return
    }

    pd := float64(p) / 100.0
    pp := float64(100-p) / 100.0

    const MAX = 310
    dp := make([][]float64, 2)
    dp[0] = make([]float64, MAX)
    dp[1] = make([]float64, MAX)
    dp[0][0] = 1.0
    next := 1
    ans := 0.0

    for i := 0; i < k; i++ {
        prev := 1 - next
        for j := 0; j < MAX; j++ {
            dp[next][j] = 0.0
        }
        for j := 0; j < 300; j++ {
            dp[next][j+1] += dp[prev][j] * pp
        }
        for j := 0; j < 300; j += 2 {
            ans += dp[prev][j] * pd
            dp[next][j/2] += dp[prev][j] * pd
        }
        next = prev
    }
    prev := 1 - next
    for i := 0; i < 300; i++ {
        r := i + x
        for r%2 == 0 {
            r /= 2
            ans += dp[prev][i]
        }
    }

    fmt.Printf("%f\n", ans)
}
