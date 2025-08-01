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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		for i := 0; i < m; i++ {
			var a, b, c int
			fmt.Fscan(reader, &a, &b, &c)
			// TODO: implement full algorithm
		}
		// output placeholder zeros
		for i := 0; i < m; i++ {
			if i > 0 {
				fmt.Fprint(writer, " ")
			}
			fmt.Fprint(writer, 0)
		}
		if t > 0 {
			fmt.Fprintln(writer)
		}
	}
}
