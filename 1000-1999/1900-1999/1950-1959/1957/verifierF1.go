package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const testcasesRaw = `
4 4 5 2 3 1 2 1 3 2 4 1 1 4 4 2 3
1 2 3 1 1 1 1 1 1 1 1 1 3 1 1 1 1 1
5 1 5 5 3 2 1 2 2 3 1 4 4 5 2 2 4 4 5 2 2 5 3 1 1
2 3 5 1 2 1 2 2 2 2 3
1 1 2 1 1 1 1 3 1 1 1 1 2
4 1 3 3 3 1 2 1 3 2 4 3 1 2 4 1 3 3 3 2 3 1 4 1 3 1 3
1 5 1 1 1 1 1 1
5 2 2 2 4 2 1 2 2 3 1 4 3 5 3 3 5 2 2 2 3 2 3 5 3 1 1 3 3 3
3 1 1 5 1 2 2 3 1 1 3 2 1 3
5 3 4 1 4 5 1 2 1 3 1 4 4 5 2 2 3 2 4 3 3 1 1 1 3
3 1 3 3 1 2 1 3 3 1 3 3 1 2 3 1 2 2 2 3 1 2 2 3
4 4 1 4 3 1 2 1 3 1 4 3 4 4 1 2 1 3 1 2 2 2 1 3 4 1 3
2 4 3 1 2 2 2 2 1 1 3 2 1 1 1 1
3 2 2 5 1 2 1 3 1 1 3 2 3 1
2 5 2 1 2 3 1 2 1 2 1 1 2 1 2 3 2 1 2 2 1
2 1 5 1 2 1 1 1 2 2 3
3 4 3 4 1 2 1 3 1 3 1 3 1 1
1 5 1 1 1 1 1 3
6 5 3 1 5 2 3 1 2 2 3 2 4 3 5 2 6 1 4 6 2 5 2
3 4 3 4 1 2 2 3 2 3 2 2 3 3 2 2 3 2 3
5 3 5 2 2 4 1 2 2 3 3 4 2 5 3 2 5 5 5 3 4 3 3 4 2 4 5 1 2 2
4 1 1 2 5 1 2 2 3 1 4 1 2 2 3 3 1
6 1 1 5 2 3 2 1 2 2 3 3 4 3 5 3 6 1 5 4 1 6 1
2 1 5 1 2 2 1 2 1 2 1 1 2 1 2 3
6 2 1 1 4 5 1 1 2 1 3 2 4 3 5 3 6 3 6 5 2 6 1 3 2 3 6 3 3 4 4 5 2
3 1 4 2 1 2 2 3 3 3 2 3 2 2 3 1 3 2 1 3 3 1 3 1
2 4 3 1 2 3 2 2 2 1 3 2 1 1 1 3 1 1 1 2 2
4 1 3 2 1 1 2 1 3 1 4 1 4 1 1 4 1
5 3 4 4 5 1 1 2 1 3 3 4 2 5 3 2 5 2 1 3 4 3 3 3 1 4 1 1 5 1
1 3 1 1 1 1 1 1
6 3 2 3 5 1 1 1 2 1 3 2 4 4 5 3 6 3 4 4 5 5 1 1 3 4 4 2 5 4 5 4 3
1 2 1 1 1 1 1 1
2 1 4 1 2 3 1 2 2 2 3 2 2 2 2 3 2 1 1 2 2
3 5 1 2 1 2 2 3 1 2 2 3 1 3
1 5 3 1 1 1 1 2 1 1 1 1 2 1 1 1 1 3
6 3 2 1 4 5 2 1 2 2 3 3 4 4 5 2 6 1 6 4 5 5 3
2 5 4 1 2 3 2 2 1 2 2 1 2 2 2 1 1 2 1 2 2
2 5 3 1 2 2 1 1 2 1 3 1 2 2 1 3
2 5 3 1 2 3 1 2 1 2 2 2 1 2 2 2 1 1 1 1 3
2 5 4 1 2 2 2 2 2 2 3 2 1 1 1 2
2 3 5 1 2 2 2 2 2 1 3 1 1 1 1 3
5 5 3 5 1 5 1 2 1 3 1 4 1 5 3 3 1 5 4 3 2 5 5 5 1 3 3 2 2 3
2 2 2 1 2 2 2 2 1 2 1 1 1 2 1 1
5 1 3 5 4 2 1 2 1 3 3 4 3 5 3 2 1 2 2 3 5 1 4 2 1 4 5 1 5 2
2 2 4 1 2 3 2 1 2 1 1 2 2 2 2 3 2 1 2 2 1
2 4 1 1 2 3 1 2 1 2 2 1 1 1 2 2 2 2 1 2 2
6 4 5 2 5 5 2 1 2 2 3 3 4 4 5 5 6 3 4 2 4 2 2 1 1 5 4 1 5 6 4 6 3
2 1 2 1 2 1 1 2 2 1 3
2 5 2 1 2 2 1 1 1 1 3 1 2 2 2 3
4 5 1 3 3 1 2 1 3 1 4 2 2 2 1 2 1 4 4 1 3 3
3 4 3 3 1 2 1 3 2 2 3 3 3 1 1 1 2 1 1
4 3 2 5 3 1 2 1 3 3 4 2 3 1 3 2 1 3 4 4 4 1
2 4 1 1 2 3 1 1 2 2 1 1 2 2 1 1 1 1 2 2 1
1 3 2 1 1 1 1 2 1 1 1 1 2
3 5 3 5 1 2 1 3 1 3 1 3 3 3
5 2 4 2 1 2 1 2 1 3 3 4 2 5 3 5 2 1 1 2 1 3 5 4 3 3 5 4 1 3
1 2 2 1 1 1 1 3 1 1 1 1 3
2 2 1 1 2 1 1 2 2 1 2
4 5 5 5 2 1 2 2 3 1 4 2 1 4 3 2 3 3 3 1 3 1
4 1 4 4 1 1 2 2 3 2 4 3 2 2 1 3 1 4 3 2 4 2 3 4 2 4 2
2 1 2 1 2 2 1 1 2 2 2 1 2 2 1 3
1 2 2 1 1 1 1 2 1 1 1 1 1
5 4 2 1 1 5 1 2 2 3 1 4 4 5 3 1 4 1 3 3 2 3 3 5 1 5 3 1 2 1
1 3 2 1 1 1 1 3 1 1 1 1 2
1 1 1 1 1 1 1 3
1 5 2 1 1 1 1 1 1 1 1 1 1
4 3 2 4 5 1 2 2 3 1 4 1 1 3 1 4 1
3 4 2 4 1 2 1 3 2 3 1 3 2 3 1 2 2 3 2
4 4 3 3 3 1 2 1 3 1 4 2 3 1 3 4 3 3 2 1 1 2
6 5 2 2 4 5 3 1 2 1 3 2 4 4 5 4 6 2 3 1 2 5 1 1 4 2 6 1
6 5 5 4 5 3 2 1 2 2 3 1 4 4 5 3 6 3 1 2 5 3 2 4 6 5 6 2 1 1 1 3 2
6 4 4 2 3 4 4 1 2 2 3 3 4 3 5 4 6 1 4 6 4 4 3
5 5 3 4 5 3 1 2 2 3 3 4 2 5 2 1 4 4 4 2 1 2 3 2 2
1 5 2 1 1 1 1 1 1 1 1 1 1
6 5 4 4 2 3 4 1 2 2 3 3 4 1 5 2 6 2 2 1 6 1 2 6 3 6 3 2
4 1 5 3 3 1 2 1 3 1 4 3 3 2 4 3 3 2 4 1 4 1 1 1 3 1 2
3 1 1 1 1 2 1 3 3 1 2 3 1 3 1 1 3 3 1 1 1 3 1 3
6 1 2 1 4 4 1 1 2 2 3 3 4 4 5 5 6 3 2 5 1 2 2 6 4 5 2 1 6 5 1 6 3
2 1 4 1 2 3 2 2 2 1 3 1 1 2 2 3 2 2 2 1 1
1 2 2 1 1 1 1 2 1 1 1 1 1
1 2 3 1 1 1 1 3 1 1 1 1 1 1 1 1 1 1
3 4 5 3 1 2 2 3 1 2 1 2 3 2
2 3 1 1 2 3 1 2 2 1 3 2 1 1 2 2 1 1 2 1 3
5 5 1 4 3 1 1 2 1 3 1 4 2 5 2 2 4 2 5 3 2 1 5 2 2
5 4 5 5 3 3 1 2 1 3 1 4 4 5 2 2 4 2 4 3 1 2 5 4 2
4 2 2 4 4 1 2 1 3 2 4 1 2 1 3 2 1
2 3 2 1 2 2 2 2 2 2 2 1 2 2 2 2
6 3 4 5 5 1 4 1 2 2 3 2 4 4 5 1 6 1 4 5 1 3 3
1 1 2 1 1 1 1 3 1 1 1 1 1
3 4 1 2 1 2 1 3 3 2 3 1 2 2 2 3 1 2 2 2 2 3 3 3
2 1 3 1 2 2 1 2 1 1 1 1 2 1 1 1
4 4 1 5 5 1 2 1 3 3 4 2 2 3 1 3 2 2 4 3 1 2
1 1 3 1 1 1 1 1 1 1 1 1 3 1 1 1 1 2
4 4 4 2 4 1 2 2 3 3 4 1 3 1 2 4 2
3 3 4 3 1 2 1 3 2 1 3 2 1 2 1 1 3 3 3
4 3 2 1 3 1 2 2 3 2 4 3 1 1 4 3 2 3 3 1 2 3 3 2 2 2 2
2 1 5 1 2 3 1 1 2 1 2 2 2 1 1 1 1 1 1 2 3
1 3 1 1 1 1 1 1
5 4 3 3 2 4 1 2 1 3 1 4 1 5 1 5 3 4 1 2
5 4 3 4 1 2 1 2 1 3 2 4 2 5 1 1 2 4 4 1

`

