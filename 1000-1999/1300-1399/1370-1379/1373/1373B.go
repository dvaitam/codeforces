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
		var s string
		fmt.Fscan(reader, &s)
		zero, one := 0, 0
		for _, c := range s {
			if c == '0' {
				zero++
			} else if c == '1' {
				one++
			}
		}
		if min(zero, one)%2 == 1 {
			fmt.Fprintln(writer, "DA")
		} else {
			fmt.Fprintln(writer, "NET")
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
