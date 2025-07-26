package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)

type edgeE struct {
	to int
	id int
}

type testE struct {
	n     int
	m     int
	w     []int64
	edges [][2]int
	s     int
}

func genTestsE() []testE {
	rand.Seed(122005)
	tests := make([]testE, 100)
	for i := range tests {
		n := rand.Intn(5) + 1
		maxEdges := n * (n - 1) / 2
		m := rand.Intn(maxEdges-n+1) + n - 1 // at least tree
		w := make([]int64, n)
		for j := range w {
			w[j] = int64(rand.Intn(10) + 1)
		}
		edges := make([][2]int, 0, m)
		// generate tree first
		for v := 2; v <= n; v++ {
			u := rand.Intn(v-1) + 1
			edges = append(edges, [2]int{u, v})
		}
		exist := map[[2]int]bool{}
		for _, e := range edges {
			a, b := e[0], e[1]
			if a > b {
				a, b = b, a
			}
			exist[[2]int{a, b}] = true
		}
		for len(edges) < m {
			u := rand.Intn(n) + 1
			v := rand.Intn(n) + 1
			if u == v {
				continue
			}
			a, b := u, v
			if a > b {
				a, b = b, a
			}
			if exist[[2]int{a, b}] {
				continue
			}
			exist[[2]int{a, b}] = true
			edges = append(edges, [2]int{u, v})
		}
		s := rand.Intn(n) + 1
		tests[i] = testE{n: n, m: m, w: w, edges: edges, s: s}
	}
	return tests
}

func solveE(tc testE) int64 {
	n, m := tc.n, tc.m
	w := make([]int64, n+1)
	copy(w[1:], tc.w)
	g := make([][]edgeE, n+1)
	for idx, e := range tc.edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], edgeE{v, idx})
		g[v] = append(g[v], edgeE{u, idx})
	}

	disc := make([]int, n+1)
	low := make([]int, n+1)
	isBridge := make([]bool, m)
	timer := 0
	var dfsBridge func(int, int)
	dfsBridge = func(u, pe int) {
		timer++
		disc[u] = timer
		low[u] = timer
		for _, e := range g[u] {
			v := e.to
			if e.id == pe {
				continue
			}
			if disc[v] == 0 {
				dfsBridge(v, e.id)
				if low[v] < low[u] {
					low[u] = low[v]
				}
				if low[v] > disc[u] {
					isBridge[e.id] = true
				}
			} else if disc[v] < low[u] {
				low[u] = disc[v]
			}
		}
	}
	for i := 1; i <= n; i++ {
		if disc[i] == 0 {
			dfsBridge(i, -1)
		}
	}

	comp := make([]int, n+1)
	compWeight := make([]int64, n+1)
	compSize := make([]int, n+1)
	comps := 0
	var dfsComp func(int)
	dfsComp = func(u int) {
		comp[u] = comps
		compWeight[comps] += w[u]
		compSize[comps]++
		for _, e := range g[u] {
			if isBridge[e.id] {
				continue
			}
			if comp[e.to] == 0 {
				dfsComp(e.to)
			}
		}
	}
	for i := 1; i <= n; i++ {
		if comp[i] == 0 {
			comps++
			dfsComp(i)
		}
	}

	compWeight = compWeight[:comps+1]
	compSize = compSize[:comps+1]
	good := make([]bool, comps+1)
	for i := 1; i <= comps; i++ {
		if compSize[i] > 1 {
			good[i] = true
		}
	}

	tree := make([][]int, comps+1)
	for i, e := range tc.edges {
		_ = i
		u, v := e[0], e[1]
		cu, cv := comp[u], comp[v]
		if cu != cv {
			tree[cu] = append(tree[cu], cv)
			tree[cv] = append(tree[cv], cu)
		}
	}

	dpRet := make([]int64, comps+1)
	dpBest := make([]int64, comps+1)
	cycleSub := make([]bool, comps+1)
	var dfs1 func(int, int)
	dfs1 = func(u, p int) {
		dpRet[u] = compWeight[u]
		cycleSub[u] = good[u]
		for _, v := range tree[u] {
			if v == p {
				continue
			}
			dfs1(v, u)
			if cycleSub[v] {
				dpRet[u] += dpRet[v]
				cycleSub[u] = true
			}
		}
	}
	var dfs2 func(int, int)
	dfs2 = func(u, p int) {
		dpBest[u] = dpRet[u]
		for _, v := range tree[u] {
			if v == p {
				continue
			}
			dfs2(v, u)
			var cand int64
			if cycleSub[v] {
				cand = dpRet[u] - dpRet[v] + dpBest[v]
			} else {
				cand = dpRet[u] + dpBest[v]
			}
			if cand > dpBest[u] {
				dpBest[u] = cand
			}
		}
	}

	root := comp[tc.s]
	dfs1(root, 0)
	dfs2(root, 0)
	return dpBest[root]
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	tests := genTestsE()

	var input bytes.Buffer
	fmt.Fprintln(&input, len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&input, "%d %d\n", tc.n, tc.m)
		for i := 0; i < tc.n; i++ {
			if i > 0 {
				input.WriteByte(' ')
			}
			fmt.Fprint(&input, tc.w[i])
		}
		input.WriteByte('\n')
		for _, e := range tc.edges {
			fmt.Fprintf(&input, "%d %d\n", e[0], e[1])
		}
		fmt.Fprintln(&input, tc.s)
	}

	expected := make([]int64, len(tests))
	for i, tc := range tests {
		expected[i] = solveE(tc)
	}

	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input.Bytes())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, out.String())
		os.Exit(1)
	}

	scanner := bufio.NewScanner(bytes.NewReader(out.Bytes()))
	scanner.Split(bufio.ScanWords)
	for i, exp := range expected {
		if !scanner.Scan() {
			fmt.Fprintf(os.Stderr, "wrong output format on test %d\n", i+1)
			os.Exit(1)
		}
		val, err := strconv.ParseInt(scanner.Text(), 10, 64)
		if err != nil {
			fmt.Fprintf(os.Stderr, "non-integer output on test %d\n", i+1)
			os.Exit(1)
		}
		if val != exp {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\n", i+1)
			os.Exit(1)
		}
	}
	if scanner.Scan() {
		fmt.Fprintln(os.Stderr, "extra output")
		os.Exit(1)
	}
	fmt.Println("Accepted")
}
