package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}
	a := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &a[i])
	}

	g := gcd(n, k)
	sums := make([]int, g)
	for i, v := range a {
		sums[i%g] = (sums[i%g] + v) % m
	}
	for i := 1; i < g; i++ {
		if sums[i] != sums[0] {
			fmt.Fprintln(out, -1)
			return
		}
	}
	if sum(a)%gcd(k, m) != 0 {
		fmt.Fprintln(out, -1)
		return
	}

	b := make([]int, n)
	b[0] = a[0] - a[n-1]
	for i := 1; i < n; i++ {
		b[i] = a[i] - a[i-1]
	}

	L := n / g
	ans := int64(0)
	for r := 0; r < g; r++ {
		pre := make([]int, L)
		cur := 0
		for t := 0; t < L; t++ {
			idx := (r + t*k) % n
			cur += b[idx]
			pre[t] = cur
		}
		sort.Ints(pre)
		med := pre[L/2]
		for _, v := range pre {
			if v > med {
				ans += int64(v - med)
			} else {
				ans += int64(med - v)
			}
		}
	}
	fmt.Fprintln(out, ans)
}

func sum(arr []int) int {
	s := 0
	for _, v := range arr {
		s += v
	}
	return s
}
