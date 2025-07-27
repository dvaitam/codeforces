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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}
	g := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		g[i] = []byte(s)
	}

	total := 1 << (n - 1)
	cnt := make([]int64, total)

	perm := make([]int, n)
	used := make([]bool, n)

	var dfs func(int)
	dfs = func(pos int) {
		if pos == n {
			var mask int
			for i := 0; i < n-1; i++ {
				if g[perm[i]][perm[i+1]] == '1' {
					mask |= 1 << i
				}
			}
			cnt[mask]++
			return
		}
		for i := 0; i < n; i++ {
			if !used[i] {
				used[i] = true
				perm[pos] = i
				dfs(pos + 1)
				used[i] = false
			}
		}
	}

	dfs(0)

	for i := 0; i < total; i++ {
		if i > 0 {
			fmt.Fprint(out, " ")
		}
		fmt.Fprint(out, cnt[i])
	}
	fmt.Fprintln(out)
}
