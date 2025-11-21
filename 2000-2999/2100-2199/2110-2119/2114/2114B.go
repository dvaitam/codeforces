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
		var n, k int
		fmt.Fscan(in, &n, &k)
		var s string
		fmt.Fscan(in, &s)
		c0 := 0
		for i := 0; i < n; i++ {
			if s[i] == '0' {
				c0++
			}
		}
		c1 := n - c0
		minGood := abs(c0-c1) / 2
		maxGood := c0/2 + c1/2
		if minGood <= k && k <= maxGood {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}
