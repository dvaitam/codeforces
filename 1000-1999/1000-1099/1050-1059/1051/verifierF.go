package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Edge struct {
	to   int
	cost int64
}

type Item struct {
	node int
	dist int64
	idx  int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int           { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i]; pq[i].idx = i; pq[j].idx = j }

func (pq *PriorityQueue) Push(x interface{}) {
	item := x.(Item)
	item.idx = len(*pq)
	*pq = append(*pq, item)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func dijkstra(graph [][]Edge, src, dst int) int64 {
	const INF int64 = 1<<63 - 1
	n := len(graph) - 1
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	dist[src] = 0
	pq := &PriorityQueue{}
	heap.Push(pq, Item{node: src, dist: 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		if cur.dist != dist[cur.node] {
			continue
		}
		if cur.node == dst {
			return cur.dist
		}
		for _, e := range graph[cur.node] {
			nd := cur.dist + e.cost
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{node: e.to, dist: nd})
			}
		}
	}
	return dist[dst]
}

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func solveCase(n, m int, edges [][3]int64, queries [][2]int) []int64 {
	graph := make([][]Edge, n+1)
	for _, e := range edges {
		u := int(e[0])
		v := int(e[1])
		w := e[2]
		graph[u] = append(graph[u], Edge{v, w})
		graph[v] = append(graph[v], Edge{u, w})
	}
	res := make([]int64, len(queries))
	for i, q := range queries {
		res[i] = dijkstra(graph, q[0], q[1])
	}
	return res
}

func genCase(rng *rand.Rand) (int, int, [][3]int64, [][2]int) {
	n := rng.Intn(6) + 2
	extra := rng.Intn(3)
	m := n - 1 + extra
	edges := make([][3]int64, 0, m)
	// create tree first
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		w := rng.Int63n(20) + 1
		edges = append(edges, [3]int64{int64(p), int64(i), w})
	}
	for i := 0; i < extra; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for v == u {
			v = rng.Intn(n) + 1
		}
		w := rng.Int63n(20) + 1
		edges = append(edges, [3]int64{int64(u), int64(v), w})
	}
	q := rng.Intn(5) + 1
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		queries[i] = [2]int{u, v}
	}
	return n, m, edges, queries
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, edges, queries := genCase(rng)
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
		sb.WriteString(fmt.Sprintf("%d\n", len(queries)))
		for _, q := range queries {
			sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
		}
		expect := solveCase(n, m, edges, queries)
		out, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != len(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\n", i+1, len(expect), len(fields))
			os.Exit(1)
		}
		for j, f := range fields {
			var val int64
			if _, err := fmt.Sscan(f, &val); err != nil || val != expect[j] {
				fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s at index %d\n", i+1, expect[j], f, j)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
