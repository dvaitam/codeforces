package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return out.String() + stderr.String(), err
	}
	return out.String(), nil
}

type Edge struct {
	to  int
	w   int64
	idx int
}

type Item struct {
	v    int
	dist int64
}

type MinHeap []Item

func (h MinHeap) Len() int            { return len(h) }
func (h MinHeap) Less(i, j int) bool  { return h[i].dist < h[j].dist }
func (h MinHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) { *h = append(*h, x.(Item)) }
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	item := old[n-1]
	*h = old[:n-1]
	return item
}

const INF int64 = 1 << 60

func dijkstra(start int, adj [][]Edge) []int64 {
	n := len(adj) - 1
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	pq := &MinHeap{}
	heap.Init(pq)
	dist[start] = 0
	heap.Push(pq, Item{v: start, dist: 0})
	for pq.Len() > 0 {
		cur := heap.Pop(pq).(Item)
		if cur.dist != dist[cur.v] {
			continue
		}
		for _, e := range adj[cur.v] {
			nd := cur.dist + e.w
			if nd < dist[e.to] {
				dist[e.to] = nd
				heap.Push(pq, Item{v: e.to, dist: nd})
			}
		}
	}
	return dist
}

func solve(t Test) string {
	n := t.n
	adj := make([][]Edge, n+1)
	radj := make([][]Edge, n+1)
	for i, e := range t.edges {
		a, b, w := int(e[0]), int(e[1]), e[2]
		adj[a] = append(adj[a], Edge{to: b, w: w, idx: i})
		radj[b] = append(radj[b], Edge{to: a, w: w, idx: i})
	}
	distS := dijkstra(t.s, adj)
	distT := dijkstra(t.t, radj)
	L := distS[t.t]
	dag := make([][]Edge, n+1)
	revDag := make([][]Edge, n+1)
	for _, e := range t.edges {
		u, v, w := int(e[0]), int(e[1]), e[2]
		if distS[u] != INF && distT[v] != INF && distS[u]+w+distT[v] == L {
			dag[u] = append(dag[u], Edge{to: v, w: w})
			revDag[v] = append(revDag[v], Edge{to: u, w: w})
		}
	}
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i + 1
	}
	sort.Slice(order, func(i, j int) bool { return distS[order[i]] < distS[order[j]] })
	f := make([]*big.Int, n+1)
	for i := 1; i <= n; i++ {
		f[i] = big.NewInt(0)
	}
	f[t.s].SetInt64(1)
	for _, u := range order {
		for _, e := range dag[u] {
			f[e.to].Add(f[e.to], f[u])
		}
	}
	revOrder := make([]int, n)
	copy(revOrder, order)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		revOrder[i], revOrder[j] = revOrder[j], revOrder[i]
	}
	g := make([]*big.Int, n+1)
	for i := 1; i <= n; i++ {
		g[i] = big.NewInt(0)
	}
	g[t.t].SetInt64(1)
	for _, u := range revOrder {
		for _, e := range revDag[u] {
			g[e.to].Add(g[e.to], g[u])
		}
	}
	total := new(big.Int).Set(f[t.t])
	var sb strings.Builder
	for _, e := range t.edges {
		u, v, w := int(e[0]), int(e[1]), e[2]
		if distS[u] == INF || distT[v] == INF {
			sb.WriteString("NO\n")
			continue
		}
		pathLen := distS[u] + w + distT[v]
		if pathLen == L {
			prod := new(big.Int).Mul(f[u], g[v])
			if prod.Cmp(total) == 0 {
				sb.WriteString("YES\n")
				continue
			}
		}
		T := L - distS[u] - distT[v]
		if T > 1 {
			newW := T - 1
			if newW >= 1 && newW < w {
				sb.WriteString("CAN " + strconv.FormatInt(w-newW, 10) + "\n")
				continue
			}
		}
		sb.WriteString("NO\n")
	}
	return strings.TrimSpace(sb.String())
}

type Test struct {
	n     int
	s     int
	t     int
	edges [][3]int64
	input string
}

func genTest(rng *rand.Rand) Test {
	n := rng.Intn(4) + 2 // 2..5
	s := 1
	t := n
	edges := make([][3]int64, 0)
	// ensure a path chain 1->2->...->n
	for i := 1; i < n; i++ {
		edges = append(edges, [3]int64{int64(i), int64(i + 1), int64(rng.Intn(5) + 1)})
	}
	mExtra := rng.Intn(5)
	for i := 0; i < mExtra; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			if a < n {
				b = a + 1
			} else {
				b = a - 1
			}
		}
		w := int64(rng.Intn(5) + 1)
		edges = append(edges, [3]int64{int64(a), int64(b), w})
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d %d\n", n, len(edges), s, t))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	return Test{n: n, s: s, t: t, edges: edges, input: sb.String()}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		t := genTest(rng)
		expected := solve(t)
		out, err := run(bin, t.input)
		if err != nil {
			fmt.Printf("test %d runtime error: %v\noutput:\n%s", i+1, err, out)
			os.Exit(1)
		}
		out = strings.TrimSpace(out)
		if out != expected {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, t.input, expected, out)
			os.Exit(1)
		}
	}
	fmt.Println("ok 100 tests")
}
