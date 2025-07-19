package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 1000000007

// Node represents a segment tree node
type Node struct {
	prod     int    // product of segment
	mask     uint64 // bitmask of primes in segment
	lazyMul  int    // pending multiplication
	lazyMask uint64 // pending mask OR
}

var (
	n, Q    int
	primes  []int
	invs    []int
	maskArr []uint64
	a       []int
	tree    []Node
)

// fast exponentiation
func fpow(a, b int) int {
	res := 1
	a %= mod
	for b > 0 {
		if b&1 != 0 {
			res = int((int64(res) * int64(a)) % mod)
		}
		a = int((int64(a) * int64(a)) % mod)
		b >>= 1
	}
	return res
}

// build the segment tree
func build(idx, l, r int) {
	tree[idx].lazyMul = 1
	tree[idx].lazyMask = 0
	if l == r {
		tree[idx].prod = a[l]
		tree[idx].mask = maskArr[a[l]]
		return
	}
	mid := (l + r) >> 1
	build(idx<<1, l, mid)
	build(idx<<1|1, mid+1, r)
	tree[idx].prod = int((int64(tree[idx<<1].prod) * int64(tree[idx<<1|1].prod)) % mod)
	tree[idx].mask = tree[idx<<1].mask | tree[idx<<1|1].mask
}

// push down lazy values
func push(idx, l, r int) {
	lm, lmask := tree[idx].lazyMul, tree[idx].lazyMask
	if lm == 1 && lmask == 0 {
		return
	}
	mid := (l + r) >> 1
	// left child
	left := idx << 1
	tree[left].prod = int((int64(tree[left].prod) * int64(fpow(lm, mid-l+1))) % mod)
	tree[left].lazyMul = int((int64(tree[left].lazyMul) * int64(lm)) % mod)
	tree[left].mask |= lmask
	tree[left].lazyMask |= lmask
	// right child
	right := idx<<1 | 1
	tree[right].prod = int((int64(tree[right].prod) * int64(fpow(lm, r-mid))) % mod)
	tree[right].lazyMul = int((int64(tree[right].lazyMul) * int64(lm)) % mod)
	tree[right].mask |= lmask
	tree[right].lazyMask |= lmask
	// clear
	tree[idx].lazyMul = 1
	tree[idx].lazyMask = 0
}

// update multiplies [ql,qr] by x, mask m
func update(idx, l, r, ql, qr, x int, m uint64) {
	if ql <= l && r <= qr {
		tree[idx].prod = int((int64(tree[idx].prod) * int64(fpow(x, r-l+1))) % mod)
		tree[idx].lazyMul = int((int64(tree[idx].lazyMul) * int64(x)) % mod)
		tree[idx].mask |= m
		tree[idx].lazyMask |= m
		return
	}
	push(idx, l, r)
	mid := (l + r) >> 1
	if ql <= mid {
		update(idx<<1, l, mid, ql, qr, x, m)
	}
	if qr > mid {
		update(idx<<1|1, mid+1, r, ql, qr, x, m)
	}
	tree[idx].prod = int((int64(tree[idx<<1].prod) * int64(tree[idx<<1|1].prod)) % mod)
	tree[idx].mask = tree[idx<<1].mask | tree[idx<<1|1].mask
}

// query product on [ql,qr]
func queryProd(idx, l, r, ql, qr int) int {
	if ql <= l && r <= qr {
		return tree[idx].prod
	}
	push(idx, l, r)
	mid := (l + r) >> 1
	res := 1
	if ql <= mid {
		res = int((int64(res) * int64(queryProd(idx<<1, l, mid, ql, qr))) % mod)
	}
	if qr > mid {
		res = int((int64(res) * int64(queryProd(idx<<1|1, mid+1, r, ql, qr))) % mod)
	}
	return res
}

// query mask on [ql,qr]
func queryMask(idx, l, r, ql, qr int) uint64 {
	if ql <= l && r <= qr {
		return tree[idx].mask
	}
	push(idx, l, r)
	mid := (l + r) >> 1
	var res uint64
	if ql <= mid {
		res |= queryMask(idx<<1, l, mid, ql, qr)
	}
	if qr > mid {
		res |= queryMask(idx<<1|1, mid+1, r, ql, qr)
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &Q)
	// sieve primes up to 300
	maxv := 300
	isComp := make([]bool, maxv+1)
	for i := 2; i <= maxv; i++ {
		if !isComp[i] {
			primes = append(primes, i)
		}
		for _, p := range primes {
			if i*p > maxv {
				break
			}
			isComp[i*p] = true
			if i%p == 0 {
				break
			}
		}
	}
	m := len(primes)
	invs = make([]int, m)
	for i := 0; i < m; i++ {
		invs[i] = fpow(primes[i], mod-2)
	}
	maskArr = make([]uint64, maxv+1)
	for v := 2; v <= maxv; v++ {
		var mm uint64
		for i, p := range primes {
			if v%p == 0 {
				mm |= 1 << uint(i)
			}
		}
		maskArr[v] = mm
	}
	// read initial array
	a = make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	tree = make([]Node, 4*(n+1))
	build(1, 1, n)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	for Q > 0 {
		Q--
		var op string
		l, r := 0, 0
		fmt.Fscan(reader, &op, &l, &r)
		if op[0] == 'M' {
			var x int
			fmt.Fscan(reader, &x)
			update(1, 1, n, l, r, x, maskArr[x])
		} else {
			pm := queryMask(1, 1, n, l, r)
			res := queryProd(1, 1, n, l, r)
			for i, p := range primes {
				if pm&(1<<uint(i)) != 0 {
					res = int((int64(res) * int64(invs[i]) % mod) * int64(p-1) % mod)
				}
			}
			fmt.Fprintln(writer, res)
		}
	}
}
