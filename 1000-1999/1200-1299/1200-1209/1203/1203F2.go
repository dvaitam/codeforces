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
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, r int
	if _, err := fmt.Fscan(reader, &n, &r); err != nil {
		return
	}

	pos := make([]project, 0)
	neg := make([]project, 0)
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(reader, &a, &b)
		if b >= 0 {
			pos = append(pos, project{a, b})
		} else {
			neg = append(neg, project{a, b})
		}
	}

	// complete all achievable positive projects greedily
	sort.Slice(pos, func(i, j int) bool {
		if pos[i].a == pos[j].a {
			return pos[i].b > pos[j].b
		}
		return pos[i].a < pos[j].a
	})

	count := 0
	curr := r
	for _, p := range pos {
		if curr >= p.a {
			curr += p.b
			count++
		}
	}

	// sort negative projects by a+b descending
	sort.Slice(neg, func(i, j int) bool {
		left := neg[i].a + neg[i].b
		right := neg[j].a + neg[j].b
		if left == right {
			return neg[i].a > neg[j].a
		}
		return left > right
	})

	maxR := curr
	const negInf = -1000000000
	dp := make([]int, maxR+1)
	for i := range dp {
		dp[i] = negInf
	}
	dp[curr] = 0

	for _, p := range neg {
		for rating := maxR; rating >= p.a; rating-- {
			if dp[rating] == negInf {
				continue
			}
			newR := rating + p.b
			if newR < 0 {
				continue
			}
			if dp[rating]+1 > dp[newR] {
				dp[newR] = dp[rating] + 1
			}
		}
	}

	best := 0
	for _, v := range dp {
		if v > best {
			best = v
		}
	}

	fmt.Fprintln(writer, best+count)
}
