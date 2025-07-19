package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var m, n int
	if _, err := fmt.Fscan(os.Stdin, &m, &n); err != nil {
		return
	}
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	// Total steps: go across first row, then down first column (excluding starting cell twice)
	fmt.Fprintln(w, m+n-1)
	for i := 1; i <= n; i++ {
		fmt.Fprintln(w, 1, i)
	}
	for i := 2; i <= m; i++ {
		fmt.Fprintln(w, i, 1)
	}
}