type node struct {
	left, right int
	sum         int
}

type testCase struct {
	n       int
	values  []int
	edges   [][2]int
	queries [][5]int
}

func parseLine(line string) (testCase, error) {
	fields := strings.Fields(line)
	if len(fields) < 1 {
		return testCase{}, fmt.Errorf("empty line")
	}
	idx := 0
	n, err := strconv.Atoi(fields[idx])
	if err != nil {
		return testCase{}, fmt.Errorf("bad n")
	}
	idx++
	if len(fields) < idx+n {
		return testCase{}, fmt.Errorf("missing values")
	}
	vals := make([]int, n+1)
	for i := 1; i <= n; i++ {
		v, err := strconv.Atoi(fields[idx])
		if err != nil {
			return testCase{}, fmt.Errorf("bad value")
		}
		vals[i] = v
		idx++
	}
	edges := make([][2]int, n-1)
	for i := 0; i < n-1; i++ {
		if idx+1 >= len(fields) {
			return testCase{}, fmt.Errorf("missing edges")
		}
		u, _ := strconv.Atoi(fields[idx])
		v, _ := strconv.Atoi(fields[idx+1])
		idx += 2
		edges[i] = [2]int{u, v}
	}
	queries := make([][5]int, 0)
	if idx < len(fields) {
		q, err := strconv.Atoi(fields[idx])
		if err != nil {
			return testCase{}, fmt.Errorf("bad q")
		}
		idx++
		for i := 0; i < q && idx+4 < len(fields); i++ {
			u1, _ := strconv.Atoi(fields[idx])
			v1, _ := strconv.Atoi(fields[idx+1])
			u2, _ := strconv.Atoi(fields[idx+2])
			v2, _ := strconv.Atoi(fields[idx+3])
			k, _ := strconv.Atoi(fields[idx+4])
			idx += 5
			queries = append(queries, [5]int{u1, v1, u2, v2, k})
		}
	}
	return testCase{n: n, values: vals, edges: edges, queries: queries}, nil
}

