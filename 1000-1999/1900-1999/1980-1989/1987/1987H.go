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
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				var x int
				fmt.Fscan(reader, &x)
			}
		}
		var i0, j0 int
		fmt.Fscan(reader, &i0, &j0)
		fmt.Fprintf(writer, "%d %d\n", i0, j0)
	}
}
