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

	var T int
	fmt.Fscan(reader, &T)
	for ; T > 0; T-- {
		var n, m int
		fmt.Fscan(reader, &n, &m)
		type Node struct {
			t, h int
			id   int
		}
		nodes := make([]Node, 0)
		idx := make(map[[2]int]int)
		for i := 1; i <= n; i++ {
			for j := 1; j <= m; j++ {
				var x int
				fmt.Fscan(reader, &x)
				if x != 0 {
					idx[[2]int{i, j}] = len(nodes)
					nodes = append(nodes, Node{i, j, x})
				}
			}
		}
		k := len(nodes)
		if k == 0 {
			fmt.Fprintln(writer, -1)
			continue
		}
		adj := make([][]int, k)
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		for id, node := range nodes {
			for _, d := range dirs {
				ni := node.t + d[0]
				nj := node.h + d[1]
				if v, ok := idx[[2]int{ni, nj}]; ok {
					adj[id] = append(adj[id], v)
				}
			}
		}
		// BFS to check connectivity
		vis := make([]bool, k)
		queue := []int{0}
		vis[0] = true
		for len(queue) > 0 {
			v := queue[0]
			queue = queue[1:]
			for _, to := range adj[v] {
				if !vis[to] {
					vis[to] = true
					queue = append(queue, to)
				}
			}
		}
		connected := true
		for _, b := range vis {
			if !b {
				connected = false
				break
			}
		}
		if !connected {
			fmt.Fprintln(writer, -1)
			continue
		}
		if k == 1 {
			fmt.Fprintln(writer, 1, 1)
			fmt.Fprintln(writer, nodes[0].id)
			continue
		}
		// DFS to build path
		vis = make([]bool, k)
		path := make([]int, 0, 2*k)
		var dfs func(int)
		dfs = func(v int) {
			vis[v] = true
			path = append(path, nodes[v].id)
			for _, to := range adj[v] {
				if !vis[to] {
					dfs(to)
					path = append(path, nodes[v].id)
				}
			}
		}
		dfs(0)
		// append one more step to make length 2k
		path = append(path, nodes[adj[0][0]].id)
		w := k
		h := 2
		index := 0
		fmt.Fprintln(writer, h, w)
		grid := make([][]int, h)
		for i := 0; i < h; i++ {
			grid[i] = make([]int, w)
			if i%2 == 0 {
				for j := 0; j < w; j++ {
					grid[i][j] = path[index]
					index++
				}
			} else {
				for j := w - 1; j >= 0; j-- {
					grid[i][j] = path[index]
					index++
				}
			}
		}
		for i := 0; i < h; i++ {
			for j := 0; j < w; j++ {
				if j > 0 {
					writer.WriteByte(' ')
				}
				fmt.Fprint(writer, grid[i][j])
			}
			writer.WriteByte('\n')
		}
	}
}
