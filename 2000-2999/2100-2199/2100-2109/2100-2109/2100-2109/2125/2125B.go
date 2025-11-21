package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(a, b, k int64) int64 {
	if a <= k && b <= k {
		return 1
	}
	if a <= k || b <= k {
		if a <= k {
			steps := (b + k - 1) / k
			return steps
		}
		steps := (a + k - 1) / k
		return steps
	}
	if a != b {
		return 2
	}
	return 1
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var a, b, k int64
		fmt.Fscan(in, &a, &b, &k)
		if a < b {
			a, b = b, a
		}
		fmt.Fprintln(out, solve(a, b, k))
	}
}
