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
		pos := 0
		double := false
		ok := true
		for pos < n {
			if double {
				if pos+1 >= n || s[pos] != s[pos+1] {
					ok = false
					break
				}
				pos += 2
			} else {
				pos++
			}
			double = !double
		}
		if pos != n {
			ok = false
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
