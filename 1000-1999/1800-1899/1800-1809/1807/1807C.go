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
		var pos [26]int
		for i := range pos {
			pos[i] = -1
		}
		ok := true
		for i, ch := range s {
			parity := i % 2
			idx := ch - 'a'
			if pos[idx] == -1 {
				pos[idx] = parity
			} else if pos[idx] != parity {
				ok = false
				break
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
