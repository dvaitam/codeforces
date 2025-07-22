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
	n   int
	adj [][]Edge
	lvl []int
	it  []int
}

func NewDinic(n int) *Dinic {
	d := &Dinic{n: n, adj: make([][]Edge, n), lvl: make([]int, n), it: make([]int, n)}
	return d
}

func (d *Dinic) AddEdge(u, v, c int) {
	d.adj[u] = append(d.adj[u], Edge{to: v, rev: len(d.adj[v]), cap: c})
	d.adj[v] = append(d.adj[v], Edge{to: u, rev: len(d.adj[u]) - 1, cap: 0})
}

func (d *Dinic) bfs(s, t int) bool {
	for i := range d.lvl {
		d.lvl[i] = -1
	}
	queue := make([]int, 0, d.n)
	d.lvl[s] = 0
	queue = append(queue, s)
	for qi := 0; qi < len(queue); qi++ {
		u := queue[qi]
		for _, e := range d.adj[u] {
			if e.cap > 0 && d.lvl[e.to] < 0 {
				d.lvl[e.to] = d.lvl[u] + 1
				queue = append(queue, e.to)
				if e.to == t {
					return true
				}
			}
		}
	}
	return d.lvl[t] >= 0
}

func (d *Dinic) dfs(u, t, f int) int {
	if u == t {
		return f
	}
	for i := d.it[u]; i < len(d.adj[u]); i++ {
		d.it[u] = i
		e := &d.adj[u][i]
		if e.cap > 0 && d.lvl[e.to] == d.lvl[u]+1 {
			ret := d.dfs(e.to, t, min(f, e.cap))
			if ret > 0 {
				e.cap -= ret
				d.adj[e.to][e.rev].cap += ret
				return ret
			}
		}
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
			pushed := d.dfs(s, t, 1<<30)
			if pushed == 0 {
				break
			}
			flow += pushed
		}
	}
	return flow
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func factor(x int) map[int]int {
	m := make(map[int]int)
	d := 2
	for d*d <= x {
		for x%d == 0 {
			m[d]++
			x /= d
		}
		d++
	}
	if x > 1 {
		m[x]++
	}
	return m
}

func solveCase(n, m int, a []int, pairs [][2]int) string {
	cnt := make([]map[int]int, n+1)
	primes := make(map[int]struct{})
	for i := 1; i <= n; i++ {
		cnt[i] = factor(a[i])
		for p := range cnt[i] {
			primes[p] = struct{}{}
		}
	}
	total := 0
	const INF = 1000000000
	for p := range primes {
		oddIds := make(map[int]int)
		evenIds := make(map[int]int)
		oid, eid := 0, 0
		for i := 1; i <= n; i++ {
			if e, ok := cnt[i][p]; ok && e > 0 {
				if i%2 == 1 {
					oid++
					oddIds[i] = oid
				} else {
					eid++
					evenIds[i] = eid
				}
			}
		}
		if oid == 0 || eid == 0 {
			continue
		}
		N := 1 + oid + eid + 1
		src, sink := 0, N-1
		din := NewDinic(N)
		for i, id := range oddIds {
			din.AddEdge(src, id, cnt[i][p])
		}
		for i, id := range evenIds {
			din.AddEdge(oid+id, sink, cnt[i][p])
		}
		for _, pr := range pairs {
			u, v := pr[0], pr[1]
			ui, uok := oddIds[u]
			vi, vok := evenIds[v]
			if uok && vok {
				din.AddEdge(ui, oid+vi, INF)
			}
		}
		flow := din.MaxFlow(src, sink)
		total += flow
	}
	return fmt.Sprintf("%d", total)
}

func runBinary(bin, input string) (string, error) {
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
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateCase(rng *rand.Rand) (string, int, []int, [][2]int) {
	n := rng.Intn(6) + 2
	m := rng.Intn(n)
	a := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(20) + 1
	}
	pairs := make([][2]int, 0, m)
	used := make(map[[2]int]struct{})
	for len(pairs) < m {
		u := rng.Intn(n) + 1
		v := rng.Intn(n) + 1
		if u == v || (u+v)%2 != 1 {
			continue
		}
		if u%2 == 0 {
			u, v = v, u
		}
		pr := [2]int{u, v}
		if _, ok := used[pr]; ok {
			continue
		}
		used[pr] = struct{}{}
		pairs = append(pairs, pr)
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 1; i <= n; i++ {
		sb.WriteString(fmt.Sprintf("%d ", a[i]))
	}
	sb.WriteString("\n")
	for _, pr := range pairs {
		sb.WriteString(fmt.Sprintf("%d %d\n", pr[0], pr[1]))
	}
	return sb.String(), n, a, pairs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	inputs := []string{"2 1\n2 2\n1 2\n"}
	ns := []int{2}
	arrays := [][]int{{0, 2, 2}}
	pairLists := [][][2]int{{{1, 2}}}
	for i := 0; i < 100; i++ {
		in, n, arr, pairs := generateCase(rng)
		inputs = append(inputs, in)
		ns = append(ns, n)
		arrays = append(arrays, arr)
		pairLists = append(pairLists, pairs)
	}
	for i := range inputs {
		out, err := runBinary(bin, inputs[i])
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		exp := solveCase(ns[i], len(pairLists[i]), arrays[i], pairLists[i])
		if strings.TrimSpace(out) != exp {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected: %s\nfound: %s\n", i+1, inputs[i], exp, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
