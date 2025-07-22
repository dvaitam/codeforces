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

type testCase struct {
	input  string
	expect int
}

const INF = 1000000000

type Edge struct{ to, rev, cap int }

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func solve(sets [][]int, costs []int) int {
	n := len(sets)
	// bipartite matching
	p := make([]int, n)
	q := make([]int, n)
	for i := 0; i < n; i++ {
		p[i], q[i] = -1, -1
	}
	was := make([]bool, n)
	var dfs func(int) bool
	dfs = func(v int) bool {
		was[v] = true
		for _, u := range sets[v] {
			if q[u] == -1 {
				q[u] = v
				p[v] = u
				return true
			}
		}
		for _, u := range sets[v] {
			w := q[u]
			if !was[w] && dfs(w) {
				q[u] = v
				p[v] = u
				return true
			}
		}
		return false
	}
	now := 0
	for now < n {
		for i := range was {
			was[i] = false
		}
		for i := 0; i < n; i++ {
			if p[i] == -1 && !was[i] {
				if dfs(i) {
					now++
				}
			}
		}
	}
	N := n + 2
	S, T := n, n+1
	G := make([][]Edge, N)
	addEdge := func(u, v, c int) {
		G[u] = append(G[u], Edge{v, len(G[v]), c})
		G[v] = append(G[v], Edge{u, len(G[u]) - 1, 0})
	}
	sumNeg := 0
	for i := 0; i < n; i++ {
		if costs[i] >= 0 {
			addEdge(i, T, costs[i])
		} else {
			sumNeg += costs[i]
			addEdge(S, i, -costs[i])
		}
	}
	for i := 0; i < n; i++ {
		for _, j := range sets[i] {
			if q[j] != i {
				addEdge(i, q[j], INF)
			}
		}
	}
	level := make([]int, N)
	iter := make([]int, N)
	bfs := func(mn int) bool {
		for i := range level {
			level[i] = -1
		}
		q := []int{S}
		level[S] = 0
		for qi := 0; qi < len(q); qi++ {
			v := q[qi]
			for _, e := range G[v] {
				if e.cap >= mn && level[e.to] < 0 {
					level[e.to] = level[v] + 1
					q = append(q, e.to)
				}
			}
		}
		return level[T] >= 0
	}
	var dfsf func(int, int, int) int
	dfsf = func(v, up, mn int) int {
		if v == T {
			return up
		}
		res := 0
		for ; up >= mn && iter[v] < len(G[v]); iter[v]++ {
			e := &G[v][iter[v]]
			if e.cap < mn || level[e.to] != level[v]+1 {
				continue
			}
			f := dfsf(e.to, min(up, e.cap), mn)
			if f > 0 {
				e.cap -= f
				G[e.to][e.rev].cap += f
				up -= f
				res += f
			}
		}
		return res
	}
	flow := 0
	for mn := 1 << 20; mn > 0; mn >>= 1 {
		for bfs(mn) {
			for i := range iter {
				iter[i] = 0
			}
			for {
				f := dfsf(S, INF, mn)
				if f == 0 {
					break
				}
				flow += f
			}
		}
	}
	return flow + sumNeg
}

func buildCase(sets [][]int, costs []int) testCase {
	var sb strings.Builder
	n := len(sets)
	fmt.Fprintf(&sb, "%d\n", n)
	for i := 0; i < n; i++ {
		fmt.Fprintf(&sb, "%d ", len(sets[i]))
		for _, x := range sets[i] {
			fmt.Fprintf(&sb, "%d ", x+1)
		}
		sb.WriteByte('\n')
	}
	for i, c := range costs {
		if i > 0 {
			sb.WriteByte(' ')
		}
		fmt.Fprintf(&sb, "%d", c)
	}
	sb.WriteByte('\n')
	return testCase{input: sb.String(), expect: solve(sets, costs)}
}

func generateRandomCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 1
	sets := make([][]int, n)
	for i := 0; i < n; i++ {
		m := rng.Intn(n) + 1
		used := map[int]bool{i: true}
		sets[i] = []int{i}
		for len(sets[i]) < m {
			v := rng.Intn(n)
			if !used[v] {
				used[v] = true
				sets[i] = append(sets[i], v)
			}
		}
	}
	costs := make([]int, n)
	for i := range costs {
		costs[i] = rng.Intn(201) - 100
	}
	return buildCase(sets, costs)
}

func runCase(bin string, tc testCase) error {
	cmd := exec.Command(bin)
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
	if got != tc.expect {
		return fmt.Errorf("expected %d got %d", tc.expect, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	sets := [][]int{{0}, {1}}
	costs := []int{5, -3}
	cases = append(cases, buildCase(sets, costs))

	for i := 0; i < 100; i++ {
		cases = append(cases, generateRandomCase(rng))
	}

	for i, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
