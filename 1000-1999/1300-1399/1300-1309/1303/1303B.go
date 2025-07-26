package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    for ; t > 0; t-- {
        var n, g, b int64
        fmt.Fscan(in, &n, &g, &b)
        // minimum high-quality segments needed
        need := (n + 1) / 2
        // number of full good periods before the last
        periods := (need + g - 1) / g
        full := periods - 1
        // remaining high segments in last good period
        rem := need - full*g
        // days to get required high-quality segments
        daysHigh := full*(g+b) + rem
        // answer is max of daysHigh and total segments n
        var ans int64
        if daysHigh < n {
            ans = n
        } else {
            ans = daysHigh
        }
        fmt.Fprintln(out, ans)
    }
}
