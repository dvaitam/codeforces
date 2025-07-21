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
		var n, x, y int
		fmt.Fscan(reader, &n, &x, &y)
		for i := 1; i < y; i++ {
			if (y-i)%2 == 1 {
				fmt.Fprint(writer, "-1 ")
			} else {
				fmt.Fprint(writer, "1 ")
			}
		}
		for i := y; i <= x; i++ {
			fmt.Fprint(writer, "1 ")
		}
		for i := x + 1; i <= n; i++ {
			if (i-x)%2 == 1 {
				fmt.Fprint(writer, "-1 ")
			} else {
				fmt.Fprint(writer, "1 ")
			}
		}
		fmt.Fprintln(writer)
	}
}
