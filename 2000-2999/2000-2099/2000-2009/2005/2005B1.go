package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int64
		var m, q int
		fmt.Fscan(in, &n, &m, &q)

		teachers := make([]int64, m)
		for i := 0; i < m; i++ {
			fmt.Fscan(in, &teachers[i])
		}
		if teachers[0] > teachers[1] {
			teachers[0], teachers[1] = teachers[1], teachers[0]
		}

		for ; q > 0; q-- {
			var a int64
			fmt.Fscan(in, &a)

			var ans int64
			if a < teachers[0] {
				ans = teachers[0] - 1
			} else if a > teachers[1] {
				ans = n - teachers[1]
			} else {
				ans = min(a-teachers[0], teachers[1]-a)
			}
			fmt.Fprintln(out, ans)
		}
	}
}
