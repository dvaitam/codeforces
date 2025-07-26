package main

import (
	"bufio"
	"fmt"
	"os"
)

func maxCagesKnown(n int) int {
	if n == 0 {
		return 0
	}
	if n%2 == 0 {
		return n/2 + 1
	}
	return (n + 1) / 2
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(in, &n)
		plan := make([]int, n)
		for i := range plan {
			fmt.Fscan(in, &plan[i])
		}

		known := 0
		unknown := 0
		ans := 0

		current := func() int { return maxCagesKnown(known) + unknown }
		ans = current()
		for _, b := range plan {
			if b == 1 {
				unknown++
			} else {
				known += unknown
				unknown = 0
			}
			if c := current(); c > ans {
				ans = c
			}
		}
		fmt.Fprintln(out, ans)
	}
}
