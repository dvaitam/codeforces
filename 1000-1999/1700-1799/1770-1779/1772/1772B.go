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
		var a, b, c, d int
		fmt.Fscan(reader, &a, &b)
		fmt.Fscan(reader, &c, &d)
		ok := false
		for i := 0; i < 4; i++ {
			if a < b && c < d && a < c && b < d {
				ok = true
				break
			}
			// rotate 90 degrees clockwise
			a, b, c, d = c, a, d, b
		}
		if ok {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
