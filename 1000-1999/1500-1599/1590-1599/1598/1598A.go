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
	dirs := [8][2]int{{-1, -1}, {-1, 0}, {-1, 1}, {0, -1}, {0, 1}, {1, -1}, {1, 0}, {1, 1}}
	for ; t > 0; t-- {
		var n int
		fmt.Fscan(reader, &n)
		grid := make([]string, 2)
		fmt.Fscan(reader, &grid[0])
		fmt.Fscan(reader, &grid[1])

		vis := make([][]bool, 2)
		for i := 0; i < 2; i++ {
			vis[i] = make([]bool, n)
		}
		type node struct{ r, c int }
		q := []node{{0, 0}}
		vis[0][0] = true
		for len(q) > 0 {
			cur := q[0]
			q = q[1:]
			if cur.r == 1 && cur.c == n-1 {
				break
			}
			for _, d := range dirs {
				nr, nc := cur.r+d[0], cur.c+d[1]
				if nr < 0 || nr >= 2 || nc < 0 || nc >= n {
					continue
				}
				if grid[nr][nc] == '1' || vis[nr][nc] {
					continue
				}
				vis[nr][nc] = true
				q = append(q, node{nr, nc})
			}
		}
		if vis[1][n-1] {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
