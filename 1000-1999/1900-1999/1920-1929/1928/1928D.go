package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		var b, x int64
		fmt.Fscan(reader, &n, &b, &x)

		counts := make([]int, n)
		maxC := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &counts[i])
			if counts[i] > maxC {
				maxC = counts[i]
			}
		}

		// frequency of each count
		freq := make(map[int]int)
		for _, v := range counts {
			freq[v]++
		}

		// difference arrays for constant and slope
		diffConst := make([]int64, maxC+2)
		diffSlope := make([]int64, maxC+2)

		for c, f := range freq {
			choose := int64(c*(c-1)) / 2
			limit := c
			if limit > maxC {
				limit = maxC
			}
			l := 1
			for l <= limit {
				q := c / l
				r := c / q
				if r > limit {
					r = limit
				}
				constPart := choose - int64(c)*int64(q)
				slopePart := int64(q*(q+1)) / 2
				diffConst[l] += int64(f) * constPart
				diffConst[r+1] -= int64(f) * constPart
				diffSlope[l] += int64(f) * slopePart
				diffSlope[r+1] -= int64(f) * slopePart
				l = r + 1
			}
			if c < maxC {
				diffConst[c+1] += int64(f) * choose
				diffConst[maxC+1] -= int64(f) * choose
			}
		}

		var constAcc, slopeAcc int64
		best := int64(-1 << 63)
		for k := 1; k <= maxC; k++ {
			constAcc += diffConst[k]
			slopeAcc += diffSlope[k]
			val := constAcc + slopeAcc*int64(k)
			strength := val*b - int64(k-1)*x
			if strength > best {
				best = strength
			}
		}
		fmt.Fprintln(writer, best)
	}
}
