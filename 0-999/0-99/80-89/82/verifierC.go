package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strings"
	"time"
)

type Div struct {
	ai int
	t  int
	id int
}

type DivHeap []Div

func (h DivHeap) Len() int            { return len(h) }
func (h DivHeap) Less(i, j int) bool  { return h[i].ai < h[j].ai }
func (h DivHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *DivHeap) Push(x interface{}) { *h = append(*h, x.(Div)) }
func (h *DivHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func schedule(divs []Div, cap int) []Div {
	sort.Slice(divs, func(i, j int) bool { return divs[i].t < divs[j].t })
	var pq DivHeap
	heap.Init(&pq)
	res := make([]Div, 0, len(divs))
	idx := 0
	day := 0
	for len(res) < len(divs) {
		if pq.Len() == 0 && idx < len(divs) && divs[idx].t > day {
			day = divs[idx].t
		}
		for idx < len(divs) && divs[idx].t <= day {
			heap.Push(&pq, divs[idx])
			idx++
		}
		for k := 0; k < cap && pq.Len() > 0; k++ {
			d := heap.Pop(&pq).(Div)
			d.t = day + 1
			res = append(res, d)
		}
		day++
	}
	return res
}

func solveCaseC(n int, a []int, edges [][3]int) []int {
	adj := make([][]struct{ to, cap int }, n+1)
	for _, e := range edges {
		u, v, c := e[0], e[1], e[2]
		adj[u] = append(adj[u], struct{ to, cap int }{v, c})
		adj[v] = append(adj[v], struct{ to, cap int }{u, c})
	}
	children := make([][]int, n+1)
	capToParent := make([]int, n+1)
	parent := make([]int, n+1)
	q := []int{1}
	for len(q) > 0 {
		u := q[0]
		q = q[1:]
		for _, e := range adj[u] {
			if e.to == parent[u] {
				continue
			}
			parent[e.to] = u
			capToParent[e.to] = e.cap
			children[u] = append(children[u], e.to)
			q = append(q, e.to)
		}
	}
	var dfs func(int) []Div
	dfs = func(u int) []Div {
		divs := []Div{{ai: a[u], t: 0, id: u}}
		for _, v := range children[u] {
			sub := dfs(v)
			sched := schedule(sub, capToParent[v])
			divs = append(divs, sched...)
		}
		return divs
	}
	rootDivs := dfs(1)
	ans := make([]int, n+1)
	for _, d := range rootDivs {
		ans[d.id] = d.t
	}
	return ans[1:]
}

func generateCase(rng *rand.Rand) (string, []int) {
	n := rng.Intn(5) + 1 // 1..5
	a := rng.Perm(n)
	for i := range a {
		a[i] = a[i] + 1 + rng.Intn(5)
	}
	edges := make([][3]int, 0, n-1)
	for v := 2; v <= n; v++ {
		u := rng.Intn(v-1) + 1
		c := rng.Intn(n) + 1
		edges = append(edges, [3]int{u, v, c})
	}
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", n)
	for i, val := range a {
		if i > 0 {
			fmt.Fprint(&b, " ")
		}
		fmt.Fprint(&b, val)
	}
	fmt.Fprintln(&b)
	for _, e := range edges {
		fmt.Fprintf(&b, "%d %d %d\n", e[0], e[1], e[2])
	}
	exp := solveCaseC(n, append([]int{0}, a...), edges)
	return b.String(), exp
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	var errBuf bytes.Buffer
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func parseOutput(out string, n int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != n {
		return nil, fmt.Errorf("expected %d numbers", n)
	}
	res := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Sscan(fields[i], &res[i])
	}
	return res, nil
}

func equalInts(a, b []int) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 20; i++ {
		input, expected := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d runtime error: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := parseOutput(out, len(expected))
		if err != nil {
			fmt.Printf("case %d invalid output: %v\n", i+1, err)
			os.Exit(1)
		}
		if !equalInts(expected, got) {
			fmt.Printf("case %d failed\ninput:\n%sexpected:%v\ngot:%v\n", i+1, input, expected, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
