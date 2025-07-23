package main

import (
	"bufio"
	"fmt"
	"math/big"
	"math/bits"
	"os"
)

var (
	n         int
	M         int
	g         [][]int
	start     []int
	endt      []int
	flat      []int
	timer     int
	mask      *big.Int
	primeMask *big.Int
)

type SegmentTree struct {
	n    int
	tree []*big.Int
	lazy []int
}

func NewSegmentTree(vals []int) *SegmentTree {
	st := &SegmentTree{
		n:    len(vals),
		tree: make([]*big.Int, len(vals)*4),
		lazy: make([]int, len(vals)*4),
	}
	st.build(1, 0, st.n-1, vals)
	return st
}

func (st *SegmentTree) build(node, l, r int, vals []int) {
	if l == r {
		st.tree[node] = new(big.Int).SetBit(new(big.Int), vals[l], 1)
		return
	}
	mid := (l + r) / 2
	st.build(node*2, l, mid, vals)
	st.build(node*2+1, mid+1, r, vals)
	st.pull(node)
}

func (st *SegmentTree) pull(node int) {
	if st.tree[node] == nil {
		st.tree[node] = new(big.Int)
	}
	st.tree[node].Or(st.tree[node*2], st.tree[node*2+1])
}

func rotateInPlace(b *big.Int, k int) {
	k %= M
	if k == 0 {
		return
	}
	tmp := new(big.Int).Set(b)
	b.Lsh(b, uint(k))
	tmp.Rsh(tmp, uint(M-k))
	b.Or(b, tmp)
	b.And(b, mask)
}

func (st *SegmentTree) apply(node, k int) {
	rotateInPlace(st.tree[node], k)
	st.lazy[node] = (st.lazy[node] + k) % M
}

func (st *SegmentTree) push(node int) {
	if st.lazy[node] != 0 {
		k := st.lazy[node]
		st.apply(node*2, k)
		st.apply(node*2+1, k)
		st.lazy[node] = 0
	}
}

func (st *SegmentTree) update(node, l, r, ql, qr, k int) {
	if ql <= l && r <= qr {
		st.apply(node, k)
		return
	}
	st.push(node)
	mid := (l + r) / 2
	if ql <= mid {
		st.update(node*2, l, mid, ql, qr, k)
	}
	if qr > mid {
		st.update(node*2+1, mid+1, r, ql, qr, k)
	}
	st.pull(node)
}

func (st *SegmentTree) Update(l, r, k int) {
	k %= M
	if k < 0 {
		k += M
	}
	st.update(1, 0, st.n-1, l, r, k)
}

func (st *SegmentTree) query(node, l, r, ql, qr int, res *big.Int) {
	if ql <= l && r <= qr {
		res.Or(res, st.tree[node])
		return
	}
	st.push(node)
	mid := (l + r) / 2
	if ql <= mid {
		st.query(node*2, l, mid, ql, qr, res)
	}
	if qr > mid {
		st.query(node*2+1, mid+1, r, ql, qr, res)
	}
}

func (st *SegmentTree) Query(l, r int) *big.Int {
	res := new(big.Int)
	st.query(1, 0, st.n-1, l, r, res)
	return res
}

func sievePrimes(m int) []int {
	isPrime := make([]bool, m)
	for i := 2; i < m; i++ {
		isPrime[i] = true
	}
	for i := 2; i*i < m; i++ {
		if isPrime[i] {
			for j := i * i; j < m; j += i {
				isPrime[j] = false
			}
		}
	}
	primes := []int{}
	for i := 2; i < m; i++ {
		if isPrime[i] {
			primes = append(primes, i)
		}
	}
	return primes
}

func bitCount(b *big.Int) int {
	cnt := 0
	for _, w := range b.Bits() {
		cnt += bits.OnesCount(uint(w))
	}
	return cnt
}

func dfs(u, p int) {
	start[u] = timer
	flat[timer] = u
	timer++
	for _, v := range g[u] {
		if v != p {
			dfs(v, u)
		}
	}
	endt[u] = timer - 1
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &M)
	valsOrig := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &valsOrig[i])
	}

	g = make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	start = make([]int, n+1)
	endt = make([]int, n+1)
	flat = make([]int, n)
	timer = 0
	dfs(1, 0)

	values := make([]int, n)
	for i := 0; i < n; i++ {
		node := flat[i]
		values[i] = valsOrig[node] % M
	}

	// prepare masks
	mask = new(big.Int).Lsh(big.NewInt(1), uint(M))
	mask.Sub(mask, big.NewInt(1))

	primeMask = new(big.Int)
	for _, p := range sievePrimes(M) {
		primeMask.SetBit(primeMask, p, 1)
	}

	st := NewSegmentTree(values)

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var typ, v int
		fmt.Fscan(reader, &typ, &v)
		if typ == 1 {
			var x int
			fmt.Fscan(reader, &x)
			st.Update(start[v], endt[v], x%M)
		} else if typ == 2 {
			res := st.Query(start[v], endt[v])
			res.And(res, primeMask)
			cnt := bitCount(res)
			fmt.Fprintln(writer, cnt)
		}
	}
}
