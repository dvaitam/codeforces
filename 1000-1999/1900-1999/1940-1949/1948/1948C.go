package main

import (
	"bufio"
	"fmt"
	"os"
)

type Node struct {
	r, c int
}

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
		if _, err := fmt.Fscan(reader, &n); err != nil {
			return
		}
		var s1, s2 string
		fmt.Fscan(reader, &s1)
		fmt.Fscan(reader, &s2)
		grid := [][]byte{[]byte(s1), []byte(s2)}

		visited := make([][]bool, 2)
		for i := 0; i < 2; i++ {
			visited[i] = make([]bool, n)
		}
		queue := []Node{{0, 0}}
		visited[0][0] = true
		reachable := false

		dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for len(queue) > 0 && !reachable {
			cur := queue[0]
			queue = queue[1:]

			if cur.r == 1 && cur.c == n-1 {
				reachable = true
				break
			}
			for _, d := range dirs {
				nr, nc := cur.r+d[0], cur.c+d[1]
				if nr < 0 || nr >= 2 || nc < 0 || nc >= n {
					continue
				}
				var tr, tc int
				if grid[nr][nc] == '>' {
					tr, tc = nr, nc+1
				} else {
					tr, tc = nr, nc-1
				}
				if tr < 0 || tr >= 2 || tc < 0 || tc >= n {
					continue
				}
				if !visited[tr][tc] {
					visited[tr][tc] = true
					queue = append(queue, Node{tr, tc})
				}
			}
		}

		if reachable {
			fmt.Fprintln(writer, "YES")
		} else {
			fmt.Fprintln(writer, "NO")
		}
	}
}
