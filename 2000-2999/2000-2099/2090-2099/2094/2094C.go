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
		size := 2 * n
		ans := make([]int, size+1)
		usedVal := make([]bool, size+1)

		// Since G[i][j] = p_{i+j}, read once, assign p[i+j] = value.
		for i := 0; i < n; i++ {
			for j := 0; j < n; j++ {
				var val int
				fmt.Fscan(in, &val)
				index := i + j + 2
				if ans[index] == 0 {
					ans[index] = val
					usedVal[val] = true
				}
			}
		}

		// Find missing number
		missing := 0
		for v := 1; v <= size; v++ {
			if !usedVal[v] {
				missing = v
				break
			}
		}
		for i := 1; i <= size; i++ {
			if ans[i] == 0 {
				ans[i] = missing
			}
		}

		for i := 1; i <= size; i++ {
			if i > 1 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, ans[i])
		}
		fmt.Fprintln(out)
	}
}
