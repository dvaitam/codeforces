package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type edge struct{ to, id int }

type testcase struct {
	n     int
	m     int
	k     int
	edges [][2]int
	pairs [][2]int
}

// Embedded copy of testcasesF.txt so the verifier is self-contained.
const testcasesRaw = `5 3 5 1 2 2 3 2 4 3 5 5 3 4 5 3 4
4 4 3 1 2 1 3 3 4 4 3 4 2 2 4 2 1
4 5 2 1 2 2 3 1 4 4 2 4 2 1 2 2 3 3 1
2 4 2 1 2 2 1 1 2 2 1 1 2
2 6 1 1 2 1 2 2 1 2 1 1 2 1 2 2 1
4 4 1 1 2 1 3 2 4 4 3 2 4 1 3 3 2
2 5 2 1 2 1 2 2 1 2 1 2 1 2 1
3 3 2 1 2 1 3 2 1 3 2 1 2
2 3 1 1 2 1 2 2 1 2 1
6 3 5 1 2 1 3 1 4 2 5 1 6 2 4 6 2 1 4
4 3 3 1 2 2 3 3 4 1 3 4 2 4 1
6 5 6 1 2 2 3 1 4 3 5 3 6 2 6 6 3 2 1 3 6 4 2
4 3 2 1 2 1 3 3 4 1 2 4 1 1 2
6 1 2 1 2 2 3 3 4 1 5 2 6 2 5
6 5 1 1 2 2 3 1 4 1 5 4 6 1 6 1 5 2 4 6 2 5 3
6 6 5 1 2 1 3 2 4 1 5 5 6 3 6 3 5 2 5 3 1 5 4 3 1
4 3 3 1 2 2 3 3 4 2 3 1 4 1 4
3 4 3 1 2 2 3 1 3 3 2 3 1 3 1
2 1 1 1 2 2 1
3 6 3 1 2 1 3 1 2 2 1 2 1 3 2 3 2 3 1
6 3 5 1 2 1 3 2 4 4 5 5 6 6 1 2 3 4 6
3 5 1 1 2 1 3 1 2 2 3 2 1 1 3 3 2
4 1 4 1 2 2 3 1 4 4 2
5 4 2 1 2 2 3 1 4 2 5 2 5 4 2 5 1 5 1
3 2 2 1 2 2 3 1 3 3 1
6 4 3 1 2 2 3 1 4 1 5 5 6 2 6 4 6 5 3 4 1
3 5 2 1 2 2 3 3 1 3 2 2 1 2 3 1 2
4 2 3 1 2 1 3 1 4 2 4 1 3
3 3 2 1 2 1 3 3 2 3 1 2 3
3 5 2 1 2 1 3 3 2 3 2 1 2 3 2 1 2
4 2 2 1 2 2 3 1 4 3 4 3 1
3 4 3 1 2 1 3 1 2 2 3 2 1 2 1
5 6 1 1 2 1 3 1 4 1 5 4 2 3 1 4 5 3 5 1 5 1 4
5 5 3 1 2 1 3 1 4 1 5 1 5 5 2 3 4 5 2 4 2
6 5 2 1 2 1 3 3 4 1 5 3 6 2 1 5 3 5 2 1 6 6 5
2 6 2 1 2 2 1 1 2 2 1 2 1 2 1 1 2
4 3 3 1 2 1 3 2 4 3 2 3 4 2 1
5 6 5 1 2 2 3 3 4 1 5 5 1 3 1 4 2 1 2 1 5 3 2
4 2 1 1 2 2 3 2 4 3 1 1 4
4 6 4 1 2 1 3 3 4 4 1 2 1 2 1 1 3 1 2 2 4
4 3 3 1 2 2 3 3 4 1 2 2 1 2 1
2 5 1 1 2 2 1 2 1 1 2 2 1 2 1
5 3 5 1 2 1 3 2 4 3 5 5 3 4 1 5 3
4 6 1 1 2 2 3 3 4 1 3 1 2 1 3 4 1 1 4 2 3
6 3 3 1 2 2 3 3 4 4 5 4 6 3 1 5 3 1 3
2 2 2 1 2 2 1 1 2
5 3 4 1 2 1 3 3 4 2 5 5 1 3 1 4 3
3 3 2 1 2 2 3 2 3 3 2 1 2
2 4 1 1 2 2 1 2 1 1 2 1 2
5 4 5 1 2 2 3 1 4 1 5 4 3 2 4 1 4 5 1
4 1 2 1 2 1 3 2 4 4 3
2 3 2 1 2 2 1 1 2 1 2
4 1 2 1 2 1 3 1 4 3 2
3 5 2 1 2 2 3 1 2 2 1 1 2 3 2 3 2
4 2 4 1 2 1 3 3 4 1 2 3 2
2 1 1 1 2 1 2
5 1 3 1 2 2 3 3 4 1 5 2 5
2 3 2 1 2 1 2 2 1 2 1
4 5 3 1 2 1 3 3 4 2 1 3 1 2 3 3 4 3 4
3 4 2 1 2 1 3 2 1 2 1 2 1 2 3
4 2 4 1 2 1 3 1 4 4 2 1 3
3 1 2 1 2 1 3 1 2
4 5 2 1 2 1 3 3 4 2 4 4 1 3 2 2 3 1 4
3 5 2 1 2 1 3 2 1 2 3 1 2 3 1 2 3
3 3 1 1 2 1 3 2 3 2 3 1 3
5 2 3 1 2 1 3 2 4 4 5 3 4 5 2
2 6 2 1 2 2 1 2 1 2 1 1 2 2 1 2 1
6 3 6 1 2 1 3 2 4 3 5 2 6 2 3 2 3 2 5
6 2 1 1 2 1 3 2 4 3 5 3 6 2 1 3 5
4 1 4 1 2 2 3 3 4 3 1
6 5 1 1 2 2 3 3 4 4 5 5 6 4 3 6 5 2 1 4 6 1 5
4 5 3 1 2 1 3 1 4 2 4 3 1 3 4 1 2 2 3
6 4 2 1 2 1 3 1 4 3 5 2 6 4 5 6 1 4 2 3 1
2 6 1 1 2 1 2 2 1 1 2 2 1 1 2 2 1
5 5 2 1 2 1 3 2 4 4 5 1 3 3 1 5 3 1 2 3 5
2 1 2 1 2 1 2
3 2 1 1 2 1 3 3 1 2 3
6 2 4 1 2 1 3 3 4 1 5 3 6 2 6 3 1
5 4 2 1 2 2 3 2 4 2 5 4 2 5 2 5 4 3 4
4 4 2 1 2 2 3 1 4 2 3 3 1 3 4 2 4
3 3 3 1 2 1 3 2 1 2 3 1 3
6 1 6 1 2 1 3 2 4 3 5 4 6 3 6
2 5 2 1 2 2 1 2 1 1 2 1 2 2 1
2 2 1 1 2 1 2 2 1
2 5 2 1 2 1 2 2 1 1 2 2 1 1 2
3 6 2 1 2 1 3 2 3 2 3 2 3 3 2 2 3 3 1
2 3 1 1 2 2 1 2 1 1 2
4 2 1 1 2 2 3 3 4 2 3 2 4
5 2 3 1 2 2 3 1 4 4 5 4 1 2 1
6 3 3 1 2 2 3 2 4 2 5 1 6 5 2 1 4 4 3
5 6 4 1 2 2 3 2 4 1 5 1 2 5 2 2 3 4 2 5 3 4 2
2 3 1 1 2 1 2 2 1 2 1
2 2 1 1 2 2 1 2 1
5 6 5 1 2 2 3 2 4 1 5 4 2 2 1 5 1 5 2 5 4 3 2
6 1 3 1 2 1 3 3 4 2 5 1 6 4 1
5 3 3 1 2 1 3 1 4 2 5 2 4 2 3 5 1
2 2 1 1 2 2 1 1 2
6 1 6 1 2 1 3 1 4 3 5 4 6 2 5
4 6 3 1 2 2 3 3 4 1 4 1 4 4 1 1 4 2 4 2 4
6 1 3 1 2 2 3 2 4 2 5 1 6 6 5`

