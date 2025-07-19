package main

import (
    "fmt"
    "os"
    "sort"
)

func main() {
    in, out := os.Stdin, os.Stdout
    var n int
    if _, err := fmt.Fscan(in, &n); err != nil {
        return
    }
    a := make([]int64, n)
    var sum int64
    for i := 0; i < n; i++ {
        fmt.Fscan(in, &a[i])
        sum += a[i]
    }
    sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
    var m int
    fmt.Fscan(in, &m)
    for i := 0; i < m; i++ {
        var t int
        fmt.Fscan(in, &t)
        idx := n - t
        ans := sum - a[idx]
        fmt.Fprintln(out, ans)
    }
}
