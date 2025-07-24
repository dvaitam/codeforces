package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	n    int
	a    []int64
	memo map[int]int
	full int
)

func abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func dfs(mask int) int {
	if mask == full {
		return 0
	}
	if v, ok := memo[mask]; ok {
		return v
	}
	// find first unpaired index
	i := 0
	for ; i < n; i++ {
		if mask>>i&1 == 0 {
			break
		}
	}
	best := dfs(mask | 1<<i) // leave i unmatched
	for j := i + 1; j < n; j++ {
		if mask>>j&1 == 0 && abs(a[i]-a[j]) == 1 {
			val := 1 + dfs(mask|1<<i|1<<j)
			if val > best {
				best = val
			}
		}
	}
	memo[mask] = best
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	a = make([]int64, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &a[i])
	}
	full = (1 << n) - 1
	memo = make(map[int]int)
	pairs := dfs(0)
	ans := n - pairs
	fmt.Println(ans)
}
