package main

import (
	"bufio"
	"fmt"
	"os"
)

var bad []int64
var limit int64 = 1_000_000_000_000_000_000

func generate() {
	primes := []int64{2, 3, 5, 7}
	var dfs func(idx int, val int64)
	dfs = func(idx int, val int64) {
		if val > limit {
			return
		}
		if idx == len(primes) {
			if val > 1 {
				bad = append(bad, val)
			}
			return
		}
		dfs(idx+1, val)
		mul := val
		for {
			mul *= primes[idx]
			if mul > limit || mul == 0 {
				break
			}
			dfs(idx+1, mul)
		}
	}
	dfs(0, 1)
}

func countBad(n int64) int64 {
	if n <= 0 {
		return 0
	}
	var ans int64
	for _, v := range bad {
		ans += n / v
	}
	return ans
}

func main() {
	generate()
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var l, r int64
		fmt.Fscan(in, &l, &r)
		ans := (r - l + 1) - (countBad(r) - countBad(l-1))
		fmt.Fprintln(out, ans)
	}
}
