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

type edge struct {
	u, v int
	w    int64
}

type item struct {
	v    int
	dist int64
}

type minHeap []item

func (h minHeap) Len() int            { return len(h) }
func (h minHeap) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h minHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}) { *h = append(*h, x.(item)) }
func (h *minHeap) Pop() interface{} {
	old := *h
	n := len(old)
	it := old[n-1]
	*h = old[:n-1]
	return it
}

func expected(n, m, s int, L int64, edges []edge) string {
	adj := make([][]struct {
		to int
		w  int64
	}, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], struct {
			to int
			w  int64
		}{e.v, e.w})
		adj[e.v] = append(adj[e.v], struct {
			to int
			w  int64
		}{e.u, e.w})
	}
	const INF = int64(1 << 60)
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	pq := &minHeap{}
	heap.Init(pq)
	dist[s] = 0
	heap.Push(pq, item{v: s, dist: 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(item)
		if cur.dist != dist[cur.v] {
			continue
		}
		for _, e := range adj[cur.v] {
			nd := cur.dist + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, item{v: e.to, dist: nd})
			}
		}
	}
	var ans int64
	for i := 1; i <= n; i++ {
		if dist[i] == L {
			ans++
		}
	}
	for _, e := range edges {
		du := dist[e.u]
		dv := dist[e.v]
		w := e.w
		flagU := du < L && du+w > L && du+dv+w >= 2*L
		flagV := dv < L && dv+w > L && du+dv+w >= 2*L
		if flagU && flagV && du+dv+w == 2*L {
			ans++
		} else {
			if flagU {
				ans++
			}
			if flagV {
				ans++
			}
		}
	}
	return fmt.Sprintf("%d", ans)
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges-(n-1)+1) + (n - 1)
	edges := make([]edge, 0, m)
	parent := make([]int, n+1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		w := int64(rng.Intn(10) + 1)
		edges = append(edges, edge{p, i, w})
		parent[i] = p
	}
	existing := map[[2]int]bool{}
	for _, e := range edges {
		if e.u < e.v {
			existing[[2]int{e.u, e.v}] = true
		} else {
			existing[[2]int{e.v, e.u}] = true
		}
	}
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		key := [2]int{u, v}
		if u > v {
			key = [2]int{v, u}
		}
		if existing[key] {
			continue
		}
		existing[key] = true
		w := int64(rng.Intn(10) + 1)
		edges = append(edges, edge{u, v, w})
	}
	s := rng.Intn(n) + 1
	L := int64(rng.Intn(20) + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, len(edges), s))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.w))
	}
	sb.WriteString(fmt.Sprintf("%d\n", L))
	exp := expected(n, len(edges), s, L, edges)
	return sb.String(), exp
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
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
