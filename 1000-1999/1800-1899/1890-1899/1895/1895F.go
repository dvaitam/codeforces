package main

import (
	"bufio"
	"fmt"
	"os"
)

type state struct {
	pos  int
	val  int64
	seen bool
}

const MOD int64 = 1_000_000_007

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		var x, k int64
		fmt.Fscan(in, &n, &x, &k)
		ans := solve(int64(n), x, k)
		fmt.Fprintln(out, ans)
	}
}

func solve(n, x, k int64) int64 {
	maxVal := x + k*(n-1)
	memo := make(map[state]int64)
	var dfs func(pos int, val int64, seen bool) int64
	dfs = func(pos int, val int64, seen bool) int64 {
		if val < 0 || val > maxVal {
			return 0
		}
		if pos == int(n) {
			if seen {
				return 1
			}
			return 0
		}
		st := state{pos, val, seen}
		if v, ok := memo[st]; ok {
			return v
		}
		var res int64
		for d := -k; d <= k; d++ {
			nv := val + d
			ns := seen || (nv >= x && nv <= x+k-1)
			res = (res + dfs(pos+1, nv, ns)) % MOD
		}
		memo[st] = res
		return res
	}

	var total int64
	for start := int64(0); start <= maxVal; start++ {
		total = (total + dfs(1, start, start >= x && start <= x+k-1)) % MOD
	}
	return total
}
