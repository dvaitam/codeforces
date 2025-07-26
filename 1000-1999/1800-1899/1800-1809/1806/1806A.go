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
		var a, b, c, d int
		fmt.Fscan(reader, &a, &b, &c, &d)

		if d < b || c > a+d-b {
			fmt.Fprintln(writer, -1)
			continue
		}
		steps := (d - b) + (a + d - b - c)
		fmt.Fprintln(writer, steps)
	}
}
