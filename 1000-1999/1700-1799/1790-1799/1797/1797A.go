package main

import (
	"bufio"
	"fmt"
	"os"
)

func degree(n, m, x, y int) int {
	d := 4
	if x == 1 {
		d--
	}
	if x == n {
		d--
	}
	if y == 1 {
		d--
	}
	if y == m {
		d--
	}
	return d
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		var x1, y1, x2, y2 int
		fmt.Fscan(reader, &x1, &y1, &x2, &y2)

		d1 := degree(n, m, x1, y1)
		d2 := degree(n, m, x2, y2)
		if d1 < d2 {
			fmt.Fprintln(writer, d1)
		} else {
			fmt.Fprintln(writer, d2)
		}
	}
}
