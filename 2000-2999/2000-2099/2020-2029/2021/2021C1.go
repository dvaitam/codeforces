package main

import (
    "bufio"
    "fmt"
    "os"
)

const inf = int(1e9)

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n, m, q int
        fmt.Fscan(in, &n, &m, &q)

        a := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
        }

        first := make([]int, n+1)
        for i := range first {
            first[i] = inf
        }

        for i := 0; i < m; i++ {
            var x int
            fmt.Fscan(in, &x)
            if first[x] == inf {
                first[x] = i
            }
        }

        prev := -1
        seenMissing := false
        good := true
        for _, val := range a {
            occ := first[val]
            if occ == inf {
                seenMissing = true
            } else {
                if seenMissing || occ < prev {
                    good = false
                    break
                }
                prev = occ
            }
        }

        if good {
            fmt.Fprintln(out, "YA")
        } else {
            fmt.Fprintln(out, "TIDAK")
        }
    }
}
