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
		arr := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}

		maxPrefix := arr[0]
		var maxDiff int64
		for i := 1; i < n; i++ {
			if arr[i] < maxPrefix {
				diff := maxPrefix - arr[i]
				if diff > maxDiff {
					maxDiff = diff
				}
			} else {
				maxPrefix = arr[i]
			}
		}

		ans := 0
		for maxDiff > 0 {
			ans++
			maxDiff >>= 1
		}
		fmt.Fprintln(out, ans)
	}
}
