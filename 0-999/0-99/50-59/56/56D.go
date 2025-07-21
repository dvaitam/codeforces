package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()

    var s, t string
    fmt.Fscan(reader, &s)
    fmt.Fscan(reader, &t)
    n, m := len(s), len(t)

    // dp[i][j]: min operations to convert s[:i] to t[:j]
    dp := make([][]int, n+1)
    for i := 0; i <= n; i++ {
        dp[i] = make([]int, m+1)
    }
    for i := 0; i <= n; i++ {
        dp[i][0] = i
    }
    for j := 0; j <= m; j++ {
        dp[0][j] = j
    }
    for i := 1; i <= n; i++ {
        for j := 1; j <= m; j++ {
            cost := 0
            if s[i-1] != t[j-1] {
                cost = 1
            }
            // replace or match
            dp[i][j] = dp[i-1][j-1] + cost
            // delete
            if dp[i-1][j]+1 < dp[i][j] {
                dp[i][j] = dp[i-1][j] + 1
            }
            // insert
            if dp[i][j-1]+1 < dp[i][j] {
                dp[i][j] = dp[i][j-1] + 1
            }
        }
    }

    // backtrack operations
    type op struct { typ string; pos int; ch byte }
    ops := make([]op, 0)
    i, j := n, m
    for i > 0 || j > 0 {
        if i > 0 && j > 0 && s[i-1] == t[j-1] && dp[i][j] == dp[i-1][j-1] {
            i--
            j--
        } else if i > 0 && j > 0 && dp[i][j] == dp[i-1][j-1]+1 {
            ops = append(ops, op{"REPLACE", i, t[j-1]})
            i--
            j--
        } else if j > 0 && dp[i][j] == dp[i][j-1]+1 {
            // insert at position i+1
            ops = append(ops, op{"INSERT", i + 1, t[j-1]})
            j--
        } else {
            // delete at position i
            ops = append(ops, op{"DELETE", i, 0})
            i--
        }
    }

    // output
    k := len(ops)
    fmt.Fprintln(writer, k)
    for idx := k - 1; idx >= 0; idx-- {
        o := ops[idx]
        switch o.typ {
        case "DELETE":
            fmt.Fprintf(writer, "%s %d\n", o.typ, o.pos)
        case "INSERT", "REPLACE":
            fmt.Fprintf(writer, "%s %d %c\n", o.typ, o.pos, o.ch)
        }
    }
}
