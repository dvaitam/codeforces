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
		grid := make([][]int, n)
		for i := 0; i < n; i++ {
			grid[i] = make([]int, m)
			for j := 0; j < m; j++ {
				fmt.Fscan(reader, &grid[i][j])
			}
		}

		dist := make([][]int, n)
		for i := range dist {
			dist[i] = make([]int, m)
			for j := range dist[i] {
				dist[i][j] = -1
			}
		}

		type node struct{ r, c int }
		queue := make([]node, 0, n*m)
		dist[0][0] = 0
		queue = append(queue, node{0, 0})
		for head := 0; head < len(queue); head++ {
			cur := queue[head]
			d := dist[cur.r][cur.c]

			// move down
			r1 := (cur.r + 1) % n
			r2 := (cur.r + 2) % n
			if grid[r1][cur.c] == 0 && grid[r2][cur.c] == 0 {
				if dist[r2][cur.c] == -1 {
					dist[r2][cur.c] = d + 1
					queue = append(queue, node{r2, cur.c})
				}
			}

			// move right
			if cur.c+1 < m {
				rNext := (cur.r + 1) % n
				if grid[rNext][cur.c+1] == 0 && dist[rNext][cur.c+1] == -1 {
					dist[rNext][cur.c+1] = d + 1
					queue = append(queue, node{rNext, cur.c + 1})
				}
			}
		}

		ans := -1
		for r := 0; r < n; r++ {
			if dist[r][m-1] != -1 {
				time := dist[r][m-1]
				row := r - time%n
				if row < 0 {
					row += n
				}
				if row == n-1 {
					if ans == -1 || time < ans {
						ans = time
					}
				}
			}
		}

		fmt.Fprintln(writer, ans)
	}
}
