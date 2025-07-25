package main

import (
	"bufio"
	"fmt"
	"os"
)

func canMake(a []int) bool {
	const INF int = int(1e9 + 7)
	min0 := -INF // minimal last value with no element removed yet
	min1 := INF  // minimal last value after exactly one removal
	for _, x := range a {
		next0 := INF
		if min0 != INF {
			cand := min(aMax(min0+1, x), x+1)
			if cand >= aMax(min0+1, x) {
				next0 = cand
			}
		}

		next1 := INF
		if min1 != INF {
			cand := min(aMax(min1+1, x), x+1)
			if cand >= aMax(min1+1, x) {
				next1 = cand
			}
		}
		if min0 != INF {
			if min0 < next1 {
				next1 = min0
			}
		}
		min0 = next0
		min1 = next1
	}
	return min1 != INF
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func aMax(a, b int) int {
	if a > b {
		return a
	}
	return b
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
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if canMake(a) {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}
