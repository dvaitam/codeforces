package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var s string
    fmt.Fscan(reader, &s)
    n := len(s)
    // Compute odd palindrome radii
    rad := make([]int, n)
    for i := 0; i < n; i++ {
        l, r := i, i
        for l-1 >= 0 && r+1 < n && s[l-1] == s[r+1] {
            l--
            r++
        }
        rad[i] = (r - l) / 2
    }
    // dp[l][r] Grundy for s[l:r+1]
    dp := make([][]uint16, n)
    for i := range dp {
        dp[i] = make([]uint16, n)
    }
    // temporary seen for mex
    const maxG = 512
    seen := make([]int, maxG)
    iter := 1
    // DP by increasing length
    for length := 1; length <= n; length++ {
        for l := 0; l+length-1 < n; l++ {
            r := l + length - 1
            var g uint16 = 0
            if length >= 3 {
                // collect reachable
                for k := l + 1; k <= r-1; k++ {
                    if rad[k] > 0 {
                        x := dp[l][k-1] ^ dp[k+1][r]
                        if int(x) < maxG {
                            seen[x] = iter
                        }
                    }
                }
                // mex
                for j := 0; j < maxG; j++ {
                    if seen[j] != iter {
                        g = uint16(j)
                        break
                    }
                }
                iter++
            }
            dp[l][r] = g
        }
    }
    full := dp[0][n-1]
    if full == 0 {
        fmt.Println("Second")
        return
    }
    fmt.Println("First")
    // find minimal winning move
    for i := 1; i+1 < n; i++ {
        if rad[i] > 0 {
            if dp[0][i-1]^dp[i+1][n-1] == 0 {
                fmt.Println(i + 1)
                return
            }
        }
    }
}
