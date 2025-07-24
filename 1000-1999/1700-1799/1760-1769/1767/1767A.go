package main

import (
	"bufio"
	"fmt"
	"os"
)

func between(a, b, c int64) bool {
	return (b < a && a < c) || (c < a && a < b)
}

func canCut(x1, y1, x2, y2, x3, y3 int64) bool {
	if between(x1, x2, x3) || between(y1, y2, y3) {
		return true
	}
	if between(x2, x1, x3) || between(y2, y1, y3) {
		return true
	}
	if between(x3, x1, x2) || between(y3, y1, y2) {
		return true
	}
	return false
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var x1, y1, x2, y2, x3, y3 int64
		fmt.Fscan(reader, &x1, &y1)
		fmt.Fscan(reader, &x2, &y2)
		fmt.Fscan(reader, &x3, &y3)
		if canCut(x1, y1, x2, y2, x3, y3) {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
