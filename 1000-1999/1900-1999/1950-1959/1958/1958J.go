package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}

	a := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &a[i])
	}

	b := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &b[i])
	}

	prefixB := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		prefixB[i] = prefixB[i-1] + b[i]
	}

	for ; q > 0; q-- {
		var l, r int
		fmt.Fscan(reader, &l, &r)
		var moves int64
		if l < r {
			for j := l + 1; j <= r; j++ {
				power := prefixB[j-1] - prefixB[l-1]
				moves += (a[j] + power - 1) / power
			}
		}
		fmt.Fprintln(writer, moves)
	}
}
