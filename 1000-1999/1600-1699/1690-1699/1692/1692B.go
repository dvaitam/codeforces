package main

import (
	"bufio"
	"fmt"
	"os"
)

// Solution to problemB.txt for 1692B (All Distinct).
// We remove elements in pairs to maximize the length of the
// resulting array with all unique values. The final size is the
// count of unique elements minus one if the number of duplicates is odd.
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
		uniq := make(map[int]struct{})
		for i := 0; i < n; i++ {
			var x int
			fmt.Fscan(in, &x)
			uniq[x] = struct{}{}
		}
		duplicates := n - len(uniq)
		if duplicates%2 == 1 {
			fmt.Fprintln(out, len(uniq)-1)
		} else {
			fmt.Fprintln(out, len(uniq))
		}
	}
}
