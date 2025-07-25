package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

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
		var c int
		fmt.Fscan(in, &n, &c)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		ans := make([]int, n)
		for target := 0; target < n; target++ {
			best := n + 1
			for mask := 0; mask < 1<<n; mask++ {
				if mask&(1<<target) != 0 {
					continue
				}
				undecided := c
				votes := make([]int, n)
				minIdx := -1
				for i := 0; i < n; i++ {
					if mask&(1<<i) != 0 {
						undecided += a[i]
					} else {
						votes[i] = a[i]
						if minIdx == -1 {
							minIdx = i
						}
					}
				}
				if minIdx != -1 {
					votes[minIdx] += undecided
				}
				win := -1
				bestVal := -1
				for i := 0; i < n; i++ {
					if mask&(1<<i) != 0 {
						continue
					}
					v := votes[i]
					if v > bestVal || (v == bestVal && i < win) {
						bestVal = v
						win = i
					}
				}
				if win == target {
					removed := bits.OnesCount(uint(mask))
					if removed < best {
						best = removed
					}
				}
			}
			ans[target] = best
		}
		for i, v := range ans {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
