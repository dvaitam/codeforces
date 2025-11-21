package main

import (
	"bufio"
	"fmt"
	"os"
)

func min(a, b int) int {
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
		var n, x int
		fmt.Fscan(in, &n, &x)
		var s string
		fmt.Fscan(in, &s)

		left := 0
		for i := x - 2; i >= 0 && s[i] == '.'; i-- {
			left++
		}
		right := 0
		for i := x; i < n && s[i] == '.'; i++ {
			right++
		}

		ans := (left + right) / 2
		if left == right {
			ans++
		}
		if ans == 0 {
			ans = 1
		}
		fmt.Fprintln(out, ans)
	}
}

