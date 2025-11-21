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
		ans := make([]byte, n)
		curMin := int(1e9 + 10)
		for i := 0; i < n; i++ {
			if a[i] < curMin {
				ans[i] = '1'
				curMin = a[i]
			} else {
				ans[i] = '0'
			}
		}
		curMax := -1
		for i := n - 1; i >= 0; i-- {
			if a[i] > curMax {
				ans[i] = '1'
				curMax = a[i]
			}
		}
		fmt.Fprintln(out, string(ans))
	}
}
