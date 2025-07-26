package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

type edge struct {
	to   int
	rev  int
	cap  int
	cost int
}

type MCMF struct {
	n     int
	graph [][]edge
	dist  []int
	prevv []int
	preve []int
}

func NewMCMF(n int) *MCMF {
	g := make([][]edge, n)
	return &MCMF{n: n, graph: g, dist: make([]int, n), prevv: make([]int, n), preve: make([]int, n)}
}

func (f *MCMF) AddEdge(u, v, cap, cost int) {
	f.graph[u] = append(f.graph[u], edge{to: v, rev: len(f.graph[v]), cap: cap, cost: cost})
	f.graph[v] = append(f.graph[v], edge{to: u, rev: len(f.graph[u]) - 1, cap: 0, cost: -cost})
}

const INF = int(1e18)

func (f *MCMF) MinCostMaxFlow(s, t int) int {
	res := 0
	for {
		for i := 0; i < f.n; i++ {
			f.dist[i] = INF
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
		if f.dist[t] >= 0 || f.dist[t] == INF {
			break
		}
		d := 1
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
	}
	return -res
}

func solve(n int, s string, b []int) int {
	freq := make([]int, 26)
	for i := 0; i < n; i++ {
		freq[int(s[i]-'a')]++
	}
	letters := make([]int, 0)
	idx := make([]int, 26)
	for i := 0; i < 26; i++ {
		if freq[i] > 0 {
			idx[i] = len(letters)
			letters = append(letters, i)
		} else {
			idx[i] = -1
		}
	}
	L := len(letters)
	pairs := n / 2
	total := 1 + L + pairs*L + n + 1
	source := 0
	letterStart := 1
	pairLetterStart := letterStart + L
	posStart := pairLetterStart + pairs*L
	sink := total - 1

	flow := NewMCMF(total)
	for k, l := range letters {
		flow.AddEdge(source, letterStart+k, freq[l], 0)
	}
	for p := 0; p < pairs; p++ {
		i := p
		j := n - 1 - p
		for k, l := range letters {
			node := pairLetterStart + p*L + k
			flow.AddEdge(letterStart+k, node, 1, 0)
			if int(s[i]-'a') == l {
				flow.AddEdge(node, posStart+i, 1, -b[i])
			} else {
				flow.AddEdge(node, posStart+i, 1, 0)
			}
			if int(s[j]-'a') == l {
				flow.AddEdge(node, posStart+j, 1, -b[j])
			} else {
				flow.AddEdge(node, posStart+j, 1, 0)
			}
		}
	}
	for i := 0; i < n; i++ {
		flow.AddEdge(posStart+i, sink, 1, 0)
	}

	return flow.MinCostMaxFlow(source, sink)
}

func generateString(n int) (string, string) {
	t := make([]byte, n)
	for i := 0; i < n/2; i++ {
		a := byte('a' + rand.Intn(26))
		b := byte('a' + rand.Intn(26))
		for b == a {
			b = byte('a' + rand.Intn(26))
		}
		t[i] = a
		t[n-1-i] = b
	}
	perm := rand.Perm(n)
	s := make([]byte, n)
	for i := 0; i < n; i++ {
		s[i] = t[perm[i]]
	}
	return string(s), string(t)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(6)
	for tc := 0; tc < 100; tc++ {
		n := rand.Intn(10)*2 + 2 // even >=2
		s, _ := generateString(n)
		b := make([]int, n)
		for i := 0; i < n; i++ {
			b[i] = rand.Intn(100) + 1
		}
		var input bytes.Buffer
		fmt.Fprintf(&input, "%d\n", n)
		fmt.Fprintln(&input, s)
		for i := 0; i < n; i++ {
			fmt.Fprintf(&input, "%d ", b[i])
		}
		fmt.Fprintln(&input)
		expected := solve(n, s, b)
		cmd := exec.Command(bin)
		cmd.Stdin = bytes.NewReader(input.Bytes())
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println("error running binary:", err)
			os.Exit(1)
		}
		got := strings.TrimSpace(out.String())
		if got != fmt.Sprint(expected) {
			fmt.Println("wrong answer on test", tc+1)
			fmt.Println("input:\n" + input.String())
			fmt.Println("expected:", expected)
			fmt.Println("got:", got)
			os.Exit(1)
		}
	}
	fmt.Println("ok")
}
