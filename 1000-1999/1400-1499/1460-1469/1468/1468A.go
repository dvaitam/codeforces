package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	l, r uint32
	val  int
}

type PersistentSegTree struct {
	nodes []Node
	m     int
}

func NewPST(m, capNodes int) *PersistentSegTree {
	nodes := make([]Node, 1, capNodes+1)
	return &PersistentSegTree{nodes: nodes, m: m}
}

func (pst *PersistentSegTree) query(idx uint32, l, r, ql, qr int) int {
	if idx == 0 || ql > r || qr < l {
		return 0
	}
	if ql <= l && r <= qr {
		return pst.nodes[idx].val
	}
	mid := (l + r) >> 1
	res := 0
	if ql <= mid {
		v := pst.query(pst.nodes[idx].l, l, mid, ql, qr)
		if v > res {
			res = v
		}
	}
	if qr > mid {
		v := pst.query(pst.nodes[idx].r, mid+1, r, ql, qr)
		if v > res {
			res = v
		}
	}
	return res
}

func (pst *PersistentSegTree) update(idx uint32, l, r, pos, val int) uint32 {
	if l == r {
		old := 0
		if idx != 0 {
			old = pst.nodes[idx].val
		}
		if val < old {
			val = old
		}
		pst.nodes = append(pst.nodes, Node{0, 0, val})
		return uint32(len(pst.nodes) - 1)
	}
	mid := (l + r) >> 1
	var oldL, oldR uint32
	if idx != 0 {
		n := pst.nodes[idx]
		oldL, oldR = n.l, n.r
	}
	var newL, newR uint32
	if pos <= mid {
		newL = pst.update(oldL, l, mid, pos, val)
		newR = oldR
	} else {
		newL = oldL
		newR = pst.update(oldR, mid+1, r, pos, val)
	}
	leftVal := 0
	if newL != 0 {
		leftVal = pst.nodes[newL].val
	}
	rightVal := 0
	if newR != 0 {
		rightVal = pst.nodes[newR].val
	}
	newVal := leftVal
	if rightVal > newVal {
		newVal = rightVal
	}
	pst.nodes = append(pst.nodes, Node{newL, newR, newVal})
	return uint32(len(pst.nodes) - 1)
}

type MaxSeg struct {
	n    int
	data []int
}

func NewMaxSeg(size int) *MaxSeg {
	n := 1
	for n < size {
		n <<= 1
	}
	return &MaxSeg{n: n, data: make([]int, n<<1)}
}

func (st *MaxSeg) pointUpdate(pos, val int) {
	i := st.n + pos - 1
	if val > st.data[i] {
		st.data[i] = val
	}
	for i >>= 1; i > 0; i >>= 1 {
		v := st.data[i<<1]
		if st.data[i<<1|1] > v {
			v = st.data[i<<1|1]
		}
		st.data[i] = v
	}
}

func (st *MaxSeg) rangeMax(l, r int) int {
	if l > r {
		return 0
	}
	l += st.n - 1
	r += st.n - 1
	res := 0
	for l <= r {
		if l&1 == 1 {
			if st.data[l] > res {
				res = st.data[l]
			}
			l++
		}
		if r&1 == 0 {
			if st.data[r] > res {
				res = st.data[r]
			}
			if r == 0 {
				break
			}
			r--
		}
		l >>= 1
		r >>= 1
		if l == 0 || r == 0 {
			break
		}
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

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
		m := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > m {
				m = a[i]
			}
		}
		if m == 0 {
			fmt.Fprintln(out, 0)
			continue
		}
		lastGt := NewMaxSeg(max(m, 1))
		depth := 0
		p := 1
		for p < max(m, 1) {
			p <<= 1
			depth++
		}
		capNodes := (depth + 4) * n
		pst := NewPST(max(m, 1), capNodes)
		roots := make([]uint32, n+1)

		ans := 0
		for i := 1; i <= n; i++ {
			v := a[i-1]
			g := 0
			if v < m {
				g = lastGt.rangeMax(v+1, m)
			}
			up := 1 + pst.query(roots[i-1], 1, pst.m, 1, v)
			down := 0
			if g > 0 {
				down = 2 + pst.query(roots[g-1], 1, pst.m, 1, v)
			}
			dp := up
			if down > dp {
				dp = down
			}
			if dp > ans {
				ans = dp
			}
			roots[i] = pst.update(roots[i-1], 1, pst.m, v, dp)
			lastGt.pointUpdate(v, i)
		}
		fmt.Fprintln(out, ans)
	}
}
