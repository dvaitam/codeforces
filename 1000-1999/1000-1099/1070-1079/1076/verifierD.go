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
	to, w, idx int
}

type item struct {
	node int
	dist int64
}

type pq []item

func (p pq) Len() int            { return len(p) }
func (p pq) Less(i, j int) bool  { return p[i].dist < p[j].dist }
func (p pq) Swap(i, j int)       { p[i], p[j] = p[j], p[i] }
func (p *pq) Push(x interface{}) { *p = append(*p, x.(item)) }
func (p *pq) Pop() interface{} {
	old := *p
	v := old[len(old)-1]
	*p = old[:len(old)-1]
	return v
}

func solve(n, m, k int, edges [][4]int) string {
	adj := make([][]Edge, n+1)
	for _, e := range edges {
		a, b, c, idx := e[0], e[1], e[2], e[3]
		adj[a] = append(adj[a], Edge{b, c, idx})
		adj[b] = append(adj[b], Edge{a, c, idx})
	}
	const inf int64 = 1 << 60
	dist := make([]int64, n+1)
	for i := range dist {
		dist[i] = inf
	}
	dist[1] = 0
	parent := make([]int, n+1)
	parentEdge := make([]int, n+1)
	q := &pq{}
	heap.Push(q, item{1, 0})
	for q.Len() > 0 {
		it := heap.Pop(q).(item)
		if it.dist != dist[it.node] {
			continue
		}
		u := it.node
		for _, e := range adj[u] {
			nd := it.dist + int64(e.w)
			if nd < dist[e.to] {
				dist[e.to] = nd
				parent[e.to] = u
				parentEdge[e.to] = e.idx
				heap.Push(q, item{e.to, nd})
			}
		}
	}
	if k > n-1 {
		k = n - 1
	}
	children := make([][]int, n+1)
	edgesIdx := make([][]int, n+1)
	for v := 2; v <= n; v++ {
		p := parent[v]
		if p != 0 {
			children[p] = append(children[p], v)
			edgesIdx[p] = append(edgesIdx[p], parentEdge[v])
		}
	}
	ans := make([]int, 0, k)
	type stackEntry struct{ u, next int }
	stack := []stackEntry{{1, 0}}
	for len(stack) > 0 && len(ans) < k {
		top := &stack[len(stack)-1]
		if top.next < len(children[top.u]) {
			v := children[top.u][top.next]
			id := edgesIdx[top.u][top.next]
			top.next++
			ans = append(ans, id)
			stack = append(stack, stackEntry{v, 0})
		} else {
			stack = stack[:len(stack)-1]
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(ans)))
	for i, v := range ans {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", v))
	}
	if len(ans) > 0 {
		sb.WriteByte('\n')
	}
	return strings.TrimSpace(sb.String())
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	// create tree first
	edges := make([][4]int, 0)
	idx := 1
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		w := rng.Intn(20) + 1
		edges = append(edges, [4]int{p, i, w, idx})
		idx++
	}
	m := len(edges)
	// add extra edges
	extra := rng.Intn(3)
	for e := 0; e < extra; e++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(n) + 1
		if a == b {
			b = (b % n) + 1
		}
		w := rng.Intn(20) + 1
		edges = append(edges, [4]int{a, b, w, idx})
		idx++
	}
	m = len(edges)
	k := rng.Intn(m + 1)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", e[0], e[1], e[2]))
	}
	expect := solve(n, m, k, edges)
	return sb.String(), expect
}

func runCase(bin, in, expect string) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(in)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	got := strings.TrimSpace(out.String())
	expect = strings.TrimSpace(expect)
	if got != expect {
		return fmt.Errorf("expected:\n%s\n\ngot:\n%s", expect, got)
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
		in, expect := generateCase(rng)
		if err := runCase(bin, in, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
