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
        var x, n int
        fmt.Fscan(in, &x, &n)
        if n%2 == 0 {
            fmt.Fprintln(out, 0)
        } else {
            fmt.Fprintln(out, x)
        }
    }
}
