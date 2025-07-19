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
	for i := 0; i < t; i++ {
		var x int64
		fmt.Fscan(reader, &x)
		if x%2 == 1 {
			fmt.Fprintln(writer, -1)
		} else {
			a := x / 2
			b := 2*x - a
			if a^b == x {
				fmt.Fprintln(writer, a, b)
			} else {
				fmt.Fprintln(writer, -1)
			}
		}
	}
}
