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
		var n, k int
		fmt.Fscan(reader, &n, &k)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &grid[i])
		}
		for i := 0; i < n; i += k {
			for j := 0; j < n; j += k {
				fmt.Fprintf(writer, "%c", grid[i][j])
			}
			fmt.Fprintln(writer)
		}
	}
}
