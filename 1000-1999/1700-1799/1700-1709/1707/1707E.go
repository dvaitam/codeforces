package main

import (
	"bufio"
	"fmt"
	"os"
)

// segment tree for range min and max

type SegTree struct {
	n  int
	mn []int
	mx []int
}

func NewSegTree(a []int) *SegTree {
	n := 1
	for n < len(a) {
		n <<= 1
	}
	mn := make([]int, 2*n)
	mx := make([]int, 2*n)
	for i := 0; i < len(a); i++ {
		mn[n+i] = a[i]
		mx[n+i] = a[i]
	}
	for i := n - 1; i > 0; i-- {
		left := i << 1
		right := left | 1
		mn[i] = mn[left]
		if mn[right] < mn[i] {
			mn[i] = mn[right]
		}
		mx[i] = mx[left]
		if mx[right] > mx[i] {
			mx[i] = mx[right]
		}
	}
	return &SegTree{n: n, mn: mn, mx: mx}
}

func (t *SegTree) Query(l, r int) (int, int) {
	l += t.n
	r += t.n
	minVal := int(1e9)
	maxVal := 0
	for l <= r {
		if l&1 == 1 {
			if t.mn[l] < minVal {
				minVal = t.mn[l]
			}
			if t.mx[l] > maxVal {
				maxVal = t.mx[l]
			}
			l++
		}
		if r&1 == 0 {
			if t.mn[r] < minVal {
				minVal = t.mn[r]
			}
			if t.mx[r] > maxVal {
				maxVal = t.mx[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return minVal, maxVal
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
		a[i]-- // work with 0-indexed internally
	}

	segMinMax := NewSegTree(a)

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		l--
		r--
		steps := 0
		L := l
		R := r
		visited := make(map[[2]int]bool)
		for steps <= n*2 {
			if L == 0 && R == n-1 {
				break
			}
			key := [2]int{L, R}
			if visited[key] {
				break
			}
			visited[key] = true
			mn, mx := segMinMax.Query(L, R)
			L = mn
			R = mx
			steps++
		}
		if L == 0 && R == n-1 {
			fmt.Fprintln(writer, steps)
		} else {
			fmt.Fprintln(writer, -1)
		}
	}
}
