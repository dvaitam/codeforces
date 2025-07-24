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
		maxDelay := 0
		cur := 0
		seenA := false
		for _, ch := range s {
			if ch == 'A' {
				seenA = true
				cur = 0
			} else if seenA {
				cur++
				if cur > maxDelay {
					maxDelay = cur
				}
			}
		}
		fmt.Fprintln(writer, maxDelay)
	}
}
