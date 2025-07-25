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
		var n, a, b int
		fmt.Fscan(in, &n, &a, &b)
		// Cost if we buy two with promotion vs individually
		pairCost := b
		if 2*a < b {
			pairCost = 2 * a
		}
		cost := (n/2)*pairCost + (n%2)*a
		fmt.Fprintln(out, cost)
	}
}
