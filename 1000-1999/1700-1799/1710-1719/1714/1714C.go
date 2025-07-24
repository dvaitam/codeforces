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
		var s int
		fmt.Fscan(reader, &s)
		var digits []int
		for d := 9; d >= 1; d-- {
			if s >= d {
				digits = append(digits, d)
				s -= d
			}
		}
		for i := len(digits) - 1; i >= 0; i-- {
			fmt.Fprint(writer, digits[i])
		}
		fmt.Fprintln(writer)
	}
}
