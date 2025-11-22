package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const maxLog = 20

type graph struct {
	n    int
	g    [][]int
	tin  []int
	tout []int
	up   [][]int
	dep  []int
	timer int
}

func newGraph(n int) *graph {
	return &graph{
		n:    n,
		g:    make([][]int, n),
		tin:  make([]int, n),
		tout: make([]int, n),
		up:   make([][]int, maxLog),
		dep:  make([]int, n),
	}
}

func (gr *graph) addEdge(u, v int) {
	gr.g[u] = append(gr.g[u], v)
	gr.g[v] = append(gr.g[v], u)
}

func (gr *graph) dfs(v, p int) {
	gr.tin[v] = gr.timer
	gr.timer++
	gr.up[0][v] = p
	for i := 1; i < maxLog; i++ {
		gr.up[i][v] = gr.up[i-1][gr.up[i-1][v]]
	}
	for _, to := range gr.g[v] {
		if to == p {
			continue
		}
		gr.dep[to] = gr.dep[v] + 1
		gr.dfs(to, v)
	}
	gr.tout[v] = gr.timer
}

func (gr *graph) prepare(root int) {
	for i := 0; i < maxLog; i++ {
		gr.up[i] = make([]int, gr.n)
	}
	gr.timer = 0
	// iterative DFS to avoid stack issues
	type nodeState struct {
		v   int
		p   int
		vis bool
	}
	st := []nodeState{{root, root, false}}
	for len(st) > 0 {
		cur := st[len(st)-1]
		st = st[:len(st)-1]
		if cur.vis {
			gr.tout[cur.v] = gr.timer
			continue
		}
		// entry
		gr.tin[cur.v] = gr.timer
		gr.timer++
		gr.up[0][cur.v] = cur.p
		for i := 1; i < maxLog; i++ {
			gr.up[i][cur.v] = gr.up[i-1][gr.up[i-1][cur.v]]
		}
		gr.dep[cur.v] = 0
		if cur.v != cur.p {
			gr.dep[cur.v] = gr.dep[cur.p] + 1
		}
		st = append(st, nodeState{cur.v, cur.p, true})
		for i := len(gr.g[cur.v]) - 1; i >= 0; i-- {
			to := gr.g[cur.v][i]
			if to == cur.p {
				continue
			}
			st = append(st, nodeState{to, cur.v, false})
		}
	}
}

func (gr *graph) isAnc(a, b int) bool {
	return gr.tin[a] <= gr.tin[b] && gr.tout[b] <= gr.tout[a]
}

func (gr *graph) lca(a, b int) int {
	if gr.isAnc(a, b) {
		return a
	}
	if gr.isAnc(b, a) {
		return b
	}
	for i := maxLog - 1; i >= 0; i-- {
		if !gr.isAnc(gr.up[i][a], b) {
			a = gr.up[i][a]
		}
	}
	return gr.up[0][a]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for ; T > 0; T-- {
		var n, k int
		fmt.Fscan(in, &n, &k)
		w := make([]int64, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &w[i])
		}
		c := make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &c[i])
		}

		gr := newGraph(n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			gr.addEdge(u, v)
		}
		gr.prepare(0)

		colorNodes := make(map[int][]int)
		for i, col := range c {
			if col != 0 {
				colorNodes[col] = append(colorNodes[col], i)
			}
		}

		countMulti := make([]int, n)
		ownMulti := make([]bool, n)
		repeatColor := make([]int, n) // 0 none, -1 multiple, >0 unique

		processColor := func(col int, nodes []int) {
			if len(nodes) < 2 {
				return
			}
			// unique and sort by tin
			sort.Slice(nodes, func(i, j int) bool { return gr.tin[nodes[i]] < gr.tin[nodes[j]] })
			uniq := make([]int, 0, len(nodes))
			prev := -1
			for _, v := range nodes {
				if prev == -1 || v != prev {
					uniq = append(uniq, v)
					prev = v
				}
			}
			nodes = uniq
			if len(nodes) < 2 {
				return
			}

			// virtual tree
			stack := make([]int, 0, len(nodes))
			stack = append(stack, nodes[0])
			deg := make(map[int]int)

			addEdge := func(p, ch int) {
				deg[p]++
			}

			for i := 1; i < len(nodes); i++ {
				u := nodes[i]
				l := gr.lca(u, stack[len(stack)-1])
				for len(stack) >= 2 && gr.dep[stack[len(stack)-2]] >= gr.dep[l] {
					addEdge(stack[len(stack)-2], stack[len(stack)-1])
					stack = stack[:len(stack)-1]
				}
				if stack[len(stack)-1] != l {
					addEdge(l, stack[len(stack)-1])
					stack = stack[:len(stack)-1]
					if len(stack) == 0 || stack[len(stack)-1] != l {
						stack = append(stack, l)
					}
				}
				stack = append(stack, u)
			}
			for len(stack) > 1 {
				addEdge(stack[len(stack)-2], stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}

			for v, d := range deg {
				if d >= 2 {
					countMulti[v]++
					if c[v] == col {
						ownMulti[v] = true
					}
					if repeatColor[v] == 0 {
						repeatColor[v] = col
					} else if repeatColor[v] != col {
						repeatColor[v] = -1
					}
				}
			}
		}

		for col, nodes := range colorNodes {
			processColor(col, nodes)
		}

		var cost int64
		for i := 0; i < n; i++ {
			mult := countMulti[i]
			if c[i] != 0 {
				if ownMulti[i] {
					mult--
				}
				if mult > 0 {
					cost += w[i]
				}
			} else {
				if mult >= 2 {
					cost += w[i]
				}
			}
		}

		// assign colors top-down
		if c[0] == 0 {
			if countMulti[0] == 1 && repeatColor[0] > 0 {
				c[0] = repeatColor[0]
			} else {
				c[0] = 1
			}
		}

		type stItem struct {
			v int
			p int
		}
		stack := []stItem{{0, 0}}
		for len(stack) > 0 {
			cur := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			for _, to := range gr.g[cur.v] {
				if to == cur.p {
					continue
				}
				if c[to] == 0 {
					if countMulti[to] == 1 && repeatColor[to] > 0 {
						c[to] = repeatColor[to]
					} else {
						c[to] = c[cur.v]
					}
				}
				stack = append(stack, stItem{to, cur.v})
			}
		}

		fmt.Fprintln(out, cost)
		for i := 0; i < n; i++ {
			if i > 0 {
				fmt.Fprint(out, " ")
			}
			fmt.Fprint(out, c[i])
		}
		fmt.Fprintln(out)
	}
}
