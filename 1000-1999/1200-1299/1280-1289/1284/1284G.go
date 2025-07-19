package main

import (
	"bufio"
	"fmt"
	"os"
)

// Edge represents a directed connection with a direction value
type Edge struct {
	to, val int
}

var dx = [4]int{1, -1, 0, 0}
var dy = [4]int{0, 0, 1, -1}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var T int
	fmt.Fscan(reader, &T)
	for t := 0; t < T; t++ {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		c := make([]string, n+1)
		for i := 1; i <= n; i++ {
			fmt.Fscan(reader, &c[i])
		}

		total := n * m
		graph := make([][]Edge, total+1)
		parent := make([]int, total+1)
		for i := 1; i <= total; i++ {
			parent[i] = i
		}
		num := make([]int, total+1)
		dir := make([]int, total+1)
		flag := make([]bool, total+1)
		// build graph: cells with (x+y) odd
		for x := 1; x <= n; x++ {
			startY := (x & 1) + 1
			for y := startY; y <= m; y += 2 {
				if c[x][y] != 'O' {
					continue
				}
				u := (x-1)*m + y
				for k := 0; k < 4; k++ {
					nx, ny := x+dx[k], y+dy[k]
					if nx < 1 || ny < 1 || nx > n || ny > m {
						continue
					}
					if nx == 1 && ny == 1 {
						continue
					}
					if c[nx][ny] != 'O' {
						continue
					}
					v := (nx-1)*m + ny
					graph[u] = append(graph[u], Edge{v, k})
				}
			}
		}
		var find func(int) int
		find = func(x int) int {
			if parent[x] != x {
				parent[x] = find(parent[x])
			}
			return parent[x]
		}
		// matching dfs
		vis := make([]bool, total+1)
		var dfs func(int) bool
		dfs = func(u int) bool {
			if vis[u] {
				return false
			}
			vis[u] = true
			for _, e := range graph[u] {
				v := e.to
				if !vis[v] && (num[v] == 0 || dfs(num[v])) {
					num[v] = u
					dir[v] = e.val ^ 1
					return true
				}
			}
			return false
		}
		// find maximum matching on odd cells
		for x := 1; x <= n; x++ {
			startY := (x & 1) + 1
			for y := startY; y <= m; y += 2 {
				if c[x][y] != 'O' {
					continue
				}
				for i := 1; i <= total; i++ {
					vis[i] = false
				}
				u := (x-1)*m + y
				flag[u] = dfs(u)
			}
		}
		// prepare answer grid
		H, W := 2*n-1, 2*m-1
		ans := make([][]byte, H)
		for i := 0; i < H; i++ {
			ans[i] = make([]byte, W)
			for j := 0; j < W; j++ {
				ans[i][j] = ' '
			}
		}
		// union-find helper
		union := func(a, b int) {
			ra, rb := find(a), find(b)
			if ra != rb {
				parent[ra] = rb
			}
		}
		// match even cells (x+y even)
		failed := false
		for x := 1; x <= n && !failed; x++ {
			startY := 2 - (x & 1)
			for y := startY; y <= m; y += 2 {
				if x+y <= 2 {
					continue
				}
				if c[x][y] != 'O' {
					continue
				}
				v := (x-1)*m + y
				if num[v] == 0 {
					fmt.Fprintln(writer, "NO")
					failed = true
					break
				}
				u := num[v]
				union(v, u)
				ci := 2*x - 2 + dx[dir[v]]
				cj := 2*y - 2 + dy[dir[v]]
				ans[ci][cj] = 'O'
			}
		}
		if failed {
			continue
		}
		// add special edges
		if n >= 2 && c[2][1] == 'O' {
			graph[m+1] = append(graph[m+1], Edge{1, 1})
		}
		if m >= 2 && c[1][2] == 'O' {
			graph[2] = append(graph[2], Edge{1, 3})
		}
		// dfs2 to connect remaining odd cells
		var dfs2 func(int, int)
		dfs2 = func(u, f int) {
			for _, e := range graph[u] {
				v := e.to
				if v == f {
					continue
				}
				if find(u) != find(v) {
					union(u, v)
					ci := (u-1)/m*2 + dx[e.val]
					cj := (u-1)%m*2 + dy[e.val]
					ans[ci][cj] = 'O'
					if v != 1 {
						dfs2(num[v], v)
					}
				}
			}
		}
		for x := 1; x <= n; x++ {
			startY := (x & 1) + 1
			for y := startY; y <= m; y += 2 {
				if c[x][y] != 'O' {
					continue
				}
				u := (x-1)*m + y
				if !flag[u] {
					dfs2(u, -1)
				}
			}
		}
		// connectivity check
		root1 := find(1)
		for x := 1; x <= n && !failed; x++ {
			for y := 1; y <= m; y++ {
				if c[x][y] == 'O' {
					if find((x-1)*m+y) != root1 {
						fmt.Fprintln(writer, "NO")
						failed = true
						break
					}
				}
			}
		}
		if failed {
			continue
		}
		// set cell centers
		for x := 1; x <= n; x++ {
			for y := 1; y <= m; y++ {
				if c[x][y] == 'O' {
					ans[2*x-2][2*y-2] = 'O'
				}
			}
		}
		// output result
		fmt.Fprintln(writer, "YES")
		for i := 0; i < H; i++ {
			writer.Write(ans[i])
			writer.WriteByte('\n')
		}
	}
}
