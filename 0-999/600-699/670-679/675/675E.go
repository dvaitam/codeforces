package main

import (
	"bufio"
	"fmt"
	"os"
)

type SegTree struct {
	size int
	data []int
}

func NewSegTree(a []int) *SegTree {
	n := len(a) - 1 // array is 1-indexed
	size := 1
	for size < n {
		size <<= 1
	}
	data := make([]int, 2*size)
	for i := 1; i <= n; i++ {
		data[size+i-1] = a[i]
	}
	for i := size - 1; i >= 1; i-- {
		left, right := data[2*i], data[2*i+1]
		if left > right {
			data[i] = left
		} else {
			data[i] = right
		}
	}
	return &SegTree{size, data}
}

func (st *SegTree) Query(l, r int) int {
	if l > r {
		return 0
	}
	l += st.size - 1
	r += st.size - 1
	res := 0
	for l <= r {
		if l%2 == 1 {
			if st.data[l] > res {
				res = st.data[l]
			}
			l++
		}
		if r%2 == 0 {
			if st.data[r] > res {
				res = st.data[r]
			}
			r--
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	a[n] = n
	st := NewSegTree(a)
	var ans int64
	for i := 1; i < n; i++ {
		cur := i
		reach := a[i]
		dist := 1
		for {
			if reach >= n {
				ans += int64(dist) * int64(n-cur)
				break
			}
			ans += int64(dist) * int64(reach-cur)
			newReach := st.Query(cur+1, reach)
			cur = reach
			reach = newReach
			dist++
		}
	}
	fmt.Println(ans)
}
