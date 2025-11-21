package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type SegmentTree struct {
	n    int
	tree []int
}

func NewSegmentTree(n int) *SegmentTree {
	size := 1
	for size < n {
		size <<= 1
	}
	return &SegmentTree{n: size, tree: make([]int, 2*size)}
}

func (st *SegmentTree) Update(pos, val int) {
	idx := pos + st.n
	st.tree[idx] = val
	idx >>= 1
	for idx > 0 {
		left := st.tree[idx<<1]
		right := st.tree[idx<<1|1]
		st.tree[idx] = max(left, right)
		idx >>= 1
	}
}

func (st *SegmentTree) Query(l, r int) int {
	l += st.n
	r += st.n
	res := 0
	for l <= r {
		if l&1 == 1 {
			res = max(res, st.tree[l])
			l++
		}
		if r&1 == 0 {
			res = max(res, st.tree[r])
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func canAchieve(mid int, n, k int, s string, a []int) bool {
	bad := make([]bool, n)
	for i := 0; i < n; i++ {
		if a[i] > mid {
			bad[i] = true
		}
	}

	intervals := make([][2]int, 0)
	start := -1
	for i := 0; i < n; i++ {
		if bad[i] {
			if start == -1 {
				start = i
			}
		} else {
			if start != -1 {
				intervals = append(intervals, [2]int{start, i - 1})
				start = -1
			}
		}
	}
	if start != -1 {
		intervals = append(intervals, [2]int{start, n - 1})
	}

	st := NewSegmentTree(n)
	for i := 0; i < n; i++ {
		st.Update(i, 0)
	}

	for i := 0; i < n; i++ {
		if s[i] == 'B' {
			st.Update(i, 1)
		} else {
			st.Update(i, 0)
		}
	}

	remaining := k
	for _, interval := range intervals {
		l, r := interval[0], interval[1]
		maxBlue := st.Query(l, r)
		if maxBlue == r-l+1 {
			if remaining == 0 {
				return false
			}
			remaining--
			for j := l; j <= r; j++ {
				st.Update(j, 0)
			}
		} else {
			return false
		}
	}
	return true
}

func solveCase(n, k int, s string, a []int) int {
	unique := make(map[int]struct{})
	for _, val := range a {
		unique[val] = struct{}{}
	}
	values := make([]int, 0, len(unique))
	for val := range unique {
		values = append(values, val)
	}
	sort.Ints(values)

	lo, hi := 0, len(values)-1
	ans := values[hi]
	for lo <= hi {
		mid := (lo + hi) / 2
		if canAchieve(values[mid], n, k, s, a) {
			ans = values[mid]
			hi = mid - 1
		} else {
			lo = mid + 1
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		result := solveCase(n, k, s, a)
		fmt.Fprintln(out, result)
	}
}
