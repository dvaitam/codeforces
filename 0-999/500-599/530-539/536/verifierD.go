package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

// ── embedded solver (CF-accepted 536D) ──────────────────────────────

type Edge struct {
	to     int
	weight int64
}

type Item struct {
	v int
	d int64
}

type MinHeap []Item

func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i].d < h[j].d }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x interface{}) {
	*h = append(*h, x.(Item))
}
func (h *MinHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func dijkstra(start int, n int, adj [][]Edge) []int64 {
	dist := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		dist[i] = -1
	}
	dist[start] = 0
	h := &MinHeap{}
	heap.Init(h)
	heap.Push(h, Item{v: start, d: 0})
	for h.Len() > 0 {
		curr := heap.Pop(h).(Item)
		u := curr.v
		d := curr.d
		if d > dist[u] {
			continue
		}
		for _, edge := range adj[u] {
			v := edge.to
			w := edge.weight
			if dist[v] == -1 || dist[u]+w < dist[v] {
				dist[v] = dist[u] + w
				heap.Push(h, Item{v: v, d: dist[v]})
			}
		}
	}
	return dist
}

func uniqueSorted(a []int64) []int64 {
	sort.Slice(a, func(i, j int) bool { return a[i] < a[j] })
	res := make([]int64, 0, len(a))
	for i, v := range a {
		if i == 0 || v != a[i-1] {
			res = append(res, v)
		}
	}
	return res
}

func getRank(a []int64, val int64) int {
	l, r := 0, len(a)-1
	for l <= r {
		m := l + (r-l)/2
		if a[m] == val {
			return m
		} else if a[m] < val {
			l = m + 1
		} else {
			r = m - 1
		}
	}
	return -1
}

