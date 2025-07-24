package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func choose(k int) int64 {
	if k <= 0 {
		return 0
	}
	return int64(k) * int64(k+1) / 2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	p := make([]int, n)
	q := make([]int, n)
	posP := make([]int, n+2)
	posQ := make([]int, n+2)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &p[i])
		posP[p[i]] = i + 1
	}
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &q[i])
		posQ[q[i]] = i + 1
	}

	ans := int64(0)

	ip := posP[1]
	iq := posQ[1]
	a := min(ip, iq)
	b := max(ip, iq)
	ans += choose(a - 1)
	ans += choose(b - a - 1)
	ans += choose(n - b)

	L := n + 1
	R := 0
	for x := 2; x <= n+1; x++ {
		t := x - 1
		if posP[t] < L {
			L = posP[t]
		}
		if posQ[t] < L {
			L = posQ[t]
		}
		if posP[t] > R {
			R = posP[t]
		}
		if posQ[t] > R {
			R = posQ[t]
		}
		ip, iq = n+1, n+1
		if x <= n {
			ip = posP[x]
			iq = posQ[x]
		}
		a = min(ip, iq)
		b = max(ip, iq)
		segs := [][2]int{{1, a - 1}, {a + 1, b - 1}, {b + 1, n}}
		for _, s := range segs {
			l, r := s[0], s[1]
			if l <= r && l <= L && r >= R {
				ans += int64(L-l+1) * int64(r-R+1)
			}
		}
	}

	fmt.Fprintln(out, ans)
}
