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
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		bestL, bestR := 1, 1
		bestDelta := 0
		for l := 0; l < n; l++ {
			cntLess := 0
			cntGreater := 0
			for r := l + 1; r < n; r++ {
				if a[r] < a[l] {
					cntLess++
				} else if a[r] > a[l] {
					cntGreater++
				}
				delta := cntGreater - cntLess
				if delta < bestDelta {
					bestDelta = delta
					bestL = l + 1
					bestR = r + 1
				}
			}
		}
		fmt.Fprintln(out, bestL, bestR)
	}
}
