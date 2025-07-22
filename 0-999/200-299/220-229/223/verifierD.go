package main

import (
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

type edge struct {
	to   int
	dist float64
}

type item struct {
	node int
	dist float64
	idx  int
}

type pq []*item

func (p pq) Len() int           { return len(p) }
func (p pq) Less(i, j int) bool { return p[i].dist < p[j].dist }
func (p pq) Swap(i, j int)      { p[i], p[j] = p[j], p[i]; p[i].idx = i; p[j].idx = j }
func (p *pq) Push(x interface{}) {
	it := x.(*item)
	it.idx = len(*p)
	*p = append(*p, it)
}
func (p *pq) Pop() interface{} {
	old := *p
	n := len(old)
	it := old[n-1]
	*p = old[:n-1]
	return it
}

func solvePolygon(xs, ys []int, s, t int) float64 {
	n := len(xs)
	adj := make([][]edge, n)
	for i := 0; i < n; i++ {
		j := (i + 1) % n
		dx := float64(xs[i] - xs[j])
		dy := float64(ys[i] - ys[j])
		d := math.Hypot(dx, dy)
		adj[i] = append(adj[i], edge{j, d})
		adj[j] = append(adj[j], edge{i, d})
	}
	mp := make(map[int][]int)
	for i, x := range xs {
		mp[x] = append(mp[x], i)
	}
	for _, vec := range mp {
		if len(vec) < 2 {
			continue
		}
		sort.Slice(vec, func(i, j int) bool { return ys[vec[i]] > ys[vec[j]] })
		for k := 0; k+1 < len(vec); k++ {
			u := vec[k]
			v := vec[k+1]
			dy := float64(ys[u] - ys[v])
			adj[u] = append(adj[u], edge{v, dy})
		}
	}
	N := n
	const inf = 1e100
	dist := make([]float64, N)
	for i := range dist {
		dist[i] = inf
	}
	dist[s] = 0
	visited := make([]bool, N)
	q := &pq{}
	heap.Init(q)
	heap.Push(q, &item{node: s, dist: 0})
	for q.Len() > 0 {
		it := heap.Pop(q).(*item)
		u := it.node
		if visited[u] {
			continue
		}
		visited[u] = true
		if u == t {
			break
		}
		for _, e := range adj[u] {
			v := e.to
			nd := dist[u] + e.dist
			if nd < dist[v] {
				dist[v] = nd
				heap.Push(q, &item{node: v, dist: nd})
			}
		}
	}
	return dist[t]
}

func genCase(rng *rand.Rand) (xs, ys []int, s, t int) {
	w := rng.Intn(9) + 1
	h := rng.Intn(9) + 1
	xs = []int{0, w, w, 0}
	ys = []int{0, 0, h, h}
	s = rng.Intn(4)
	t = rng.Intn(4)
	return
}

func runCase(bin string, xs, ys []int, s, t int, expected float64) error {
	var sb strings.Builder
	n := len(xs)
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", xs[i], ys[i]))
	}
	sb.WriteString(fmt.Sprintf("%d %d\n", s+1, t+1))
	input := sb.String()
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got float64
	if _, err := fmt.Sscan(strings.TrimSpace(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if math.Abs(got-expected) > 1e-6 {
		return fmt.Errorf("expected %.6f got %.6f", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		xs, ys, s, t := genCase(rng)
		exp := solvePolygon(xs, ys, s, t)
		if err := runCase(bin, xs, ys, s, t, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
