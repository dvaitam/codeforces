package main

import (
	"bufio"
	"fmt"
	"os"
)

var (
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
)

func main() {
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for ; t > 0; t-- {
		var a, b, c, x, y int
		fmt.Fscan(reader, &a, &b, &c, &x, &y)
		remX := x - a
		if remX < 0 {
			remX = 0
		}
		remY := y - b
		if remY < 0 {
			remY = 0
		}
		if remX+remY <= c {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
