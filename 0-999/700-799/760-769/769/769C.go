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

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	grid := make([][]byte, n)
	var sx, sy int
	for i := 0; i < n; i++ {
		var row string
		fmt.Fscan(reader, &row)
		grid[i] = []byte(row)
		for j := 0; j < m; j++ {
			if grid[i][j] == 'X' {
				sx, sy = i, j
				grid[i][j] = '.'
			}
		}
	}

	if k%2 == 1 {
		fmt.Fprintln(writer, "IMPOSSIBLE")
		return
	}

	type pair struct{ x, y int }
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, m)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	q := make([]pair, 0, n*m)
	q = append(q, pair{sx, sy})
	dist[sx][sy] = 0
	for head := 0; head < len(q); head++ {
		p := q[head]
		d := dist[p.x][p.y] + 1
		if p.x+1 < n && grid[p.x+1][p.y] != '*' && dist[p.x+1][p.y] == -1 {
			dist[p.x+1][p.y] = d
			q = append(q, pair{p.x + 1, p.y})
		}
		if p.y-1 >= 0 && grid[p.x][p.y-1] != '*' && dist[p.x][p.y-1] == -1 {
			dist[p.x][p.y-1] = d
			q = append(q, pair{p.x, p.y - 1})
		}
		if p.y+1 < m && grid[p.x][p.y+1] != '*' && dist[p.x][p.y+1] == -1 {
			dist[p.x][p.y+1] = d
			q = append(q, pair{p.x, p.y + 1})
		}
		if p.x-1 >= 0 && grid[p.x-1][p.y] != '*' && dist[p.x-1][p.y] == -1 {
			dist[p.x-1][p.y] = d
			q = append(q, pair{p.x - 1, p.y})
		}
	}

	dirs := []byte{'D', 'L', 'R', 'U'}
	dx := []int{1, 0, 0, -1}
	dy := []int{0, -1, 1, 0}
	ans := make([]byte, 0, k)
	x, y := sx, sy
	for step := 0; step < k; step++ {
		rem := k - step - 1
		moved := false
		for i := 0; i < 4; i++ {
			nx, ny := x+dx[i], y+dy[i]
			if nx < 0 || nx >= n || ny < 0 || ny >= m || grid[nx][ny] == '*' {
				continue
			}
			d := dist[nx][ny]
			if d != -1 && d <= rem && (rem-d)%2 == 0 {
				ans = append(ans, dirs[i])
				x, y = nx, ny
				moved = true
				break
			}
		}
		if !moved {
			fmt.Fprintln(writer, "IMPOSSIBLE")
			return
		}
	}
	fmt.Fprintln(writer, string(ans))
}
