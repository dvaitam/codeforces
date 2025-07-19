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
	for i := 0; i < t; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
        if (x%2 == 1 && y%2 == 0) || (x%2 == 0 && y%2 == 1) {
			fmt.Fprintln(writer, "-1 -1")
		} else {
			if x%2 == 1 {
				a := x / 2
				b := (y + 1) / 2
				fmt.Fprintln(writer, a, b)
			} else {
				a := x / 2
				b := y / 2
				fmt.Fprintln(writer, a, b)
			}
		}
	}
}
