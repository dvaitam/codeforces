package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var n, m, k int
	if _, err := fmt.Fscan(os.Stdin, &n, &m, &k); err != nil {
		return
	}
	total := m * (m - 1) / 2
	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()
	fmt.Fprintln(w, total)
	if k == 0 {
		for i := 1; i <= m; i++ {
			for j := i + 1; j <= m; j++ {
				fmt.Fprintln(w, i, j)
			}
		}
	} else {
		for i := 1; i <= m; i++ {
			for j := i + 1; j <= m; j++ {
				fmt.Fprintln(w, j, i)
			}
		}
	}
}
