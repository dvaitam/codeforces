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
		var c1, c2, c3 int64
		fmt.Fscan(reader, &c1, &c2, &c3)
		var a1, a2, a3, a4, a5 int64
		fmt.Fscan(reader, &a1, &a2, &a3, &a4, &a5)

		if c1 < a1 || c2 < a2 {
			fmt.Fprintln(writer, "NO")
			continue
		}
		c1 -= a1
		c2 -= a2

		needed := a3
		if a4 > c1 {
			needed += a4 - c1
		}
		if a5 > c2 {
			needed += a5 - c2
		}

		if c3 >= needed {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
