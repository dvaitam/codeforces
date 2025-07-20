package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const inf = 2_000_000_005 // 2 e 9 + 5

// ---------- tiny helpers ----------
func min(a, b int) int { if a < b { return a }; return b }
func max(a, b int) int { if a > b { return a }; return b }

// ---------- matrix type & “multiplication” ----------
type Mat struct {
	a [2][2]int
}

func newMat() Mat {
	var m Mat
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			m.a[i][j] = inf
		}
	}
	return m
}

var ident = Mat{a: [2][2]int{{0, inf}, {inf, 0}}}

func mul(u, v Mat) Mat {
	w := newMat()
	for i := 0; i < 2; i++ {
		for j := 0; j < 2; j++ {
			w.a[i][j] = min(
				max(u.a[i][0], v.a[0][j]),
				max(u.a[i][1], v.a[1][j]),
			)
		}
	}
	return w
}

// ---------- segment‑tree ----------
type SegTree struct {
	N  int
	tr []Mat // 1‑based: indices [1 … 2N]
}

func NewSegTree(n int, leaf []Mat) *SegTree {
	N := 1
	for N <= n {
		N <<= 1
	}
	st := &SegTree{
		N:  N,
		tr: make([]Mat, 2*N+2), // +2 = safe guard
	}
	// fill with identity
	for i := 1; i < 2*N; i++ {
		st.tr[i] = ident
	}
	// copy leaves
	for i := 1; i <= n; i++ {
		st.tr[i+N] = leaf[i]
	}
	// build
	for i := N - 1; i >= 1; i-- {
		st.tr[i] = mul(st.tr[i<<1], st.tr[i<<1|1])
	}
	return st
}

func (st *SegTree) Update(pos int, val Mat) {
	idx := pos + st.N
	st.tr[idx] = val
	for idx >>= 1; idx > 0; idx >>= 1 {
		st.tr[idx] = mul(st.tr[idx<<1], st.tr[idx<<1|1])
	}
}

// ---------- per‑test‑case solver ----------
type Op struct{ val, idx, x, y int }

func solve(in *bufio.Reader, out *bufio.Writer) {
	var n int
	fmt.Fscan(in, &n)

	// input arrays (1‑based, size n+2 for a[n+1] access)
	a := make([][2]int, n+2)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i][0])
	}
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i][1])
	}

	// f[1..n] matrices
	f := make([]Mat, n+2)
	for i := 1; i <= n; i++ {
		f[i] = newMat()
	}

	var ops []Op

	if n%2 == 0 {
		for i := 1; i <= n; i++ {
			if i%2 == 0 { // even index → zero‑matrix
				f[i].a = [2][2]int{{0, 0}, {0, 0}}
				continue
			}
			for x := 0; x < 2; x++ {
				for y := 0; y < 2; y++ {
					l := a[i][x] + a[i+1][y]
					r := a[i][x^1] + a[i+1][y^1]
					f[i].a[x][y] = max(l, r)
					ops = append(ops, Op{min(l, r), i, x, y})
				}
			}
		}
	} else {
		a[n+1][0], a[n+1][1] = a[1][1], a[1][0]
		for i := 1; i <= n; i++ {
			for x := 0; x < 2; x++ {
				for y := 0; y < 2; y++ {
					if i&1 == 1 { // odd
						f[i].a[x][y] = a[i][x] + a[i+1][y]
					} else { // even
						f[i].a[x][y] = a[i][x^1] + a[i+1][y^1]
					}
					ops = append(ops, Op{f[i].a[x][y], i, x, y})
				}
			}
		}
	}

	st := NewSegTree(n, f)

	ans := inf
	sort.Slice(ops, func(i, j int) bool { return ops[i].val < ops[j].val })

	for _, op := range ops {
		root := st.tr[1]
		v := min(root.a[0][0], root.a[1][1])
		if v == inf {
			break
		}
		if cur := v - op.val; cur < ans {
			ans = cur
		}
		// invalidate the chosen entry & update seg‑tree
		f[op.idx].a[op.x][op.y] = inf
		st.Update(op.idx, f[op.idx])
	}

	fmt.Fprintln(out, ans)
}

// ---------- main ----------
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var tc int
	fmt.Fscan(in, &tc)
	for ; tc > 0; tc-- {
		solve(in, out)
	}
}

