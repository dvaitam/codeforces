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

type DSU struct {
	p, sz []int
}

func NewDSU(n int) *DSU {
	p := make([]int, n+1)
	sz := make([]int, n+1)
	for i := 1; i <= n; i++ {
		p[i] = i
		sz[i] = 1
	}
	return &DSU{p: p, sz: sz}
}

func (d *DSU) Find(x int) int {
	for d.p[x] != x {
		d.p[x] = d.p[d.p[x]]
		x = d.p[x]
	}
	return x
}

func (d *DSU) Union(a, b int) {
	ra := d.Find(a)
	rb := d.Find(b)
	if ra == rb {
		return
	}
	if d.sz[ra] < d.sz[rb] {
		ra, rb = rb, ra
	}
	d.p[rb] = ra
	d.sz[ra] += d.sz[rb]
}

type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func solve(n, m, k int, edges [][2]int) int {
	dsu := NewDSU(n)
	for _, e := range edges {
		dsu.Union(e[0], e[1])
	}
	compCap := make(map[int]int)
	for i := 1; i <= n; i++ {
		r := dsu.Find(i)
		compCap[r]++
	}
	h := &IntHeap{}
	heap.Init(h)
	sumCaps := 0
	D := 0
	for _, sz := range compCap {
		cap := sz
		if cap > k {
			cap = k
		}
		heap.Push(h, cap)
		sumCaps += cap
		D++
	}
	if D <= 1 {
		return 0
	}
	roads := 0
	for sumCaps < 2*(D-1) {
		c1 := heap.Pop(h).(int)
		c2 := heap.Pop(h).(int)
		newC := c1 + c2
		loss := 0
		if newC > k {
			loss = newC - k
			newC = k
		}
		sumCaps -= loss
		D--
		roads++
		heap.Push(h, newC)
	}
	return roads
}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n := rng.Intn(5) + 1
		maxEdges := n * (n - 1) / 2
		m := rng.Intn(maxEdges + 1)
		edges := make([][2]int, 0, m)
		used := make(map[[2]int]bool)
		for len(edges) < m {
			u := rng.Intn(n) + 1
			v := rng.Intn(n) + 1
			if u == v {
				continue
			}
			e := [2]int{u, v}
			if u > v {
				e = [2]int{v, u}
			}
			if used[e] {
				continue
			}
			used[e] = true
			edges = append(edges, e)
		}
		k := rng.Intn(5) + 1
		expected := fmt.Sprintf("%d", solve(n, m, k, edges))
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", n, m, k))
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
		}
		input := sb.String()
		got, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, expected, got, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
