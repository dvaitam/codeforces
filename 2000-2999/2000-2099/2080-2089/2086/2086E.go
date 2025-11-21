package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	coins []int64
	maxK  int
)

func buildCoins(limit int64) {
	val := int64(1)
	for val <= limit {
		coins = append(coins, val)
		val = 4*val + 1
	}
	maxK = len(coins) * 4
}

type state struct {
	idx  int
	need int
	rem  int64
}

func countUpTo(limit int64, k int64) int64 {
	if limit <= 0 || k < 0 {
		return 0
	}
	if k > int64(maxK) {
		return 0
	}
	memo := make(map[state]int64)
	var dfs func(int, int64, int) int64
	dfs = func(idx int, rem int64, need int) int64 {
		if need < 0 {
			return 0
		}
		if idx < 0 {
			if need == 0 {
				return 1
			}
			return 0
		}
		key := state{idx: idx, need: need, rem: rem}
		if val, ok := memo[key]; ok {
			return val
		}
		coin := coins[idx]
		maxD := rem / coin
		if maxD > 4 {
			maxD = 4
		}
		total := int64(0)
		var cap int64
		if idx > 0 {
			cap = coin - 1
		}
		for d := int64(0); d <= maxD; d++ {
			newRem := rem - d*coin
			if idx > 0 {
				if newRem > cap {
					newRem = cap
				}
			} else {
				newRem = 0
			}
			total += dfs(idx-1, newRem, need-int(d))
		}
		memo[key] = total
		return total
	}
	return dfs(len(coins)-1, limit, int(k))
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	buildCoins(1_000_000_000_000_000_000)

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for ; t > 0; t-- {
		var l, r, k int64
		fmt.Fscan(in, &l, &r, &k)
		res := countUpTo(r, k) - countUpTo(l-1, k)
		fmt.Fprintln(out, res)
	}
}
