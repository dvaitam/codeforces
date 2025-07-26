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
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		if n == 3 && a[1]%2 == 1 {
			fmt.Fprintln(out, -1)
			continue
		}
		allOne := true
		for i := 1; i < n-1; i++ {
			if a[i] != 1 {
				allOne = false
				break
			}
		}
		if allOne {
			fmt.Fprintln(out, -1)
			continue
		}
		var res int64
		for i := 1; i < n-1; i++ {
			res += (a[i] + 1) / 2
		}
		fmt.Fprintln(out, res)
	}
}
