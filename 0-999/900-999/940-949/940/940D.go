package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	var s string
	fmt.Fscan(reader, &s)
	// 1-indexed b
	b := make([]byte, n+1)
	for i := 1; i <= n; i++ {
		b[i] = s[i-1]
	}
	const INF = 1000000000
	l, r := -INF, INF
	// iterate from i=5 to n
	for i := 5; i <= n; i++ {
		prev, curr := b[i-1], b[i]
		switch {
		case prev == '1' && curr == '1':
			// l = max(getMin(i), l)
			m := a[i]
			for j := i - 4; j <= i; j++ {
				if j >= 1 && a[j] < m {
					m = a[j]
				}
			}
			if m > l {
				l = m
			}
		case prev == '1' && curr == '0':
			// r = min(getMin(i)-1, r)
			m := a[i]
			for j := i - 4; j <= i; j++ {
				if j >= 1 && a[j] < m {
					m = a[j]
				}
			}
			if m-1 < r {
				r = m - 1
			}
			i += 3
		case prev == '0' && curr == '1':
			// l = max(getMax(i)+1, l)
			mx := a[i]
			for j := i - 4; j <= i; j++ {
				if j >= 1 && a[j] > mx {
					mx = a[j]
				}
			}
			if mx+1 > l {
				l = mx + 1
			}
			i += 3
		default:
			// prev=='0' && curr=='0'
			mx := a[i]
			for j := i - 4; j <= i; j++ {
				if j >= 1 && a[j] > mx {
					mx = a[j]
				}
			}
			if mx < r {
				r = mx
			}
		}
	}
	fmt.Printf("%d %d\n", l, r)
}
