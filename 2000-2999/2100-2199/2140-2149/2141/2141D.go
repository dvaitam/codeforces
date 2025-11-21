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
		var k int64
		fmt.Fscan(in, &n, &k)
		a := make([]int64, n)
		var sum int64
		var minVal int64 = 1<<62 - 1
		var maxVal int64
		cntMin := 0
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			sum += a[i]
			if a[i] < minVal {
				minVal = a[i]
				cntMin = 1
			} else if a[i] == minVal {
				cntMin++
			}
			if a[i] > maxVal {
				maxVal = a[i]
			}
		}

		target := (sum + k) / int64(n)
		if target < maxVal {
			fmt.Fprintln(out, -1)
			continue
		}
		coins := int64(n-cntMin)*target - sum + int64(cntMin)*minVal
		fmt.Fprintln(out, coins)
	}
}