func solve(input string) string {
	data := []byte(input)
	pos := 0

	nextInt := func() int {
		for pos < len(data) && data[pos] <= ' ' {
			pos++
		}
		if pos >= len(data) {
			return 0
		}
		sign := 1
		if data[pos] == '-' {
			sign = -1
			pos++
		}
		res := 0
		for pos < len(data) && data[pos] > ' ' {
			res = res*10 + int(data[pos]-'0')
			pos++
		}
		return res * sign
	}

	nextInt64 := func() int64 {
		for pos < len(data) && data[pos] <= ' ' {
			pos++
		}
		if pos >= len(data) {
			return 0
		}
		sign := int64(1)
		if data[pos] == '-' {
			sign = -1
			pos++
		}
		res := int64(0)
		for pos < len(data) && data[pos] > ' ' {
			res = res*10 + int64(data[pos]-'0')
			pos++
		}
		return res * sign
	}

	n := nextInt()
	if n == 0 {
		return ""
	}
	m := nextInt()
	s := nextInt()
	t := nextInt()

	P := make([]int64, n+1)
	for i := 1; i <= n; i++ {
		P[i] = nextInt64()
	}

	adj := make([][]Edge, n+1)
	for i := 0; i < m; i++ {
		u := nextInt()
		v := nextInt()
		w := nextInt64()
		adj[u] = append(adj[u], Edge{to: v, weight: w})
		adj[v] = append(adj[v], Edge{to: u, weight: w})
	}

	distS := dijkstra(s, n, adj)
	distT := dijkstra(t, n, adj)

	distS_copy := make([]int64, 0, n+1)
	distT_copy := make([]int64, 0, n+1)
	distS_copy = append(distS_copy, -1)
	distT_copy = append(distT_copy, -1)

	for i := 1; i <= n; i++ {
		distS_copy = append(distS_copy, distS[i])
		distT_copy = append(distT_copy, distT[i])
	}

	distS_distinct := uniqueSorted(distS_copy)
	distT_distinct := uniqueSorted(distT_copy)

	ks := len(distS_distinct) - 1
	kt := len(distT_distinct) - 1

	Rs := make([]int, n+1)
	Rt := make([]int, n+1)
	for i := 1; i <= n; i++ {
		Rs[i] = getRank(distS_distinct, distS[i])
		Rt[i] = getRank(distT_distinct, distT[i])
	}

	cols := kt + 1
	sumP := make([]int64, (ks+1)*cols)
	for v := 1; v <= n; v++ {
		sumP[Rs[v]*cols+Rt[v]] += P[v]
	}

	F := make([]int64, (ks+1)*cols)
	row_suf := make([]int64, kt+2)
	for i := 0; i <= ks; i++ {
		for j := 0; j <= kt+1; j++ {
			row_suf[j] = 0
		}
		for j := kt; j >= 0; j-- {
			row_suf[j] = row_suf[j+1] + sumP[i*cols+j]
		}
		for j := 0; j <= kt; j++ {
			if i > 0 {
				F[i*cols+j] = F[(i-1)*cols+j] + row_suf[j+1]
			} else {
				F[i*cols+j] = row_suf[j+1]
			}
		}
	}

	G := make([]int64, (ks+1)*cols)
	col_suf := make([]int64, ks+2)
	for j := 0; j <= kt; j++ {
		for i := 0; i <= ks+1; i++ {
			col_suf[i] = 0
		}
		for i := ks; i >= 0; i-- {
			col_suf[i] = col_suf[i+1] + sumP[i*cols+j]
		}
		for i := 0; i <= ks; i++ {
			if j > 0 {
				G[i*cols+j] = G[i*cols+(j-1)] + col_suf[i+1]
			} else {
				G[i*cols+j] = col_suf[i+1]
			}
		}
	}

	max_Rt := make([]int, ks+1)
	max_Rs := make([]int, kt+1)
	for i := 0; i <= ks; i++ {
		max_Rt[i] = -1
	}
	for j := 0; j <= kt; j++ {
		max_Rs[j] = -1
	}
	for v := 1; v <= n; v++ {
		rS := Rs[v]
		rT := Rt[v]
		if rT > max_Rt[rS] {
			max_Rt[rS] = rT
		}
		if rS > max_Rs[rT] {
			max_Rs[rT] = rS
		}
	}

	next_i := make([]int32, (ks+1)*cols)
	for j := 0; j <= kt; j++ {
		next_i[ks*cols+j] = int32(ks + 1)
		for i := ks - 1; i >= 0; i-- {
			if max_Rt[i+1] > j {
				next_i[i*cols+j] = int32(i + 1)
			} else {
				next_i[i*cols+j] = next_i[(i+1)*cols+j]
			}
		}
	}

	next_j := make([]int32, (ks+1)*cols)
	for i := 0; i <= ks; i++ {
		next_j[i*cols+kt] = int32(kt + 1)
		for j := kt - 1; j >= 0; j-- {
			if max_Rs[j+1] > i {
				next_j[i*cols+j] = int32(j + 1)
			} else {
				next_j[i*cols+j] = next_j[i*cols+(j+1)]
			}
		}
	}

	suf_T := make([]int64, (ks+2)*cols)
	suf_N := make([]int64, kt+2)
	const INF = int64(2e18)

	for i := 0; i < len(suf_T); i++ {
		suf_T[i] = -INF
	}

	var dp0_00 int64

	for i := ks; i >= 0; i-- {
		for j := 0; j <= kt+1; j++ {
			suf_N[j] = -INF
		}

		for j := kt; j >= 0; j-- {
			idx := i*cols + j

			var dp0 int64
			ni := int(next_i[idx])
			if ni > ks {
				dp0 = 0
			} else {
				dp0 = -F[idx] + suf_T[ni*cols+j]
			}

			var dp1 int64
			nj := int(next_j[idx])
			if nj > kt {
				dp1 = 0
			} else {
				dp1 = -G[idx] + suf_N[nj]
			}

			if i == 0 && j == 0 {
				dp0_00 = dp0
			}

			val_T := F[idx] - dp1
			if val_T > suf_T[(i+1)*cols+j] {
				suf_T[i*cols+j] = val_T
			} else {
				suf_T[i*cols+j] = suf_T[(i+1)*cols+j]
			}

			val_N := G[idx] - dp0
			if val_N > suf_N[j+1] {
				suf_N[j] = val_N
			} else {
				suf_N[j] = suf_N[j+1]
			}
		}
	}

	if dp0_00 > 0 {
		return "Break a heart"
	} else if dp0_00 < 0 {
		return "Cry"
	}
	return "Flowers"
}

