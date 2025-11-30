package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const rawTestcasesData = `
4 4 1 3 4 4
3 4 3 2 3 2 1 3 2 2
2 3 1 1 2 2
1 3 2 3 1 2 1 3
3 1 3 1 3 1 1 1 2 1
3 1 1 2 1 1
4 1 1 3 4 1
3 3 1 3 3 1
3 4 1 4 2 2
3 2 2 2 1 2 2 1
1 2 1 1 1 2
3 2 2 4 3 2 2 2
3 1 2 1 1 1 2 1
1 3 1 2 1 1
3 4 1 1 1 2
1 1 1 1 1 1
4 1 3 1 3 1 1 1 2 1
2 1 2 1 1 1 2 1
2 2 1 2 1 2
3 1 3 4 3 1 1 1 2 1
2 1 2 1 1 1 2 1
3 3 3 2 2 3 3 2 1 1
2 3 2 2 1 3 2 2
4 1 4 4 3 1 1 1 4 1 2 1
2 1 2 4 1 1 2 1
1 1 1 1 1 1
3 2 1 4 2 1
1 3 2 1 1 1 1 2
1 1 1 2 1 1
3 1 3 4 3 1 1 1 2 1
3 1 1 2 2 1
3 3 1 3 3 1
1 3 1 2 1 2
4 2 3 1 1 2 4 1 2 1
3 4 3 3 2 3 2 1 1 4
3 1 2 1 3 1 2 1
3 3 2 2 2 3 3 3
1 1 1 2 1 1
4 1 1 4 3 1
4 2 4 1 3 1 1 1 1 2 4 1
4 3 2 1 2 3 2 1
2 4 4 3 2 3 1 1 2 4 2 2
4 4 2 2 3 3 1 3
1 2 2 2 1 1 1 2
1 1 1 2 1 1
3 1 2 3 3 1 1 1
3 4 1 1 3 1
4 2 3 4 1 1 2 1 4 2
3 3 3 1 3 1 1 1 2 2
2 3 3 1 2 3 1 2 2 1
1 2 2 4 1 1 1 2
4 2 1 1 2 2
1 3 1 4 1 1
1 4 3 1 1 2 1 3 1 4
4 3 2 4 1 1 3 3
1 4 3 1 1 1 1 2 1 4
2 4 4 1 1 1 2 4 1 4 2 2
1 4 3 3 1 1 1 2 1 3
2 1 1 1 1 1
3 1 3 2 3 1 1 1 2 1
2 1 1 3 2 1
3 3 3 2 3 1 3 2 1 3
2 4 2 2 1 2 1 3
2 3 3 4 1 1 1 2 1 3
1 2 2 3 1 1 1 2
3 1 1 4 3 1
1 4 4 1 1 1 1 2 1 3 1 4
4 2 4 4 4 1 2 1 4 2 2 2
4 2 3 1 3 1 1 1 4 2
3 2 4 4 3 1 1 2 2 1 2 2
4 1 3 3 3 1 1 1 2 1
4 3 3 2 2 3 1 2 4 3
3 3 4 3 1 2 1 1 3 3 1 3
4 4 2 3 3 4 4 2
4 4 3 4 3 2 4 1 1 4
2 1 2 4 1 1 2 1
4 1 3 1 3 1 1 1 2 1
3 2 1 4 2 2
4 2 3 4 3 1 4 1 2 1
4 3 4 4 1 2 4 1 2 1 4 2
1 2 1 2 1 1
2 2 3 1 1 2 2 1 2 2
1 1 1 3 1 1
1 2 1 4 1 2
2 4 2 3 2 3 1 1
4 3 3 1 3 1 1 1 2 3
3 2 2 3 3 2 2 2
3 4 3 2 3 1 1 3 1 4
4 4 1 4 4 4
4 1 1 1 2 1
2 4 2 4 2 1 1 4
2 2 4 2 1 1 1 2 2 1 2 2
1 2 2 2 1 1 1 2
4 3 1 3 2 3
3 2 2 1 3 2 1 2
4 3 2 2 3 2 1 3
1 2 2 2 1 1 1 2
3 2 4 1 3 1 3 2 2 1 2 2
2 2 4 4 1 1 1 2 2 1 2 2
3 2 3 1 3 1 1 1 2 1
`

