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
        var n int
        fmt.Fscan(in, &n)
        fmt.Fprintln(out, 1)
        for i := 0; i < n; i++ {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, 1)
        }
        fmt.Fprintln(out)
    }
}
