package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	var t int
	fmt.Fscan(in, &t)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		maxPrefix := a[0]
		var ans int64
		for i := 1; i < n; i++ {
			if a[i] < maxPrefix {
				if maxPrefix > ans {
					ans = maxPrefix
				}
			} else {
				maxPrefix = a[i]
			}
		}
		fmt.Fprintln(out, ans)
	}
}
