package main

import (
    "bufio"
    "fmt"
    "os"
    "sort"
)

const mod = 1000000007

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var T int
    fmt.Fscan(in, &T)
    for ; T > 0; T-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int, n)
        vals := make([]int, n)
        for i := 0; i < n; i++ {
            fmt.Fscan(in, &a[i])
            vals[i] = a[i]
        }
        sort.Ints(vals)
        vals = unique(vals)
        m := len(vals)
        bit1 := make([]int64, m+2)
        bit2 := make([]int64, m+2)
        ans := int64(1)
        for i := 0; i < n; i++ {
            p := lowerBound(vals, a[i])
            c1 := (sum(bit1, p) + 1) % mod
            c2 := (sum(bit2, p) + sum(bit1, p-1)) % mod
            ans = (ans + c2) % mod
            add(bit1, p, c1)
            add(bit2, p, c2)
        }
        fmt.Fprintln(out, ans%mod)
    }
}

func unique(a []int) []int {
    if len(a) == 0 {
        return a
    }
    idx := 1
    for i := 1; i < len(a); i++ {
        if a[i] != a[idx-1] {
            a[idx] = a[i]
            idx++
        }
    }
    return a[:idx]
}

func lowerBound(a []int, x int) int {
    l, r := 0, len(a)
    for l < r {
        mid := (l + r) >> 1
        if a[mid] < x {
            l = mid + 1
        } else {
            r = mid
        }
    }
    return l
}

func add(bit []int64, idx int, val int64) {
    idx++
    for idx < len(bit) {
        bit[idx] = (bit[idx] + val) % mod
        idx += idx & -idx
    }
}

func sum(bit []int64, idx int) int64 {
    if idx < 0 {
        return 0
    }
    idx++
    res := int64(0)
    for idx > 0 {
        res = (res + bit[idx]) % mod
        idx -= idx & -idx
    }
    return res
}
