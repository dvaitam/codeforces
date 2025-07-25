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

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}

		prefixMax := make([]int, n)
		suffixMin := make([]int, n)
		maxVal := a[0]
		for i := 0; i < n; i++ {
			if a[i] > maxVal {
				maxVal = a[i]
			}
			prefixMax[i] = maxVal
		}
		minVal := a[n-1]
		for i := n - 1; i >= 0; i-- {
			if a[i] < minVal {
				minVal = a[i]
			}
			suffixMin[i] = minVal
		}

		res := "NO"
		for i := 0; i < n-1; i++ {
			if prefixMax[i] > suffixMin[i+1] {
				res = "YES"
				break
			}
		}
		fmt.Fprintln(out, res)
	}
}
