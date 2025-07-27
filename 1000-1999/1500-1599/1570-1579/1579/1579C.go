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
		var n, m, k int
		fmt.Fscan(reader, &n, &m, &k)
		grid := make([][]byte, n)
		for i := 0; i < n; i++ {
			var s string
			fmt.Fscan(reader, &s)
			grid[i] = []byte(s)
		}
		covered := make([][]bool, n)
		for i := 0; i < n; i++ {
			covered[i] = make([]bool, m)
		}
		for i := n - 1; i >= 0; i-- {
			for j := 0; j < m; j++ {
				if grid[i][j] != '*' {
					continue
				}
				size := 0
				for {
					x := i - size - 1
					y1 := j - size - 1
					y2 := j + size + 1
					if x < 0 || y1 < 0 || y2 >= m {
						break
					}
					if grid[x][y1] == '*' && grid[x][y2] == '*' {
						size++
					} else {
						break
					}
				}
				if size >= k {
					covered[i][j] = true
					for d := 1; d <= size; d++ {
						covered[i-d][j-d] = true
						covered[i-d][j+d] = true
					}
				}
			}
		}
		ok := true
		for i := 0; i < n && ok; i++ {
			for j := 0; j < m; j++ {
				if grid[i][j] == '*' && !covered[i][j] {
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