func parseTests(raw string) ([]testCase, error) {
	lines := strings.Split(raw, "\n")
	cases := make([]testCase, 0)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		tc, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		cases = append(cases, tc)
	}
	return cases, nil
}

func solveCase(tc testCase) string {
	n := tc.n
	a := tc.values
	g := make([][]int, n+1)
	for _, e := range tc.edges {
		u, v := e[0], e[1]
		g[u] = append(g[u], v)
		g[v] = append(g[v], u)
	}
	LOG := 0
	for (1 << LOG) <= n {
		LOG++
	}
	parent := make([][]int, LOG)
	for i := 0; i < LOG; i++ {
		parent[i] = make([]int, n+1)
	}
	depth := make([]int, n+1)
	root := make([]int, n+1)
	seg := make([]node, 1)

	var newNode func() int
	newNode = func() int {
		seg = append(seg, node{})
		return len(seg) - 1
	}

	var update func(prev, l, r, pos int) int
	update = func(prev, l, r, pos int) int {
		idx := newNode()
		seg[idx] = seg[prev]
		seg[idx].sum++
		if l != r {
			mid := (l + r) >> 1
			if pos <= mid {
				seg[idx].left = update(seg[prev].left, l, mid, pos)
			} else {
				seg[idx].right = update(seg[prev].right, mid+1, r, pos)
			}
		}
		return idx
	}

	var dfs func(u, p int)
	dfs = func(u, p int) {
		parent[0][u] = p
		depth[u] = depth[p] + 1
		root[u] = update(root[p], 1, 100000, a[u])
		for _, v := range g[u] {
			if v == p {
				continue
			}
			dfs(v, u)
		}
	}
	dfs(1, 0)
	for k := 1; k < LOG; k++ {
		for v := 1; v <= n; v++ {
			parent[k][v] = parent[k-1][parent[k-1][v]]
		}
	}

	lca := func(u, v int) int {
		if depth[u] < depth[v] {
			u, v = v, u
		}
		diff := depth[u] - depth[v]
		for k := LOG - 1; k >= 0; k-- {
			if diff>>uint(k)&1 == 1 {
				u = parent[k][u]
			}
		}
		if u == v {
			return u
		}
		for k := LOG - 1; k >= 0; k-- {
			if parent[k][u] != parent[k][v] {
				u = parent[k][u]
				v = parent[k][v]
			}
		}
		return parent[0][u]
	}

	var diff func(ru1, rv1, rl1, ru2, rv2, rl2 int, valL1, valL2, l, r int) int
	diff = func(ru1, rv1, rl1, ru2, rv2, rl2 int, valL1, valL2, l, r int) int {
		c1 := seg[ru1].sum + seg[rv1].sum - 2*seg[rl1].sum
		c2 := seg[ru2].sum + seg[rv2].sum - 2*seg[rl2].sum
		if valL1 >= l && valL1 <= r {
			c1++
		}
		if valL2 >= l && valL2 <= r {
			c2++
		}
		if c1 == c2 {
			return -1
		}
		if l == r {
			return l
		}
		mid := (l + r) >> 1
		res := diff(seg[ru1].left, seg[rv1].left, seg[rl1].left,
			seg[ru2].left, seg[rv2].left, seg[rl2].left,
			valL1, valL2, l, mid)
		if res != -1 {
			return res
		}
		return diff(seg[ru1].right, seg[rv1].right, seg[rl1].right,
			seg[ru2].right, seg[rv2].right, seg[rl2].right,
			valL1, valL2, mid+1, r)
	}

	var out strings.Builder
	for _, q := range tc.queries {
		u1, v1, u2, v2 := q[0], q[1], q[2], q[3]
		l1 := lca(u1, v1)
		l2 := lca(u2, v2)
		ans := diff(root[u1], root[v1], root[l1], root[u2], root[v2], root[l2], a[l1], a[l2], 1, 100000)
		if ans == -1 {
			out.WriteString("0\n")
		} else {
			out.WriteString("1 ")
			out.WriteString(strconv.Itoa(ans))
			out.WriteByte('\n')
		}
	}
	return strings.TrimRight(out.String(), "\n")
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.values[1:] {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, e := range tc.edges {
		sb.WriteString(fmt.Sprintf("%d %d\n", e[0], e[1]))
	}
	sb.WriteString(strconv.Itoa(len(tc.queries)))
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", q[0], q[1], q[2], q[3], q[4]))
	}
	return sb.String()
}

func run(bin, input string) (string, error) {
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
		fmt.Println("usage: go run verifierF1.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	cases, err := parseTests(testcasesRaw)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse testcases: %v\n", err)
		os.Exit(1)
	}

	for idx, tc := range cases {
		input := buildInput(tc)
		expected := solveCase(tc)
		got, err := run(bin, input)
		if err != nil {
			fmt.Printf("case %d: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expected {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
