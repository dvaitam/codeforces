package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Edge struct {
	to int
	w  int64
	id int
}

type Item struct {
	node int
	dist int64
	idx  int
}

type PQ []*Item

func (pq PQ) Len() int           { return len(pq) }
func (pq PQ) Less(i, j int) bool { return pq[i].dist < pq[j].dist }
func (pq PQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i]; pq[i].idx = i; pq[j].idx = j }

func (pq *PQ) Push(x interface{}) {
	n := len(*pq)
	it := x.(*Item)
	it.idx = n
	*pq = append(*pq, it)
}

func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[:n-1]
	return it
}

func runBinary(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func solveE(input string) string {
	reader := bufio.NewReader(strings.NewReader(input))
	var n, m int
	fmt.Fscan(reader, &n, &m)
	graph := make([][]Edge, n+1)
	edges := make([][3]int64, m)
	for i := 0; i < m; i++ {
		var u, v int
		var w int64
		fmt.Fscan(reader, &u, &v, &w)
		edges[i] = [3]int64{int64(u), int64(v), w}
		graph[u] = append(graph[u], Edge{to: v, w: w, id: i + 1})
		graph[v] = append(graph[v], Edge{to: u, w: w, id: i + 1})
	}
	var start int
	fmt.Fscan(reader, &start)
	const INF = int64(1) << 62
	dist := make([]int64, n+1)
	for i := range dist {
		dist[i] = INF
	}
	dist[start] = 0
	visited := make([]bool, n+1)
	pq := &PQ{}
	heap.Push(pq, &Item{node: start, dist: 0})
	for pq.Len() > 0 {
		it := heap.Pop(pq).(*Item)
		u := it.node
		if visited[u] {
			continue
		}
		visited[u] = true
		for _, e := range graph[u] {
			if dist[e.to] > dist[u]+e.w {
				dist[e.to] = dist[u] + e.w
				heap.Push(pq, &Item{node: e.to, dist: dist[e.to]})
			}
		}
	}
	res := make([]int, 0, n-1)
	var total int64
	for v := 1; v <= n; v++ {
		if v == start {
			continue
		}
		bestW := INF
		bestID := 0
		for _, e := range graph[v] {
			if dist[v] == dist[e.to]+e.w && e.w < bestW {
				bestW = e.w
				bestID = e.id
			}
		}
		total += bestW
		res = append(res, bestID)
	}
	var sb strings.Builder
	sb.WriteString(strconv.FormatInt(total, 10))
	sb.WriteByte('\n')
	for i, id := range res {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(id))
	}
	return sb.String()
}

func genTests() []string {
	rand.Seed(5)
	tests := make([]string, 0, 100)
	for t := 0; t < 100; t++ {
		n := rand.Intn(5) + 2
		maxM := n * (n - 1) / 2
		m := (n - 1) + rand.Intn(maxM-(n-1)+1) // ensure connectivity
		type pair struct{ u, v int }
		used := make(map[pair]bool)
		edges := make([][3]int64, 0, m)
		// first create a random tree to ensure connectivity
		for i := 2; i <= n; i++ {
			v := rand.Intn(i-1) + 1
			w := rand.Int63n(9) + 1
			edges = append(edges, [3]int64{int64(i), int64(v), w})
			used[pair{i, v}] = true
			used[pair{v, i}] = true
		}
		// add remaining edges
		for len(edges) < m {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			if u == v || used[pair{u, v}] {
				continue
			}
			w := rand.Int63n(9) + 1
			edges = append(edges, [3]int64{int64(u), int64(v), w})
			used[pair{u, v}] = true
			used[pair{v, u}] = true
		}
		start := rand.Intn(n) + 1
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
		}
		sb.WriteString(fmt.Sprintf("%d", start))
		tests = append(tests, sb.String())
	}
	return tests
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: verifierE <binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		expected := solveE(tc)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(expected) != strings.TrimSpace(got) {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, tc, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("all tests passed")
}
