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
		var a, b int
		fmt.Fscan(reader, &a, &b)
		diff := a - b
		if diff < 0 {
			diff = -diff
		}
		ops := diff / 5
		diff %= 5
		ops += diff / 2
		diff %= 2
		ops += diff
		fmt.Fprintln(writer, ops)
	}
}
