package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

type project struct {
	a int
	b int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, r int
	if _, err := fmt.Fscan(in, &n, &r); err != nil {
		return
	}
	pos := make([]project, 0)
	neg := make([]project, 0)
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		if b >= 0 {
			pos = append(pos, project{a, b})
		} else {
			neg = append(neg, project{a, b})
		}
	}

	sort.Slice(pos, func(i, j int) bool { return pos[i].a < pos[j].a })
	for _, p := range pos {
		if r < p.a {
			fmt.Fprintln(out, "NO")
			return
		}
		r += p.b
	}

	sort.Slice(neg, func(i, j int) bool { return neg[i].a+neg[i].b > neg[j].a+neg[j].b })
	m := len(neg)
	dp := make([]int, m+1)
	const negInf = -1 << 30
	for i := range dp {
		dp[i] = negInf
	}
	dp[0] = r
	for i, p := range neg {
		for j := i + 1; j >= 1; j-- {
			if dp[j-1] >= p.a && dp[j-1]+p.b >= 0 {
				cand := dp[j-1] + p.b
				if cand > dp[j] {
					dp[j] = cand
				}
			}
		}
	}
	if dp[m] >= 0 {
		fmt.Fprintln(out, "YES")
	} else {
		fmt.Fprintln(out, "NO")
	}
}
