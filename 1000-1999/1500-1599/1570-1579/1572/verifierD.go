package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"math/bits"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type edge struct {
	to   int
	rev  int
	cap  int
	cost int
}

type mcmf struct {
	n     int
	graph [][]edge
	dist  []int
	prevv []int
	preve []int
}

func newMCMF(n int) *mcmf {
	g := make([][]edge, n)
	return &mcmf{n: n, graph: g, dist: make([]int, n), prevv: make([]int, n), preve: make([]int, n)}
}

func (f *mcmf) addEdge(u, v, cap, cost int) {
	f.graph[u] = append(f.graph[u], edge{to: v, rev: len(f.graph[v]), cap: cap, cost: cost})
	f.graph[v] = append(f.graph[v], edge{to: u, rev: len(f.graph[u]) - 1, cap: 0, cost: -cost})
}

const inf = int(1e18)

func (f *mcmf) minCostFlow(s, t, maxf int) int {
	res := 0
	flow := 0
	for flow < maxf {
		for i := 0; i < f.n; i++ {
			f.dist[i] = inf
		}
		inq := make([]bool, f.n)
		q := make([]int, 0)
		f.dist[s] = 0
		inq[s] = true
		q = append(q, s)
		for idx := 0; idx < len(q); idx++ {
			v := q[idx]
			inq[v] = false
			for i, e := range f.graph[v] {
				if e.cap > 0 && f.dist[e.to] > f.dist[v]+e.cost {
					f.dist[e.to] = f.dist[v] + e.cost
					f.prevv[e.to] = v
					f.preve[e.to] = i
					if !inq[e.to] {
						q = append(q, e.to)
						inq[e.to] = true
					}
				}
			}
		}
		if f.dist[t] == inf {
			break
		}
		d := maxf - flow
		for v := t; v != s; v = f.prevv[v] {
			if f.graph[f.prevv[v]][f.preve[v]].cap < d {
				d = f.graph[f.prevv[v]][f.preve[v]].cap
			}
		}
		for v := t; v != s; v = f.prevv[v] {
			e := &f.graph[f.prevv[v]][f.preve[v]]
			e.cap -= d
			rev := &f.graph[v][e.rev]
			rev.cap += d
		}
		res += d * f.dist[t]
		flow += d
	}
	return res
}

type testCaseD struct {
	n int
	k int
	a []int
}

func generateCaseD(rng *rand.Rand) testCaseD {
	n := rng.Intn(3) + 1 // n 1..3
	m := 1 << n
	k := rng.Intn(m/2) + 1
	a := make([]int, m)
	for i := range a {
		a[i] = rng.Intn(20)
	}
	return testCaseD{n: n, k: k, a: a}
}

func buildInputD(t testCaseD) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", t.n, t.k))
	for i, v := range t.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(fmt.Sprint(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func solveD(reader *bufio.Reader) string {
	var n, k int
	fmt.Fscan(reader, &n, &k)
	m := 1 << n
	a := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(reader, &a[i])
	}
	type node struct{ idx, val int }
	nodes := make([]node, m)
	for i := 0; i < m; i++ {
		nodes[i] = node{i, a[i]}
	}
	sort.Slice(nodes, func(i, j int) bool { return nodes[i].val > nodes[j].val })
	tsize := 4 * k
	if tsize > m {
		tsize = m
	}
	selected := make([]bool, m)
	for i := 0; i < tsize; i++ {
		selected[nodes[i].idx] = true
	}
	for i := 0; i < tsize; i++ {
		v := nodes[i].idx
		for b := 0; b < n; b++ {
			selected[v^(1<<b)] = true
		}
	}
	leftNodes := make([]int, 0)
	rightNodes := make([]int, 0)
	idxL := make(map[int]int)
	idxR := make(map[int]int)
	for i := 0; i < m; i++ {
		if selected[i] {
			if bits.OnesCount(uint(i))%2 == 0 {
				idxL[i] = len(leftNodes)
				leftNodes = append(leftNodes, i)
			} else {
				idxR[i] = len(rightNodes)
				rightNodes = append(rightNodes, i)
			}
		}
	}
	total := 1 + len(leftNodes) + len(rightNodes) + 1
	source := 0
	sink := total - 1
	offsetR := 1 + len(leftNodes)
	f := newMCMF(total)
	for i := 0; i < len(leftNodes); i++ {
		f.addEdge(source, 1+i, 1, 0)
	}
	for i := 0; i < len(rightNodes); i++ {
		f.addEdge(offsetR+i, sink, 1, 0)
	}
	for _, u := range leftNodes {
		lu := 1 + idxL[u]
		for b := 0; b < n; b++ {
			v := u ^ (1 << b)
			if ri, ok := idxR[v]; ok {
				w := -(a[u] + a[v])
				f.addEdge(lu, offsetR+ri, 1, w)
			}
		}
	}
	cost := f.minCostFlow(source, sink, k)
	return fmt.Sprint(-cost)
}

func expectedD(t testCaseD) string {
	input := buildInputD(t)
	return solveD(bufio.NewReader(strings.NewReader(input)))
}

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return out.String(), fmt.Errorf("timeout")
	}
	if err != nil {
		return out.String(), err
	}
	return out.String(), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tc := generateCaseD(rng)
		input := buildInputD(tc)
		expect := expectedD(tc)
		out, err := runBinary(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d: runtime error: %v\nOutput:%s", i+1, err, out)
			os.Exit(1)
		}
		got := strings.TrimSpace(out)
		exp := strings.TrimSpace(expect)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
