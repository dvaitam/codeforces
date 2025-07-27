package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type testCase struct {
	input  string
	expect string
}

type edge struct {
	u, v int
	t    int
}

type pair struct {
	to int
	t  int
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func pathInfo(start, end int, adj [][]pair, visited []bool) (length int, stones int, ok bool) {
	if start == end {
		return 0, 0, true
	}
	visited[start] = true
	for _, e := range adj[start] {
		if visited[e.to] {
			continue
		}
		if l, s, f := pathInfo(e.to, end, adj, visited); f {
			return l + 1, s + e.t, true
		}
	}
	return 0, 0, false
}

func maxEvenPath(n int, edges []edge) int {
	adj := make([][]pair, n+1)
	for _, e := range edges {
		adj[e.u] = append(adj[e.u], pair{e.v, e.t})
		adj[e.v] = append(adj[e.v], pair{e.u, e.t})
	}
	best := 0
	visited := make([]bool, n+1)
	for i := 1; i <= n; i++ {
		for j := i; j <= n; j++ {
			for k := range visited {
				visited[k] = false
			}
			l, s, ok := pathInfo(i, j, adj, visited)
			if ok && s%2 == 0 && l > best {
				best = l
			}
		}
	}
	return best
}

func solveF(input string) string {
	in := bufio.NewScanner(strings.NewReader(input))
	in.Split(bufio.ScanWords)
	nextInt := func() int {
		if !in.Scan() {
			return 0
		}
		v, _ := strconv.Atoi(in.Text())
		return v
	}
	n := nextInt()
	edges := make([]edge, n-1)
	for i := 0; i < n-1; i++ {
		u := nextInt()
		v := nextInt()
		t := nextInt()
		edges[i] = edge{u, v, t}
	}
	m := nextInt()
	var sb strings.Builder
	for i := 0; i < m; i++ {
		id := nextInt() - 1
		edges[id].t ^= 1
		ans := maxEvenPath(n, edges)
		sb.WriteString(strconv.Itoa(ans))
		if i+1 < m {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func genTests() []testCase {
	rng := rand.New(rand.NewSource(6))
	tests := make([]testCase, 100)
	for i := 0; i < 100; i++ {
		n := rng.Intn(4) + 2
		edges := make([]edge, n-1)
		for v := 2; v <= n; v++ {
			p := rng.Intn(v-1) + 1
			t := rng.Intn(2)
			edges[v-2] = edge{p, v, t}
		}
		m := rng.Intn(3) + 1
		queries := make([]int, m)
		for j := 0; j < m; j++ {
			queries[j] = rng.Intn(n-1) + 1
			edges[queries[j]-1].t ^= 0 // no flip yet
		}
		// build input
		var sb strings.Builder
		sb.WriteString(strconv.Itoa(n))
		sb.WriteByte('\n')
		for _, e := range edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.t))
		}
		sb.WriteString(strconv.Itoa(m))
		sb.WriteByte('\n')
		// We also need to recompute initial edge states for solver
		edgeCopy := make([]edge, len(edges))
		copy(edgeCopy, edges)
		for j := 0; j < m; j++ {
			sb.WriteString(strconv.Itoa(queries[j]))
			if j+1 < m {
				sb.WriteByte('\n')
			}
		}
		sb.WriteByte('\n')
		input := sb.String()
		expect := solveF(input)
		tests[i] = testCase{input: input, expect: expect}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	tests := genTests()
	for i, tc := range tests {
		out, err := runBinary(bin, tc.input)
		if err != nil {
			fmt.Printf("test %d: runtime error: %v\n", i+1, err)
			return
		}
		if out != tc.expect {
			fmt.Printf("test %d failed\ninput:\n%s\nexpected:\n%s\nactual:\n%s\n", i+1, tc.input, tc.expect, out)
			return
		}
	}
	fmt.Println("All tests passed")
}
