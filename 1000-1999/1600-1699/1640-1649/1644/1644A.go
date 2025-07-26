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
		var s string
		fmt.Fscan(reader, &s)
		keys := make(map[rune]bool)
		ok := true
		for _, ch := range s {
			if ch == 'r' || ch == 'g' || ch == 'b' {
				keys[ch] = true
			} else {
				if !keys[ch+32] {
					ok = false
					break
				}
			}
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
