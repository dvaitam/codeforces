package main

import (
	"bufio"
	"fmt"
	"os"
)

type Bitset []uint64

var (
	n, q   int
	tree   [][]int
	result Bitset
	words  int
)

func newBitset() Bitset {
	return make(Bitset, words)
}

func (b Bitset) copy() Bitset {
	nb := make(Bitset, len(b))
	copy(nb, b)
	return nb
}

func (b Bitset) or(other Bitset) {
	for i := range b {
		b[i] |= other[i]
	}
}

func (b Bitset) shiftLeft(x int) Bitset {
	res := make(Bitset, len(b))
	ws := x >> 6
	bs := uint(x & 63)
	if ws >= len(b) {
		return res
	}
	if bs == 0 {
		for i := len(b) - 1; i >= ws; i-- {
			res[i] = b[i-ws]
		}
	} else {
		for i := len(b) - 1; i > ws; i-- {
			res[i] = b[i-ws]<<bs | b[i-ws-1]>>(64-bs)
		}
		res[ws] = b[0] << bs
	}
	// mask bits beyond n
	last := n >> 6
	rem := uint(n & 63)
	mask := uint64(^uint64(0))
	if rem != 63 {
		mask = (uint64(1) << (rem + 1)) - 1
	}
	for i := last + 1; i < len(res); i++ {
		res[i] = 0
	}
	res[last] &= mask
	return res
}

func add(idx, l, r, ql, qr, val int) {
	if ql <= l && r <= qr {
		tree[idx] = append(tree[idx], val)
		return
	}
	mid := (l + r) >> 1
	if ql <= mid {
		add(idx<<1, l, mid, ql, qr, val)
	}
	if qr > mid {
		add(idx<<1|1, mid+1, r, ql, qr, val)
	}
}

func dfs(idx, l, r int, bit Bitset) {
	cur := bit.copy()
	for _, v := range tree[idx] {
		shifted := cur.shiftLeft(v)
		for i := range cur {
			cur[i] |= shifted[i]
		}
	}
	if l == r {
		result.or(cur)
		return
	}
	mid := (l + r) >> 1
	dfs(idx<<1, l, mid, cur)
	dfs(idx<<1|1, mid+1, r, cur)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n, &q)
	words = (n+64)/64 + 1
	tree = make([][]int, 4*n+5)
	for i := 0; i < q; i++ {
		var l, r, x int
		fmt.Fscan(reader, &l, &r, &x)
		if x > n {
			// values larger than n are useless for our range
			continue
		}
		add(1, 1, n, l, r, x)
	}
	result = newBitset()
	root := newBitset()
	root[0] = 1
	dfs(1, 1, n, root)
	var ans []int
	for i := 1; i <= n; i++ {
		if (result[i>>6]>>(uint(i)&63))&1 == 1 {
			ans = append(ans, i)
		}
	}
	fmt.Println(len(ans))
	for i, v := range ans {
		if i > 0 {
			fmt.Print(" ")
		}
		fmt.Print(v)
	}
	if len(ans) > 0 {
		fmt.Println()
	}
}