var rawTestcases = func() []string {
	lines := strings.Split(strings.TrimSpace(rawTestcasesData), "\n")
	out := make([]string, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		out = append(out, line)
	}
	return out
}()

// Dinic maxflow implementation (embedded from solution)
type edge struct{ to, rev, cap int }

type Dinic struct {
	N     int
	G     [][]edge
	level []int
	it    []int
}

func NewDinic(N int) *Dinic {
	return &Dinic{N: N, G: make([][]edge, N), level: make([]int, N), it: make([]int, N)}
}

func (d *Dinic) AddEdge(u, v, c int) {
	d.G[u] = append(d.G[u], edge{v, len(d.G[v]), c})
	d.G[v] = append(d.G[v], edge{u, len(d.G[u]) - 1, 0})
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *Dinic) bfs(s, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	queue := make([]int, 0, d.N)
	d.level[s] = 0
	queue = append(queue, s)
	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		for _, e := range d.G[u] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[u] + 1
				queue = append(queue, e.to)
			}
		}
	}
	return d.level[t] >= 0
}

func (d *Dinic) dfs(u, t, f int) int {
	if u == t {
		return f
	}
	for i := d.it[u]; i < len(d.G[u]); i++ {
		e := &d.G[u][i]
		if e.cap > 0 && d.level[u] < d.level[e.to] {
			ret := d.dfs(e.to, t, min(f, e.cap))
			if ret > 0 {
				e.cap -= ret
				d.G[e.to][e.rev].cap += ret
				return ret
			}
		}
		d.it[u]++
	}
	return 0
}

