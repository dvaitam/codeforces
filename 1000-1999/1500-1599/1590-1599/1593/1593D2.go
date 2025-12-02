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

	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}

	for i := 0; i < t; i++ {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		counts := make(map[int]int)
		maxCount := 0
		for j := 0; j < n; j++ {
			fmt.Fscan(in, &a[j])
			counts[a[j]]++
			if counts[a[j]] > maxCount {
				maxCount = counts[a[j]]
			}
		}

		if maxCount >= n/2 {
			fmt.Fprintln(out, -1)
			continue
		}

		maxK := 1
		for j := 0; j < n; j++ {
			pivot := a[j]
			needed := n/2 - counts[pivot]
			
			divCounts := make(map[int]int)
			for k := 0; k < n; k++ {
				if a[k] == pivot {
					continue
				}
				diff := a[k] - pivot
				if diff < 0 {
					diff = -diff
				}
				for d := 1; d*d <= diff; d++ {
					if diff%d == 0 {
						divCounts[d]++
						if d*d != diff {
							divCounts[diff/d]++
						}
					}
				}
			}

			for k, count := range divCounts {
				if count >= needed {
					if k > maxK {
						maxK = k
					}
				}
			}
		}
		fmt.Fprintln(out, maxK)
	}
}