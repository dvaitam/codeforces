package main

import (
    "bufio"
    "fmt"
    "os"
)

func main() {
    in := bufio.NewReader(os.Stdin)
    var t, p string
    if _, err := fmt.Fscan(in, &t); err != nil {
        return
    }
    fmt.Fscan(in, &p)
    n := len(t)
    pos := make([]int, n)
    for i := 0; i < n; i++ {
        var x int
        fmt.Fscan(in, &x)
        pos[x-1] = i + 1
    }
    can := func(k int) bool {
        j := 0
        m := len(p)
        for i := 0; i < n && j < m; i++ {
            if pos[i] <= k {
                continue
            }
            if t[i] == p[j] {
                j++
            }
        }
        return j == m
    }
    l, r := 0, n
    for l < r {
        mid := (l + r + 1) / 2
        if can(mid) {
            l = mid
        } else {
            r = mid - 1
        }
    }
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()
    fmt.Fprintln(out, l)
}
