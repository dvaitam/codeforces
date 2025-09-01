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
)

type maxHeap77 []int

func (h maxHeap77) Len() int           { return len(h) }
func (h maxHeap77) Less(i, j int) bool { return h[i] < h[j] } // min-heap by depth; drop nearest first
func (h maxHeap77) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxHeap77) Push(x interface{}) {
	*h = append(*h, x.(int))
}
func (h *maxHeap77) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[:n-1]
	return x
}
func (h *maxHeap77) Peek() int { return (*h)[0] }

type state77 struct {
	m     map[int]int64 // depth -> count
	h     *maxHeap77    // depths present in m (max-heap)
	total int64         // total count
}

func (s *state77) merge(o *state77) {
	if o == nil || o.m == nil {
		return
	}
	// small-to-large
	if len(s.m) < len(o.m) {
		s.m, o.m = o.m, s.m
		s.h, o.h = o.h, s.h
		s.total, o.total = o.total, s.total
	}
	for d, cnt := range o.m {
		if cnt == 0 {
			continue
		}
		if _, ok := s.m[d]; ok {
			s.m[d] += cnt
		} else {
			s.m[d] = cnt
			heap.Push(s.h, d)
		}
		s.total += cnt
	}
	o.m = nil
	o.h = nil
	o.total = 0
}

// Oracle computing the objective using the accepted greedy from C++ reference
func solve77COracle(n int, k0 []int64, edges [][2]int, root int) int64 {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	k := make([]int64, n+1)
	copy(k, k0)
	var dfs func(u, p int) int64
	dfs = func(u, p int) int64 {
		r := make([]int64, 0)
		var fb int64 = 0
		var res int64 = 0
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			k[v]--
			val := -dfs(v, u)
			r = append(r, val)
			fb += k[v]
		}
		sort.Slice(r, func(i, j int) bool { return r[i] < r[j] })
		for _, pv := range r {
			if k[u] == 0 {
				break
			}
			res -= pv - 2
			k[u]--
		}
		if fb > k[u] {
			fb = k[u]
		}
		k[u] -= fb
		return res + fb*2
	}
	return dfs(root, 0)
}

func run77C(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errBuf.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rand.Seed(42)
	tests := 100
	for t := 1; t <= tests; t++ {
		n := rand.Intn(5) + 3 // 3..7
		k := make([]int64, n+1)
		for i := 1; i <= n; i++ {
			k[i] = int64(rand.Intn(5))
		}
		edges := make([][2]int, n-1)
		for i := 2; i <= n; i++ {
			p := rand.Intn(i-1) + 1
			edges[i-2] = [2]int{i, p}
		}
		root := rand.Intn(n) + 1
		var input bytes.Buffer
		fmt.Fprintln(&input, n)
		for i := 1; i <= n; i++ {
			fmt.Fprintf(&input, "%d ", k[i])
		}
		input.WriteByte('\n')
		for _, e := range edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
		fmt.Fprintln(&input, root)
		exp := solve77COracle(n, k, edges, root)
		out, err := run77C(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\nInput:\n%s", t, err, input.String())
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(out, &got)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\nInput:\n%s", t, exp, out, input.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
