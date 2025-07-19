package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n int
	fmt.Fscan(reader, &n)
	a := make([]int, n)
	b := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &b[i])
	}
	m := n >> 1
	st := make([]int, n+5)
	tp := 0
	ans := make([]float64, n)
	vb := make([]float64, n)
	va := make([]float64, n)
	o := 0.0

	calc := func(x1, x2, y1, y2 int) float64 {
		return float64(n-1-2*x2)*float64(y1-y2)/float64(x1-x2) + float64(y2+y2)
	}

	// forward pass
	for i := n - 1; i >= 0; i-- {
		for tp > 1 && int64(a[i]-a[st[tp-1]])*int64(st[tp]-i) <= int64(a[i]-a[st[tp]])*int64(st[tp-1]-i) {
			tp--
		}
		if o > ans[i] {
			ans[i] = o
		}
		if tp > 0 && i < m && st[tp] >= m {
			o = calc(i, st[tp], a[i], a[st[tp]])
		}
		tp++
		st[tp] = i
		if i == m {
			for j := 0; j < m; j++ {
				l, r := 1, tp
				for l < r {
					mid := (l + r + 1) >> 1
					if int64(a[st[mid-1]]-b[j])*int64(st[mid]-j) > int64(a[st[mid]]-b[j])*int64(st[mid-1]-j) {
						r = mid - 1
					} else {
						l = mid
					}
				}
				vb[j] = calc(j, st[l], b[j], a[st[l]])
			}
		}
	}

	// backward pass
	tp = 0
	o = 0.0
	for i := 0; i < n; i++ {
		for tp > 1 && int64(b[i]-b[st[tp-1]])*int64(i-st[tp]) <= int64(b[i]-b[st[tp]])*int64(i-st[tp-1]) {
			tp--
		}
		if tp > 0 && st[tp] < m && i >= m {
			o = calc(st[tp], i, b[st[tp]], b[i])
		}
		if o > ans[i] {
			ans[i] = o
		}
		tp++
		st[tp] = i
		if i == m-1 {
			for j := m; j < n; j++ {
				l, r := 1, tp
				for l < r {
					mid := (l + r) >> 1
					if int64(a[j]-b[st[mid+1]])*int64(j-st[mid]) < int64(a[j]-b[st[mid]])*int64(j-st[mid+1]) {
						l = mid + 1
					} else {
						r = mid
					}
				}
				va[j] = calc(j, st[l], a[j], b[st[l]])
			}
		}
	}

	for i := 1; i < m; i++ {
		if vb[i-1] > vb[i] {
			vb[i] = vb[i-1]
		}
	}
	for i := n - 2; i >= m; i-- {
		if va[i+1] > va[i] {
			va[i] = va[i+1]
		}
	}
	for i := 0; i < n; i++ {
		var comp float64
		if i < m {
			comp = vb[i]
		} else if i+1 < n {
			comp = va[i+1]
		}
		res := comp
		if ans[i] > res {
			res = ans[i]
		}
		if i+1 == n {
			fmt.Fprintf(writer, "%.9f\n", res)
		} else {
			fmt.Fprintf(writer, "%.9f ", res)
		}
	}
}
