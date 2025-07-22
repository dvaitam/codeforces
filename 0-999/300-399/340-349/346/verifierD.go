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
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

const INF = 1000000007

type Item struct{ v, f int }
type PQ []Item

func (pq PQ) Len() int            { return len(pq) }
func (pq PQ) Less(i, j int) bool  { return pq[i].f < pq[j].f }
func (pq PQ) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PQ) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PQ) Pop() interface{} {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[:n-1]
	return x
}

func expectedOrders(n int, edges [][2]int, s, t int) int {
	fwd := make([][]int, n+1)
	rev := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		fwd[u] = append(fwd[u], v)
		rev[v] = append(rev[v], u)
	}
	reachS := make([]bool, n+1)
	q := []int{s}
	reachS[s] = true
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for _, v := range fwd[u] {
			if !reachS[v] {
				reachS[v] = true
				q = append(q, v)
			}
		}
	}
	reachT := make([]bool, n+1)
	q = []int{t}
	reachT[t] = true
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for _, v := range rev[u] {
			if !reachT[v] {
				reachT[v] = true
				q = append(q, v)
			}
		}
	}
	if !reachS[t] {
		return -1
	}
	good := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		good[i] = reachS[i] && reachT[i]
	}
	deg := make([]int, n+1)
	revGood := make([][]int, n+1)
	for u := 1; u <= n; u++ {
		if !good[u] {
			continue
		}
		for _, v := range fwd[u] {
			if good[v] {
				deg[u]++
				revGood[v] = append(revGood[v], u)
			}
		}
	}

	f := make([]int, n+1)
	minF := make([]int, n+1)
	maxF := make([]int, n+1)
	cnt := make([]int, n+1)
	for i := 1; i <= n; i++ {
		f[i] = INF
		minF[i] = INF
	}
	f[t] = 0
	pq := &PQ{}
	heap.Init(pq)
	heap.Push(pq, Item{v: t, f: 0})
	for pq.Len() > 0 {
		it := heap.Pop(pq).(Item)
		v := it.v
		fv := it.f
		if fv != f[v] {
			continue
		}
		for _, u := range revGood[v] {
			cnt[u]++
			if fv < minF[u] {
				minF[u] = fv
			}
			if fv > maxF[u] {
				maxF[u] = fv
			}
			if cnt[u] == deg[u] {
				noOrder := maxF[u]
				order := minF[u] + 1
				cand := noOrder
				if order < cand {
					cand = order
				}
				if cand < f[u] {
					f[u] = cand
					heap.Push(pq, Item{v: u, f: cand})
				}
			}
		}
	}
	if f[s] >= INF {
		return -1
	}
	return f[s]
}

func generateCase(rng *rand.Rand) string {
	n := rng.Intn(5) + 2 // 2..6
	maxEdges := n * (n - 1)
	m := rng.Intn(maxEdges + 1)
	used := make(map[[2]int]bool)
	edges := make([][2]int, 0, m)
	for len(edges) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		e := [2]int{u, v}
		if used[e] {
			continue
		}
		used[e] = true
		edges = append(edges, e)
	}
	s := rng.Intn(n) + 1
	t := rng.Intn(n) + 1
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	fmt.Fprintf(&sb, "%d %d\n", s, t)
	return sb.String()
}

func parseCase(input string) (int, [][2]int, int, int) {
	r := strings.NewReader(input)
	var n, m int
	fmt.Fscan(r, &n, &m)
	edges := make([][2]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(r, &edges[i][0], &edges[i][1])
	}
	var s, t int
	fmt.Fscan(r, &s, &t)
	return n, edges, s, t
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	cases := []string{}
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for idx, tc := range cases {
		n, edges, s, t := parseCase(tc)
		expect := expectedOrders(n, edges, s, t)
		got, err := runBinary(bin, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, tc)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != fmt.Sprint(expect) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\ninput:\n%s", idx+1, expect, got, tc)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
