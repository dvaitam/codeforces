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
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return
	}
	for ; t > 0; t-- {
		var a, b, c, d int64
		fmt.Fscan(reader, &a, &b, &c, &d)
		var x, y, x1, y1, x2, y2 int64
		fmt.Fscan(reader, &x, &y, &x1, &y1, &x2, &y2)

		if x1 == x2 && (a > 0 || b > 0) {
			fmt.Fprintln(writer, "NO")
			continue
		}
		if y1 == y2 && (c > 0 || d > 0) {
			fmt.Fprintln(writer, "NO")
			continue
		}

		nx := x + b - a
		ny := y + d - c
		if nx < x1 || nx > x2 || ny < y1 || ny > y2 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")
	}
}
