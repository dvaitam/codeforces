package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

type TestCase struct {
	n     int
	edges [][2]int
}

type pair struct {
	d int
	v int
}

var (
	n  int
	K  int
	id int
	g  [][]int
)

func dfs(x, fa int) pair {
	if id != -1 {
		return pair{-n, -1}
	}
	st := make([]pair, 0)
	for _, u := range g[x] {
		if u == fa {
			continue
		}
		st = append(st, dfs(u, x))
	}
	if id != -1 {
		return pair{-n, -1}
	}
	if len(st) == 0 {
		return pair{1, x}
	}
	sort.Slice(st, func(i, j int) bool {
		if st[i].d == st[j].d {
			return st[i].v < st[j].v
		}
		return st[i].d < st[j].d
	})
	if len(st) >= 2 && st[0].d < K && st[1].d < K {
		id = st[0].v
		return pair{-n, -1}
	}
	return pair{st[0].d + 1, st[0].v}
}

func chk(rd int) bool {
	K = rd
	id = -1
	o := dfs(1, 0)
	if o.d >= K {
		return true
	}
	if id == -1 {
		id = o.v
	}
	u := id
	id = -1
	o = dfs(u, 0)
	return o.d >= K
}

func solve(tc TestCase) int {
	n = tc.n
	g = make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}

	l, r, ret := 1, n, 1
	for l <= r {
		md := (l + r) >> 1
		if chk(md) {
			ret = md
			l = md + 1
		} else {
			r = md - 1
		}
	}
	return ret
}

func genTests() []TestCase {
	rng := rand.New(rand.NewSource(1))
	tests := make([]TestCase, 0, 140)

	// Deterministic chains.
	for i := 0; i < 30; i++ {
		n := 3 + i%8
		edges := make([][2]int, n-1)
		for j := 2; j <= n; j++ {
			edges[j-2] = [2]int{j - 1, j}
		}
		tests = append(tests, TestCase{n, edges})
	}

	// Deterministic stars.
	for i := 0; i < 20; i++ {
		n := 3 + i%8
		edges := make([][2]int, 0, n-1)
		for v := 2; v <= n; v++ {
			edges = append(edges, [2]int{1, v})
		}
		tests = append(tests, TestCase{n, edges})
	}

	// Random trees.
	for i := 0; i < 90; i++ {
		n := rng.Intn(30) + 3
		edges := make([][2]int, 0, n-1)
		for v := 2; v <= n; v++ {
			p := rng.Intn(v-1) + 1
			edges = append(edges, [2]int{p, v})
		}
		tests = append(tests, TestCase{n, edges})
	}
	return tests
}

func run(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTests()

	var input strings.Builder
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintln(&input, tc.n)
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
	}

	expected := make([]string, len(tests))
	for i := range tests {
		expected[i] = fmt.Sprintf("%d", solve(tests[i]))
	}

	out, err := run(bin, input.String())
	if err != nil {
		fmt.Fprintln(os.Stderr, "error running binary:", err)
		fmt.Print(out)
		os.Exit(1)
	}

	tokens := strings.Fields(out)
	if len(tokens) != len(expected) {
		fmt.Printf("wrong number of answers: got %d want %d\n", len(tokens), len(expected))
		os.Exit(1)
	}
	for i, tk := range tokens {
		got, err := strconv.Atoi(strings.TrimSpace(tk))
		if err != nil {
			fmt.Printf("test %d failed: non-integer output %q\n", i+1, tk)
			os.Exit(1)
		}
		want, _ := strconv.Atoi(expected[i])
		if got != want {
			fmt.Printf("test %d failed expected %d got %d\n", i+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed.")
}
