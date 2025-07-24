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
        var n, d int64
        fmt.Fscan(in, &n, &d)
        possible := false
        for x := int64(0); x*x <= d; x++ {
            val := x + (d+x)/(x+1)
            if val <= n {
                possible = true
                break
            }
        }
        if possible {
            fmt.Fprintln(out, "YES")
        } else {
            fmt.Fprintln(out, "NO")
        }
    }
}
