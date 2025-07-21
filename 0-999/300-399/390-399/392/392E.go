package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var n int
    fmt.Fscan(in, &n)
    v := make([]int, n+1)
    for i := 1; i <= n; i++ {
        fmt.Fscan(in, &v[i])
    }
    w := make([]int, n+1)
    for i := 1; i <= n; i++ {
        fmt.Fscan(in, &w[i])
    }
    // Precompute which original intervals are mountain-contiguous
    good := make([][]bool, n+1)
    for i := 0; i <= n; i++ {
        good[i] = make([]bool, n+1)
    }
    for i := 1; i <= n; i++ {
        good[i][i] = true
        for j := i + 1; j <= n; j++ {
            ok := true
            for k := i; k < j; k++ {
                if abs(w[k]-w[k+1]) != 1 {
                    ok = false
                    break
                }
            }
            if ok {
                for k := i + 1; k < j; k++ {
                    if 2*w[k] < w[k-1]+w[k+1] {
                        ok = false
                        break
                    }
                }
            }
            good[i][j] = ok
        }
    }
    // dp[l][r]: max score to delete all in [l..r]
    dp := make([][]int, n+2)
    for i := range dp {
        dp[i] = make([]int, n+2)
    }
    // len=1..n
    for length := 1; length <= n; length++ {
        for l := 1; l+length-1 <= n; l++ {
            r := l + length - 1
            best := 0
            // split
            for k := l; k < r; k++ {
                s := dp[l][k] + dp[k+1][r]
                if s > best {
                    best = s
                }
            }
            // delete full segment
            if good[l][r] {
                if v[length] > best {
                    best = v[length]
                }
            }
            dp[l][r] = best
        }
    }
    fmt.Println(dp[1][n])
}

func abs(x int) int {
    if x < 0 {
        return -x
    }
    return x
}
