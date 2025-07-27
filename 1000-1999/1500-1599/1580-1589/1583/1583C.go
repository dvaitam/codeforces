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

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}

	grid := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(reader, &grid[i])
	}

	bad := make([]int, m+1)
	for i := 1; i < n; i++ {
		for j := 1; j < m; j++ {
			if grid[i-1][j] == 'X' && grid[i][j-1] == 'X' {
				bad[j+1] = 1
			}
		}
	}

	pref := make([]int, m+1)
	for j := 1; j <= m; j++ {
		pref[j] = pref[j-1] + bad[j]
	}

	var q int
	fmt.Fscan(reader, &q)
	for ; q > 0; q-- {
		var x1, x2 int
		fmt.Fscan(reader, &x1, &x2)
		if pref[x2]-pref[x1] == 0 {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
