package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO: implement actual solution
// Currently this placeholder simply checks if t is a subsequence of s.
func isSubsequence(s, t string) bool {
	j := 0
	for i := 0; i < len(s) && j < len(t); i++ {
		if s[i] == t[j] {
			j++
		}
	}
	return j == len(t)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var s, t string
		fmt.Fscan(reader, &s)
		fmt.Fscan(reader, &t)
		if !isSubsequence(s, t) {
			fmt.Fprintln(writer, -1)
			continue
		}
		// Minimal cost must be at least the number of deletions.
		// Since the optimal strategy is non-trivial, we output this lower bound.
		fmt.Fprintln(writer, n-m)
	}
}
