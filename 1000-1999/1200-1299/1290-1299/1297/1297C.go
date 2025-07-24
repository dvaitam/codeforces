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
		a := make([]int, n)
		sumPos := 0
		minAbs := 10001
		idx := -1
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
			if a[i] > 0 {
				sumPos += a[i]
			}
			if a[i] != 0 {
				x := a[i]
				if x < 0 {
					x = -x
				}
				if x < minAbs {
					minAbs = x
					idx = i
				}
			}
		}
		ans := make([]byte, n)
		for i := 0; i < n; i++ {
			if a[i] > 0 {
				ans[i] = '1'
			} else {
				ans[i] = '0'
			}
		}
		if a[idx] > 0 {
			ans[idx] = '0'
			sumPos -= a[idx]
		} else {
			ans[idx] = '1'
			sumPos += a[idx]
		}
		fmt.Fprintln(out, sumPos)
		fmt.Fprintln(out, string(ans))
	}
}
