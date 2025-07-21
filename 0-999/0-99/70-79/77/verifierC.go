package main

import (
	"bytes"
	"container/heap"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
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
func (h *IntHeap) Peek() int { return (*h)[0] }

type State struct {
	m     map[int]int64
	heap  *IntHeap
	total int64
}

func (s *State) merge(o *State) {
	if o == nil || o.m == nil {
		return
	}
	if len(s.m) < len(o.m) {
		s.m, o.m = o.m, s.m
		s.heap, o.heap = o.heap, s.heap
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
			heap.Push(s.heap, d)
		}
		s.total += cnt
	}
	o.m = nil
	o.heap = nil
	o.total = 0
}

func solveC(n int, k []int64, edges [][2]int, sRoot int) int64 {
	adj := make([][]int, n+1)
	for _, e := range edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}
	depth := make([]int, n+1)
	var dfs func(u, p int) *State
	dfs = func(u, p int) *State {
		st := &State{m: make(map[int]int64), heap: &IntHeap{}, total: 0}
		heap.Init(st.heap)
		for _, v := range adj[u] {
			if v == p {
				continue
			}
			depth[v] = depth[u] + 1
			child := dfs(v, u)
			st.merge(child)
		}
		if u != sRoot {
			d := depth[u]
			cnt := k[u]
			if cnt > 0 {
				if _, ok := st.m[d]; ok {
					st.m[d] += cnt
				} else {
					st.m[d] = cnt
					heap.Push(st.heap, d)
				}
				st.total += cnt
			}
		}
		capc := k[u]
		for st.total > capc {
			d := st.heap.Peek()
			cnt := st.m[d]
			rem := st.total - capc
			if cnt <= rem {
				heap.Pop(st.heap)
				delete(st.m, d)
				st.total -= cnt
			} else {
				st.m[d] = cnt - rem
				st.total -= rem
				break
			}
		}
		return st
	}
	depth[sRoot] = 0
	st := dfs(sRoot, 0)
	var ans int64
	for d, cnt := range st.m {
		ans += int64(d) * cnt * 2
	}
	return ans
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
		exp := solveC(n, k, edges, root)
		out, err := run(bin, input.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", t, err)
			os.Exit(1)
		}
		var got int64
		fmt.Sscan(out, &got)
		if got != exp {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d got %s\n", t, exp, out)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", tests)
}
