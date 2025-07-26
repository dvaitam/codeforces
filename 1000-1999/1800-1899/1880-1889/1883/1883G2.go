package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	if _, err := fmt.Fscan(in, &T); err != nil {
		return
	}
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		a := make([]int, n-1)
		for i := 0; i < n-1; i++ {
			fmt.Fscan(in, &a[i])
		}
		b := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &b[i])
		}
		sort.Ints(a)
		sort.Ints(b)
		pre := make([]int, n)
		for j := 0; j < n; j++ {
			pre[j] = sort.SearchInts(a, b[j])
		}
		hist := make([]int, n+1)
		used := 0
		for j := 0; j < n; j++ {
			if pre[j] > used {
				used++
			}
			hist[j+1] = used
		}
		k0 := hist[n]
		good := make([]int, n+1)
		exist := false
		for p := n; p >= 0; p-- {
			if p < n && pre[p] == hist[p] {
				exist = true
			}
			if exist {
				good[p] = 1
			}
		}
		counts := make([]int64, n+1)
		prev := 1
		for p := 0; p < n; p++ {
			nxt := m
			if b[p]-1 < m {
				nxt = b[p] - 1
			}
			if nxt >= prev {
				counts[p] = int64(nxt - prev + 1)
			}
			if b[p] > prev {
				prev = b[p]
			}
		}
		if prev <= m {
			counts[n] = int64(m - prev + 1)
		}
		var ans int64
		for p := 0; p <= n; p++ {
			if counts[p] > 0 {
				k := k0 + good[p]
				ops := n - k
				ans += counts[p] * int64(ops)
			}
		}
		fmt.Fprintln(out, ans)
	}
}
