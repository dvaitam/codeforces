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
    fmt.Fscan(in, &t)
    for t > 0 {
        t--
        var n int64
        fmt.Fscan(in, &n)
        n0 := n
        var p1, p2, p3 int64
        // find first factor p1
        for i := int64(2); i*i <= n0; i++ {
            if n0%i == 0 {
                p1 = i
                break
            }
        }
        if p1 == 0 {
            fmt.Fprintln(out, "NO")
            continue
        }
        // reduce and find second factor p2
        n1 := n0 / p1
        for j := p1 + 1; j*j <= n1; j++ {
            if n1%j == 0 {
                p2 = j
                break
            }
        }
        if p2 == 0 {
            fmt.Fprintln(out, "NO")
            continue
        }
        // third factor is remaining
        p3 = n1 / p2
        if p3 <= 1 || p3 == p1 || p3 == p2 {
            fmt.Fprintln(out, "NO")
        } else {
            fmt.Fprintln(out, "YES")
            fmt.Fprintln(out, p1, p2, p3)
        }
    }
}
