package main

import (
    "bufio"
    "fmt"
    "os"
)

const mod = 1000000007

func modPow(base, exp int64) int64 {
    res := int64(1)
    b := base % mod
    for exp > 0 {
        if exp&1 == 1 {
            res = res * b % mod
        }
        b = b * b % mod
        exp >>= 1
    }
    return res
}

type dsu struct {
    parent []int
    xorVal []int
    size   []int
    comps  int
}

func newDSU(n int) *dsu {
    parent := make([]int, n)
    xorVal := make([]int, n)
    size := make([]int, n)
    for i := 0; i < n; i++ {
        parent[i] = i
        xorVal[i] = 0
        size[i] = 1
    }
    return &dsu{parent: parent, xorVal: xorVal, size: size, comps: n}
}

func (d *dsu) find(x int) (int, int) {
    if d.parent[x] == x {
        return x, 0
    }
    root, xr := d.find(d.parent[x])
    d.xorVal[x] ^= xr
    d.parent[x] = root
    return root, d.xorVal[x]
}

func (d *dsu) unite(x, y, val int) bool {
    rx, vx := d.find(x)
    ry, vy := d.find(y)
    if rx == ry {
        return (vx ^ vy) == val
    }
    if d.size[rx] < d.size[ry] {
        rx, ry = ry, rx
        vx, vy = vy, vx
    }
    d.parent[ry] = rx
    d.xorVal[ry] = vx ^ vy ^ val
    d.size[rx] += d.size[ry]
    d.comps--
    return true
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n, m, k, q int
        fmt.Fscan(in, &n, &m, &k, &q)
        total := n + m
        d := newDSU(total)
        ok := true
        for i := 0; i < k; i++ {
            var r, c int
            var v int
            fmt.Fscan(in, &r, &c, &v)
            if ok {
                nodeR := r - 1
                nodeC := n + c - 1
                if !d.unite(nodeR, nodeC, v) {
                    ok = false
                }
            }
        }
        for i := 0; i < q; i++ {
            var r, c, v int
            fmt.Fscan(in, &r, &c, &v)
            // q = 0 in this version, but consume input just in case
            if ok {
                nodeR := r - 1
                nodeC := n + c - 1
                if !d.unite(nodeR, nodeC, v) {
                    ok = false
                }
            }
        }
        if !ok {
            fmt.Fprintln(out, 0)
            continue
        }
        dof := d.comps - 1
        if dof < 0 {
            dof = 0
        }
        ans := modPow(2, int64(dof)*30)
        fmt.Fprintln(out, ans)
    }
}
