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
		arr := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &arr[i])
		}
		target := arr[n-1]
		maxBefore := arr[0]
		for i := 1; i < n-1; i++ {
			if arr[i] > maxBefore {
				maxBefore = arr[i]
			}
		}
		if maxBefore == target-1 {
			fmt.Fprintln(out, target-1)
		} else {
			fmt.Fprintln(out, "Ambiguous")
		}
	}
}
