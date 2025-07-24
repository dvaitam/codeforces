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

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var b, g, n int
	if _, err := fmt.Fscan(in, &b, &g, &n); err != nil {
		return
	}

	l := max(0, n-g)
	r := min(n, b)

	fmt.Fprintln(out, r-l+1)
}
