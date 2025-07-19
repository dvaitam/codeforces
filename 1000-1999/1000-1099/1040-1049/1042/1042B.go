package main

import (
    "bufio"
    "fmt"
    "os"
)

const INF = 1000000000

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n int
    if _, err := fmt.Fscan(reader, &n); err != nil {
        return
    }
    // dp[mask] = minimum cost to get vitamins in mask
    dp := [8]int{}
    for i := 1; i < 8; i++ {
        dp[i] = INF
    }
    for i := 0; i < n; i++ {
        var cost int
        var s string
        fmt.Fscan(reader, &cost, &s)
        mask := 0
        for _, ch := range s {
            switch ch {
            case 'A':
                mask |= 1
            case 'B':
                mask |= 2
            case 'C':
                mask |= 4
            }
        }
        if cost < dp[mask] {
            dp[mask] = cost
        }
    }
    ans := dp[7]
    // consider using two or three packages
    for i := 1; i < 8; i++ {
        for j := 1; j < 8; j++ {
            if i|j == 7 {
                sum := dp[i] + dp[j]
                if sum < ans {
                    ans = sum
                }
            }
            for k := 1; k < 8; k++ {
                if i|j|k == 7 {
                    sum := dp[i] + dp[j] + dp[k]
                    if sum < ans {
                        ans = sum
                    }
                }
            }
        }
    }
    if ans >= INF {
        fmt.Println(-1)
    } else {
        fmt.Println(ans)
    }
}
