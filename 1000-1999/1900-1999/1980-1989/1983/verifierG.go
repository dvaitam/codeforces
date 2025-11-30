package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type Edge struct{ u, v int }

type testCase struct {
	n       int
	edges   []Edge
	vals    []int
	queries [][2]int
}

const testcasesRaw = `
100
2
1 2
6 3
2
1 1
1 1
3
1 2
1 3
7 6 8
2
2 2
3 2
5
1 2
1 3
3 4
4 5
3 0 1 2 9
2
3 4
2 4
3
1 2
1 3
10 6 8
2
2 2
1 2
2
1 2
7 0
3
2 1
2 2
2 2
5
1 2
1 3
2 4
2 5
2 6 8 3 7
1
4 4
2
1 2
2 0
2
1 2
1 2
3
1 2
1 3
1 4 4
1
2 1
4
1 2
2 3
1 4
1 4 6 6
2
1 2
1 4
2
1 2
7 0
2
1 2
1 1
2
1 2
2 0
2
2 1
1 2
4
1 2
1 3
2 4
8 3 8 6
2
2 2
3 3
5
1 2
1 3
1 4
2 5
8 6 5 7 10
3
5 4
2 2
2 3
1
9
3
1 1
1 1
1 1
5
1 2
1 3
3 4
2 5
6 4 9 9 9
1
4 3
5
1 2
2 3
1 4
3 5
8 1 5 3 8
3
3 2
2 5
1 2
2
1 2
3 3
1
2 1
5
1 2
1 3
3 4
3 5
4 8 0 5 6
1
1 4
2
1 2
9 8
1
2 1
4
1 2
1 3
3 4
1 5 9 3
1
3 2
1
8
1
1 1
5
1 2
1 3
1 4
1 5
1 0 4 7 7
2
2 4
2 3
4
1 2
1 3
1 4
1 7 6 9
2
2 2
3 2
5
1 2
2 3
2 4
4 5
8 6 2 2 10
2
4 4
4 3
1
6
1
1 1
1
0
2
1 1
1 1
2
1 2
5 0
3
2 2
1 2
1 1
4
1 2
1 3
3 4
6 1 10 5
3
4 4
4 3
1 1
1
0
3
1 1
1 1
1 1
4
1 2
1 3
2 4
8 8 4 4
1
1 4
1
6
1
1 1
1
8
1
1 1
3
1 2
2 3
4 5 2
3
3 2
3 3
1 3
2
1 2
6 6
2
1 1
2 1
2
1 2
7 8
2
1 1
1 1
4
1 2
1 3
1 4
2 4 4 2
1
4 2
4
1 2
2 3
2 4
7 10 7 5
3
4 4
1 1
4 2
1
6
3
1 1
1 1
1 1
4
1 2
2 3
2 4
5 5 6 3
2
4 1
2 3
1
0
1
1 1
5
1 2
2 3
3 4
3 5
7 6 1 3 3
3
2 4
2 4
1 2
2
1 2
9 7
3
2 2
1 1
1 2
2
1 2
2 9
2
2 2
1 1
1
9
1
1 1
4
1 2
1 3
1 4
5 2 3 0
3
2 4
1 4
2 3
2
1 2
0 8
1
1 1
5
1 2
1 3
1 4
3 5
4 10 4 3 0
1
5 1
1
6
1
1 1
3
1 2
2 3
7 0 2
3
2 2
1 3
1 1
3
1 2
2 3
5 6 3
2
1 2
3 1
4
1 2
2 3
2 4
2 0 0 6
3
3 1
4 1
1 2
5
1 2
1 3
1 4
2 5
0 2 2 2 10
2
2 4
2 5
1
8
2
1 1
1 1
2
1 2
10 2
3
2 2
1 1
2 2
2
1 2
9 1
1
1 1
5
1 2
1 3
2 4
3 5
4 5 7 4 10
2
4 3
2 4
1
10
2
1 1
1 1
5
1 2
1 3
1 4
1 5
5 1 5 1 8
2
2 5
5 3
5
1 2
1 3
3 4
1 5
6 1 8 10 10
1
5 1
5
1 2
1 3
2 4
3 5
3 5 5 5 1
2
2 4
4 1
1
8
1
1 1
3
1 2
1 3
5 2 7
1
3 3
1
1
3
1 1
1 1
1 1
3
1 2
1 3
0 1 1
3
2 2
1 1
3 3
3
1 2
2 3
10 5 4
3
3 2
3 1
3 3
5
1 2
1 3
1 4
2 5
7 4 1 2 1
1
1 2
4
1 2
2 3
1 4
6 7 1 0
2
4 1
3 3
4
1 2
1 3
3 4
8 1 5 10
1
3 3
3
1 2
2 3
8 4 2
1
1 2
1
9
2
1 1
1 1
4
1 2
1 3
3 4
6 8 9 8
2
4 3
2 2
2
1 2
6 1
3
1 2
2 1
1 2
5
1 2
1 3
3 4
1 5
0 5 6 0 5
3
3 3
3 5
3 4
4
1 2
2 3
2 4
10 6 2 1
1
2 3
5
1 2
1 3
2 4
1 5
1 10 8 10 8
3
2 4
4 2
2 1
2
1 2
6 3
3
2 2
2 2
1 1
4
1 2
2 3
2 4
5 1 7 3
3
2 1
3 3
1 2
3
1 2
1 3
9 1 3
3
3 3
1 1
1 1
4
1 2
2 3
3 4
10 7 6 1
2
3 3
4 2
1
5
1
1 1
3
1 2
1 3
3 1 1
3
3 1
3 1
3 3
5
1 2
1 3
1 4
4 5
3 4 1 5 10
2
1 2
5 1
1
10
1
1 1
3
1 2
2 3
2 0 0
2
2 2
1 2
5
1 2
2 3
1 4
1 5
6 8 2 1 0
3
2 3
4 2
4 2
4
1 2
2 3
3 4
6 7 6 1
3
3 1
3 4
1 1
2
1 2
6 6
1
1 1
5
1 2
1 3
3 4
2 5
9 5 8 7 8
1
2 1
4
1 2
1 3
1 4
7 4 1 5
3
3 2
4 1
1 3
3
1 2
1 3
1 1 2
1
2 3
5
1 2
1 3
3 4
4 5
1 10 4 5 5
2
1 5
1 3
5
1 2
1 3
3 4
2 5
9 3 6 0 9
2
1 2
1 4
5
1 2
1 3
3 4
2 5
1 6 10 10 0
3
3 3
5 1
4 5
5
1 2
1 3
3 4
4 5
3 1 5 10 4
3
4 2
3 5
2 2
3
1 2
2 3
8 3 5
3
3 2
1 2
2 2
2
1 2
8 6
3
1 1
2 1
1 1
3
1 2
1 3
2 4 0
2
2 3
3 2
3
1 2
2 3
5 10 2
2
1 1
1 2
3
1 2
2 3
3 3 10
3
2 2
3 3
3 2
3
1 2
1 3
10 4 6
2
3 1
3 2
`