func (d *Dinic) MaxFlow(s, t int) int {
	flow := 0
	for d.bfs(s, t) {
		for i := range d.it {
			d.it[i] = 0
		}
		for {
			f := d.dfs(s, t, 1e9)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func roundMatrix(d []int, ksum []int, t int) [][]int {
	n := len(d)
	lower := make([][]int, n)
	upper := make([][]int, n)
	for i := 0; i < n; i++ {
		fl := d[i] / t
		ce := fl
		if d[i]%t != 0 {
			ce = fl + 1
		}
		lower[i] = make([]int, t)
		upper[i] = make([]int, t)
		for j := 0; j < t; j++ {
			lower[i][j] = fl
			upper[i][j] = ce
		}
	}
	U := n
	S := 0
	T := U + t + 1
	N := U + t + 2
	demand := make([]int, N)
	dinic := NewDinic(N + 2)
	SS := N
	TT := N + 1
	for i := 0; i < U; i++ {
		dinic.AddEdge(S, 1+i, d[i])
	}
	for j := 0; j < t; j++ {
		dinic.AddEdge(1+U+j, T, ksum[j])
	}
	for i := 0; i < U; i++ {
		for j := 0; j < t; j++ {
			l := lower[i][j]
			ucap := upper[i][j]
			dinic.AddEdge(1+i, 1+U+j, ucap-l)
			demand[1+i] -= l
			demand[1+U+j] += l
		}
	}
	dinic.AddEdge(T, S, 1e9)
	totalDemand := 0
	for v := 0; v < N; v++ {
		if demand[v] > 0 {
			dinic.AddEdge(SS, v, demand[v])
			totalDemand += demand[v]
		} else if demand[v] < 0 {
			dinic.AddEdge(v, TT, -demand[v])
		}
	}
	if dinic.MaxFlow(SS, TT) != totalDemand {
		panic("roundMatrix infeasible")
	}
	A := make([][]int, n)
	for i := range A {
		A[i] = make([]int, t)
	}
	for i := 0; i < U; i++ {
		for _, e := range dinic.G[1+i] {
			if e.to >= 1+U && e.to < 1+U+t {
				j := e.to - (1 + U)
				l := lower[i][j]
				used := upper[i][j] - l - e.cap
				A[i][j] = l + used
			}
		}
	}
	return A
}

func runReference(reader *bufio.Reader, writer *bufio.Writer) error {
	var n, m, k, t int
	if _, err := fmt.Fscan(reader, &n, &m, &k, &t); err != nil {
		return fmt.Errorf("failed to read header: %w", err)
	}
	xs := make([]int, k)
	ys := make([]int, k)
	degL := make([]int, n)
	degR := make([]int, m)
	for i := 0; i < k; i++ {
		if _, err := fmt.Fscan(reader, &xs[i], &ys[i]); err != nil {
			return fmt.Errorf("failed to read edge %d: %w", i, err)
		}
		xs[i]--
		ys[i]--
		degL[xs[i]]++
		degR[ys[i]]++
	}
	ksum := make([]int, t)
	base := k / t
	rem := k % t
	for j := 0; j < t; j++ {
		ksum[j] = base
		if j < rem {
			ksum[j]++
		}
	}
	A := roundMatrix(degL, ksum, t)
	B := roundMatrix(degR, ksum, t)
	adj := make([][]int, n)
	for i := 0; i < k; i++ {
		adj[xs[i]] = append(adj[xs[i]], i)
	}
	col := make([]int, k)
	used := make([]bool, k)
	for j := 0; j < t; j++ {
		N := 2 + n + m
		S := n + m
		T := n + m + 1
		din := NewDinic(N)
		type emap struct{ u, pos, idx int }
		emaps := make([]emap, 0, k)
		for u := 0; u < n; u++ {
			if A[u][j] > 0 {
				din.AddEdge(S, u, A[u][j])
			}
			for _, ei := range adj[u] {
				if used[ei] {
					continue
				}
				v := ys[ei]
				pos := len(din.G[u])
				din.AddEdge(u, n+v, 1)
				emaps = append(emaps, emap{u, pos, ei})
			}
		}
		for v := 0; v < m; v++ {
			if B[v][j] > 0 {
				din.AddEdge(n+v, T, B[v][j])
			}
		}
		din.MaxFlow(S, T)
		for _, em := range emaps {
			e := din.G[em.u][em.pos]
			if e.cap == 0 && !used[em.idx] {
				col[em.idx] = j + 1
				used[em.idx] = true
			}
		}
	}
	uneven := 0
	cntL := make([][]int, n)
	for i := 0; i < n; i++ {
		cntL[i] = make([]int, t)
	}
	cntR := make([][]int, m)
	for i := 0; i < m; i++ {
		cntR[i] = make([]int, t)
	}
	for i := 0; i < k; i++ {
		cj := col[i] - 1
		cntL[xs[i]][cj]++
		cntR[ys[i]][cj]++
	}
	for i := 0; i < n; i++ {
		mn, mx := cntL[i][0], cntL[i][0]
		for j := 1; j < t; j++ {
			if cntL[i][j] < mn {
				mn = cntL[i][j]
			}
			if cntL[i][j] > mx {
				mx = cntL[i][j]
			}
		}
		if mx-mn > uneven {
			uneven = mx - mn
		}
	}
	for i := 0; i < m; i++ {
		mn, mx := cntR[i][0], cntR[i][0]
		for j := 1; j < t; j++ {
			if cntR[i][j] < mn {
				mn = cntR[i][j]
			}
			if cntR[i][j] > mx {
				mx = cntR[i][j]
			}
		}
		if mx-mn > uneven {
			uneven = mx - mn
		}
	}
	fmt.Fprintln(writer, uneven)
	for i := 0; i < k; i++ {
		fmt.Fprint(writer, col[i], " ")
	}
	return writer.Flush()
}

func solveCase(line string) (string, error) {
	reader := bufio.NewReader(strings.NewReader(line))
	var sb strings.Builder
	writer := bufio.NewWriter(&sb)
	if err := runReference(reader, writer); err != nil {
		return "", err
	}
	return strings.TrimSpace(sb.String()), nil
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("runtime error: %v, output: %s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for idx, line := range rawTestcases {
		expected, err := solveCase(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d invalid: %v\n", idx+1, err)
			os.Exit(1)
		}
		got, err := run(bin, line+"\n")
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if got != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(rawTestcases))
}
