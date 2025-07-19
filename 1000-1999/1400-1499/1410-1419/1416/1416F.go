package main

import (
	"bufio"
	"fmt"
	"os"
)

var dx = [4]int{-1, 1, 0, 0}
var dy = [4]int{0, 0, -1, 1}
var dir = [4]byte{'U', 'D', 'L', 'R'}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	var t int
	fmt.Fscan(reader, &t)
	for tc := 0; tc < t; tc++ {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		size := n * m
		a := make([]int, size)
		good := make([]bool, size)
		g := make([][]int, size)
		p := make([]int, size)
		u := make([]int, size)
		for i := 0; i < size; i++ {
			p[i] = -1
		}
		for i := 0; i < size; i++ {
			fmt.Fscan(reader, &a[i])
		}
		// build graph and good
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				idx := i*m + j
				good[idx] = false
				// adjacency reset by new slice
				for k := 0; k < 4; k++ {
					x := i + dx[k]
					y := j + dy[k]
					if x >= 0 && x < n && y >= 0 && y < m {
						ni := x*m + y
						if a[idx] > a[ni] {
							good[idx] = true
						} else if a[idx] == a[ni] {
							g[idx] = append(g[idx], ni)
						}
					}
				}
			}
		}
		// list of nodes needing match (no smaller neighbor)
		li := make([]int, 0, size)
		for i := 0; i < size; i++ {
			if !good[i] {
				li = append(li, i)
			}
		}
		// dfs for augmenting
		T := 0
		var dfs func(int) bool
		dfs = func(v int) bool {
			if u[v] == T {
				return false
			}
			u[v] = T
			for _, to := range g[v] {
				if p[to] == -1 {
					p[to] = v
					p[v] = to
					return true
				}
			}
			for _, to := range g[v] {
				if good[p[to]] {
					p[p[to]] = -1
					p[to] = v
					p[v] = to
					return true
				}
			}
			for _, to := range g[v] {
				if dfs(p[to]) {
					p[to] = v
					p[v] = to
					return true
				}
			}
			return false
		}
		for {
			ok := false
			T++
			for _, v := range li {
				if p[v] == -1 && dfs(v) {
					ok = true
				}
			}
			if !ok {
				break
			}
		}
		// check feasibility
		feasible := true
		for i := 0; i < size; i++ {
			if !good[i] && p[i] == -1 {
				feasible = false
				break
			}
		}
		if !feasible {
			fmt.Fprintln(writer, "NO")
			continue
		}
		fmt.Fprintln(writer, "YES")
		// build b and c
		b := make([]int, size)
		c := make([]byte, size)
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				idx := i*m + j
				if p[idx] == -1 {
					// has smaller neighbor
					for k := 0; k < 4; k++ {
						x := i + dx[k]
						y := j + dy[k]
						if x >= 0 && x < n && y >= 0 && y < m {
							ni := x*m + y
							if a[idx] > a[ni] {
								b[idx] = a[idx] - a[ni]
								c[idx] = dir[k]
								break
							}
						}
					}
				} else {
					// matched equal neighbor
					for k := 0; k < 4; k++ {
						x := i + dx[k]
						y := j + dy[k]
						if x >= 0 && x < n && y >= 0 && y < m {
							ni := x*m + y
							if a[idx] == a[ni] && p[idx] == ni {
								if idx < ni {
									b[idx] = 1
								} else {
									b[idx] = a[idx] - 1
								}
								c[idx] = dir[k]
								break
							}
						}
					}
				}
			}
		}
		// output b
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				writer.WriteString(fmt.Sprintf("%d ", b[i*m+j]))
			}
			writer.WriteByte('\n')
		}
		// output c
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				writer.WriteByte(c[i*m+j])
				writer.WriteByte(' ')
			}
			writer.WriteByte('\n')
		}
	}
}