// ── test generation ─────────────────────────────────────────────────

type edge struct {
	u, v int
	w    int64
}

func buildInput(n int, s, t int, values []int64, edges []edge) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", n, len(edges))
	fmt.Fprintf(&sb, "%d %d\n", s, t)
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", values[i])
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d %d\n", e.u, e.v, e.w)
	}
	return sb.String()
}

func deterministicTests() []string {
	var tests []string
	tests = append(tests, buildInput(
		2, 1, 2,
		[]int64{5, -3},
		[]edge{{1, 2, 1}},
	))
	tests = append(tests, buildInput(
		3, 1, 3,
		[]int64{10, 20, 30},
		[]edge{{1, 2, 2}, {2, 3, 2}, {1, 3, 5}},
	))
	tests = append(tests, buildInput(
		4, 2, 4,
		[]int64{1, 2, 3, 4},
		[]edge{{1, 1, 0}, {1, 2, 1}, {2, 3, 2}, {3, 4, 3}, {4, 2, 1}},
	))
	return tests
}

func randIntR(rnd *rand.Rand, l, r int) int {
	return rnd.Intn(r-l+1) + l
}

func randomValues(n int, rnd *rand.Rand) []int64 {
	vals := make([]int64, n)
	for i := range vals {
		vals[i] = rnd.Int63n(2_000_000_001) - 1_000_000_000
	}
	return vals
}

func randomGraph(n int, rnd *rand.Rand) []edge {
	maxEdges := 100000
	edges := make([]edge, 0, n-1)
	for i := 2; i <= n; i++ {
		parent := randIntR(rnd, 1, i-1)
		w := rnd.Int63n(1_000_000_000 + 1)
		edges = append(edges, edge{parent, i, w})
	}
	remaining := maxEdges - len(edges)
	if remaining < 0 {
		remaining = 0
	}
	additionalLimit := minInt(remaining, n*4)
	extraCount := 0
	if additionalLimit > 0 {
		extraCount = rnd.Intn(additionalLimit + 1)
	}
	for i := 0; i < extraCount; i++ {
		u := randIntR(rnd, 1, n)
		v := randIntR(rnd, 1, n)
		w := rnd.Int63n(1_000_000_000 + 1)
		edges = append(edges, edge{u, v, w})
	}
	return edges
}

func randomTests(count int) []string {
	rnd := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 0, count)
	for i := 0; i < count; i++ {
		var n int
		switch {
		case i%25 == 0:
			n = randIntR(rnd, 500, 1000)
		case i%5 == 0:
			n = randIntR(rnd, 100, 500)
		default:
			n = randIntR(rnd, 2, 80)
		}
		values := randomValues(n, rnd)
		edges := randomGraph(n, rnd)
		s := randIntR(rnd, 1, n)
		t := randIntR(rnd, 1, n)
		for t == s {
			t = randIntR(rnd, 1, n)
		}
		tests = append(tests, buildInput(n, s, t, values, edges))
	}
	n := 2000
	values := randomValues(n, rnd)
	edges := randomGraph(n, rnd)
	tests = append(tests, buildInput(n, 1, n, values, edges))
	return tests
}

func normalizeOutput(out string) (string, error) {
	out = strings.TrimSpace(out)
	if out == "" {
		return "", fmt.Errorf("empty output")
	}
	switch out {
	case "Break a heart", "Cry", "Flowers":
		return out, nil
	default:
		return "", fmt.Errorf("invalid outcome %q", out)
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

// suppress unused import
var _ = io.Discard

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := deterministicTests()
	tests = append(tests, randomTests(150)...)

	for idx, input := range tests {
		exp := solve(input)
		gotOut, err := runProgram(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		expNorm, err := normalizeOutput(exp)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: oracle output invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := normalizeOutput(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if expNorm != got {
			fmt.Fprintf(os.Stderr, "case %d mismatch: expected %q got %q\n", idx+1, expNorm, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
