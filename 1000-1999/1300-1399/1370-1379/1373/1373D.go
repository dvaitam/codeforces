package main

import (
	"bufio"
	"fmt"
	"os"
)

func kadaneDiff(diffs []int64) int64 {
	best := int64(0)
	cur := int64(0)
	for _, v := range diffs {
		cur += v
		if cur < 0 {
			cur = 0
		}
		if cur > best {
			best = cur
		}
	}
	return best
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &a[i])
		}

		base := int64(0)
		for i := 0; i < n; i += 2 {
			base += a[i]
		}

		diffs1 := make([]int64, 0, n/2)
		for i := 0; i+1 < n; i += 2 {
			diffs1 = append(diffs1, a[i+1]-a[i])
		}
		best1 := kadaneDiff(diffs1)

		diffs2 := make([]int64, 0, n/2)
		for i := 1; i+1 < n; i += 2 {
			diffs2 = append(diffs2, a[i]-a[i+1])
		}
		best2 := kadaneDiff(diffs2)

		best := best1
		if best2 > best {
			best = best2
		}

		fmt.Fprintln(writer, base+best)
	}
}
