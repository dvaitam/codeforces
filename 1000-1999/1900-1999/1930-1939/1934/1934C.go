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

        fmt.Fprintf(out, "? %d %d\n", 1, 1)
        out.Flush()
        var a int
        fmt.Fscan(in, &a)

        fmt.Fprintf(out, "? %d %d\n", 1, m)
        out.Flush()
        var b int
        fmt.Fscan(in, &b)

        fmt.Fprintf(out, "? %d %d\n", n, 1)
        out.Flush()
        var c int
        fmt.Fscan(in, &c)

        fmt.Fprintf(out, "? %d %d\n", n, m)
        out.Flush()
        var d int
        fmt.Fscan(in, &d)

        x := (a + b - m + 3) / 2
        y := (a + m - b + 1) / 2
        if x < 1 || x > n || y < 1 || y > m {
            x = (2*n + m - d - c - 1) / 2
            y = (m - d + c + 1) / 2
        }

        fmt.Fprintf(out, "! %d %d\n", x, y)
        out.Flush()
        var verdict string
        fmt.Fscan(in, &verdict)
        if verdict != "Ok" && verdict != "ok" {
            return
        }
    }
}

