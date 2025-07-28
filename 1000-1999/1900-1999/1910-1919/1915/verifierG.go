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

func runBinary(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

type Edge struct {
	to   int
	cost int64
}

type Item struct {
	dist int64
	node int
	slow int
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

const INF int64 = 1 << 60

func solveCase(n int, edges [][3]int, s []int) string {
	g := make([][]Edge, n+1)
	for _, e := range edges {
		u, v, w := e[0], e[1], int64(e[2])
		g[u] = append(g[u], Edge{to: v, cost: w})
		g[v] = append(g[v], Edge{to: u, cost: w})
	}
	maxS := 1000
	dist := make([][]int64, n+1)
	for i := 0; i <= n; i++ {
		dist[i] = make([]int64, maxS+1)
		for j := 0; j <= maxS; j++ {
			dist[i][j] = INF
		}
	}
	pq := &PriorityQueue{}
	startSlow := s[1]
	dist[1][startSlow] = 0
	heap.Push(pq, Item{dist: 0, node: 1, slow: startSlow})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		if cur.dist != dist[cur.node][cur.slow] {
			continue
		}
		u := cur.node
		d := cur.dist
		slow := cur.slow
		for _, e := range g[u] {
			v := e.to
			ns := slow
			if s[v] < ns {
				ns = s[v]
			}
			nd := d + e.cost*int64(slow)
			if nd < dist[v][ns] {
				dist[v][ns] = nd
				heap.Push(pq, Item{dist: nd, node: v, slow: ns})
			}
		}
	}
	ans := INF
	for i := 1; i <= maxS; i++ {
		if dist[n][i] < ans {
			ans = dist[n][i]
		}
	}
	return fmt.Sprintf("%d", ans)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(4) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-n+1) + n - 1
	edges := make([][3]int, 0, m)
	// ensure connectivity by chain
	for i := 1; i < n; i++ {
		w := rng.Intn(10) + 1
		edges = append(edges, [3]int{i, i + 1, w})
	}
	existing := map[[2]int]bool{}
	for i := 1; i < n; i++ {
		existing[[2]int{i, i + 1}] = true
		existing[[2]int{i + 1, i}] = true
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		key := [2]int{u, v}
		if existing[key] {
			continue
		}
		existing[[2]int{u, v}] = true
		existing[[2]int{v, u}] = true
		w := rng.Intn(10) + 1
		edges = append(edges, [3]int{u, v, w})
	}
	s := make([]int, n+1)
	for i := 1; i <= n; i++ {
		s[i] = rng.Intn(10) + 1
	}
	var sb strings.Builder
	sb.WriteString("1\n")
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
	}
	for i := 1; i <= n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", s[i])
	}
	sb.WriteByte('\n')
	input := sb.String()
	expected := solveCase(n, edges, s)
	return input, expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expected := generateCase(rng)
		got, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
