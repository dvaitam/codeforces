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

type edge struct {
	to, rev int
	cap     int64
	cost    int64
}

type mcmf struct {
	n     int
	graph [][]edge
	h     []int64
	dist  []int64
	prevV []int
	prevE []int
}

func newMCMF(n int) *mcmf {
	return &mcmf{n: n, graph: make([][]edge, n), h: make([]int64, n), dist: make([]int64, n), prevV: make([]int, n), prevE: make([]int, n)}
}

func (m *mcmf) addEdge(from, to int, cap, cost int64) {
	m.graph[from] = append(m.graph[from], edge{to, len(m.graph[to]), cap, cost})
	m.graph[to] = append(m.graph[to], edge{from, len(m.graph[from]) - 1, 0, -cost})
}

func (m *mcmf) minCostFlow(s, t int, maxf int64) (int64, int64) {
	const inf = int64(4e18)
	var flow, cost int64
	for i := 0; i < m.n; i++ {
		m.h[i] = 0
	}
	for flow < maxf {
		for i := 0; i < m.n; i++ {
			m.dist[i] = inf
		}
		m.dist[s] = 0
		type item struct {
			v    int
			dist int64
		}
		pq := []item{{s, 0}}
		for len(pq) > 0 {
			it := pq[0]
			pq = pq[1:]
			v := it.v
			if it.dist != m.dist[v] {
				continue
			}
			for ei, e := range m.graph[v] {
				if e.cap > 0 {
					nd := it.dist + e.cost + m.h[v] - m.h[e.to]
					if nd < m.dist[e.to] {
						m.dist[e.to] = nd
						m.prevV[e.to] = v
						m.prevE[e.to] = ei
						pq = append(pq, item{e.to, nd})
					}
				}
			}
		}
		if m.dist[t] == inf {
			break
		}
		for v := 0; v < m.n; v++ {
			if m.dist[v] < inf {
				m.h[v] += m.dist[v]
			}
		}
		d := maxf - flow
		for v := t; v != s; v = m.prevV[v] {
			e := m.graph[m.prevV[v]][m.prevE[v]]
			if e.cap < d {
				d = e.cap
			}
		}
		flow += d
		cost += d * m.h[t]
		for v := t; v != s; v = m.prevV[v] {
			e := &m.graph[m.prevV[v]][m.prevE[v]]
			e.cap -= d
			m.graph[v][e.rev].cap += d
		}
	}
	return flow, cost
}

func solveCase(n, m int, u, v []int, c, f []int64) int64 {
	res := int64(0)
	for i := 0; i < m; i++ {
		if f[i] > c[i] {
			res += f[i] - c[i]
			c[i] = f[i]
		}
	}
	b := make([]int64, n)
	for i := 0; i < m; i++ {
		b[u[i]] -= f[i]
		b[v[i]] += f[i]
	}
	var totalDemand int64
	for i := 1; i < n-1; i++ {
		if b[i] > 0 {
			totalDemand += b[i]
		}
	}
	if totalDemand > 0 {
		SS := n
		TT := n + 1
		mvc := newMCMF(n + 2)
		for i := 0; i < m; i++ {
			ui, vi := u[i], v[i]
			if f[i] > 0 {
				mvc.addEdge(vi, ui, f[i], 1)
			}
			capInc := c[i] - f[i]
			if capInc > 0 {
				mvc.addEdge(ui, vi, capInc, 1)
			}
			mvc.addEdge(ui, vi, totalDemand, 2)
		}
		for i := 1; i < n-1; i++ {
			if b[i] > 0 {
				mvc.addEdge(i, TT, b[i], 0)
			} else if b[i] < 0 {
				mvc.addEdge(SS, i, -b[i], 0)
			}
		}
		_, cost := mvc.minCostFlow(SS, TT, totalDemand)
		res += cost
	}
	return res
}

func runCase(bin string, n, m int, u, v []int, c, f []int64) error {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < m; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", u[i]+1, v[i]+1, c[i], f[i]))
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(sb.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expect := fmt.Sprintf("%d", solveCase(n, m, u, v, c, f))
	if got != expect {
		return fmt.Errorf("expected %s got %s", expect, got)
	}
	return nil
}

func randomGraph(rng *rand.Rand) (int, int, []int, []int, []int64, []int64) {
	n := rng.Intn(5) + 2
	m := rng.Intn(6)
	u := make([]int, m)
	v := make([]int, m)
	c := make([]int64, m)
	f := make([]int64, m)
	for i := 0; i < m; i++ {
		for {
			a := rng.Intn(n)
			b := rng.Intn(n)
			if a != b && b != 0 && a != n-1 {
				u[i] = a
				v[i] = b
				break
			}
		}
		c[i] = int64(rng.Intn(5))
		f[i] = int64(rng.Intn(5))
	}
	return n, m, u, v, c, f
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		n, m, u, v, c, f := randomGraph(rng)
		if err := runCase(bin, n, m, u, v, c, f); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
