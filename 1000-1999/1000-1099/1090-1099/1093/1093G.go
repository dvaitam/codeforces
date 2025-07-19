package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

var (
	reader *bufio.Reader
	writer *bufio.Writer
	n, k   int
	M      int
	data   []int64
	st     []int64
	qmax   []int64
)

func maxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func build(o, L, R int) {
	if L == R {
		base := o * M
		row := L * k
		for S := 0; S < M; S++ {
			var sum int64
			for j := 0; j < k; j++ {
				v := data[row+j]
				if (S>>j)&1 == 1 {
					sum += v
				} else {
					sum -= v
				}
			}
			st[base+S] = sum
		}
	} else {
		mid := (L + R) >> 1
		left := o << 1
		right := left | 1
		build(left, L, mid)
		build(right, mid+1, R)
		base := o * M
		lbase := left * M
		rbase := right * M
		for S := 0; S < M; S++ {
			st[base+S] = maxInt64(st[lbase+S], st[rbase+S])
		}
	}
}

func modify(o, L, R, idx int) {
	if L == R {
		base := o * M
		row := L * k
		for S := 0; S < M; S++ {
			var sum int64
			for j := 0; j < k; j++ {
				v := data[row+j]
				if (S>>j)&1 == 1 {
					sum += v
				} else {
					sum -= v
				}
			}
			st[base+S] = sum
		}
	} else {
		mid := (L + R) >> 1
		left := o << 1
		right := left | 1
		if idx <= mid {
			modify(left, L, mid, idx)
		} else {
			modify(right, mid+1, R, idx)
		}
		base := o * M
		lbase := left * M
		rbase := right * M
		for S := 0; S < M; S++ {
			st[base+S] = maxInt64(st[lbase+S], st[rbase+S])
		}
	}
}

func query(o, L, R, ql, qr int) {
	if ql <= L && R <= qr {
		base := o * M
		for S := 0; S < M; S++ {
			if st[base+S] > qmax[S] {
				qmax[S] = st[base+S]
			}
		}
	} else {
		mid := (L + R) >> 1
		left := o << 1
		right := left | 1
		if ql <= mid {
			query(left, L, mid, ql, qr)
		}
		if qr > mid {
			query(right, mid+1, R, ql, qr)
		}
	}
}

func main() {
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	fmt.Fscan(reader, &n, &k)
	M = 1 << k
	data = make([]int64, n*k)
	for i := 0; i < n; i++ {
		for j := 0; j < k; j++ {
			var x int64
			fmt.Fscan(reader, &x)
			data[i*k+j] = x
		}
	}
	st = make([]int64, 4*n*M)
	build(1, 0, n-1)
	qmax = make([]int64, M)
	var q int
	fmt.Fscan(reader, &q)
	inf := int64(math.MinInt64 / 2)
	for qi := 0; qi < q; qi++ {
		var op int
		fmt.Fscan(reader, &op)
		if op == 1 {
			var idx1 int
			fmt.Fscan(reader, &idx1)
			idx := idx1 - 1
			for j := 0; j < k; j++ {
				var x int64
				fmt.Fscan(reader, &x)
				data[idx*k+j] = x
			}
			modify(1, 0, n-1, idx)
		} else {
			var l1, r1 int
			fmt.Fscan(reader, &l1, &r1)
			l, r := l1-1, r1-1
			for S := 0; S < M; S++ {
				qmax[S] = inf
			}
			query(1, 0, n-1, l, r)
			ans := inf
			full := M - 1
			for S := 0; S < M; S++ {
				v := qmax[S] + qmax[full^S]
				if v > ans {
					ans = v
				}
			}
			fmt.Fprintln(writer, ans)
		}
	}
}
