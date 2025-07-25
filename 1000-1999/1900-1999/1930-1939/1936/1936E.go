package main

import (
	"bufio"
	"fmt"
	"os"
)

const mod = 998244353

func countPerm(p []int) int64 {
	n := len(p)
	A := make([]int, n)
	m := 0
	for i, v := range p {
		if v > m {
			m = v
		}
		A[i] = m
	}
	used := make([]bool, n+1)
	var ans int64
	var dfs func(int, int)
	dfs = func(pos, mx int) {
		if pos == n {
			ans++
			if ans >= mod {
				ans -= mod
			}
			return
		}
		for x := 1; x <= n; x++ {
			if used[x] {
				continue
			}
			used[x] = true
			nmx := mx
			if x > nmx {
				nmx = x
			}
			if pos == n-1 || nmx != A[pos] {
				dfs(pos+1, nmx)
			}
			used[x] = false
		}
	}
	dfs(0, 0)
	return ans % mod
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		p := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &p[i])
		}
		if n > 8 {
			fmt.Fprintln(out, 0)
			continue
		}
		res := countPerm(p)
		fmt.Fprintln(out, res)
	}
}
