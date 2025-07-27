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
		fmt.Fscan(reader, &n)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &grid[i])
		}
		ok := true
		for i := 0; i < n && ok; i++ {
			for j := 0; j < n; j++ {
				if grid[i][j] == '1' {
					if i == n-1 || j == n-1 {
						continue
					}
					if grid[i+1][j] != '1' && grid[i][j+1] != '1' {
						ok = false
						break
					}
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
