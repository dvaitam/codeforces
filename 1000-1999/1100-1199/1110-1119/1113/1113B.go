package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	var n int
	fmt.Fscan(reader, &n)
	b := [101]int{}
	var ans int
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(reader, &v)
		if v >= 1 && v <= 100 {
			b[v]++
		}
		ans += v
	}
	// Precompute divisors (excluding 1 and itself)
	adj := make([][]int, 101)
	for j := 1; j <= 100; j++ {
		for d := 2; d*d <= j; d++ {
			if j%d == 0 {
				adj[j] = append(adj[j], d)
				if d*d != j {
					adj[j] = append(adj[j], j/d)
				}
			}
		}
	}
	// initial best is no operation
	best := ans
	// try all possible operations
	for j := 1; j <= 100; j++ {
		if b[j] == 0 || len(adj[j]) == 0 {
			continue
		}
		for _, d := range adj[j] {
			// move factor d from j to i
			for i := 1; i <= 100; i++ {
				if b[i] == 0 {
					continue
				}
				pre := i + j
				cur := i*d + j/d
				if cur < pre {
					cand := ans - pre + cur
					if cand < best {
						best = cand
					}
				}
			}
		}
	}
	fmt.Println(best)
}
