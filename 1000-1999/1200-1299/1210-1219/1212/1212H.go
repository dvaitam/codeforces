package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1 << 30

type Edge struct {
	to  int
	idx int
}

var (
	n, k      int
	g         [][]Edge
	parent    []int
	pEdge     []int
	dp        [][]int
	trace     [][][]int
	size      []int
	children  [][]int
	childEdge [][]int
)

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func dfs(v, p, idx int) {
	parent[v] = p
	pEdge[v] = idx
	size[v] = 1
	var ch []int
	var ce []int
	for _, e := range g[v] {
		if e.to == p {
			continue
		}
		dfs(e.to, v, e.idx)
		size[v] += size[e.to]
		ch = append(ch, e.to)
		ce = append(ce, e.idx)
	}
	children[v] = ch
	childEdge[v] = ce

	dp[v] = make([]int, k+1)
	for i := 0; i <= k; i++ {
		dp[v][i] = INF
	}
	trace[v] = make([][]int, k+1)

	dp[v][1] = 0
	trace[v][1] = make([]int, 0)
	curSize := 1
	curDP := make([]int, k+1)
	for i := 0; i <= k; i++ {
		curDP[i] = INF
	}
	curDP[1] = 0
	curChoice := make([][]int, k+1)
	curChoice[1] = []int{}

	for _, u := range ch {
		nxtDP := make([]int, k+1)
		for i := 0; i <= k; i++ {
			nxtDP[i] = INF
		}
		nxtChoice := make([][]int, k+1)
		for s1 := 1; s1 <= min(k, curSize); s1++ {
			if curDP[s1] == INF {
				continue
			}
			base := curChoice[s1]
			// exclude child
			if curDP[s1]+1 < nxtDP[s1] {
				nxtDP[s1] = curDP[s1] + 1
				tmp := append([]int(nil), base...)
				tmp = append(tmp, 0)
				nxtChoice[s1] = tmp
			}
			// include child with w nodes
			for w := 1; w <= min(size[u], k-s1); w++ {
				if dp[u][w] == INF {
					continue
				}
				c := curDP[s1] + dp[u][w]
				if c < nxtDP[s1+w] {
					nxtDP[s1+w] = c
					tmp := append([]int(nil), base...)
					tmp = append(tmp, w)
					nxtChoice[s1+w] = tmp
				}
			}
		}
		curSize += size[u]
		curDP = nxtDP
		curChoice = nxtChoice
	}

	for s := 1; s <= min(k, curSize); s++ {
		dp[v][s] = curDP[s]
		if curDP[s] != INF {
			trace[v][s] = curChoice[s]
		}
	}
}

var result []int
var rootNode int

func collect(v, need int) {
	ch := children[v]
	ce := childEdge[v]
	choices := trace[v][need]
	idx := 0
	for i, u := range ch {
		w := 0
		if idx < len(choices) {
			w = choices[idx]
		}
		idx++
		if w == 0 {
			result = append(result, ce[i])
		} else {
			collect(u, w)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
	defer writer.Flush()

	if _, err := fmt.Fscan(reader, &n, &k); err != nil {
		return
	}
	g = make([][]Edge, n+1)
	parent = make([]int, n+1)
	pEdge = make([]int, n+1)
	dp = make([][]int, n+1)
	trace = make([][][]int, n+1)
	size = make([]int, n+1)
	children = make([][]int, n+1)
	childEdge = make([][]int, n+1)

	for i := 1; i <= n-1; i++ {
		var x, y int
		fmt.Fscan(reader, &x, &y)
		g[x] = append(g[x], Edge{y, i})
		g[y] = append(g[y], Edge{x, i})
	}

	dfs(1, 0, 0)

	bestCost := INF
	rootNode = 1
	for v := 1; v <= n; v++ {
		if k >= len(dp[v]) {
			continue
		}
		if dp[v][k] == INF {
			continue
		}
		cost := dp[v][k]
		if v != 1 {
			cost++
		}
		if cost < bestCost {
			bestCost = cost
			rootNode = v
		}
	}

	if rootNode != 1 {
		result = append(result, pEdge[rootNode])
	}
	collect(rootNode, k)

	fmt.Fprintln(writer, len(result))
	if len(result) > 0 {
		for i, id := range result {
			if i > 0 {
				writer.WriteByte(' ')
			}
			fmt.Fprint(writer, id)
		}
		writer.WriteByte('\n')
	}
}
