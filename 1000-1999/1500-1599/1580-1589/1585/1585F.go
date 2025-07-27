package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const MOD int64 = 998244353

type Node struct {
	sum int64
	len int64
	mul int64
	add int64
}

var tree []Node
var size int

func buildLens(m int, lens []int64) {
	size = 1
	for size < m {
		size <<= 1
	}
	tree = make([]Node, size<<1)
	for i := 0; i < m; i++ {
		tree[size+i].len = lens[i]
	}
	for i := size - 1; i > 0; i-- {
		tree[i].len = tree[i<<1].len + tree[i<<1|1].len
	}
	for i := range tree {
		if tree[i].mul == 0 {
			tree[i].mul = 1
		}
	}
}

func apply(i int, mul, add int64) {
	tree[i].sum = (tree[i].sum*mul + add*tree[i].len) % MOD
	tree[i].mul = tree[i].mul * mul % MOD
	tree[i].add = (tree[i].add*mul + add) % MOD
}

func push(i int) {
	if tree[i].mul != 1 || tree[i].add != 0 {
		apply(i<<1, tree[i].mul, tree[i].add)
		apply(i<<1|1, tree[i].mul, tree[i].add)
		tree[i].mul = 1
		tree[i].add = 0
	}
}

func rangeApply(l, r int, mul, add int64) {
	var rec func(i, tl, tr int)
	rec = func(i, tl, tr int) {
		if l >= tr || r <= tl {
			return
		}
		if l <= tl && tr <= r {
			apply(i, mul, add)
			return
		}
		push(i)
		mid := (tl + tr) >> 1
		rec(i<<1, tl, mid)
		rec(i<<1|1, mid, tr)
		tree[i].sum = (tree[i<<1].sum + tree[i<<1|1].sum) % MOD
	}
	rec(1, 0, size)
}

func query(l, r int) int64 {
	var rec func(i, tl, tr int) int64
	rec = func(i, tl, tr int) int64 {
		if l >= tr || r <= tl {
			return 0
		}
		if l <= tl && tr <= r {
			return tree[i].sum
		}
		push(i)
		mid := (tl + tr) >> 1
		res := rec(i<<1, tl, mid) + rec(i<<1|1, mid, tr)
		if res >= MOD {
			res %= MOD
		}
		return res
	}
	return rec(1, 0, size)
}

func solve(a []int) int64 {
	n := len(a)
	valsMap := make(map[int]bool)
	for _, v := range a {
		valsMap[v] = true
	}
	vals := make([]int, 0, len(valsMap))
	for v := range valsMap {
		vals = append(vals, v)
	}
	sort.Ints(vals)
	val := make([]int, len(vals)+1)
	for i, v := range vals {
		val[i+1] = v
	}
	m := len(vals)
	lens := make([]int64, m)
	for i := 0; i < m; i++ {
		lens[i] = int64(val[i+1] - val[i])
	}
	buildLens(m, lens)
	// map from value to index
	idxMap := make(map[int]int)
	for i, v := range vals {
		idxMap[v] = i
	}
	firstIdx := idxMap[a[0]]
	rangeApply(0, firstIdx+1, 0, 1)
	dpTotal := query(0, m)
	for i := 1; i < n; i++ {
		idx := idxMap[a[i]]
		prefix := query(0, idx+1)
		newTotal := (dpTotal*int64(a[i]) - prefix) % MOD
		if newTotal < 0 {
			newTotal += MOD
		}
		if idx+1 < m {
			rangeApply(idx+1, m, 0, 0)
		}
		rangeApply(0, idx+1, MOD-1, dpTotal%MOD)
		dpTotal = newTotal % MOD
	}
	return dpTotal % MOD
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	res := solve(a)
	fmt.Println(res)
}
