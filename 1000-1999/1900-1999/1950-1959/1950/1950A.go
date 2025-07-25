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
		var a, b, c int
		fmt.Fscan(reader, &a, &b, &c)
		if a < b && b < c {
			fmt.Fprintln(writer, "STAIR")
		} else if a < b && b > c {
			fmt.Fprintln(writer, "PEAK")
		} else {
			fmt.Fprintln(writer, "NONE")
		}
	}
}
