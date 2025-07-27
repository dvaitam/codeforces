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
	const offset = 100000
	for ; t > 0; t-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		g := make([][]int, n)
		for i := 0; i < m; i++ {
			var u, v int
			fmt.Fscan(reader, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		feasible := true
		for i := 0; i < n && feasible; i++ {
			if len(g[i]) > 3 {
				feasible = false
			}
		}
		color := make([]int, n)
		for i := range color {
			color[i] = -1
		}
		queue := make([]int, 0)
		for i := 0; i < n && feasible; i++ {
			if color[i] != -1 {
				continue
			}
			color[i] = 0
			queue = append(queue, i)
			for len(queue) > 0 && feasible {
				u := queue[0]
				queue = queue[1:]
				for _, v := range g[u] {
					if color[v] == -1 {
						color[v] = color[u] ^ 1
						queue = append(queue, v)
					} else if color[v] == color[u] {
						feasible = false
						break
					}
				}
			}
		}
		if !feasible {
			fmt.Fprintln(writer, "No")
			continue
		}
		used0 := make(map[int]bool)
		used1 := make(map[int]bool)
		x := make([]int, n)
		y := make([]int, n)
		visited := make([]bool, n)
		queue = []int{0}
		visited[0] = true
		x[0] = color[0]
		y[0] = 0
		if color[0] == 0 {
			used0[0] = true
		} else {
			used1[0] = true
		}
		nextRow := []int{1, 1}
		for len(queue) > 0 {
			u := queue[0]
			queue = queue[1:]
			for _, v := range g[u] {
				if visited[v] {
					continue
				}
				visited[v] = true
				row := color[v]
				if row != color[u] {
					if row == 0 {
						if !used0[y[u]] {
							y[v] = y[u]
						} else {
							y[v] = nextRow[row]
							nextRow[row]++
						}
					} else {
						if !used1[y[u]] {
							y[v] = y[u]
						} else {
							y[v] = nextRow[row]
							nextRow[row]++
						}
					}
				} else {
					target := y[u] + 1
					if row == 0 {
						if !used0[target] {
							y[v] = target
						} else if !used0[y[u]-1] {
							y[v] = y[u] - 1
						} else {
							y[v] = nextRow[row]
							nextRow[row]++
						}
					} else {
						if !used1[target] {
							y[v] = target
						} else if !used1[y[u]-1] {
							y[v] = y[u] - 1
						} else {
							y[v] = nextRow[row]
							nextRow[row]++
						}
					}
				}
				x[v] = row
				if row == 0 {
					used0[y[v]] = true
				} else {
					used1[y[v]] = true
				}
				queue = append(queue, v)
			}
		}
		fmt.Fprintln(writer, "Yes")
		for i := 0; i < n; i++ {
			fmt.Fprintln(writer, x[i], y[i]-offset)
		}
	}
}
