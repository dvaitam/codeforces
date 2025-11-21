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
		var w, h, a, b int64
		fmt.Fscan(reader, &w, &h, &a, &b)
		var x1, y1, x2, y2 int64
		fmt.Fscan(reader, &x1, &y1, &x2, &y2)

		dx := x1 - x2
		if dx < 0 {
			dx = -dx
		}
		dy := y1 - y2
		if dy < 0 {
			dy = -dy
		}

		if dx%a == 0 && dy%b == 0 {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
