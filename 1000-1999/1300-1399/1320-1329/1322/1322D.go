package main

import (
	"bufio"
	"fmt"
	"os"
)

// Greedy attempt for problem described in problemD.txt.
// We process candidates in order and recruit a candidate only
// if their addition immediately increases total profit while
// respecting the non-increasing aggressiveness constraint.

func simulateGain(cnt []int, level int, c []int) int {
	gain := c[level]
	cnt[level]++
	for lvl := level; cnt[lvl] == 2; lvl++ {
		cnt[lvl] = 0
		if lvl+1 >= len(cnt) {
			cnt = append(cnt, 0)
		}
		cnt[lvl+1]++
		gain += c[lvl+1]
	}
	return gain
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	l := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &l[i])
	}
	s := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &s[i])
	}
	c := make([]int, n+m+5)
	for i := 1; i <= n+m; i++ {
		fmt.Fscan(in, &c[i])
	}

	cnt := make([]int, n+m+5)
	profit := 0
	cost := 0
	maxAllowed := n + m + 5

	for i := 0; i < n; i++ {
		if l[i] > maxAllowed {
			continue
		}
		tmp := make([]int, len(cnt))
		copy(tmp, cnt)
		g := simulateGain(tmp, l[i], c)
		if profit+g-(cost+s[i]) > profit-cost {
			cnt = tmp
			profit += g
			cost += s[i]
			if l[i] < maxAllowed {
				maxAllowed = l[i]
			}
		}
	}

	fmt.Fprintln(out, profit-cost)
}
