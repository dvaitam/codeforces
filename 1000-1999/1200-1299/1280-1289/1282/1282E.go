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

	var Cases int
	fmt.Fscan(reader, &Cases)
	for Cases > 0 {
		Cases--
		var n int
		fmt.Fscan(reader, &n)
		// triangles t[1..n-2][0..2]
		t := make([][3]int, n)
		// adjacency from node to triangles
		g := make([][]int, n+1)
		d := make([]int, n+1)
		// bij[u] holds up to two connected nodes
		bij := make([][2]int, n+1)
		// visited triangles and nodes
		visTri := make([]bool, n)
		visNode := make([]bool, n+1)
		// union-find parent
		f := make([]int, n+1)
		for i := 1; i <= n; i++ {
			f[i] = i
		}
		// read triangles
		for i := 1; i <= n-2; i++ {
			fmt.Fscan(reader, &t[i][0], &t[i][1], &t[i][2])
			for j := 0; j < 3; j++ {
				u := t[i][j]
				g[u] = append(g[u], i)
				d[u]++
			}
		}
		// union-find find with path compression
		var find func(int) int
		find = func(u int) int {
			v := u
			for f[v] != v {
				v = f[v]
			}
			for f[u] != u {
				w := f[u]
				f[u] = v
				u = w
			}
			return v
		}
		// add neighbor link
		add := func(u, v int) {
			if bij[u][0] == 0 {
				bij[u][0] = v
			} else {
				bij[u][1] = v
			}
		}
		// connect u and v in bij and union-find
		connect := func(u, v int) {
			if find(u) == find(v) {
				return
			}
			add(u, v)
			add(v, u)
			f[find(v)] = u
		}
		// queue of nodes with degree 1
		nodeQ := make([]int, n+1)
		head, tail := 0, 0
		for i := 1; i <= n; i++ {
			if d[i] == 1 {
				nodeQ[tail] = i
				tail++
			}
		}
		q := make([]int, n)
		qnt := 0
		// process
		for head < tail {
			u := nodeQ[head]
			head++
			for _, w := range g[u] {
				if !visTri[w] {
					qnt++
					q[qnt] = w
					// examine triangle w
					for j := 0; j < 3; j++ {
						v := t[w][j]
						if v != u {
							connect(u, v)
							d[v]--
							if d[v] == 1 {
								nodeQ[tail] = v
								tail++
							}
						}
					}
					visTri[w] = true
				}
			}
		}
		// build node order p
		p := make([]int, n+1)
		// find start: node with only one neighbor in bij
		for i := 1; i <= n; i++ {
			if bij[i][1] == 0 {
				p[1] = i
				break
			}
		}
		visNode[p[1]] = true
		k := 1
		for i := 1; i <= n; i++ {
			for j := 0; j < 2; j++ {
				v := bij[p[i]][j]
				if v != 0 && !visNode[v] {
					k++
					p[k] = v
					visNode[v] = true
				}
			}
		}
		// output p[1..n]
		for i := 1; i <= n; i++ {
			if i > 1 {
				writer.WriteByte(' ')
			}
			fmt.Fprintf(writer, "%d", p[i])
		}
		writer.WriteByte('\n')
		// output q[1..n-2]
		for i := 1; i <= n-2; i++ {
			if i > 1 {
				writer.WriteByte(' ')
			}
			fmt.Fprintf(writer, "%d", q[i])
		}
		writer.WriteByte('\n')
	}
}