func parseTestcases() ([]testcase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testcase, 0, len(lines))
	for idx, ln := range lines {
		ln = strings.TrimSpace(ln)
		if ln == "" {
			continue
		}
		fields := strings.Fields(ln)
		if len(fields) < 3 {
			return nil, fmt.Errorf("line %d: too few fields", idx+1)
		}
		pos := 0
		nextInt := func() (int, error) {
			if pos >= len(fields) {
				return 0, fmt.Errorf("line %d: unexpected EOF", idx+1)
			}
			v, err := strconv.Atoi(fields[pos])
			pos++
			return v, err
		}
		n, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		m, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", idx+1, err)
		}
		k, err := nextInt()
		if err != nil {
			return nil, fmt.Errorf("line %d: parse k: %v", idx+1, err)
		}
		expected := 3 + 2*(n-1) + 2*m
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d fields, got %d", idx+1, expected, len(fields))
		}
		edges := make([][2]int, n-1)
		for i := 0; i < n-1; i++ {
			u, err := nextInt()
			if err != nil {
				return nil, err
			}
			v, err := nextInt()
			if err != nil {
				return nil, err
			}
			edges[i] = [2]int{u, v}
		}
		pairs := make([][2]int, m)
		for i := 0; i < m; i++ {
			s, err := nextInt()
			if err != nil {
				return nil, err
			}
			t, err := nextInt()
			if err != nil {
				return nil, err
			}
			pairs[i] = [2]int{s, t}
		}
		cases = append(cases, testcase{n: n, m: m, k: k, edges: edges, pairs: pairs})
	}
	return cases, nil
}

