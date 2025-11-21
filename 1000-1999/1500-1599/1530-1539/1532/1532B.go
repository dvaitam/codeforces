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
    for i := 0; i < t; i++ {
        var a, b, k int64
        fmt.Fscan(in, &a, &b, &k)
        pairs := k / 2
        res := (a - b) * pairs
        if k%2 == 1 {
            res += a
        }
        fmt.Fprintln(out, res)
    }
}
