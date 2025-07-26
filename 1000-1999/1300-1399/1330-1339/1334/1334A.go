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
		var n int
		fmt.Fscan(reader, &n)
		prevP, prevC := 0, 0
		ok := true
		for i := 0; i < n; i++ {
			var p, c int
			fmt.Fscan(reader, &p, &c)
			if p < c {
				ok = false
			}
			if p < prevP || c < prevC {
				ok = false
			}
			if p-prevP < c-prevC {
				ok = false
			}
			prevP = p
			prevC = c
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
