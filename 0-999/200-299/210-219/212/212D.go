package main

import (
	"bufio"
	"fmt"
	"os"
)

type FastScanner struct {
	r *bufio.Reader
}

func NewFastScanner() *FastScanner {
	return &FastScanner{r: bufio.NewReaderSize(os.Stdin, 1<<20)}
}

func (fs *FastScanner) NextInt() int {
	sign := 1
	val := 0
	c, err := fs.r.ReadByte()
	for (c < '0' || c > '9') && c != '-' {
		if err != nil {
			return 0
		}
		c, err = fs.r.ReadByte()
	}
	if c == '-' {
		sign = -1
		c, err = fs.r.ReadByte()
	}
	for c >= '0' && c <= '9' {
		val = val*10 + int(c-'0')
		c, err = fs.r.ReadByte()
		if err != nil {
			break
		}
	}
	return sign * val
}

func addRange(l, r int, slope, intercept float64, diffSlope, diffIntercept []float64) {
	if l > r {
		return
	}
	diffSlope[l] += slope
	diffSlope[r+1] -= slope
	diffIntercept[l] += intercept
	diffIntercept[r+1] -= intercept
}

func main() {
	fs := NewFastScanner()
	n := fs.NextInt()
	if n == 0 {
		return
	}
	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		a[i] = int64(fs.NextInt())
	}

	prev := make([]int, n+1)
	stack := make([]int, 0, n)
	for i := 1; i <= n; i++ {
		for len(stack) > 0 && a[stack[len(stack)-1]] >= a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			prev[i] = 0
		} else {
			prev[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}

	next := make([]int, n+1)
	stack = stack[:0]
	for i := n; i >= 1; i-- {
		for len(stack) > 0 && a[stack[len(stack)-1]] > a[i] {
			stack = stack[:len(stack)-1]
		}
		if len(stack) == 0 {
			next[i] = n + 1
		} else {
			next[i] = stack[len(stack)-1]
		}
		stack = append(stack, i)
	}

	diffSlope := make([]float64, n+3)
	diffIntercept := make([]float64, n+3)

	for i := 1; i <= n; i++ {
		left := i - prev[i]
		right := next[i] - i
		small := left
		large := right
		if small > large {
			small, large = large, small
		}
		total := left + right - 1
		val := float64(a[i])

		addRange(1, small, val, 0, diffSlope, diffIntercept)
		if large > small {
			addRange(small+1, large, 0, val*float64(small), diffSlope, diffIntercept)
		}
		if total > large {
			addRange(large+1, total, -val, val*float64(left+right), diffSlope, diffIntercept)
		}
	}

	ans := make([]float64, n+1)
	curSlope := 0.0
	curIntercept := 0.0
	for k := 1; k <= n; k++ {
		curSlope += diffSlope[k]
		curIntercept += diffIntercept[k]
		ans[k] = curSlope*float64(k) + curIntercept
	}

	m := fs.NextInt()
	if m == 0 {
		return
	}
	out := bufio.NewWriterSize(os.Stdout, 1<<20)
	defer out.Flush()
	for i := 0; i < m; i++ {
		k := fs.NextInt()
		exp := ans[k] / float64(n-k+1)
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprintf(out, "%.15f", exp)
	}
}
