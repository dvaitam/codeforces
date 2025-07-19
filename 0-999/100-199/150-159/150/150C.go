package main

import (
	"bufio"
	"fmt"
	"os"
)

// Sum holds segment data for maximum subarray calculations
type Sum struct {
	l, r, w, sum int64
}

// meg merges two Sum nodes
func meg(a, b Sum) Sum {
	var ans Sum
	ans.w = a.w
	if b.w > ans.w {
		ans.w = b.w
	}
	if a.r+b.l > ans.w {
		ans.w = a.r + b.l
	}
	ans.l = a.l
	if a.sum+b.l > ans.l {
		ans.l = a.sum + b.l
	}
	ans.r = b.r
	if b.sum+a.r > ans.r {
		ans.r = b.sum + a.r
	}
	ans.sum = a.sum + b.sum
	return ans
}

var (
	n, m, c int
	xx      []int64
	p       []int64
	arr     []int64
	st      []Sum
)

// build constructs the segment tree
func build(u, L, R int) {
	if L == R {
		x := arr[L]
		if x < 0 {
			st[u] = Sum{0, 0, 0, x}
		} else {
			st[u] = Sum{x, x, x, x}
		}
		return
	}
	mid := (L + R) >> 1
	build(u<<1, L, mid)
	build(u<<1|1, mid+1, R)
	st[u] = meg(st[u<<1], st[u<<1|1])
}

// query returns the merged Sum for range [l, r]
func query(u, L, R, l, r int) Sum {
	if l <= L && R <= r {
		return st[u]
	}
	mid := (L + R) >> 1
	if r <= mid {
		return query(u<<1, L, mid, l, r)
	}
	if l > mid {
		return query(u<<1|1, mid+1, R, l, r)
	}
	left := query(u<<1, L, mid, l, r)
	right := query(u<<1|1, mid+1, R, l, r)
	return meg(left, right)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	// Read inputs
	var tmpN, tmpM, tmpC int
	fmt.Fscan(reader, &tmpN, &tmpM, &tmpC)
	n, m, c = tmpN, tmpM, tmpC
	xx = make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &xx[i])
	}
	p = make([]int64, n+1)
	for i := 1; i < n; i++ {
		fmt.Fscan(reader, &p[i])
	}
	// Prepare transformed array
	if n > 1 {
		arr = make([]int64, n)
		for i := 1; i < n; i++ {
			diff := xx[i+1] - xx[i]
			arr[i] = (diff-int64(2*c))*p[i] + diff*(100-p[i])
		}
		st = make([]Sum, 4*n)
		build(1, 1, n-1)
	}
	// Process queries
	var ans int64
	for i := 0; i < m; i++ {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		if l < r {
			res := query(1, 1, n-1, l, r-1)
			ans += res.w
		}
	}
	// Output result
	fmt.Fprintf(writer, "%.12f\n", float64(ans)/200.0)
}
