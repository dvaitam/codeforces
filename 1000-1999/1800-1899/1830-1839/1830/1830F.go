package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxCostNaive(n, m int, intervals [][2]int, p []int) int64 {
	if m > 20 {
		// Too large for naive enumeration
		return 0
	}
	var best int64
	totalMasks := 1 << m
	for mask := 0; mask < totalMasks; mask++ {
		var sum int64
		for _, inter := range intervals {
			l, r := inter[0], inter[1]
			val := 0
			for x := r; x >= l; x-- {
				if mask&(1<<(x-1)) != 0 {
					val = p[x-1]
					break
				}
			}
			sum += int64(val)
		}
		if sum > best {
			best = sum
		}
	}
	return best
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(in, &n, &m)
		intervals := make([][2]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &intervals[i][0], &intervals[i][1])
		}
		p := make([]int, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &p[i])
		}
		ans := maxCostNaive(n, m, intervals, p)
		fmt.Fprintln(out, ans)
	}
}
