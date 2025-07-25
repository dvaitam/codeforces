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
		var ones, cost int64
		for _, ch := range s {
			if ch == '1' {
				ones++
			} else if ones > 0 {
				cost += ones + 1
			}
		}
		fmt.Fprintln(writer, cost)
	}
}
