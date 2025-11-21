package main

import (
	"bufio"
	"fmt"
	"os"
)

const inf int64 = 1 << 60

type segTree struct {
	size int
	data []int64
}

func newSegTree(n int) *segTree {
	size := 1
	for size < n {
		size <<= 1
	}
	data := make([]int64, 2*size)
	for i := range data {
		data[i] = inf
	}
	return &segTree{size: size, data: data}
}

func (st *segTree) update(pos int, val int64) {
	idx := pos + st.size
	st.data[idx] = val
	for idx >>= 1; idx > 0; idx >>= 1 {
		left := st.data[idx<<1]
		right := st.data[idx<<1|1]
		if left < right {
			st.data[idx] = left
		} else {
			st.data[idx] = right
		}
	}
}

func (st *segTree) query(l, r int) int64 {
	if l >= r {
		return inf
	}
	l += st.size
	r += st.size
	res := inf
	for l < r {
		if l&1 == 1 {
			if st.data[l] < res {
				res = st.data[l]
			}
			l++
		}
		if r&1 == 1 {
			r--
			if st.data[r] < res {
				res = st.data[r]
			}
		}
		l >>= 1
		r >>= 1
	}
	return res
}

func computeNext(a []int64, limit int64) []int {
	n := len(a)
	res := make([]int, n)
	r := 0
	var sum int64
	for i := 0; i < n; i++ {
		if r < i {
			r = i
			sum = 0
		}
		for r < n && sum+a[r] <= limit {
			sum += a[r]
			r++
		}
		if r == i {
			res[i] = -1
		} else {
			res[i] = r
			sum -= a[i]
		}
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &b[i])
		}

		nextPos := make([][]int, m)
		for k := 0; k < m; k++ {
			nextPos[k] = computeNext(a, b[k])
		}

		dpNext := make([]int64, n+1)
		for i := 0; i <= n; i++ {
			dpNext[i] = inf
		}
		dpNext[n] = 0

		for k := m - 1; k >= 0; k-- {
			dpCur := make([]int64, n+1)
			for i := 0; i <= n; i++ {
				dpCur[i] = inf
			}
			dpCur[n] = 0
			st := newSegTree(n + 1)
			st.update(n, 0)
			costPer := int64(m - (k + 1))
			for i := n - 1; i >= 0; i-- {
				best := inf
				if k+1 < m && dpNext[i] < best {
					best = dpNext[i]
				}
				nxt := nextPos[k][i]
				if nxt != -1 {
					minVal := st.query(i+1, nxt+1)
					if minVal < inf {
						val := costPer + minVal
						if val < best {
							best = val
						}
					}
				}
				dpCur[i] = best
				st.update(i, best)
			}
			dpNext = dpCur
		}

		if dpNext[0] >= inf {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, dpNext[0])
		}
	}
}
