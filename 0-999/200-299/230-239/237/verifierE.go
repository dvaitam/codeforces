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

type Edge struct {
	to, rev, cap int
	cost         int
}

type Graph [][]*Edge

func (g Graph) AddEdge(u, v, cap, cost int) {
	g[u] = append(g[u], &Edge{to: v, rev: len(g[v]), cap: cap, cost: cost})
	g[v] = append(g[v], &Edge{to: u, rev: len(g[u]) - 1, cap: 0, cost: -cost})
}

type Item struct {
	v    int
	dist int
}

type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	it := old[n-1]
	*pq = old[0 : n-1]
	return it
}

func minCostFlow(g Graph, s, t, maxf int) (int, int) {
	n := len(g)
	h := make([]int, n)
	prevv := make([]int, n)
	preve := make([]int, n)
	flow, cost := 0, 0
	const INF = int(1 << 60)
	for flow < maxf {
		dist := make([]int, n)
		for i := range dist {
			dist[i] = INF
		}
		dist[s] = 0
		pq := &PriorityQueue{}
		heap.Init(pq)
		heap.Push(pq, Item{v: s, dist: 0})
		for pq.Len() > 0 {
			it := heap.Pop(pq).(Item)
			v := it.v
			if dist[v] < it.dist {
				continue
			}
			for i, e := range g[v] {
				if e.cap > 0 && dist[e.to] > dist[v]+e.cost+h[v]-h[e.to] {
					dist[e.to] = dist[v] + e.cost + h[v] - h[e.to]
					prevv[e.to] = v
					preve[e.to] = i
					heap.Push(pq, Item{v: e.to, dist: dist[e.to]})
				}
			}
		}
		if dist[t] == INF {
			break
		}
		for v := 0; v < n; v++ {
			if dist[v] < INF {
				h[v] += dist[v]
			}
		}
		d := maxf - flow
		for v := t; v != s; v = prevv[v] {
			e := g[prevv[v]][preve[v]]
			if d > e.cap {
				d = e.cap
			}
		}
		flow += d
		cost += d * h[t]
		for v := t; v != s; v = prevv[v] {
			e := g[prevv[v]][preve[v]]
			e.cap -= d
			g[v][e.rev].cap += d
		}
	}
	return flow, cost
}

type strCase struct {
	s string
	a int
}

type testCase struct {
	input    string
	expected int
}

func solve(t string, arr []strCase) int {
	m := len(t)
	freqT := make([]int, 26)
	for i := 0; i < m; i++ {
		freqT[t[i]-'a']++
	}
	s := 0
	charBase := 1
	strBase := charBase + 26
	sink := strBase + len(arr)
	V := sink + 1
	g := make(Graph, V)
	for c := 0; c < 26; c++ {
		if freqT[c] > 0 {
			g.AddEdge(s, charBase+c, freqT[c], 0)
		}
	}
	for i, item := range arr {
		freqS := make([]int, 26)
		for j := 0; j < len(item.s); j++ {
			freqS[item.s[j]-'a']++
		}
		for c := 0; c < 26; c++ {
			if freqS[c] > 0 {
				g.AddEdge(charBase+c, strBase+i, freqS[c], i+1)
			}
		}
		if item.a > 0 {
			g.AddEdge(strBase+i, sink, item.a, 0)
		}
	}
	if m == 0 {
		return 0
	}
	flow, cost := minCostFlow(g, s, sink, m)
	if flow < m {
		return -1
	}
	return cost
}

func generateCase(rng *rand.Rand) testCase {
	tlen := rng.Intn(5) + 1
	letters := []rune("abcdefghijklmnopqrstuvwxyz")
	tb := make([]rune, tlen)
	for i := range tb {
		tb[i] = letters[rng.Intn(26)]
	}
	t := string(tb)
	n := rng.Intn(4) + 1
	arr := make([]strCase, n)
	var sb strings.Builder
	sb.WriteString(t)
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		slen := rng.Intn(5) + 1
		sbuf := make([]rune, slen)
		for j := range sbuf {
			sbuf[j] = letters[rng.Intn(26)]
		}
		arr[i].s = string(sbuf)
		arr[i].a = rng.Intn(slen + 1)
		sb.WriteString(fmt.Sprintf("%s %d\n", arr[i].s, arr[i].a))
	}
	expect := solve(t, arr)
	return testCase{input: sb.String(), expected: expect}
}

func runCase(bin string, tc testCase) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(tc.input)
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
	if got != tc.expected {
		return fmt.Errorf("expected %d got %d", tc.expected, got)
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
	var cases []testCase
	cases = append(cases, generateCase(rand.New(rand.NewSource(1))))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}
	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
