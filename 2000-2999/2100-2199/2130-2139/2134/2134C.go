package main

import (
	"bufio"
	"fmt"
	"os"
)

func max(a, b int64) int64 {
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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		var totalOps int64
		var sumOdd int64
		var S int64
		var maxPrevPrev int64 = 0
		var prevS int64 = 0

		for i := 1; i <= n; i++ {
			val := a[i-1]
			if i%2 == 1 {
				sumOdd += val
				S -= val
			} else {
				S += val
			}

			SAdj := S + totalOps
			if i >= 2 && SAdj < maxPrevPrev {
				need := maxPrevPrev - SAdj
				totalOps += need
				SAdj += need
			}

			maxPrevPrev = max(maxPrevPrev, prevS)
			prevS = SAdj
		}

		fmt.Fprintln(out, totalOps)
	}
}
