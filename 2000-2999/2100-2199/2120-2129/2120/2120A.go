package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

type rect struct {
    l, b int
}

func canFormSquare(r []rect) bool {
    area := 0
    for _, x := range r {
        area += x.l * x.b
    }
    side := 0
    for side*side < area {
        side++
    }
    if side*side != area {
        return false
    }
    // try arrange by matching widths
    if r[0].l == side && r[0].b+r[1].b+r[2].b == side && r[1].l == side && r[2].l == side {
        return true
    }
    // try stacking horizontally below first rectangle
    if r[0].l+r[1].l+r[2].l == side && r[0].b == side && r[1].b == side && r[2].b == side {
        return true
    }
    return false
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        r := make([]rect, 3)
        for i := 0; i < 3; i++ {
            fmt.Fscan(in, &r[i].l, &r[i].b)
        }
        sort.Slice(r, func(i, j int) bool {
            if r[i].l == r[j].l {
                return r[i].b > r[j].b
            }
            return r[i].l > r[j].l
        })
        if canFormSquare(r) {
            fmt.Fprintln(out, "YES")
        } else {
            fmt.Fprintln(out, "NO")
        }
    }
}
