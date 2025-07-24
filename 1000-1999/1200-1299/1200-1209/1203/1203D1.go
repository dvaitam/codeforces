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
	for i := 0; i < m; i++ {
		for idx < n && s[idx] != t[i] {
			idx++
		}
		if idx == n {
			// should not happen as t is subsequence
			pre[i] = n
		} else {
			pre[i] = idx
			idx++
		}
	}

	suf := make([]int, m)
	idx = n - 1
	for i := m - 1; i >= 0; i-- {
		for idx >= 0 && s[idx] != t[i] {
			idx--
		}
		if idx < 0 {
			suf[i] = -1
		} else {
			suf[i] = idx
			idx--
		}
	}

	maxDel := 0
	for i := 0; i <= m; i++ {
		var l, r int
		if i == 0 {
			l = 0
		} else {
			l = pre[i-1] + 1
		}
		if i == m {
			r = n - 1
		} else {
			r = suf[i] - 1
		}
		if r >= l {
			if r-l+1 > maxDel {
				maxDel = r - l + 1
			}
		}
	}

	fmt.Fprintln(out, maxDel)
}
