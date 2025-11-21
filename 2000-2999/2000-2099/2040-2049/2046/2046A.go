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
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)

		top := make([]int64, n)
		totalTop := int64(0)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &top[i])
			totalTop += top[i]
		}

		totalBottom := int64(0)
		sumMin := int64(0)
		maxMin := int64(-1 << 60)
		for i := 0; i < n; i++ {
			var val int64
			fmt.Fscan(in, &val)
			totalBottom += val
			minVal := top[i]
			if val < minVal {
				minVal = val
			}
			sumMin += minVal
			if minVal > maxMin {
				maxMin = minVal
			}
		}

		answer := totalTop + totalBottom - sumMin + maxMin
		fmt.Fprintln(out, answer)
	}
}
