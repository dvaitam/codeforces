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
    for ; t > 0; t-- {
        var n, m int
        fmt.Fscan(in, &n, &m)
        total := 0
        count := 0
        for i := 0; i < n; i++ {
            var s string
            fmt.Fscan(in, &s)
            if total+len(s) <= m {
                total += len(s)
                count = i + 1
            }
        }
        fmt.Fprintln(out, count)
    }
}
