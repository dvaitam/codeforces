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

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n, m int
        fmt.Fscan(in, &n, &m)
        a := make([]int64, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
        }
        var b int64
        if m > 0 {
            fmt.Fscan(in, &b)
            for i := 1; i < m; i++ { // read remaining if any (shouldn't happen in easy version)
                var tmp int64
                fmt.Fscan(in, &tmp)
            }
        }

        prev := int64(-1 << 60)
        possible := true
        for i := 0; i < n; i++ {
            x := a[i]
            y := b - x
            small, large := x, y
            if small > large {
                small, large = large, small
            }
            if small >= prev {
                prev = small
            } else if large >= prev {
                prev = large
            } else {
                possible = false
                break
            }
        }
        if possible {
            fmt.Fprintln(out, "YES")
        } else {
            fmt.Fprintln(out, "NO")
        }
    }
}
