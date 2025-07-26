package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type edge struct {
	u, v int
	w    int64
}
type item struct {
	node int
	dist int64
}
type priorityQueue []item

func (pq priorityQueue) Len() int            { return len(pq) }
func (pq priorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq priorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *priorityQueue) Push(x interface{}) { *pq = append(*pq, x.(item)) }
func (pq *priorityQueue) Pop() interface{} {
	old := *pq
	v := old[len(old)-1]
	*pq = old[:len(old)-1]
	return v
}

func expectedF(n, m, k, q int, edges []edge, queries [][2]int) []int64 {
	g := make([][]edge, n+1)
	for _, e := range edges {
		g[e.u] = append(g[e.u], edge{e.v, e.u, e.w})
		g[e.v] = append(g[e.v], edge{e.u, e.v, e.w})
	}
	const INF int64 = 1 << 60
	dist := make([]int64, n+1)
	belong := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	pqq := &priorityQueue{}
	heap.Init(pqq)
	for i := 1; i <= k; i++ {
		dist[i] = 0
		belong[i] = i
		heap.Push(pqq, item{i, 0})
	}
	for pqq.Len() > 0 {
		it := heap.Pop(pqq).(item)
		if it.dist != dist[it.node] {
			continue
		}
		u := it.node
		for _, e := range g[u] {
			nd := dist[u] + e.w
			if nd < dist[e.u] {
				dist[e.u] = nd
				belong[e.u] = belong[u]
				heap.Push(pqq, item{e.u, nd})
			}
		}
	}
	type e2 struct {
		u, v int
		w    int64
	}
	edges2 := make([]e2, 0)
	for _, e := range edges {
		bu := belong[e.u]
		bv := belong[e.v]
		if bu != bv {
			cost := dist[e.u] + dist[e.v] + e.w
			edges2 = append(edges2, e2{bu, bv, cost})
		}
	}
	sort.Slice(edges2, func(i, j int) bool { return edges2[i].w < edges2[j].w })
	parent := make([]int, k+1)
	for i := 1; i <= k; i++ {
		parent[i] = i
	}
	var find func(int) int
	find = func(x int) int {
		if parent[x] != x {
			parent[x] = find(parent[x])
		}
		return parent[x]
	}
	unite := func(a, b int) bool {
		fa := find(a)
		fb := find(b)
		if fa == fb {
			return false
		}
		parent[fb] = fa
		return true
	}
	adj := make([][]edge, k+1)
	for _, e := range edges2 {
		if unite(e.u, e.v) {
			adj[e.u] = append(adj[e.u], edge{e.u, e.v, e.w})
			adj[e.v] = append(adj[e.v], edge{e.v, e.u, e.w})
		}
	}
	res := make([]int64, len(queries))
	for idx, qu := range queries {
		a := qu[0]
		b := qu[1]
		// BFS on MST
		type node struct {
			id int
			mx int64
		}
		visited := make([]bool, k+1)
		queue := []node{{a, 0}}
		visited[a] = true
		ans := int64(0)
		for len(queue) > 0 {
			cur := queue[0]
			queue = queue[1:]
			if cur.id == b {
				ans = cur.mx
				break
			}
			for _, e := range adj[cur.id] {
				if !visited[e.v] {
					visited[e.v] = true
					mx := cur.mx
					if e.w > mx {
						mx = e.w
					}
					queue = append(queue, node{e.v, mx})
				}
			}
		}
		res[idx] = ans
	}
	return res
}

func generateCaseF(rng *rand.Rand) (int, int, int, int, []edge, [][2]int) {
	n := rng.Intn(6) + 1
	k := rng.Intn(n) + 1
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges) + 1
	edges := make([]edge, 0, m)
	used := make(map[[2]int]bool)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		key := [2]int{u, v}
		if used[key] {
			continue
		}
		used[key] = true
		w := int64(rng.Intn(10) + 1)
		edges = append(edges, edge{u, v, w})
	}
	q := rng.Intn(5) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		queries[i][0] = rng.Intn(k) + 1
		queries[i][1] = rng.Intn(k) + 1
	}
	return n, m, k, q, edges, queries
}

func runCaseF(bin string, n, m, k, q int, edges []edge, queries [][2]int) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, m, k, q))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	for _, qu := range queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", qu[0], qu[1]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	expected := expectedF(n, m, k, q, edges, queries)
	outLines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(outLines) != q {
		return fmt.Errorf("expected %d lines got %d", q, len(outLines))
	}
	for i := 0; i < q; i++ {
		var val int64
		if _, err := fmt.Sscan(strings.TrimSpace(outLines[i]), &val); err != nil {
			return fmt.Errorf("failed to parse line %d: %v", i+1, err)
		}
		if val != expected[i] {
			return fmt.Errorf("query %d expected %d got %d", i+1, expected[i], val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, k, q, edges, queries := generateCaseF(rng)
		if err := runCaseF(bin, n, m, k, q, edges, queries); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
