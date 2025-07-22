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
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var line string
		fmt.Fscan(reader, &line)
		grid[i] = []byte(line)
	}

	dirs := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			ch := grid[i][j]
			if ch == '*' {
				continue
			}
			count := 0
			for _, d := range dirs {
				ni, nj := i+d[0], j+d[1]
				if ni >= 0 && ni < n && nj >= 0 && nj < m && grid[ni][nj] == '*' {
					count++
				}
			}
			if ch == '.' {
				if count != 0 {
					fmt.Fprintln(writer, "NO")
					return
				}
			} else { // digit
				if int(ch-'0') != count {
					fmt.Fprintln(writer, "NO")
					return
				}
			}
		}
	}
	fmt.Fprintln(writer, "YES")
}
