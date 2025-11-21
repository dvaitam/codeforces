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
		ans := make([]int, n)
		if k%2 == 1 {
			for i := 0; i < n-1; i++ {
				ans[i] = n
			}
			ans[n-1] = n - 1
		} else {
			if n >= 3 {
				for i := 0; i < n-2; i++ {
					ans[i] = n - 1
				}
			}
			if n >= 2 {
				ans[n-2] = n
				ans[n-1] = n - 1
			}
		}
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
