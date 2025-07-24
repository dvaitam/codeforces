package main

import (
	"bufio"
	"fmt"
	"os"
)

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int64
		fmt.Fscan(reader, &n, &m)
		if n == 1 && m == 1 {
			fmt.Fprintln(writer, 0)
			continue
		}
		ans := n + m - 1 + minInt64(n-1, m-1)
		fmt.Fprintln(writer, ans)
	}
}
