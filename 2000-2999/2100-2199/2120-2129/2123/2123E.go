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
        a := make([]int, n)
        maxVal := n + 1
        freq := make([]int, maxVal)
        zeros := 0
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
            if a[i] < maxVal {
                freq[a[i]]++
            }
            if a[i] == 0 {
                zeros++
            }
        }

        res := make([]int, n+1)
        for k := 0; k <= n; k++ {
            if k <= zeros {
                res[k] = 1
            }
        }

        for i := 0; i <= n; i++ {
            if i+1 <= n {
                res[i+1]++
            }
        }

        for i := 0; i <= n; i++ {
            if i > 0 {
                fmt.Fprint(out, " ")
            }
            if res[i] > 0 {
                fmt.Fprint(out, res[i])
            } else {
                fmt.Fprint(out, 1)
            }
        }
        fmt.Fprintln(out)
    }
}
