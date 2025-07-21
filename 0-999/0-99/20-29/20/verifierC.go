package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Edge struct {
	to int
	w  int64
}

func shortestPath(n int, edges [][3]int) []int {
	adj := make([][]Edge, n+1)
	for _, e := range edges {
		u, v, w := e[0], e[1], int64(e[2])
		adj[u] = append(adj[u], Edge{v, w})
		adj[v] = append(adj[v], Edge{u, w})
	}
	const inf = int64(1<<63 - 1)
	dist := make([]int64, n+1)
	prev := make([]int, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = inf
	}
	dist[1] = 0
	pq := &itemPQ{}
	heap.Push(pq, &item{1, 0})
	for pq.Len() > 0 {
		it := heap.Pop(pq).(*item)
		if it.d != dist[it.v] {
			continue
		}
		if it.v == n {
			break
		}
		for _, e := range adj[it.v] {
			nd := it.d + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				prev[e.to] = it.v
				heap.Push(pq, &item{e.to, nd})
			}
		}
	}
	if dist[n] == inf {
		return nil
	}
	var path []int
	for v := n; v != 0; v = prev[v] {
		path = append(path, v)
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path
}

type item struct {
	v   int
	d   int64
	idx int
}

type itemPQ []*item

func (pq itemPQ) Len() int           { return len(pq) }
func (pq itemPQ) Less(i, j int) bool { return pq[i].d < pq[j].d }
func (pq itemPQ) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].idx = i
	pq[j].idx = j
}
func (pq *itemPQ) Push(x interface{}) {
	it := x.(*item)
	it.idx = len(*pq)
	*pq = append(*pq, it)
}
func (pq *itemPQ) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(7) + 2 // 2..8
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	edges := make([][3]int, m)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		w := rng.Intn(100) + 1
		edges[i] = [3]int{u, v, w}
		sb.WriteString(fmt.Sprintf("%d %d %d\n", u, v, w))
	}
	path := shortestPath(n, edges)
	return sb.String(), path
}

func runCase(bin, input string, expected []int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	outStr := strings.TrimSpace(out.String())
	if expected == nil {
		if outStr != "-1" {
			return fmt.Errorf("expected -1 got %s", outStr)
		}
		return nil
	}
	tokens := strings.Fields(outStr)
	if len(tokens) != len(expected) {
		return fmt.Errorf("expected %d nodes got %d", len(expected), len(tokens))
	}
	for i, exp := range expected {
		val, err := strconv.Atoi(tokens[i])
		if err != nil {
			return fmt.Errorf("failed to parse int '%s'", tokens[i])
		}
		if val != exp {
			return fmt.Errorf("position %d: expected %d got %d", i+1, exp, val)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
