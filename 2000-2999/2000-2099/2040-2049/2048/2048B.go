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
        var n, k int
        fmt.Fscan(in, &n, &k)
        perm := make([]int, 0, n)
        for start := k; start > 0; start-- {
            for i := start; i <= n; i += k {
                perm = append(perm, i)
            }
        }
        for i := 0; i < n; i++ {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, perm[i])
        }
        fmt.Fprintln(out)
    }
}
