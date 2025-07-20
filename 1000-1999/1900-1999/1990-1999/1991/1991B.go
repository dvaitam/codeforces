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
        b := make([]int, n-1)
        for i := 0; i < n-1; i++ {
            fmt.Fscan(in, &b[i])
        }
        possible := true
        for i := 1; i < n-1 && possible; i++ {
            if ((b[i-1] & b[i+1]) &^ b[i]) != 0 {
                possible = false
            }
        }
        if !possible {
            fmt.Fprintln(out, -1)
            continue
        }
        a := make([]int, n)
        if n == 1 {
            // not reachable since b length is n-1 but handle
            a[0] = 0
        } else {
            a[0] = b[0]
            for i := 1; i < n-1; i++ {
                a[i] = b[i-1] | b[i]
            }
            a[n-1] = b[n-2]
        }
        for i := 0; i < n; i++ {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            fmt.Fprint(out, a[i])
        }
        fmt.Fprintln(out)
    }
}

