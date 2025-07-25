package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

// Edge for Dinic
type Edge struct {
	to, cap, rev int
}

// Dinic structure
type Dinic struct {
	N     int
	G     [][]Edge
	level []int
	iter  []int
}

func NewDinic(n int) *Dinic {
	g := make([][]Edge, n)
	return &Dinic{N: n, G: g, level: make([]int, n), iter: make([]int, n)}
}

func (d *Dinic) AddEdge(u, v, c int) {
	d.G[u] = append(d.G[u], Edge{to: v, cap: c, rev: len(d.G[v])})
	d.G[v] = append(d.G[v], Edge{to: u, cap: 0, rev: len(d.G[u]) - 1})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *Dinic) bfs(s int) {
	for i := range d.level {
		d.level[i] = -1
	}
	q := make([]int, 0, d.N)
	d.level[s] = 0
	q = append(q, s)
	for i := 0; i < len(q); i++ {
		v := q[i]
		for _, e := range d.G[v] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[v] + 1
				q = append(q, e.to)
			}
		}
	}
}

func (d *Dinic) dfs(v, t, f int) int {
	if v == t {
		return f
	}
	for ; d.iter[v] < len(d.G[v]); d.iter[v]++ {
		e := &d.G[v][d.iter[v]]
		if e.cap > 0 && d.level[v] < d.level[e.to] {
			ret := d.dfs(e.to, t, min(f, e.cap))
			if ret > 0 {
				e.cap -= ret
				d.G[e.to][e.rev].cap += ret
				return ret
			}
		}
	}
	return 0
}

func (d *Dinic) MaxFlow(s, t int) int {
	flow := 0
	const INF = int(1e9)
	for {
		d.bfs(s)
		if d.level[t] < 0 {
			break
		}
		for i := range d.iter {
			d.iter[i] = 0
		}
		for {
			f := d.dfs(s, t, INF)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func solveCase(n, m int, A, B []int, edges [][2]int) (bool, [][]int) {
	sumA := 0
	sumB := 0
	for i := 0; i < n; i++ {
		sumA += A[i]
		sumB += B[i]
	}
	if sumA != sumB {
		return false, nil
	}
	s := 0
	t := 2*n + 1
	d := NewDinic(2*n + 2)
	for i := 0; i < n; i++ {
		d.AddEdge(s, i+1, A[i])
		d.AddEdge(i+1, i+1+n, A[i])
		d.AddEdge(i+1+n, t, B[i])
	}
	const INF = int(1e9)
	for _, e := range edges {
		x, y := e[0], e[1]
		d.AddEdge(x, y+n, INF)
		d.AddEdge(y, x+n, INF)
	}
	if d.MaxFlow(s, t) != sumA {
		return false, nil
	}
	ans := make([][]int, n)
	for i := range ans {
		ans[i] = make([]int, n)
	}
	for u := 1; u <= n; u++ {
		for _, e := range d.G[u] {
			if e.to >= n+1 && e.to <= n+n {
				j := e.to - (n + 1)
				revCap := d.G[e.to][e.rev].cap
				ans[u-1][j] = revCap
			}
		}
	}
	return true, ans
}

func generateCase(rng *rand.Rand) (string, bool, []int, []int, map[[2]int]struct{}) {
	n := rng.Intn(5) + 1
	maxEdges := n * (n - 1) / 2
	m := rng.Intn(maxEdges + 1)
	edgeSet := make(map[[2]int]struct{})
	for len(edgeSet) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v {
			continue
		}
		if u > v {
			u, v = v, u
		}
		edgeSet[[2]int{u, v}] = struct{}{}
	}
	edges := make([][2]int, 0, len(edgeSet))
	for e := range edgeSet {
		edges = append(edges, e)
	}
	A := make([]int, n)
	B := make([]int, n)
	for i := 0; i < n; i++ {
		A[i] = rng.Intn(5)
		B[i] = rng.Intn(5)
	}
	if rng.Intn(2) == 0 {
		sumA, sumB := 0, 0
		for _, v := range A {
			sumA += v
		}
		for _, v := range B {
			sumB += v
		}
		diff := sumA - sumB
		B[n-1] += diff
		if B[n-1] < 0 {
			B[n-1] = 0
		}
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, len(edges)))
	for i, v := range A {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for i, v := range B {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", v)
	}
	sb.WriteByte('\n')
	for _, e := range edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	input := sb.String()
	ok, _ := solveCase(n, len(edges), A, B, edges)
	return input, ok, A, B, edgeSet
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func checkOutput(out string, n int, A, B []int, edges map[[2]int]struct{}, expectOK bool) error {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}
	first := strings.TrimSpace(lines[0])
	if first == "NO" {
		if expectOK {
			return fmt.Errorf("expected YES got NO")
		}
		if len(lines) != 1 {
			return fmt.Errorf("extra output after NO")
		}
		return nil
	}
	if first != "YES" {
		return fmt.Errorf("first line must be YES or NO")
	}
	if !expectOK {
		return fmt.Errorf("expected NO got YES")
	}
	if len(lines) != n+1 {
		return fmt.Errorf("expected %d lines, got %d", n+1, len(lines))
	}
	matrix := make([][]int, n)
	for i := 0; i < n; i++ {
		parts := strings.Fields(lines[i+1])
		if len(parts) != n {
			return fmt.Errorf("line %d should have %d numbers", i+1, n)
		}
		row := make([]int, n)
		for j, p := range parts {
			val, err := strconv.Atoi(p)
			if err != nil || val < 0 {
				return fmt.Errorf("invalid number at (%d,%d)", i+1, j+1)
			}
			row[j] = val
		}
		matrix[i] = row
	}
	for i := 0; i < n; i++ {
		sumOut := 0
		sumIn := 0
		for j := 0; j < n; j++ {
			sumOut += matrix[i][j]
			sumIn += matrix[j][i]
			if i != j && matrix[i][j] > 0 {
				p := [2]int{i + 1, j + 1}
				if p[0] > p[1] {
					p[0], p[1] = p[1], p[0]
				}
				if _, ok := edges[p]; !ok {
					return fmt.Errorf("edge %d-%d not allowed", i+1, j+1)
				}
			}
		}
		if sumOut != A[i] {
			return fmt.Errorf("out sum mismatch for city %d", i+1)
		}
		if sumIn != B[i] {
			return fmt.Errorf("in sum mismatch for city %d", i+1)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, ok, A, B, edges := generateCase(rng)
		out, err := run(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if err := checkOutput(out, len(A), A, B, edges, ok); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
