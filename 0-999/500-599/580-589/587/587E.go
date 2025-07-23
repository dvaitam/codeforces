package main

import (
	"bufio"
	"fmt"
	"os"
)

type Basis struct {
	b  [31]int
	sz int
}

func (bs *Basis) Add(x int) {
	for i := 30; i >= 0; i-- {
		if (x>>i)&1 == 0 {
			continue
		}
		if bs.b[i] != 0 {
			x ^= bs.b[i]
		} else {
			bs.b[i] = x
			bs.sz++
			return
		}
	}
}

func (bs *Basis) Merge(o *Basis) {
	for i := 30; i >= 0; i-- {
		if o.b[i] != 0 {
			bs.Add(o.b[i])
		}
	}
}

// Segment tree for basis of diff array
var seg []Basis
var diff []int

func build(pos, l, r int) {
	if l == r {
		var b Basis
		if diff[l] != 0 {
			b.Add(diff[l])
		}
		seg[pos] = b
		return
	}
	mid := (l + r) >> 1
	build(pos<<1, l, mid)
	build(pos<<1|1, mid+1, r)
	seg[pos] = Basis{}
	seg[pos].Merge(&seg[pos<<1])
	seg[pos].Merge(&seg[pos<<1|1])
}

func update(pos, l, r, idx int) {
	if l == r {
		var b Basis
		if diff[l] != 0 {
			b.Add(diff[l])
		}
		seg[pos] = b
		return
	}
	mid := (l + r) >> 1
	if idx <= mid {
		update(pos<<1, l, mid, idx)
	} else {
		update(pos<<1|1, mid+1, r, idx)
	}
	seg[pos] = Basis{}
	seg[pos].Merge(&seg[pos<<1])
	seg[pos].Merge(&seg[pos<<1|1])
}

func query(pos, l, r, ql, qr int, res *Basis) {
	if ql > r || qr < l {
		return
	}
	if ql <= l && r <= qr {
		res.Merge(&seg[pos])
		return
	}
	mid := (l + r) >> 1
	if ql <= mid {
		query(pos<<1, l, mid, ql, qr, res)
	}
	if qr > mid {
		query(pos<<1|1, mid+1, r, ql, qr, res)
	}
}

// Fenwick tree for prefix xor of diff
var bit []int
var n int

func bitAdd(idx, val int) {
	for idx <= n {
		bit[idx] ^= val
		idx += idx & -idx
	}
}

func bitQuery(idx int) int {
	res := 0
	for idx > 0 {
		res ^= bit[idx]
		idx -= idx & -idx
	}
	return res
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	arr := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &arr[i])
	}
	diff = make([]int, n+2)
	diff[1] = arr[1]
	for i := 2; i <= n; i++ {
		diff[i] = arr[i] ^ arr[i-1]
	}
	bit = make([]int, n+2)
	for i := 1; i <= n; i++ {
		bitAdd(i, diff[i])
	}
	seg = make([]Basis, 4*(n+2))
	build(1, 1, n)

	for ; q > 0; q-- {
		var t int
		fmt.Fscan(reader, &t)
		if t == 1 {
			var l, r, k int
			fmt.Fscan(reader, &l, &r, &k)
			diff[l] ^= k
			bitAdd(l, k)
			update(1, 1, n, l)
			if r+1 <= n {
				diff[r+1] ^= k
				bitAdd(r+1, k)
				update(1, 1, n, r+1)
			}
		} else {
			var l, r int
			fmt.Fscan(reader, &l, &r)
			var ansBasis Basis
			if l < r {
				query(1, 1, n, l+1, r, &ansBasis)
			}
			arrL := bitQuery(l)
			if arrL != 0 {
				ansBasis.Add(arrL)
			}
			res := 1
			for i := 0; i < ansBasis.sz; i++ {
				res <<= 1
			}
			fmt.Fprintln(writer, res)
		}
	}
}
