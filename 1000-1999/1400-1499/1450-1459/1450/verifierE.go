package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Edge struct{ to, val int }

type testCase struct {
	input    string
	expected string
}

const inf = 1e9

func solveCase(n int, edges [][3]int) string {
	graph := make([][]Edge, n+1)
	addEdge := func(u, v, w int) { graph[u] = append(graph[u], Edge{v, w}) }
	for _, e := range edges {
		x, y, z := e[0], e[1], e[2]
		addEdge(x, y, 1)
		if z != 0 {
			addEdge(y, x, -1)
		} else {
			addEdge(y, x, 1)
		}
	}
	col := make([]int, n+1)
	var dfs func(int, int) bool
	dfs = func(u, c int) bool {
		col[u] = c
		for _, ed := range graph[u] {
			v := ed.to
			if col[v] != 0 {
				if col[v] != -c {
					return false
				}
			} else {
				if !dfs(v, -c) {
					return false
				}
			}
		}
		return true
	}
	dis := make([]int, n+1)
	cnt := make([]int, n+1)
	inq := make([]bool, n+1)
	SPFA := func() bool {
		for i := 1; i <= n; i++ {
			dis[i] = inf
			cnt[i] = 0
			inq[i] = false
		}
		dis[1] = 0
		q := []int{1}
		inq[1] = true
		for len(q) > 0 {
			u := q[0]
			q = q[1:]
			inq[u] = false
			for _, ed := range graph[u] {
				v := ed.to
				if dis[v] > dis[u]+ed.val {
					dis[v] = dis[u] + ed.val
					cnt[v] = cnt[u] + 1
					if cnt[v] > n {
						return true
					}
					if !inq[v] {
						inq[v] = true
						q = append(q, v)
					}
				}
			}
		}
		return false
	}
	type Item struct{ dist, node int }
	type MinHeap []Item
	var heapLess = func(h MinHeap, i, j int) bool { return h[i].dist < h[j].dist }
	var hswap = func(h MinHeap, i, j int) { h[i], h[j] = h[j], h[i] }
	var hpush = func(h *MinHeap, x Item) { *h = append(*h, x) }
	solveFrom := func(src int) (int, []int, bool) {
		dist := make([]int, n+1)
		for i := 1; i <= n; i++ {
			dist[i] = inf
		}
		dist[src] = 0
		h := MinHeap{{0, src}}
		for len(h) > 0 {
			// pop min
			bestIdx := 0
			for i := 1; i < len(h); i++ {
				if heapLess(h, i, bestIdx) {
					bestIdx = i
				}
			}
			it := h[bestIdx]
			hswap(h, bestIdx, len(h)-1)
			h = h[:len(h)-1]
			d, u := it.dist, it.node
			if d > dist[u] {
				continue
			}
			if d < 0 {
				return 0, nil, false
			}
			for _, ed := range graph[u] {
				v := ed.to
				nd := d + ed.val
				if nd < dist[v] {
					dist[v] = nd
					hpush(&h, Item{nd, v})
				}
			}
		}
		mx := 0
		for i := 1; i <= n; i++ {
			if dist[i] > mx {
				mx = dist[i]
			}
		}
		res := make([]int, n)
		for i := 1; i <= n; i++ {
			res[i-1] = dist[i]
		}
		return mx, res, true
	}

	if !dfs(1, 1) || SPFA() {
		return "NO"
	}
	best := -1
	var bestVec []int
	for i := 1; i <= n; i++ {
		mx, vec, ok := solveFrom(i)
		if ok && mx > best {
			best = mx
			bestVec = vec
		}
	}
	var sb strings.Builder
	sb.WriteString("YES\n")
	sb.WriteString(fmt.Sprintf("%d\n", best))
	for i, v := range bestVec {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	return strings.TrimSpace(sb.String())
}

func buildCase(n int, edges [][3]int) testCase {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e[0], e[1], e[2])
	}
	exp := solveCase(n, edges)
	return testCase{input: sb.String(), expected: exp}
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(4) + 2
	m := rng.Intn(n*2) + 1
	edges := make([][3]int, m)
	for i := 0; i < m; i++ {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		for v == u {
			v = rng.Intn(n) + 1
		}
		z := rng.Intn(2)
		edges[i] = [3]int{u, v, z}
	}
	return buildCase(n, edges)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(tc.input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	if got != tc.expected {
		return fmt.Errorf("expected:\n%s\n----\ngot:\n%s", tc.expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	for i := 0; i < 100; i++ {
		cases = append(cases, randomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
