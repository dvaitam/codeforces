package main

import (
	"bufio"
	"fmt"
	"os"
)

// Pos represents a coordinate in the grid
type Pos struct{ x, y int }

var dx = [4]int{1, -1, 0, 0}
var dy = [4]int{0, 0, 1, -1}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m int
	if _, err := fmt.Fscan(reader, &n, &m); err != nil {
		return
	}
	gs := make([][]byte, n)
	ge := make([][]byte, n)
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		gs[i] = []byte(s)
	}
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(reader, &s)
		ge[i] = []byte(s)
	}
	// BFS setup
	dist := make([][]int, n)
	pre := make([][]Pos, n)
	for i := 0; i < n; i++ {
		dist[i] = make([]int, m)
		pre[i] = make([]Pos, m)
		for j := 0; j < m; j++ {
			dist[i][j] = -1
		}
	}
	var q []Pos
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if gs[i][j] == '*' {
				dist[i][j] = 0
				pre[i][j] = Pos{i, j}
				q = append(q, Pos{i, j})
			}
		}
	}
	head := 0
	found := false
	for head < len(q) && !found {
		p := q[head]
		head++
		x, y := p.x, p.y
		if ge[x][y] == '*' {
			moveTo(writer, gs, ge, pre, n, m, p)
			found = true
			break
		}
		for d := 0; d < 4; d++ {
			a, b := x+dx[d], y+dy[d]
			if a >= 0 && a < n && b >= 0 && b < m && dist[a][b] == -1 && gs[a][b] != 'X' {
				dist[a][b] = dist[x][y] + 1
				pre[a][b] = p
				q = append(q, Pos{a, b})
			}
		}
	}
	if !found {
		fmt.Fprintln(writer, "NO")
	}
}

// moveTo reconstructs path and outputs moves
func moveTo(writer *bufio.Writer, gs, ge [][]byte, pre [][]Pos, n, m int, end Pos) {
	// Reconstruct BFS path
	ex, ey := end.x, end.y
	rex, rey := ex, ey
	var st []Pos
	for {
		par := pre[ex][ey]
		if par.x == ex && par.y == ey {
			break
		}
		st = append(st, Pos{ex, ey})
		ex, ey = par.x, par.y
	}
	// Collect toRemove via post-order DFS on gs
	toRemove := []Pos{}
	vis := make([][]bool, n)
	for i := range vis {
		vis[i] = make([]bool, m)
	}
	// post-order stacking
	stack := []Pos{{ex, ey}}
	var stack2 []Pos
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if vis[cur.x][cur.y] {
			continue
		}
		vis[cur.x][cur.y] = true
		stack2 = append(stack2, cur)
		for d := 0; d < 4; d++ {
			a, b := cur.x+dx[d], cur.y+dy[d]
			if a >= 0 && a < n && b >= 0 && b < m && !vis[a][b] && gs[a][b] == '*' {
				stack = append(stack, Pos{a, b})
			}
		}
	}
	for i := len(stack2) - 1; i >= 0; i-- {
		toRemove = append(toRemove, stack2[i])
	}
	// Initialize toAdd with st minus last point
	toAdd := []Pos{}
	for len(st) > 1 {
		last := st[len(st)-1]
		toAdd = append(toAdd, last)
		toRemove = append(toRemove, last)
		st = st[:len(st)-1]
	}
	// Pre-order DFS on target grid ge
	vis2 := make([][]bool, n)
	for i := range vis2 {
		vis2[i] = make([]bool, m)
	}
	stack = []Pos{{rex, rey}}
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if vis2[cur.x][cur.y] {
			continue
		}
		vis2[cur.x][cur.y] = true
		toAdd = append(toAdd, cur)
		for d := 0; d < 4; d++ {
			a, b := cur.x+dx[d], cur.y+dy[d]
			if a >= 0 && a < n && b >= 0 && b < m && !vis2[a][b] && ge[a][b] == '*' {
				stack = append(stack, Pos{a, b})
			}
		}
	}
	added := make([][]bool, n)
	for i := range added {
		added[i] = make([]bool, m)
	}
	var ans [][4]int
	// Process toAdd queue
	for len(toAdd) > 0 {
		p := toAdd[0]
		toAdd = toAdd[1:]
		x, y := p.x, p.y
		if gs[x][y] == '*' {
			added[x][y] = true
			continue
		}
		// find a valid toRemove
		var aPos Pos
		for len(toRemove) > 0 {
			aPos = toRemove[0]
			toRemove = toRemove[1:]
			if added[aPos.x][aPos.y] && ge[aPos.x][aPos.y] == '*' {
				continue
			}
			break
		}
		added[x][y] = true
		gs[aPos.x][aPos.y] = '.'
		gs[x][y] = '*'
		ans = append(ans, [4]int{aPos.x + 1, aPos.y + 1, x + 1, y + 1})
	}
	// Output result
	fmt.Fprintln(writer, "YES")
	fmt.Fprintln(writer, len(ans))
	for _, v := range ans {
		fmt.Fprintf(writer, "%d %d %d %d\n", v[0], v[1], v[2], v[3])
	}
}
