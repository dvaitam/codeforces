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

	var n int
	if _, err := fmt.Fscan(in, &n); err != nil {
		return
	}

	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &a[i])
	}

	prefPos := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefPos[i] = prefPos[i-1]
		if a[i] > 0 {
			prefPos[i] += a[i]
		}
	}

	bestSum := int64(-1 << 62)
	bestI, bestJ := 1, 2
	for i := 1; i <= n; i++ {
		for j := i + 1; j <= n; j++ {
			if a[i] == a[j] {
				cur := a[i] + a[j] + prefPos[j-1] - prefPos[i]
				if cur > bestSum {
					bestSum = cur
					bestI, bestJ = i, j
				}
			}
		}
	}

	cuts := make([]int, 0)
	for idx := 1; idx < bestI; idx++ {
		cuts = append(cuts, idx)
	}
	for idx := bestI + 1; idx < bestJ; idx++ {
		if a[idx] < 0 {
			cuts = append(cuts, idx)
		}
	}
	for idx := bestJ + 1; idx <= n; idx++ {
		cuts = append(cuts, idx)
	}

	fmt.Fprintf(out, "%d %d\n", bestSum, len(cuts))
	if len(cuts) > 0 {
		for i, idx := range cuts {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, idx)
		}
		fmt.Fprintln(out)
	} else {
		fmt.Fprintln(out)
	}
}
