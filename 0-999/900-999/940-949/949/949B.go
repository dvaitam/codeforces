package main

import (
	"bufio"
	"fmt"
	"os"
)

func solve(n, x int64) int64 {
	if x%2 == 1 {
		return (x + 1) / 2
	}
	if n%2 == 0 {
		return n/2 + solve(n/2, x/2)
	}
	m := n / 2
	i := x / 2
	if i == 1 {
		return m + 1 + solve(m, m)
	}
	return m + 1 + solve(m, i-1)
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, q int64
	if _, err := fmt.Fscan(reader, &n, &q); err != nil {
		return
	}
	for ; q > 0; q-- {
		var x int64
		fmt.Fscan(reader, &x)
		fmt.Fprintln(writer, solve(n, x))
	}
}
