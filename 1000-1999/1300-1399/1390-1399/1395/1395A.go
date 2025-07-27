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
		var r, g, b, w int
		fmt.Fscan(reader, &r, &g, &b, &w)
		odds := r&1 + g&1 + b&1 + w&1
		ok := false
		if odds <= 1 {
			ok = true
		} else if r > 0 && g > 0 && b > 0 {
			// perform one operation: flip parity of r,g,b,w
			// new odd count = 4 - odds
			if 4-odds <= 1 {
				ok = true
			}
		}
		if ok {
			fmt.Fprintln(writer, "Yes")
		} else {
			fmt.Fprintln(writer, "No")
		}
	}
}
