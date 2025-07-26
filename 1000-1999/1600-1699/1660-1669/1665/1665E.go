package main

import (
	"bufio"
	"fmt"
	"os"
)

const K = 31

type SegTree struct {
	n    int
	data [][]int
}

func merge(a, b []int) []int {
	res := make([]int, 0, len(a)+len(b))
	i, j := 0, 0
	for len(res) < K && (i < len(a) || j < len(b)) {
		if j >= len(b) || (i < len(a) && a[i] <= b[j]) {
			res = append(res, a[i])
			i++
		} else {
			res = append(res, b[j])
			j++
		}
	}
	return res
}

func NewSegTree(arr []int) *SegTree {
	n := 1
	for n < len(arr) {
		n <<= 1
	}
	data := make([][]int, 2*n)
	st := &SegTree{n, data}
	for i := 0; i < len(arr); i++ {
		data[n+i] = []int{arr[i]}
	}
	for i := n - 1; i >= 1; i-- {
		data[i] = merge(data[2*i], data[2*i+1])
	}
	return st
}

func (st *SegTree) Query(l, r int) []int {
	l += st.n
	r += st.n
	leftParts := [][]int{}
	rightParts := [][]int{}
	for l <= r {
		if l%2 == 1 {
			leftParts = append(leftParts, st.data[l])
			l++
		}
		if r%2 == 0 {
			rightParts = append(rightParts, st.data[r])
			r--
		}
		l >>= 1
		r >>= 1
	}
	res := []int{}
	for _, part := range leftParts {
		res = merge(res, part)
	}
	for i := len(rightParts) - 1; i >= 0; i-- {
		res = merge(res, rightParts[i])
	}
	return res
}

func minOr(nums []int) int {
	ans := int(^uint(0) >> 1) // max int
	for i := 0; i < len(nums); i++ {
		for j := i + 1; j < len(nums); j++ {
			v := nums[i] | nums[j]
			if v < ans {
				ans = v
			}
		}
	}
	return ans
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		st := NewSegTree(arr)
		var q int
		fmt.Fscan(in, &q)
		for ; q > 0; q-- {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			r--
			vals := st.Query(l, r)
			ans := minOr(vals)
			fmt.Fprintln(out, ans)
		}
	}
}
