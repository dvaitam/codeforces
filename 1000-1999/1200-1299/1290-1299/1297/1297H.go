package main

import (
	"bufio"
	"fmt"
	"os"
)

type pair struct{ r, b string }

func lexCompare(a, b string) int {
	if a == b {
		return 0
	}
	if a < b {
		return -1
	}
	return 1
}

func maxStr(a, b string) string {
	if lexCompare(a, b) >= 0 {
		return a
	}
	return b
}

func solveCase(s string) string {
	n := len(s)
	dp := make([][]*pair, n+1)
	col := make([][]string, n+1)
	for i := 0; i <= n; i++ {
		dp[i] = make([]*pair, n+1)
		col[i] = make([]string, n+1)
	}
	dp[0][0] = &pair{"", ""}
	col[0][0] = ""

	for i := 0; i < n; i++ {
		ch := s[i : i+1]
		for rlen := i; rlen >= 0; rlen-- {
			p := dp[i][rlen]
			if p == nil {
				continue
			}
			// assign to red
			r1 := p.r + ch
			b1 := p.b
			candMax := maxStr(r1, b1)
			if old := dp[i+1][rlen+1]; old == nil || lexCompare(candMax, maxStr(old.r, old.b)) < 0 {
				dp[i+1][rlen+1] = &pair{r1, b1}
				col[i+1][rlen+1] = col[i][rlen] + "R"
			}
			// assign to blue
			r2 := p.r
			b2 := p.b + ch
			candMax = maxStr(r2, b2)
			if old := dp[i+1][rlen]; old == nil || lexCompare(candMax, maxStr(old.r, old.b)) < 0 {
				dp[i+1][rlen] = &pair{r2, b2}
				col[i+1][rlen] = col[i][rlen] + "B"
			}
		}
	}

	bestMax := ""
	bestCol := ""
	first := true
	for rlen := 0; rlen <= n; rlen++ {
		p := dp[n][rlen]
		if p == nil {
			continue
		}
		m := maxStr(p.r, p.b)
		if first || lexCompare(m, bestMax) < 0 {
			bestMax = m
			bestCol = col[n][rlen]
			first = false
		}
	}
	return bestCol
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var s string
		fmt.Fscan(in, &s)
		fmt.Fprintln(out, solveCase(s))
	}
}
