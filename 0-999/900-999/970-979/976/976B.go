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

	var n, m, k int64
	if _, err := fmt.Fscan(in, &n, &m, &k); err != nil {
		return
	}

	if k < n {
		fmt.Fprintf(out, "%d %d\n", k+1, 1)
		return
	}
	k -= n - 1
	if k < m {
		fmt.Fprintf(out, "%d %d\n", n, k+1)
		return
	}
	k -= m
	m1 := m - 1
	rowOffset := k / m1
	rem := k % m1
	row := n - 1 - rowOffset
	var col int64
	if rowOffset%2 == 0 {
		col = m - rem
	} else {
		col = 2 + rem
	}
	fmt.Fprintf(out, "%d %d\n", row, col)
}
