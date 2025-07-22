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

func runCandidate(bin, input string) (string, error) {
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

type e struct{ to, w int }

func dijkstra(n, src, skip int, adj [][]e) []int64 {
	const INF = int64(1 << 60)
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = INF
	}
	if src == skip {
		return dist
	}
	dist[src] = 0
	h := &hp{{src, 0}}
	heap.Init(h)
	for h.Len() > 0 {
		cur := heap.Pop(h).(pr)
		if cur.d != dist[cur.u] {
			continue
		}
		for _, ed := range adj[cur.u] {
			if ed.to == skip || cur.u == skip {
				continue
			}
			nd := cur.d + int64(ed.w)
			if nd < dist[ed.to] {
				dist[ed.to] = nd
				heap.Push(h, pr{ed.to, nd})
			}
		}
	}
	return dist
}

type pr struct {
	u int
	d int64
}
type hp []pr

func (h hp) Len() int            { return len(h) }
func (h hp) Less(i, j int) bool  { return h[i].d < h[j].d }
func (h hp) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *hp) Push(x interface{}) { *h = append(*h, x.(pr)) }
func (h *hp) Pop() interface{}   { old := *h; x := old[len(old)-1]; *h = old[:len(old)-1]; return x }

func solveCase(n, m int, spec []int, edges [][3]int) string {
	adj := make([][]e, n+1)
	for _, ed := range edges {
		a, b, c := ed[0], ed[1], ed[2]
		adj[a] = append(adj[a], e{b, c})
		adj[b] = append(adj[b], e{a, c})
	}
	isSpec := make([]bool, n+1)
	for _, s := range spec {
		isSpec[s] = true
	}
	// all pairs distances
	dist := make([][]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = dijkstra(n, i, -1, adj)
	}
	far := make([][]int, n+1)
	for _, s := range spec {
		md := int64(-1)
		for _, t := range spec {
			if dist[s][t] > md {
				md = dist[s][t]
				far[s] = []int{t}
			} else if dist[s][t] == md {
				far[s] = append(far[s], t)
			}
		}
	}
	best := -1
	ways := 0
	for v := 1; v <= n; v++ {
		if isSpec[v] {
			continue
		}
		unhappy := 0
		for _, s := range spec {
			d := dijkstra(n, s, v, adj)
			unreachable := true
			for _, f := range far[s] {
				if d[f] < int64(1<<60) {
					unreachable = false
					break
				}
			}
			if unreachable {
				unhappy++
			}
		}
		if unhappy > best {
			best = unhappy
			ways = 1
		} else if unhappy == best {
			ways++
		}
	}
	return fmt.Sprintf("%d %d", best, ways)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(6) + 3
	m := rng.Intn(n-1) + 2
	edges := make([][3]int, n-1)
	for i := 2; i <= n; i++ {
		p := rng.Intn(i-1) + 1
		w := rng.Intn(10) + 1
		edges[i-2] = [3]int{p, i, w}
	}
	perm := rng.Perm(n)
	spec := perm[:m]
	for i := 0; i < m; i++ {
		spec[i]++
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i, s := range spec {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprintf("%d", s))
	}
	sb.WriteByte('\n')
	for _, ed := range edges {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", ed[0], ed[1], ed[2]))
	}
	expect := solveCase(n, m, spec, edges)
	return sb.String(), expect
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
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
