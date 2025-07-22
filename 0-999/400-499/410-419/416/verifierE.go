package main

import (
	"bufio"
	"bytes"
	"container/heap"
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

const INF = math.MaxInt64 / 4

type Edge struct {
	to int
	w  int64
}

type EdgeRec struct {
	u, v int
	w    int64
}

type Item struct {
	node int
	dist int64
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func solve(input string) string {
	in := bufio.NewReader(strings.NewReader(input))
	var n, m int
	fmt.Fscan(in, &n, &m)
	adj := make([][]Edge, n)
	edges := make([]EdgeRec, m)
	for i := 0; i < m; i++ {
		var x, y int
		var l int64
		fmt.Fscan(in, &x, &y, &l)
		x--
		y--
		adj[x] = append(adj[x], Edge{y, l})
		adj[y] = append(adj[y], Edge{x, l})
		edges[i] = EdgeRec{x, y, l}
	}
	totalPairs := n * (n - 1) / 2
	result := make([]int, 0, totalPairs)
	blocks := (n + 63) / 64
	dist := make([]int64, n)
	inDeg := make([]int, n)
	var adjDAG [][]int
	Reach := make([][]uint64, n)
	for i := 0; i < n; i++ {
		Reach[i] = make([]uint64, blocks)
	}
	for s := 0; s < n; s++ {
		for i := 0; i < n; i++ {
			dist[i] = INF
		}
		dist[s] = 0
		pq := &PriorityQueue{}
		heap.Init(pq)
		heap.Push(pq, Item{s, 0})
		for pq.Len() > 0 {
			it := heap.Pop(pq).(Item)
			u, d := it.node, it.dist
			if d != dist[u] {
				continue
			}
			for _, e := range adj[u] {
				nd := d + e.w
				if nd < dist[e.to] {
					dist[e.to] = nd
					heap.Push(pq, Item{e.to, nd})
				}
			}
		}
		adjDAG = make([][]int, n)
		for i := range inDeg {
			inDeg[i] = 0
		}
		for _, e := range edges {
			u, v, w := e.u, e.v, e.w
			if dist[u]+w == dist[v] {
				adjDAG[u] = append(adjDAG[u], v)
				inDeg[v]++
			}
			if dist[v]+w == dist[u] {
				adjDAG[v] = append(adjDAG[v], u)
				inDeg[u]++
			}
		}
		for i := 0; i < n; i++ {
			for b := 0; b < blocks; b++ {
				Reach[i][b] = 0
			}
			Reach[i][i/64] |= 1 << (uint(i) % 64)
		}
		nodes := make([]int, n)
		for i := 0; i < n; i++ {
			nodes[i] = i
		}
		sort.Slice(nodes, func(i, j int) bool { return dist[nodes[i]] > dist[nodes[j]] })
		for _, v := range nodes {
			for _, w := range adjDAG[v] {
				for b := 0; b < blocks; b++ {
					Reach[v][b] |= Reach[w][b]
				}
			}
		}
		ans := make([]int, n)
		for v := 0; v < n; v++ {
			deg := inDeg[v]
			if deg == 0 {
				continue
			}
			for t := s + 1; t < n; t++ {
				if (Reach[v][t/64]>>(uint(t)%64))&1 == 1 {
					ans[t] += deg
				}
			}
		}
		for t := s + 1; t < n; t++ {
			result = append(result, ans[t])
		}
	}
	var sb strings.Builder
	for i, v := range result {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	return sb.String()
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	used := make(map[[2]int]bool)
	for i := 0; i < m; i++ {
		var x, y int
		for {
			x = rng.Intn(n) + 1
			y = rng.Intn(n) + 1
			if x != y {
				if x > y {
					x, y = y, x
				}
				if !used[[2]int{x, y}] {
					used[[2]int{x, y}] = true
					break
				}
			}
		}
		l := rng.Intn(10) + 1
		fmt.Fprintf(&sb, "%d %d %d\n", x, y, l)
	}
	input := sb.String()
	exp := solve(input)
	return input, strings.TrimSpace(exp)
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if out != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
