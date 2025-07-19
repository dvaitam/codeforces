package main

import (
	"bufio"
	"fmt"
	"os"
)

var n, m int
var head, cur, nxt []int
var ne int
var sg [][]bool
var dfn, low, idx []int
var isBridge [][]bool
var g [][]int
var noNodes [][]int
var reader *bufio.Reader
var writer *bufio.Writer

func initGraph(N int) {
	head = make([]int, N)
	for i := range head {
		head[i] = -1
	}
	cur = make([]int, 2*m+5)
	nxt = make([]int, 2*m+5)
	ne = 0
}

func addEdge(u, v int) {
	cur[ne] = v
	nxt[ne] = head[u]
	head[u] = ne
	ne++
}

func tarjan(u, fa, d int) {
	dfn[u] = d
	low[u] = d
	for i := head[u]; i != -1; i = nxt[i] {
		v := cur[i]
		if dfn[v] == -1 {
			tarjan(v, u, d+1)
			if low[v] < low[u] {
				low[u] = low[v]
			}
			if dfn[u] < low[v] {
				isBridge[u][v] = true
				isBridge[v][u] = true
			}
		} else if v != fa {
			if dfn[v] < low[u] {
				low[u] = dfn[v]
			}
		}
	}
}

func dfsComp(u, c int) {
	if idx[u] != -1 {
		return
	}
	idx[u] = c
	for i := head[u]; i != -1; i = nxt[i] {
		v := cur[i]
		if isBridge[u][v] {
			continue
		}
		dfsComp(v, c)
	}
}

func output(a, b int) {
	for _, ua := range noNodes[a] {
		for _, ub := range noNodes[b] {
			if !sg[ua][ub] {
				fmt.Fprintln(writer, ua, ub)
				return
			}
		}
	}
}

func gao(u, fa int) (int, int) {
	a := []int{}
	b := []int{}
	for _, v := range g[u] {
		if v == fa {
			continue
		}
		f0, f1 := gao(v, u)
		b = append(b, f0)
		if f1 != -1 {
			a = append(a, f1)
			last := len(b) - 1
			b[0], b[last] = b[last], b[0]
		}
	}
	if len(b) == 0 {
		return u, -1
	}
	for i := len(b) - 1; i >= 0; i-- {
		a = append(a, b[i])
	}
	start := 2 - (len(a) % 2)
	for i := start; i < len(a); i += 2 {
		output(a[i], a[i+1])
	}
	if fa == -1 {
		output(a[0], a[1])
	}
	if len(a)%2 == 1 {
		return a[0], -1
	}
	return a[0], a[1]
}

func main() {
	reader = bufio.NewReader(os.Stdin)
	writer = bufio.NewWriter(os.Stdout)
	defer writer.Flush()
	fmt.Fscan(reader, &n, &m)
	sg = make([][]bool, n+1)
	isBridge = make([][]bool, n+1)
	for i := 0; i <= n; i++ {
		sg[i] = make([]bool, n+1)
		isBridge[i] = make([]bool, n+1)
	}
	initGraph(n + 1)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(reader, &u, &v)
		addEdge(u, v)
		addEdge(v, u)
		sg[u][v] = true
		sg[v][u] = true
	}
	if n == 2 {
		fmt.Fprintln(writer, -1)
		return
	}
	dfn = make([]int, n+1)
	low = make([]int, n+1)
	idx = make([]int, n+1)
	for i := 1; i <= n; i++ {
		dfn[i] = -1
		idx[i] = -1
	}
	tarjan(1, 0, 0)
	cid := 0
	for i := 1; i <= n; i++ {
		if idx[i] == -1 {
			dfsComp(i, cid)
			cid++
		}
	}
	g = make([][]int, cid)
	noNodes = make([][]int, cid)
	for u := 1; u <= n; u++ {
		for i := head[u]; i != -1; i = nxt[i] {
			v := cur[i]
			if idx[u] != idx[v] {
				g[idx[u]] = append(g[idx[u]], idx[v])
			}
		}
	}
	for i := 1; i <= n; i++ {
		noNodes[idx[i]] = append(noNodes[idx[i]], i)
	}
	if cid == 1 {
		fmt.Fprintln(writer, 0)
	} else if cid == 2 {
		fmt.Fprintln(writer, 1)
		for _, ua := range noNodes[0] {
			for _, ub := range noNodes[1] {
				if !sg[ua][ub] {
					fmt.Fprintln(writer, ua, ub)
					return
				}
			}
		}
	} else {
		leaf := 0
		st := -1
		for i := 0; i < cid; i++ {
			if len(g[i]) == 1 {
				leaf++
			} else if len(g[i]) > 1 && st == -1 {
				st = i
			}
		}
		fmt.Fprintln(writer, (leaf+1)/2)
		gao(st, -1)
	}
}
