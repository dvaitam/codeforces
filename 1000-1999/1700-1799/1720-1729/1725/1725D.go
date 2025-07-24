package main

import (
    "bufio"
    "fmt"
    "os"
)

// TODO: proper implementation
func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, q int
    if _, err := fmt.Fscan(in, &n, &q); err != nil {
        return
    }

    for i := 0; i < q; i++ {
        var p int
        fmt.Fscan(in, &p)
    }
    fmt.Fprintln(out, 0)
    for i := 0; i < q; i++ {
        fmt.Fprintln(out, 0)
    }
}

