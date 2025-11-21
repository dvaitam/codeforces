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
        var a, b, c int64
        fmt.Fscan(in, &a, &b, &c)
        sum := a + b + c
        if sum%3 != 0 {
            fmt.Fprintln(out, "NO")
            continue
        }
        x := sum / 3
        if x < a || x < b || x >= c {
            fmt.Fprintln(out, "NO")
            continue
        }
        fmt.Fprintln(out, "YES")
    }
}
