package main

import (
    "bufio"
    "fmt"
    "os"
)

type Fenwick struct {
    n   int
    bit []int
}

func NewFenwick(n int) *Fenwick {
    return &Fenwick{n: n, bit: make([]int, n+2)}
}

func (f *Fenwick) Update(idx, val int) {
    for idx <= f.n {
        if f.bit[idx] < val {
            f.bit[idx] = val
        }
        idx += idx & -idx
    }
}

func (f *Fenwick) Query(idx int) int {
    if idx <= 0 {
        return 0
    }
    res := 0
    for idx > 0 {
        if res < f.bit[idx] {
            res = f.bit[idx]
        }
        idx -= idx & -idx
    }
    return res
}

type pair struct {
    start int
    upper int
}

type SegTree struct {
    n    int
    seg  []pair
    lazy []pair
}

func NewSegTree(n int) *SegTree {
    size := 4 * (n + 2)
    return &SegTree{n: n, seg: make([]pair, size), lazy: make([]pair, size)}
}

func better(a, b pair) pair {
    if a.start == 0 {
        return b
    }
    if b.start == 0 {
        return a
    }
    if a.start != b.start {
        if a.start > b.start {
            return a
        }
        return b
    }
    if a.upper == 0 {
        return b
    }
    if b.upper == 0 {
        return a
    }
    if a.upper < b.upper {
        return a
    }
    return b
}

func (st *SegTree) apply(idx int, val pair) {
    if val.start == 0 {
        return
    }
    st.seg[idx] = better(st.seg[idx], val)
    st.lazy[idx] = better(st.lazy[idx], val)
}

func (st *SegTree) push(idx int) {
    if st.lazy[idx].start != 0 {
        st.apply(idx<<1, st.lazy[idx])
        st.apply(idx<<1|1, st.lazy[idx])
        st.lazy[idx] = pair{}
    }
}

func (st *SegTree) update(idx, l, r, ql, qr int, val pair) {
    if ql > r || qr < l {
        return
    }
    if ql <= l && r <= qr {
        st.apply(idx, val)
        return
    }
    st.push(idx)
    mid := (l + r) >> 1
    st.update(idx<<1, l, mid, ql, qr, val)
    st.update(idx<<1|1, mid+1, r, ql, qr, val)
    st.seg[idx] = better(st.seg[idx<<1], st.seg[idx<<1|1])
}

func (st *SegTree) Update(ql, qr int, val pair) {
    if ql > qr {
        return
    }
    st.update(1, 1, st.n, ql, qr, val)
}

func (st *SegTree) query(idx, l, r, pos int) pair {
    if l == r {
        return st.seg[idx]
    }
    st.push(idx)
    mid := (l + r) >> 1
    if pos <= mid {
        return st.query(idx<<1, l, mid, pos)
    }
    return st.query(idx<<1|1, mid+1, r, pos)
}

func (st *SegTree) Query(pos int) pair {
    return st.query(1, 1, st.n, pos)
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var t int
    fmt.Fscan(in, &t)
    for ; t > 0; t-- {
        var n int
        fmt.Fscan(in, &n)
        a := make([]int, n+1)
        for i := 1; i <= n; i++ {
            fmt.Fscan(in, &a[i])
        }

        bitFirst := NewFenwick(n)
        bitQuad := NewFenwick(n)
        seg := NewSegTree(n)
        stack := make([]int, 0, n)
        currentBest := 0
        var ans int64

        for i := 1; i <= n; i++ {
            val := a[i]
            for len(stack) > 0 && a[stack[len(stack)-1]] < val {
                stack = stack[:len(stack)-1]
            }
            leftGreater := 0
            if len(stack) > 0 {
                leftGreater = stack[len(stack)-1]
            }

            bestForR := bitQuad.Query(val - 1)
            if bestForR > currentBest {
                currentBest = bestForR
            }
            ans += int64(currentBest)

            quadPair := seg.Query(val)
            if quadPair.start > 0 {
                bitQuad.Update(quadPair.upper, quadPair.start)
            }

            start3 := bitFirst.Query(val - 1)
            if start3 > 0 {
                lowerVal := a[start3]
                upperVal := val
                if lowerVal+1 <= upperVal-1 {
                    seg.Update(lowerVal+1, upperVal-1, pair{start: start3, upper: upperVal})
                }
            }

            if leftGreater > 0 {
                bitFirst.Update(a[leftGreater], leftGreater)
            }

            stack = append(stack, i)
        }

        fmt.Fprintln(out, ans)
    }
}
