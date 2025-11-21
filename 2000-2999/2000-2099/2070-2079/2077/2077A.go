package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

		size := 2 * n
		vals := make([]int64, size)
		for i := range vals {
			fmt.Fscan(in, &vals[i])
		}

		sort.Slice(vals, func(i, j int) bool { return vals[i] < vals[j] })

		var sumSmall int64
		for i := 0; i < n-1; i++ {
			sumSmall += vals[i]
		}
		var sumRest int64
		for i := n - 1; i < size; i++ {
			sumRest += vals[i]
		}

		x := sumRest - sumSmall

		rest := make([]int64, size-(n-1))
		copy(rest, vals[n-1:])

		evenVals := make([]int64, n)
		evenVals[0] = x
		copy(evenVals[1:], vals[:n-1])

		answer := make([]int64, 2*n+1)
		oddIdx, evenIdx := 0, 0
		for i := range answer {
			if i&1 == 0 {
				answer[i] = rest[oddIdx]
				oddIdx++
			} else {
				answer[i] = evenVals[evenIdx]
				evenIdx++
			}
		}

		for i, v := range answer {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, v)
		}
		fmt.Fprintln(out)
	}
}
