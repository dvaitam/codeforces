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
		var n int
		var s string
		fmt.Fscan(reader, &n)
		fmt.Fscan(reader, &s)
		x, y := 0, 0
		reached := false
		for _, ch := range s {
			switch ch {
			case 'L':
				x--
			case 'R':
				x++
			case 'U':
				y++
			case 'D':
				y--
			}
			if x == 1 && y == 1 {
				reached = true
			}
		}
		if reached {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
