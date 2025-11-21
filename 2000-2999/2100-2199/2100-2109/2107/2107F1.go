package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    const inf int64 = 1<<60

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
        }
        type Pair struct {
            val int
            idx int
        }
        ord := make([]Pair, n)
        for i := 0; i < n; i++ {
            ord[i] = Pair{a[i], i}
        }
        sort.Slice(ord, func(i, j int) bool {
            if ord[i].val == ord[j].val {
                return ord[i].idx < ord[j].idx
            }
            return ord[i].val < ord[j].val
        })
        dp := make([]int64, n)
        for i := 0; i < n; i++ {
            if i == 0 {
                dp[i] = int64(ord[i].val)
            } else {
                dp[i] = dp[i-1] + int64(ord[i].val)
            }
        }
        fmt.Fprintln(out, dp[n-1])
    }
}
