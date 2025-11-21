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
        var w string
        fmt.Fscan(in, &w)
        n := len(w)
        if n < 2 {
            fmt.Fprintln(out, "")
            continue
        }
        root := w[:n-2]
        fmt.Fprintln(out, root+"i")
    }
}
