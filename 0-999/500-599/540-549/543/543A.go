package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    var n, m, b, mod int
    fmt.Fscan(reader, &n, &m, &b, &mod)
    a := make([]int, n)
    for i := 0; i < n; i++ {
        fmt.Fscan(reader, &a[i])
    }
    // dp[j][k]: ways to write j lines with k bugs
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, b+1)
    }
    dp[0][0] = 1
    for i := 0; i < n; i++ {
        ai := a[i]
        for j := 1; j <= m; j++ {
            for k := ai; k <= b; k++ {
                dp[j][k] += dp[j-1][k-ai]
                if dp[j][k] >= mod {
                    dp[j][k] -= mod
                }
            }
        }
    }
    result := 0
    for k := 0; k <= b; k++ {
        result += dp[m][k]
        if result >= mod {
            result -= mod
        }
    }
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    fmt.Fprintln(writer, result)
}
