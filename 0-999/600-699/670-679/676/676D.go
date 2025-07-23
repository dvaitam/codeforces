package main

import (
	"bufio"
	"fmt"
	"os"
)

var masks = map[byte]int{
	'+': 15,
	'-': 10,
	'|': 5,
	'^': 1,
	'>': 2,
	'<': 8,
	'v': 4,
	'L': 7,
	'R': 13,
	'U': 14,
	'D': 11,
	'*': 0,
}

func rotate(mask, r int) int {
	r &= 3
	return ((mask << r) | (mask >> (4 - r))) & 15
}

type state struct {
	x, y, rot int
}

func main() {
	in := bufio.NewReader(os.Stdin)
	var n, m int
	if _, err := fmt.Fscan(in, &n, &m); err != nil {
		return
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		grid[i] = []byte(s)
	}
	var sx, sy, tx, ty int
	fmt.Fscan(in, &sx, &sy)
	fmt.Fscan(in, &tx, &ty)
	sx--
	sy--
	tx--
	ty--

	dirs := [4][2]int{{-1, 0}, {0, 1}, {1, 0}, {0, -1}}
	dist := make([][][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([][]int, m)
		for j := 0; j < m; j++ {
			dist[i][j] = []int{-1, -1, -1, -1}
		}
	}

	q := make([]state, 0)
	dist[sx][sy][0] = 0
	q = append(q, state{sx, sy, 0})

	for head := 0; head < len(q); head++ {
		cur := q[head]
		d := dist[cur.x][cur.y][cur.rot]
		// rotate
		nr := (cur.rot + 1) & 3
		if dist[cur.x][cur.y][nr] == -1 {
			dist[cur.x][cur.y][nr] = d + 1
			q = append(q, state{cur.x, cur.y, nr})
		}
		// move
		mask := rotate(masks[grid[cur.x][cur.y]], cur.rot)
		for dir := 0; dir < 4; dir++ {
			if mask&(1<<dir) == 0 {
				continue
			}
			nx := cur.x + dirs[dir][0]
			ny := cur.y + dirs[dir][1]
			if nx < 0 || nx >= n || ny < 0 || ny >= m {
				continue
			}
			if masks[grid[nx][ny]] == 0 {
				continue
			}
			nmask := rotate(masks[grid[nx][ny]], cur.rot)
			if nmask&(1<<((dir+2)&3)) == 0 {
				continue
			}
			if dist[nx][ny][cur.rot] == -1 {
				dist[nx][ny][cur.rot] = d + 1
				q = append(q, state{nx, ny, cur.rot})
			}
		}
	}

	ans := -1
	for r := 0; r < 4; r++ {
		v := dist[tx][ty][r]
		if v != -1 {
			if ans == -1 || v < ans {
				ans = v
			}
		}
	}
	fmt.Println(ans)
}
