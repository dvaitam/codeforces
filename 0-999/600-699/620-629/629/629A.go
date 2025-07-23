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

	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return
	}

	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grid[i])
	}

	ans := 0
	// count pairs in rows
	for i := 0; i < n; i++ {
		cnt := 0
		for j := 0; j < n; j++ {
			if grid[i][j] == 'C' {
				cnt++
			}
		}
		ans += cnt * (cnt - 1) / 2
	}
	// count pairs in columns
	for j := 0; j < n; j++ {
		cnt := 0
		for i := 0; i < n; i++ {
			if grid[i][j] == 'C' {
				cnt++
			}
		}
		ans += cnt * (cnt - 1) / 2
	}

	fmt.Fprintln(writer, ans)
}
