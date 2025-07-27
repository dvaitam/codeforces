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

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, b, n, m int64
		fmt.Fscan(reader, &a, &b, &n, &m)
		if a+b < n+m || min(a, b) < m {
			fmt.Fprintln(writer, "No")
		} else {
			fmt.Fprintln(writer, "Yes")
		}
	}
}

func min(x, y int64) int64 {
	if x < y {
		return x
	}
	return y
}
