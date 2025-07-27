package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func compute(n int, a []int, b []int, c []int) []int64 {
	sort.Ints(a)
	sort.Ints(b)
	prefix := make([]int64, n+1)
	const INF int64 = 1 << 60
	prefix[0] = -INF
	for i := 1; i <= n; i++ {
		diff := int64(b[i-1] - a[i-1])
		if diff > prefix[i-1] {
			prefix[i] = diff
		} else {
			prefix[i] = prefix[i-1]
		}
	}
	suffix := make([]int64, n+2)
	suffix[n+1] = -INF
	for i := n; i >= 1; i-- {
		diff := int64(b[i] - a[i-1])
		if diff > suffix[i+1] {
			suffix[i] = diff
		} else {
			suffix[i] = suffix[i+1]
		}
	}
	res := make([]int64, len(c))
	for idx, x := range c {
		k := sort.SearchInts(a, x) + 1
		ans := prefix[k-1]
		diff := int64(b[k-1] - x)
		if diff > ans {
			ans = diff
		}
		if suffix[k] > ans {
			ans = suffix[k]
		}
		res[idx] = ans
	}
	return res
}
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n int
	fmt.Fscan(in, &n)
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}
	b := make([]int, n+1)
	for i := 0; i <= n; i++ {
		fmt.Fscan(in, &b[i])
	}
	var m int
	fmt.Fscan(in, &m)
	c := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &c[i])
	}
	ans := compute(n, a, b, c)
	for i, val := range ans {
		if i > 0 {
			out.WriteByte(' ')
		}
		fmt.Fprint(out, val)
	}
	out.WriteByte('\n')
}