// solve implements the logic from 1336F.go for one testcase.
func solve(tc testcase) string {
	n, k := tc.n, tc.k

	log := 0
	for (1 << log) <= n {
		log++
	}
	adj := make([][]edge, n+1)
	for i, e := range tc.edges {
		u, v := e[0], e[1]
		adj[u] = append(adj[u], edge{to: v, id: i})
		adj[v] = append(adj[v], edge{to: u, id: i})
	}

	up := make([][]int, log)
	for i := range up {
		up[i] = make([]int, n+1)
	}
	parentEdge := make([]int, n+1)
	depth := make([]int, n+1)

	var dfs func(v, p int)
	dfs = func(v, p int) {
		for _, e := range adj[v] {
			if e.to == p {
				continue
			}
			up[0][e.to] = v
			parentEdge[e.to] = e.id
			depth[e.to] = depth[v] + 1
			for i := 1; i < log; i++ {
				up[i][e.to] = up[i-1][up[i-1][e.to]]
			}
			dfs(e.to, v)
		}
	}
	dfs(1, 0)

	lca := func(a, b int) int {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		diff := depth[a] - depth[b]
		for i := 0; i < log; i++ {
			if diff&(1<<i) != 0 {
				a = up[i][a]
			}
		}
		if a == b {
			return a
		}
		for i := log - 1; i >= 0; i-- {
			if up[i][a] != up[i][b] {
				a = up[i][a]
				b = up[i][b]
			}
		}
		return up[0][a]
	}

	getPathEdges := func(u, v int) []int {
		p := lca(u, v)
		res := make([]int, 0)
		x := u
		for x != p {
			res = append(res, parentEdge[x])
			x = up[0][x]
		}
		tmp := make([]int, 0)
		x = v
		for x != p {
			tmp = append(tmp, parentEdge[x])
			x = up[0][x]
		}
		for i := len(tmp) - 1; i >= 0; i-- {
			res = append(res, tmp[i])
		}
		return res
	}

	edgeTrav := make([][]int, n-1)
	for i, pr := range tc.pairs {
		path := getPathEdges(pr[0], pr[1])
		for _, e := range path {
			edgeTrav[e] = append(edgeTrav[e], i)
		}
	}

	pairCount := make(map[uint64]int)
	ans := 0
	for _, list := range edgeTrav {
		L := len(list)
		for i := 0; i < L; i++ {
			for j := i + 1; j < L; j++ {
				a := list[i]
				b := list[j]
				if a > b {
					a, b = b, a
				}
				key := uint64(a)<<32 | uint64(b)
				pairCount[key]++
				if pairCount[key] == k {
					ans++
				}
			}
		}
	}
	return strconv.Itoa(ans)
}

func buildIfGo(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp, err := os.CreateTemp("", "solbin*")
		if err != nil {
			return "", nil, err
		}
		tmp.Close()
		out, err := exec.Command("go", "build", "-o", tmp.Name(), path).CombinedOutput()
		if err != nil {
			os.Remove(tmp.Name())
			return "", nil, fmt.Errorf("build failed: %v\n%s", err, out)
		}
		return tmp.Name(), func() { os.Remove(tmp.Name()) }, nil
	}
	return path, func() {}, nil
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}

	cases, err := parseTestcases()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	bin, cleanup, err := buildIfGo(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	for idx, tc := range cases {
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.m, tc.k)
		for _, e := range tc.edges {
			fmt.Fprintf(&sb, "%d %d\n", e[0], e[1])
		}
		for _, p := range tc.pairs {
			fmt.Fprintf(&sb, "%d %d\n", p[0], p[1])
		}
		input := sb.String()
		expect := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(expect) {
			fmt.Printf("case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", idx+1, input, expect, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
