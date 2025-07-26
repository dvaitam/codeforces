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

type IntHeap []int

func (h IntHeap) Len() int            { return len(h) }
func (h IntHeap) Less(i, j int) bool  { return h[i] < h[j] }
func (h IntHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *IntHeap) Push(x interface{}) { *h = append(*h, x.(int)) }
func (h *IntHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}

func expected(n, m, k int, a, b, c []int, edges [][2]int) string {
	last := make([]int, n+1)
	for i := 1; i <= n; i++ {
		last[i] = i
	}
	for _, e := range edges {
		u := e[0]
		v := e[1]
		if u > v {
			if u > last[v] {
				last[v] = u
			}
		}
	}
	base := make([]int, n+1)
	base[0] = k
	for i := 1; i <= n; i++ {
		if base[i-1] < a[i] {
			return "-1"
		}
		base[i] = base[i-1] + b[i]
	}
	cap := make([]int, n+1)
	cap[n] = base[n]
	for i := 1; i < n; i++ {
		cap[i] = base[i] - a[i+1]
		if cap[i] < 0 {
			return "-1"
		}
	}
	buckets := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		d := last[i]
		buckets[d] = append(buckets[d], c[i])
	}
	h := &IntHeap{}
	heap.Init(h)
	total := 0
	size := 0
	for t := 1; t <= n; t++ {
		for _, val := range buckets[t] {
			heap.Push(h, val)
			total += val
			size++
		}
		for size > cap[t] {
			removed := heap.Pop(h).(int)
			total -= removed
			size--
		}
	}
	return fmt.Sprintf("%d", total)
}

func generateCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(10) + 1
	maxM := n * (n - 1) / 2
	m := rng.Intn(maxM + 1)
	k := rng.Intn(20)
	a := make([]int, n+1)
	b := make([]int, n+1)
	c := make([]int, n+1)
	for i := 1; i <= n; i++ {
		a[i] = rng.Intn(10)
		b[i] = rng.Intn(10)
		c[i] = rng.Intn(10)
	}
	edges := make([][2]int, m)
	seen := map[[2]int]bool{}
	for i := 0; i < m; {
		u := rng.Intn(n) + 1
		v := rng.Intn(u-1) + 1
		e := [2]int{u, v}
		if !seen[e] {
			seen[e] = true
			edges[i] = e
			i++
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
	for i := 1; i <= n; i++ {
		fmt.Fprintf(&sb, "%d %d %d\n", a[i], b[i], c[i])
	}
	for _, e := range edges {
		fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
	}
	exp := expected(n, m, k, a, b, c, edges)
	return sb.String(), exp
}

func runCase(exe, input, expected string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	expected = strings.TrimSpace(expected)
	if got != expected {
		return fmt.Errorf("expected %q got %q", expected, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, expect := generateCase(rng)
		if err := runCase(exe, input, expect); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
