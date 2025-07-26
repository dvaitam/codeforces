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
		var n, m int
		fmt.Fscan(reader, &n, &m)
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(reader, &grid[i])
		}
		ok := true
		for i := 0; i < n-1 && ok; i++ {
			for j := 0; j < m-1; j++ {
				cnt := 0
				if grid[i][j] == '1' {
					cnt++
				}
				if grid[i][j+1] == '1' {
					cnt++
				}
				if grid[i+1][j] == '1' {
					cnt++
				}
				if grid[i+1][j+1] == '1' {
					cnt++
				}
				if cnt == 3 {
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
