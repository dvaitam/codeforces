package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	if _, err := fmt.Fscan(in, &s); err != nil {
		return
	}
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	n := len(s)
	m := len(t)
	pre := make([]int, m)
	idx := 0
	for i := 0; i < n && idx < m; i++ {
		if s[i] == t[idx] {
			pre[idx] = i
			idx++
		}
	}

	suf := make([]int, m)
	idx = m - 1
	for i := n - 1; i >= 0 && idx >= 0; i-- {
		if s[i] == t[idx] {
			suf[idx] = i
			idx--
		}
	}

	ans := max(suf[0], n-1-pre[m-1])
	for i := 0; i < m-1; i++ {
		if gap := suf[i+1] - pre[i] - 1; gap > ans {
			ans = gap
		}
	}

	fmt.Fprintln(out, ans)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
