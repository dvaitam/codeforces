package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1 << 30

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	var n, m, k int
	if _, err := fmt.Fscan(reader, &n, &m, &k); err != nil {
		return
	}
	color := make([]int, n+1)
	for i := 1; i <= n; i++ {
		fmt.Fscan(reader, &color[i])
	}
	type Edge struct{ to, sign int }
	g := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		sign := 0
		if color[u] != color[v] {
			sign = 1
		}
		g[u] = append(g[u], Edge{v, sign})
		g[v] = append(g[v], Edge{u, sign})
	}

	// dist[s][t][p]
	dist := make([][][]int, n+1)
	for i := 0; i <= n; i++ {
		dist[i] = make([][]int, n+1)
		for j := 0; j <= n; j++ {
			dist[i][j] = []int{INF, INF}
		}
	}

	type State struct{ v, p int }
	queue := make([]State, 0)

	for s := 1; s <= n; s++ {
		for j := 1; j <= n; j++ {
			dist[s][j][0] = INF
			dist[s][j][1] = INF
		}
		head, tail := 0, 0
		queue = queue[:0]
		queue = append(queue, State{s, 0})
		tail++
		dist[s][s][0] = 0
		for head < tail {
			st := queue[head]
			head++
			d := dist[s][st.v][st.p]
			for _, e := range g[st.v] {
				np := st.p ^ e.sign
				if dist[s][e.to][np] == INF {
					dist[s][e.to][np] = d + 1
					queue = append(queue, State{e.to, np})
					tail++
				}
			}
		}
	}

	// compute diameter
	ans := 0
	for s := 1; s <= n; s++ {
		for t := 1; t <= n; t++ {
			dist0 := dist[s][t][0]
			dist1 := dist[s][t][1]
			for d := 0; d <= k; d++ {
				best := INF
				if dist0 < INF {
					val := dist0 + d
					if val < best {
						best = val
					}
				}
				if dist1 < INF {
					val := dist1 + (k - d)
					if val < best {
						best = val
					}
				}
				if best > ans {
					ans = best
				}
			}
		}
	}
	fmt.Fprintln(writer, ans)
}
