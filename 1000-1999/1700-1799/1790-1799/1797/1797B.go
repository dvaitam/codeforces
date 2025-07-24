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
		var k int
		fmt.Fscan(reader, &n, &k)
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, n)
			for j := 0; j < n; j++ {
				fmt.Fscan(reader, &grid[i][j])
			}
		}
		mismatches := 0
		for i := 0; i < n; i++ {
			ni := n - 1 - i
			for j := 0; j < n; j++ {
				nj := n - 1 - j
				if i < ni || (i == ni && j < nj) {
					if grid[i][j] != grid[ni][nj] {
						mismatches++
					}
				}
			}
		}
		if k < mismatches {
			fmt.Fprintln(writer, "NO")
			continue
		}
		if n%2 == 1 {
			fmt.Fprintln(writer, "YES")
		} else if (k-mismatches)%2 == 0 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
