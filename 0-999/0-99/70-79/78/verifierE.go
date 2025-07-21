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

type Edge struct {
	to, rev, cap int
}

type Dinic struct {
	n     int
	adj   [][]Edge
	level []int
	prog  []int
}

func NewDinic(n int) *Dinic {
	d := &Dinic{n: n, adj: make([][]Edge, n), level: make([]int, n), prog: make([]int, n)}
	return d
}

func (d *Dinic) AddEdge(u, v, c int) {
	d.adj[u] = append(d.adj[u], Edge{to: v, rev: len(d.adj[v]), cap: c})
	d.adj[v] = append(d.adj[v], Edge{to: u, rev: len(d.adj[u]) - 1, cap: 0})
}

func (d *Dinic) bfs(s, t int) bool {
	for i := range d.level {
		d.level[i] = -1
	}
	queue := make([]int, 0, d.n)
	d.level[s] = 0
	queue = append(queue, s)
	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		for _, e := range d.adj[u] {
			if e.cap > 0 && d.level[e.to] < 0 {
				d.level[e.to] = d.level[u] + 1
				queue = append(queue, e.to)
				if e.to == t {
					return true
				}
			}
		}
	}
	return d.level[t] >= 0
}

func (d *Dinic) dfs(u, t, f int) int {
	if u == t {
		return f
	}
	for i := d.prog[u]; i < len(d.adj[u]); i++ {
		e := &d.adj[u][i]
		if e.cap > 0 && d.level[e.to] == d.level[u]+1 {
			minf := d.dfs(e.to, t, min(f, e.cap))
			if minf > 0 {
				e.cap -= minf
				d.adj[e.to][e.rev].cap += minf
				return minf
			}
		}
		d.prog[u]++
	}
	return 0
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func (d *Dinic) MaxFlow(s, t int) int {
	flow := 0
	for d.bfs(s, t) {
		for i := range d.prog {
			d.prog[i] = 0
		}
		for {
			f := d.dfs(s, t, 1<<30)
			if f == 0 {
				break
			}
			flow += f
		}
	}
	return flow
}

func solve(n, t int, gridS, gridC []string) int {
	const INF = int(1e9)
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, n)
		for j := range dist[i] {
			dist[i][j] = INF
		}
	}
	var zi, zj int
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if gridS[i][j] == 'Z' {
				zi, zj = i, j
			}
		}
	}
	type P struct{ i, j int }
	q := make([]P, 0, n*n)
	dist[zi][zj] = 0
	q = append(q, P{zi, zj})
	dirs := []P{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for qi := 0; qi < len(q); qi++ {
		u := q[qi]
		for _, ddir := range dirs {
			ni, nj := u.i+ddir.i, u.j+ddir.j
			if ni >= 0 && ni < n && nj >= 0 && nj < n && dist[ni][nj] == INF {
				if gridS[ni][nj] >= '0' && gridS[ni][nj] <= '9' {
					dist[ni][nj] = dist[u.i][u.j] + 1
					q = append(q, P{ni, nj})
				}
			}
		}
	}
	timeLayers := t + 1
	totalLabs := n * n
	baseTime := 0
	baseCap := baseTime + totalLabs*timeLayers
	source := baseCap + totalLabs
	sink := source + 1
	V := sink + 1
	dinic := NewDinic(V)
	infCap := totalLabs*10 + 5
	timeNode := func(i, j, k int) int {
		return baseTime + (i*n+j)*timeLayers + k
	}
	capNode := func(i, j int) int {
		return baseCap + (i*n + j)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if gridS[i][j] >= '0' && gridS[i][j] <= '9' {
				sCnt := int(gridS[i][j] - '0')
				if sCnt > 0 && dist[i][j] > 0 {
					dinic.AddEdge(source, timeNode(i, j, 0), sCnt)
				}
			}
		}
	}
	dirs = []P{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if gridS[i][j] >= '0' && gridS[i][j] <= '9' {
				capCnt := 0
				if gridC[i][j] >= '0' && gridC[i][j] <= '9' {
					capCnt = int(gridC[i][j] - '0')
				}
				if capCnt > 0 {
					dinic.AddEdge(capNode(i, j), sink, capCnt)
				}
				for k := 0; k <= t; k++ {
					if dist[i][j] > k {
						if capCnt > 0 {
							dinic.AddEdge(timeNode(i, j, k), capNode(i, j), infCap)
						}
						if k < t {
							for _, ddir := range dirs {
								ni, nj := i+ddir.i, j+ddir.j
								if ni >= 0 && ni < n && nj >= 0 && nj < n && dist[ni][nj] > k+1 && gridS[ni][nj] >= '0' && gridS[ni][nj] <= '9' {
									dinic.AddEdge(timeNode(i, j, k), timeNode(ni, nj, k+1), infCap)
								}
							}
						}
					}
				}
			}
		}
	}
	result := dinic.MaxFlow(source, sink)
	return result
}

func generateCase(rng *rand.Rand) (string, int) {
	n := rng.Intn(3) + 2
	tVal := rng.Intn(4) + 1
	gridS := make([]string, n)
	gridC := make([]string, n)
	zi := rng.Intn(n)
	zj := rng.Intn(n)
	for i := 0; i < n; i++ {
		var sbS, sbC strings.Builder
		for j := 0; j < n; j++ {
			if i == zi && j == zj {
				sbS.WriteByte('Z')
				sbC.WriteByte('Z')
				continue
			}
			if rng.Intn(5) == 0 {
				sbS.WriteByte('Y')
				sbC.WriteByte('Y')
			} else {
				sbS.WriteByte(byte('0' + rng.Intn(3)))
				sbC.WriteByte(byte('0' + rng.Intn(3)))
			}
		}
		gridS[i] = sbS.String()
		gridC[i] = sbC.String()
	}
	var input strings.Builder
	input.WriteString(fmt.Sprintf("%d %d\n", n, tVal))
	for i := 0; i < n; i++ {
		input.WriteString(gridS[i])
		input.WriteByte('\n')
	}
	input.WriteByte('\n')
	for i := 0; i < n; i++ {
		input.WriteString(gridC[i])
		input.WriteByte('\n')
	}
	expect := solve(n, tVal, gridS, gridC)
	return input.String(), expect
}

func runCase(bin, input string, expected int) error {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	var got int
	if _, err := fmt.Fscan(strings.NewReader(out.String()), &got); err != nil {
		return fmt.Errorf("bad output: %v", err)
	}
	if got != expected {
		return fmt.Errorf("expected %d got %d", expected, got)
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
		in, exp := generateCase(rng)
		if err := runCase(bin, in, exp); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
