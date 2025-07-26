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
	if _, err := fmt.Fscan(in, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		w := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &w[i])
		}
		i, j := 0, n-1
		left, right := 0, 0
		ans := 0
		for i <= j {
			if left <= right {
				left += w[i]
				i++
			} else {
				right += w[j]
				j--
			}
			if left == right {
				cand := i + (n - j - 1)
				if cand > ans {
					ans = cand
				}
			}
		}
		fmt.Fprintln(out, ans)
	}
}
