package main
import (
    "bufio"
    "fmt"
    "os"
)
// Segment tree with affine updates and range sum queries
type Node struct {
    x, t, s float64
    sz      int
}
var (
    tree []Node
    D    []float64
    N    int
)
// apply update: multiply by tt and add ts
func (nd *Node) upd(tt, ts float64) {
    nd.x = tt*nd.x + ts*float64(nd.sz)
    nd.t *= tt
    nd.s = nd.s*tt + ts
}
// build tree over [l, r)
func build(i, l, r int) {
    tree[i].t = 1
    tree[i].s = 0
    tree[i].sz = r - l
    if r-l == 1 {
        tree[i].x = D[l]
        return
    }
    m := (l + r) >> 1
    build(i<<1, l, m)
    build(i<<1|1, m, r)
    tree[i].x = tree[i<<1].x + tree[i<<1|1].x
}
// push lazy to children
func push(i int) {
    tt := tree[i].t
    ts := tree[i].s
    if tt != 1 || ts != 0 {
        tree[i<<1].upd(tt, ts)
        tree[i<<1|1].upd(tt, ts)
        tree[i].t = 1
        tree[i].s = 0
    }
}
// range sum query on [ql, qr)
func query(i, l, r, ql, qr int) float64 {
    if ql <= l && r <= qr {
        return tree[i].x
    }
    push(i)
    m := (l + r) >> 1
    var ans float64
    if ql < m && qr > l {
        ans += query(i<<1, l, m, ql, qr)
    }
    if qr > m && ql < r {
        ans += query(i<<1|1, m, r, ql, qr)
    }
    return ans
}
// apply affine update on [ql, qr): x->tt*x+ts
func modify(i, l, r, ql, qr int, tt, ts float64) {
    if ql <= l && r <= qr {
        tree[i].upd(tt, ts)
        return
    }
    push(i)
    m := (l + r) >> 1
    if ql < m && qr > l {
        modify(i<<1, l, m, ql, qr, tt, ts)
    }
    if qr > m && ql < r {
        modify(i<<1|1, m, r, ql, qr, tt, ts)
    }
    tree[i].x = tree[i<<1].x + tree[i<<1|1].x
}
func main() {
    reader := bufio.NewReader(os.Stdin)
    writer := bufio.NewWriter(os.Stdout)
    defer writer.Flush()
    var Q int
    fmt.Fscan(reader, &N, &Q)
    D = make([]float64, N)
    for i := 0; i < N; i++ {
        var v int
        fmt.Fscan(reader, &v)
        D[i] = float64(v)
    }
    tree = make([]Node, 4*N+5)
    build(1, 0, N)
    for Q > 0 {
        Q--
        var op, a, b int
        fmt.Fscan(reader, &op, &a, &b)
        if op == 1 {
            var c, d int
            fmt.Fscan(reader, &c, &d)
            a--
            c--
            b0 := b
            d0 := d
            len1 := b0 - a
            len2 := d0 - c
            sum2 := query(1, 0, N, c, d0)
            sum1 := query(1, 0, N, a, b0)
            s1 := sum2 / float64(len2) / float64(len1)
            s2 := sum1 / float64(len2) / float64(len1)
            t1 := 1 - 1/float64(len1)
            t2 := 1 - 1/float64(len2)
            modify(1, 0, N, a, b0, t1, s1)
            modify(1, 0, N, c, d0, t2, s2)
        } else {
            a--
            sum := query(1, 0, N, a, b)
            fmt.Fprintf(writer, "%.7f\n", sum)
        }
    }
}
