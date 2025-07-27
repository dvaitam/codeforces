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
		var n int
		var s string
		fmt.Fscan(reader, &n, &s)
		moves := 0
		bal := 0
		for _, ch := range s {
			if ch == '(' {
				bal++
			} else {
				if bal > 0 {
					bal--
				} else {
					moves++
				}
			}
		}
		fmt.Fprintln(writer, moves)
	}
}
