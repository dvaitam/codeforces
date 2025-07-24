package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const mod int64 = 998244353

func modPow(a, e int64) int64 {
	res := int64(1)
	for e > 0 {
		if e&1 == 1 {
			res = res * a % mod
		}
		a = a * a % mod
		e >>= 1
	}
	return res
}

type matrix struct {
	a [2][2]int64
}

func mul(A, B matrix) matrix {
	var C matrix
	for i := 0; i < 2; i++ {
		for k := 0; k < 2; k++ {
			var s int64
			for j := 0; j < 2; j++ {
				s = (s + A.a[i][j]*B.a[j][k]) % mod
			}
			C.a[i][k] = s
		}
	}
	return C
}

func identityMatrix() matrix {
	return matrix{[2][2]int64{{1, 0}, {0, 1}}}
}

// segment tree for matrix product
type segTree struct {
	n    int
	tree []matrix
}

func newSegTree(arr []matrix) *segTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	st := &segTree{n, make([]matrix, 2*n)}
	for i := 0; i < len(arr); i++ {
		st.tree[n+i] = arr[i]
	}
	for i := n - 1; i > 0; i-- {
		left := st.tree[i<<1]
		right := st.tree[i<<1|1]
		if left.a[0][0] == 0 && left.a[0][1] == 0 && left.a[1][0] == 0 && left.a[1][1] == 0 {
			st.tree[i] = right
		} else if right.a[0][0] == 0 && right.a[0][1] == 0 && right.a[1][0] == 0 && right.a[1][1] == 0 {
			st.tree[i] = left
		} else {
			st.tree[i] = mul(left, right)
		}
	}
	return st
}

func (st *segTree) update(pos int, val matrix) {
	idx := st.n + pos
	st.tree[idx] = val
	for idx >>= 1; idx > 0; idx >>= 1 {
		left := st.tree[idx<<1]
		right := st.tree[idx<<1|1]
		if left.a[0][0] == 0 && left.a[0][1] == 0 && left.a[1][0] == 0 && left.a[1][1] == 0 {
			st.tree[idx] = right
		} else if right.a[0][0] == 0 && right.a[0][1] == 0 && right.a[1][0] == 0 && right.a[1][1] == 0 {
			st.tree[idx] = left
		} else {
			st.tree[idx] = mul(left, right)
		}
	}
}

func (st *segTree) product() matrix {
	if len(st.tree) == 0 {
		return identityMatrix()
	}
	return st.tree[1]
}

type frac struct {
	num int64
	den int64
}

func less(a, b frac) bool {
	return a.num*b.den < b.num*a.den
}

func equal(a, b frac) bool {
	return a.num*b.den == b.num*a.den
}

type event struct {
	f   frac
	idx int
	r   int
	c   int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	x := make([]int64, n)
	v := make([]int64, n)
	pIn := make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &x[i], &v[i], &pIn[i])
	}
	inv100 := modPow(100, mod-2)
	p := make([]int64, n)
	for i := 0; i < n; i++ {
		p[i] = pIn[i] % mod * inv100 % mod
	}

	// prepare matrices for pairs
	mats := make([]matrix, max(0, n-1))
	for i := 0; i+1 < n; i++ {
		q := (1 - p[i+1] + mod) % mod
		mats[i] = matrix{[2][2]int64{{q, p[i+1]}, {q, p[i+1]}}}
	}
	st := newSegTree(mats)

	// compute probability no collision
	probNoCollision := func() int64 {
		prod := st.product()
		q1 := (1 - p[0] + mod) % mod
		res0 := (q1*prod.a[0][0] + p[0]*prod.a[1][0]) % mod
		res1 := (q1*prod.a[0][1] + p[0]*prod.a[1][1]) % mod
		return (res0 + res1) % mod
	}

	// generate events
	const inf int64 = 1 << 60
	events := make([]event, 0)
	for i := 0; i+1 < n; i++ {
		d := x[i+1] - x[i]
		// LL catch up
		if v[i+1] > v[i] {
			events = append(events, event{frac{d, v[i+1] - v[i]}, i, 0, 0})
		}
		// RL meet
		events = append(events, event{frac{d, v[i] + v[i+1]}, i, 1, 0})
		// RR catch up
		if v[i] > v[i+1] {
			events = append(events, event{frac{d, v[i] - v[i+1]}, i, 1, 1})
		}
	}

	sort.Slice(events, func(i, j int) bool { return less(events[i].f, events[j].f) })

	ans := int64(0)
	prevProb := int64(1)
	i := 0
	for i < len(events) {
		j := i
		f := events[i].f
		for j < len(events) && equal(events[j].f, f) {
			j++
		}
		// probability after processing previous events but before this time
		timeMod := f.num % mod * modPow(f.den%mod, mod-2) % mod
		// apply events in [i,j)
		for k := i; k < j; k++ {
			e := events[k]
			m := mats[e.idx]
			if m.a[e.r][e.c] != 0 {
				m.a[e.r][e.c] = 0
				mats[e.idx] = m
				st.update(e.idx, m)
			}
		}
		newProb := probNoCollision()
		delta := (prevProb - newProb) % mod
		if delta < 0 {
			delta += mod
		}
		ans = (ans + delta*timeMod) % mod
		prevProb = newProb
		i = j
	}
	fmt.Fprintln(out, ans)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
