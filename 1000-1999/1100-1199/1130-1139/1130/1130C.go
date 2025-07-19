package main

import (
	"bufio"
	"fmt"
	"os"
)

// pair holds a coordinate (x, y)
type pair struct{ x, y int }

var (
	n              int
	sx, sy, ex, ey int
	grid           [][]byte
	used           [][]bool
	vec1, vec2     []pair
	found          bool
	dx             = [4]int{1, -1, 0, 0}
	dy             = [4]int{0, 0, 1, -1}
)

func dfs1(x, y int) {
	if found {
		return
	}
	used[x][y] = true
	if x == ex && y == ey {
		found = true
		return
	}
	for i := 0; i < 4; i++ {
		tx := x + dx[i]
		ty := y + dy[i]
		if tx < 0 || tx >= n || ty < 0 || ty >= n || used[tx][ty] {
			continue
		}
		if grid[tx][ty] == '1' {
			vec1 = append(vec1, pair{x, y})
		} else {
			dfs1(tx, ty)
			if found {
				return
			}
		}
	}
}

func dfs2(x, y int) {
	used[x][y] = true
	for i := 0; i < 4; i++ {
		tx := x + dx[i]
		ty := y + dy[i]
		if tx < 0 || tx >= n || ty < 0 || ty >= n || used[tx][ty] {
			continue
		}
		if grid[tx][ty] == '1' {
			vec2 = append(vec2, pair{x, y})
		} else {
			dfs2(tx, ty)
		}
	}
}

func sqr(x int) int { return x * x }

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Fscan(reader, &n)
	fmt.Fscan(reader, &sx, &sy)
	fmt.Fscan(reader, &ex, &ey)
	sx--
	sy--
	ex--
	ey--
	grid = make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		grid[i] = []byte(s)
	}
	used = make([][]bool, n)
	for i := range used {
		used[i] = make([]bool, n)
	}
	dfs1(sx, sy)
	if found {
		fmt.Println(0)
		return
	}
	// reset used for second DFS
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			used[i][j] = false
		}
	}
	dfs2(ex, ey)
	// compute minimal squared distance between frontier points
	ans := 1 << 62
	for _, u := range vec1 {
		for _, v := range vec2 {
			dx := u.x - v.x
			dy := u.y - v.y
			d := sqr(dx) + sqr(dy)
			if d < ans {
				ans = d
			}
		}
	}
	fmt.Println(ans)
}
