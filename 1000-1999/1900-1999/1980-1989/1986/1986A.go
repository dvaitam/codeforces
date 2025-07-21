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
		var x1, x2, x3 int
		fmt.Fscan(reader, &x1, &x2, &x3)
		minv := x1
		if x2 < minv {
			minv = x2
		}
		if x3 < minv {
			minv = x3
		}
		maxv := x1
		if x2 > maxv {
			maxv = x2
		}
		if x3 > maxv {
			maxv = x3
		}
		fmt.Fprintln(writer, maxv-minv)
	}
}