func parseTests(raw string) ([]testCase, error) {
	tokens := strings.Fields(raw)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("empty testcases")
	}
	idx := 0
	t, err := strconv.Atoi(tokens[idx])
	if err != nil {
		return nil, fmt.Errorf("bad t: %v", err)
	}
	idx++
	cases := make([]testCase, 0, t)
	for c := 0; c < t; c++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("truncated before case %d", c+1)
		}
		n, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("bad n in case %d", c+1)
		}
		idx++
		edges := make([]Edge, n-1)
		for i := 0; i < n-1; i++ {
			if idx+1 >= len(tokens) {
				return nil, fmt.Errorf("missing edges in case %d", c+1)
			}
			u, _ := strconv.Atoi(tokens[idx])
			v, _ := strconv.Atoi(tokens[idx+1])
			idx += 2
			edges[i] = Edge{u: u, v: v}
		}
		vals := make([]int, n+1)
		for i := 1; i <= n; i++ {
			if idx >= len(tokens) {
				return nil, fmt.Errorf("missing vals in case %d", c+1)
			}
			v, _ := strconv.Atoi(tokens[idx])
			vals[i] = v
			idx++
		}
		if idx >= len(tokens) {
			return nil, fmt.Errorf("missing q in case %d", c+1)
		}
		q, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("bad q in case %d", c+1)
		}
		idx++
		queries := make([][2]int, q)
		for i := 0; i < q; i++ {
			if idx+1 >= len(tokens) {
				return nil, fmt.Errorf("missing queries in case %d", c+1)
			}
			x, _ := strconv.Atoi(tokens[idx])
			y, _ := strconv.Atoi(tokens[idx+1])
			idx += 2
			queries[i] = [2]int{x, y}
		}
		cases = append(cases, testCase{n: n, edges: edges, vals: vals, queries: queries})
	}
	return cases, nil
}

func solve(tc testCase) string {
	n := tc.n
	g := make([][]int, n+1)
	for _, e := range tc.edges {
		g[e.u] = append(g[e.u], e.v)
		g[e.v] = append(g[e.v], e.u)
	}
	const LOG = 20
	up := make([][]int, n+1)
	depth := make([]int, n+1)
	for i := 0; i <= n; i++ {
		up[i] = make([]int, LOG)
	}
	var dfs func(int, int)
	dfs = func(u, p int) {
		up[u][0] = p
		for i := 1; i < LOG; i++ {
			up[u][i] = up[up[u][i-1]][i-1]
		}
		for _, v := range g[u] {
			if v == p {
				continue
			}
			depth[v] = depth[u] + 1
			dfs(v, u)
		}
	}
	dfs(1, 1)
	lca := func(a, b int) int {
		if depth[a] < depth[b] {
			a, b = b, a
		}
		diff := depth[a] - depth[b]
		for i := LOG - 1; i >= 0; i-- {
			if diff&(1<<i) != 0 {
				a = up[a][i]
			}
		}
		if a == b {
			return a
		}
		for i := LOG - 1; i >= 0; i-- {
			if up[a][i] != up[b][i] {
				a = up[a][i]
				b = up[b][i]
			}
		}
		return up[a][0]
	}
	vals := tc.vals
	pathSum := func(x, y int) int64 {
		l := lca(x, y)
		idx := 0
		ans := int64(0)
		u := x
		for u != l {
			ans += int64(vals[u] ^ idx)
			idx++
			u = up[u][0]
		}
		ans += int64(vals[l] ^ idx)
		idx++
		var stack []int
		v := y
		for v != l {
			stack = append(stack, v)
			v = up[v][0]
		}
		for i := len(stack) - 1; i >= 0; i-- {
			ans += int64(vals[stack[i]] ^ idx)
			idx++
		}
		return ans
	}
	res := make([]string, len(tc.queries))
	for i, q := range tc.queries {
		res[i] = strconv.FormatInt(pathSum(q[0], q[1]), 10)
	}
	return strings.Join(res, "\n")
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e.u, e.v))
	}
	for i := 1; i <= tc.n; i++ {
		if i > 1 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(tc.vals[i]))
	}
	sb.WriteByte('\n')
	sb.WriteString(strconv.Itoa(len(tc.queries)))
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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
		fmt.Println("usage: go run verifierG.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range tests {
		input := buildInput(tc)
		want := solve(tc)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed:\nexpected:\n%s\ngot:\n%s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
