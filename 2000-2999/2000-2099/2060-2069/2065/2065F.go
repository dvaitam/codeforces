package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
)

type edge struct {
	to   int
	cost int
}

const negInf = -1 << 60

var (
	n       int
	g       [][]int
	a       []int
	tin     []int
	depth   []int
	up      [][]int
	timer   int
	tmpAdj  [][]edge
	mark    []int
	curMark int
	LOG     int
)

func dfsInit(root int) {
	type frame struct {
		node, parent, idx int
	}
	stack := []frame{{node: root, parent: root, idx: 0}}
	for len(stack) > 0 {
		f := &stack[len(stack)-1]
		u := f.node
		if f.idx == 0 {
			tin[u] = timer
			timer++
			if u == f.parent {
				depth[u] = 0
				up[0][u] = u
			} else {
				depth[u] = depth[f.parent] + 1
				up[0][u] = f.parent
			}
			for k := 1; k < LOG; k++ {
				up[k][u] = up[k-1][up[k-1][u]]
			}
		}
		if f.idx < len(g[u]) {
			v := g[u][f.idx]
			f.idx++
			if v == f.parent {
				continue
			}
			stack = append(stack, frame{node: v, parent: u, idx: 0})
		} else {
			stack = stack[:len(stack)-1]
		}
	}
}

func lca(u, v int) int {
	if depth[u] < depth[v] {
		u, v = v, u
	}
	diff := depth[u] - depth[v]
	bit := 0
	for diff > 0 {
		if diff&1 == 1 {
			u = up[bit][u]
		}
		diff >>= 1
		bit++
	}
	if u == v {
		return u
	}
	for k := LOG - 1; k >= 0; k-- {
		if up[k][u] != up[k][v] {
			u = up[k][u]
			v = up[k][v]
		}
	}
	return up[0][u]
}

func distance(u, v int) int {
	w := lca(u, v)
	return depth[u] + depth[v] - 2*depth[w]
}

func addEdge(u, v int) {
	if u == v {
		return
	}
	dist := distance(u, v)
	cost := dist - 1
	tmpAdj[u] = append(tmpAdj[u], edge{to: v, cost: cost})
	tmpAdj[v] = append(tmpAdj[v], edge{to: u, cost: cost})
}

var curColor int
var possible bool

func dfsVT(u, parent int) int {
	weight := -1
	if a[u] == curColor {
		weight = 1
	}
	bestDown := weight
	best1, best2 := negInf, negInf
	for _, e := range tmpAdj[u] {
		if e.to == parent {
			continue
		}
		childDown := dfsVT(e.to, u)
		contribution := childDown - e.cost
		if weight+contribution > bestDown {
			bestDown = weight + contribution
		}
		if weight+contribution > 0 {
			possible = true
		}
		if contribution > best1 {
			best2 = best1
			best1 = contribution
		} else if contribution > best2 {
			best2 = contribution
		}
	}
	if best2 > negInf/2 { // ensure both exist
		if weight+best1+best2 > 0 {
			possible = true
		}
	}
	return bestDown
}

func solveColor(nodes []int) bool {
	if len(nodes) < 2 {
		return false
	}
	sort.Slice(nodes, func(i, j int) bool { return tin[nodes[i]] < tin[nodes[j]] })
	curMark++
	vtNodes := make([]int, 0, len(nodes)*2)
	addNode := func(x int) {
		if mark[x] != curMark {
			mark[x] = curMark
			vtNodes = append(vtNodes, x)
		}
	}
	stack := make([]int, 0, len(nodes)*2)
	for _, v := range nodes {
		addNode(v)
		if len(stack) == 0 {
			stack = append(stack, v)
			continue
		}
		l := lca(v, stack[len(stack)-1])
		addNode(l)
		for len(stack) >= 2 && depth[stack[len(stack)-2]] >= depth[l] {
			addEdge(stack[len(stack)-2], stack[len(stack)-1])
			stack = stack[:len(stack)-1]
		}
		if stack[len(stack)-1] != l {
			addEdge(stack[len(stack)-1], l)
			stack[len(stack)-1] = l
		}
		stack = append(stack, v)
	}
	for len(stack) > 1 {
		addEdge(stack[len(stack)-2], stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}
	root := nodes[0]
	possible = false
	dfsVT(root, -1)
	for _, v := range vtNodes {
		tmpAdj[v] = tmpAdj[v][:0]
	}
	return possible
}

func main() {
	debug.SetMaxStack(1 << 30)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var t int
	fmt.Fscan(in, &t)
	for ; t > 0; t-- {
		fmt.Fscan(in, &n)
		a = make([]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &a[i])
		}
		g = make([][]int, n)
		for i := 0; i < n-1; i++ {
			var u, v int
			fmt.Fscan(in, &u, &v)
			u--
			v--
			g[u] = append(g[u], v)
			g[v] = append(g[v], u)
		}
		LOG = 0
		for (1 << LOG) <= n {
			LOG++
		}
		tin = make([]int, n)
		depth = make([]int, n)
		up = make([][]int, LOG)
		for i := 0; i < LOG; i++ {
			up[i] = make([]int, n)
		}
		timer = 0
		dfsInit(0)
		tmpAdj = make([][]edge, n)
		mark = make([]int, n)
		curMark = 0

		pos := make([][]int, n+1)
		for i := 0; i < n; i++ {
			val := a[i]
			pos[val] = append(pos[val], i)
		}
		res := make([]byte, n)
		for val := 1; val <= n; val++ {
			nodes := pos[val]
			if len(nodes) < 2 {
				res[val-1] = '0'
				continue
			}
			curColor = val
			if solveColor(nodes) {
				res[val-1] = '1'
			} else {
				res[val-1] = '0'
			}
		}
		fmt.Fprintln(out, string(res))
	}
}
