package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const mod = 1000000007

func add(a, b int64) int64 {
	x := a + b
	if x >= mod {
		x -= mod
	}
	return x
}

func sub(a, b int64) int64 {
	x := a - b
	if x < 0 {
		x += mod
	}
	return x
}

func mul(a, b int64) int64 {
	return (a * b) % mod
}

func powMod(a, e int) int64 {
	res := int64(1)
	base := int64(a % mod)
	for e > 0 {
		if e&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		e >>= 1
	}
	return res
}

// segment tree
type Node struct {
	t    [6]int64
	lazy int64
	set  bool
}

var n, m int
var a []int64
var P [6][]int64
var C [6][6]int64
var st []Node

func build(node, l, r int) {
	if l == r {
		ai := a[l]
		var p int64 = 1
		for j := 0; j <= 5; j++ {
			st[node].t[j] = (ai % mod) * p % mod
			p = p * int64(l) % mod
		}
		return
	}
	mid := (l + r) >> 1
	lc, rc := node<<1, node<<1|1
	build(lc, l, mid)
	build(rc, mid+1, r)
	for j := 0; j <= 5; j++ {
		st[node].t[j] = add(st[lc].t[j], st[rc].t[j])
	}
}

func applySet(node, l, r int, x int64) {
	st[node].set = true
	st[node].lazy = x
	for j := 0; j <= 5; j++ {
		// sum i^j from l to r
		sum := sub(P[j][r], P[j][l-1])
		st[node].t[j] = mul(x%mod, sum)
	}
}

func push(node, l, r int) {
	if !st[node].set || l == r {
		return
	}
	mid := (l + r) >> 1
	lc, rc := node<<1, node<<1|1
	applySet(lc, l, mid, st[node].lazy)
	applySet(rc, mid+1, r, st[node].lazy)
	st[node].set = false
}

func update(node, l, r, ql, qr int, x int64) {
	if ql <= l && r <= qr {
		applySet(node, l, r, x)
		return
	}
	push(node, l, r)
	mid := (l + r) >> 1
	if ql <= mid {
		update(node<<1, l, mid, ql, qr, x)
	}
	if qr > mid {
		update(node<<1|1, mid+1, r, ql, qr, x)
	}
	for j := 0; j <= 5; j++ {
		st[node].t[j] = add(st[node<<1].t[j], st[node<<1|1].t[j])
	}
}

func query(node, l, r, ql, qr int, res *[6]int64) {
	if ql <= l && r <= qr {
		for j := 0; j <= 5; j++ {
			res[j] = add(res[j], st[node].t[j])
		}
		return
	}
	push(node, l, r)
	mid := (l + r) >> 1
	if ql <= mid {
		query(node<<1, l, mid, ql, qr, res)
	}
	if qr > mid {
		query(node<<1|1, mid+1, r, ql, qr, res)
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	// read n, m
	fmt.Fscan(reader, &n, &m)
	a = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	// precompute P
	for j := 0; j <= 5; j++ {
		P[j] = make([]int64, n+1)
		for i := 1; i <= n; i++ {
			P[j][i] = add(P[j][i-1], mul(int64(powMod(i, j)), 1))
		}
	}
	// binomial
	for i := 0; i <= 5; i++ {
		C[i][0] = 1
		for j := 1; j <= i; j++ {
			C[i][j] = add(C[i-1][j-1], C[i-1][j])
		}
	}
	// build segtree
	st = make([]Node, 4*(n+5))
	build(1, 1, n)
	// process queries
	for qi := 0; qi < m; qi++ {
		var op string
		fmt.Fscan(reader, &op)
		if op == "?" {
			var l, r, k int
			fmt.Fscan(reader, &l, &r, &k)
			var res [6]int64
			query(1, 1, n, l, r, &res)
			pw := make([]int64, k+1)
			pw[0] = 1
			neg := (mod - int64(l-1)%mod) % mod
			for t := 1; t <= k; t++ {
				pw[t] = pw[t-1] * neg % mod
			}
			var ans int64
			for j := 0; j <= k; j++ {
				coef := C[k][j] * pw[k-j] % mod
				ans = (ans + res[j]*coef) % mod
			}
			fmt.Fprint(writer, ans, '\n')
		} else {
			// assignment, op is l
			l, _ := strconv.Atoi(op)
			var r int
			var x int64
			fmt.Fscan(reader, &r, &x)
			update(1, 1, n, l, r, x)
		}
	}
}
