package main

import (
	"bufio"
	"fmt"
	"os"
)

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

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
		time1 := abs(a - 1)
		time2 := abs(b-c) + abs(c-1)
		if time1 < time2 {
			fmt.Fprintln(writer, 1)
		} else if time2 < time1 {
			fmt.Fprintln(writer, 2)
		} else {
			fmt.Fprintln(writer, 3)
		}
	}
}
