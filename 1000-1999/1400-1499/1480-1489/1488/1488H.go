package main

import (
    "bufio"
    "fmt"
    "os"
)

const mod int64 = 998244353

type Mat [4][4]int64

func matMul(a, b Mat) Mat {
    var c Mat
    for i := 0; i < 4; i++ {
        for k := 0; k < 4; k++ {
            if a[i][k] == 0 {
                continue
            }
            for j := 0; j < 4; j++ {
                c[i][j] = (c[i][j] + a[i][k]*b[k][j]) % mod
            }
        }
    }
    return c
}

func identityMat() Mat {
    var m Mat
    for i := 0; i < 4; i++ {
        m[i][i] = 1
    }
    return m
}

var lessMat, geMat Mat

func init() {
    for i := 0; i < 4; i++ {
        for j := 0; j < 4; j++ {
            if i < j {
                lessMat[i][j] = 1
            }
            if i >= j {
                geMat[i][j] = 1
            }
        }
    }
}

type SegTree struct {
    n    int
    size int
    tree []Mat
}

func NewSegTree(a []int) *SegTree {
    n := len(a)
    size := 1
    for size < n {
        size <<= 1
    }
    tree := make([]Mat, 2*size)
    for i := 0; i < n; i++ {
        if a[i] == 1 {
            tree[size+i] = lessMat
        } else {
            tree[size+i] = geMat
        }
    }
    id := identityMat()
    for i := n; i < size; i++ {
        tree[size+i] = id
    }
    for i := size - 1; i > 0; i-- {
        tree[i] = matMul(tree[i<<1], tree[i<<1|1])
    }
    return &SegTree{n: n, size: size, tree: tree}
}

func (st *SegTree) Update(pos, val int) {
    idx := st.size + pos
    if val == 1 {
        st.tree[idx] = lessMat
    } else {
        st.tree[idx] = geMat
    }
    for idx >>= 1; idx > 0; idx >>= 1 {
        st.tree[idx] = matMul(st.tree[idx<<1], st.tree[idx<<1|1])
    }
}

func (st *SegTree) Root() Mat {
    return st.tree[1]
}

func main() {
    in := bufio.NewReader(os.Stdin)
    out := bufio.NewWriter(os.Stdout)
    defer out.Flush()

    var n, q int
    if _, err := fmt.Fscan(in, &n, &q); err != nil {
        return
    }
    a := make([]int, n-1)
    for i := 0; i < n-1; i++ {
        fmt.Fscan(in, &a[i])
    }

    st := NewSegTree(a)

    for ; q > 0; q-- {
        var idx int
        fmt.Fscan(in, &idx)
        idx--
        a[idx] ^= 1
        st.Update(idx, a[idx])
        m := st.Root()
        var total int64
        for i := 0; i < 4; i++ {
            for j := 0; j < 4; j++ {
                total += m[i][j]
            }
        }
        total %= mod
        fmt.Fprintln(out, total)
    }
}

