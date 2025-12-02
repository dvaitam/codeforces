package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
)

const testcasesData = `
100
5 6 4
1 3 9
1 5 10
2 4 7
2 5 10
3 4 4
4 5 18
3 3 1
1 2 3
1 2 18
2 3 4
5 4 2
1 2 7
1 4 17
1 5 18
2 4 13
3 3 2
2 3 19
1 3 8
3 2 5
5 7 2
1 2 11
1 5 10
2 3 4
2 4 18
2 5 20
3 4 10
3 5 20
3 2 1
1 2 7
2 3 9
2 1 0
1 2 5
2 1 1
1 2 18
4 4 0
1 2 8
2 3 14
2 4 15
3 4 12
4 4 2
1 3 11
2 3 9
2 4 8
3 4 6
5 7 0
1 3 5
1 5 19
3 4 4
1 4 12
4 5 4
5 3 2
3 2 20
3 2 2
1 3 4
2 3 2
2 1 0
1 2 20
4 4 3
1 2 3
1 4 6
2 3 15
2 4 4
3 3 2
1 2 16
1 3 7
3 2 2
3 2 2
1 3 9
2 3 15
3 3 2
1 2 14
1 2 12
2 3 13
4 6 2
1 2 1
1 3 3
1 4 2
2 3 5
2 4 16
3 4 10
3 3 0
1 3 14
1 3 1
3 2 20
4 4 1
1 2 8
1 4 19
2 4 19
3 4 2
5 6 1
1 2 6
1 3 16
2 3 2
2 4 10
3 4 14
3 5 3
2 1 1
1 2 14
2 1 0
1 2 1
3 2 0
1 2 7
2 3 6
5 6 0
1 3 3
1 4 15
2 3 9
2 4 17
3 5 5
4 5 1
2 1 1
1 2 9
4 4 1
1 4 20
1 2 12
2 3 18
3 4 6
5 8 3
1 2 1
1 3 9
1 4 12
2 3 20
2 4 9
2 5 19
3 4 13
4 5 4
3 3 1
1 3 17
1 2 13
2 3 11
5 6 0
1 2 18
1 4 11
1 5 4
2 4 10
2 5 5
4 5 3
2 1 0
1 2 8
5 8 4
1 2 13
1 3 10
1 4 16
2 4 3
2 5 9
3 4 6
3 5 12
4 5 1
5 4 0
1 3 16
1 4 7
3 4 15
3 5 18
5 5 2
1 2 13
1 4 17
1 5 11
3 5 8
4 5 11
2 1 0
1 2 9
5 7 4
1 2 16
1 3 3
1 4 2
1 5 5
2 3 1
2 5 6
4 5 12
5 4 3
1 3 19
1 4 17
3 4 18
4 5 16
3 2 2
1 3 6
3 2 10
4 4 1
1 2 16
2 3 14
2 4 13
3 4 5
4 3 1
1 3 2
1 4 5
2 3 20
4 4 2
1 2 20
1 3 15
1 4 16
2 4 20
2 1 0
1 2 12
5 7 3
1 4 2
2 3 11
3 5 5
1 4 17
4 2 7
2 3 2
3 5 13
4 5 3
1 3 3
2 4 9
1 3 8
3 2 3
2 4 17
4 4 0
1 3 19
1 4 3
2 4 4
3 4 5
5 8 0
1 2 14
1 3 11
1 5 3
2 3 1
2 4 6
2 5 8
3 5 5
4 5 7
2 1 0
1 2 1
3 3 1
1 2 4
1 3 9
2 3 7
4 5 1
1 2 10
1 3 18
2 3 3
2 4 19
3 4 13
3 2 1
1 3 8
2 3 6
4 4 3
1 2 2
2 3 3
2 4 10
3 4 5
4 3 2
1 3 15
2 4 13
3 4 14
5 6 3
1 2 10
1 4 6
2 3 18
2 5 15
3 5 14
4 5 8
4 3 0
1 3 15
2 3 7
2 4 15
3 3 0
1 2 10
1 3 5
2 3 16
5 8 0
1 3 13
1 4 1
1 5 2
2 3 13
2 4 19
3 4 9
3 5 7
4 5 2
4 4 1
1 2 11
1 3 8
2 4 10
3 4 5
2 1 1
1 2 18
4 5 2
1 2 12
1 3 5
1 4 4
2 3 19
2 4 18
5 6 0
1 3 16
1 5 2
2 3 5
2 5 20
3 4 12
3 5 1
2 1 0
1 2 13
4 3 2
1 2 20
1 3 9
3 4 19
3 2 1
1 2 15
2 3 11
5 8 0
1 2 5
1 3 5
1 4 6
1 5 13
2 3 16
2 4 8
3 4 7
3 5 2
3 3 2
2 3 6
1 3 20
3 2 10
2 1 1
1 2 17
4 3 2
1 2 2
1 4 18
3 4 16
5 7 4
1 2 9
1 3 2
1 4 1
1 5 1
2 3 14
3 4 8
3 5 11
2 1 1
1 2 11
4 6 2
1 2 17
1 3 12
1 4 17
2 3 7
2 4 10
3 4 13
5 9 3
1 2 8
1 3 18
1 4 15
1 5 15
2 3 14
2 5 4
3 4 8
3 5 5
4 5 7
2 1 1
1 2 18
2 1 1
1 2 8
3 2 0
1 2 9
2 3 14
4 5 0
2 4 15
3 4 9
1 4 1
4 2 4
2 3 20
3 3 1
1 3 10
1 3 18
3 2 17
2 1 0
1 2 9
4 4 2
1 2 18
1 4 8
2 3 10
3 4 16
3 3 0
1 2 1
1 3 11
2 3 1
5 6 1
1 2 6
1 3 3
1 4 7
2 4 7
2 5 5
4 5 19
5 8 3
1 3 1
1 4 12
1 5 10
2 3 8
2 4 6
3 4 16
3 5 3
4 5 20
3 3 2
2 3 2
1 3 2
3 2 11
3 3 2
1 2 15
1 3 9
3 2 12
3 2 0
1 2 12
1 3 2
4 6 3
1 2 2
1 3 2
1 4 11
2 3 2
2 4 14
3 4 7
5 7 1
1 2 10
1 3 11
2 3 1
2 4 9
2 5 15
3 5 4
4 5 13
2 1 0
1 2 20
2 1 0
1 2 6
2 1 0
1 2 6
2 1 0
1 2 20
5 6 2
1 4 11
2 4 2
2 5 12
3 4 5
3 5 17
4 5 13
4 5 0
1 4 17
2 4 2
1 2 2
2 3 15
3 4 13
5 4 2
1 2 2
1 3 2
1 4 19
3 5 13
3 3 2
1 2 11
1 3 19
3 2 18
4 6 3
1 2 5
1 3 1
1 4 2
2 3 12
2 4 10
3 4 13
3 3 2
1 2 5
1 3 15
2 3 1
5 5 4
1 5 20
2 3 18
2 5 6
3 4 3
4 5 4
5 5 1
1 3 10
1 4 9
2 3 18
2 4 11
2 5 14
5 7 3
1 2 2
1 3 14
2 5 6
1 3 19
3 4 8
4 2 9
2 5 6
2 1 1
1 2 8
`

type edgeInput struct {
	u, v, c int
}

type testCase struct {
	n     int
	m     int
	limit int
	edges []edgeInput
}

func parseTestcases() ([]testCase, error) {
	fields := strings.Fields(testcasesData)
	if len(fields) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	pos := 0
	t, err := strconv.Atoi(fields[pos])
	if err != nil {
		return nil, err
	}
	pos++
	cases := make([]testCase, 0, t)
	for i := 0; i < t; i++ {
		if pos+2 >= len(fields) {
			return nil, fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		n, _ := strconv.Atoi(fields[pos])
		m, _ := strconv.Atoi(fields[pos+1])
		limit, _ := strconv.Atoi(fields[pos+2])
		pos += 3
		edges := make([]edgeInput, m)
		for j := 0; j < m; j++ {
			if pos+2 >= len(fields) {
				return nil, fmt.Errorf("not enough edge values at case %d", i+1)
			}
			u, _ := strconv.Atoi(fields[pos])
			v, _ := strconv.Atoi(fields[pos+1])
			c, _ := strconv.Atoi(fields[pos+2])
			pos += 3
			edges[j] = edgeInput{u: u, v: v, c: c}
		}
		cases = append(cases, testCase{n: n, m: m, limit: limit, edges: edges})
	}
	return cases, nil
}

type Edge struct {
	u, v, c, idx int
	tip          bool
}

type GEdge struct {
	to        int
	c         int
	idx       int
	tip       bool
	next, rev int
}

type solver struct {
	n, m, limit int
	edges       []Edge
	uCost       []int
	numEdge     []int
	vis         []bool
	d           []int
	posEdge     []int
	first       []int
	g           []GEdge
	ds          []int
	eCnt        int
	tot         int
}

func newSolver(tc testCase) *solver {
	edges := make([]Edge, len(tc.edges))
	for i, e := range tc.edges {
		edges[i] = Edge{u: e.u, v: e.v, c: e.c, idx: i + 1}
	}
	return &solver{
		n:     tc.n,
		m:     tc.m,
		limit: tc.limit,
		edges: edges,
		uCost: make([]int, tc.n+1),
		numEdge: func() []int {
			return make([]int, tc.n+1)
		}(),
		vis: make([]bool, tc.n+1),
		d:   make([]int, tc.n+1),
		posEdge: func() []int {
			return make([]int, tc.n+1)
		}(),
		first: make([]int, tc.n+1),
		g:     make([]GEdge, 2*tc.m+2),
		ds: func() []int {
			f := make([]int, tc.n+1)
			for i := 1; i <= tc.n; i++ {
				f[i] = i
			}
			return f
		}(),
	}
}

func (s *solver) find(x int) int {
	if s.ds[x] != x {
		s.ds[x] = s.find(s.ds[x])
	}
	return s.ds[x]
}

func (s *solver) addEdge(u, v, c, idx int) {
	s.tot++
	s.g[s.tot] = GEdge{to: v, c: c, idx: idx, rev: s.tot + 1, next: s.first[u]}
	s.first[u] = s.tot
	s.tot++
	s.g[s.tot] = GEdge{to: u, c: c, idx: idx, rev: s.tot - 1, next: s.first[v]}
	s.first[v] = s.tot
}

func (s *solver) buildMST() {
	for i := range s.edges {
		e := s.edges[i]
		if e.u == 1 {
			if s.uCost[e.v] == 0 || e.c < s.uCost[e.v] {
				s.uCost[e.v] = e.c
				s.numEdge[e.v] = e.idx
			}
		}
		if e.v == 1 {
			if s.uCost[e.u] == 0 || e.c < s.uCost[e.u] {
				s.uCost[e.u] = e.c
				s.numEdge[e.u] = e.idx
			}
		}
	}
	sort.Slice(s.edges, func(i, j int) bool {
		return s.edges[i].c < s.edges[j].c
	})
	for i := 0; i < s.m; i++ {
		e := &s.edges[i]
		if e.u == 1 || e.v == 1 {
			continue
		}
		f1, f2 := s.find(e.u), s.find(e.v)
		if f1 == f2 {
			continue
		}
		e.tip = true
		s.ds[f2] = f1
	}
	for i := 0; i < s.m && s.eCnt < s.limit; i++ {
		e := &s.edges[i]
		if e.u == 1 {
			fv := s.find(e.v)
			if fv != 1 {
				e.tip = true
				s.vis[e.v] = true
				s.ds[fv] = 1
				s.eCnt++
			}
		} else if e.v == 1 {
			fu := s.find(e.u)
			if fu != 1 {
				e.tip = true
				s.vis[e.u] = true
				s.ds[fu] = 1
				s.eCnt++
			}
		}
	}
	for i := 0; i < s.m; i++ {
		e := s.edges[i]
		if e.tip && e.u != 1 && e.v != 1 {
			s.addEdge(e.u, e.v, e.c, e.idx)
		}
	}
}

func (s *solver) dfs(start int, visited []bool) {
	stack := []struct {
		v     int
		curMx int
		last  int
	}{}
	stack = append(stack, struct {
		v     int
		curMx int
		last  int
	}{start, 0, 0})
	for len(stack) > 0 {
		node := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		v := node.v
		if visited[v] {
			continue
		}
		visited[v] = true
		s.d[v] = node.curMx
		s.posEdge[v] = node.last
		for i := s.first[v]; i != 0; i = s.g[i].next {
			ge := s.g[i]
			if visited[ge.to] || s.g[i].tip {
				continue
			}
			ncur := node.curMx
			nlast := node.last
			if ge.c > ncur {
				ncur = ge.c
				nlast = i
			}
			stack = append(stack, struct {
				v     int
				curMx int
				last  int
			}{ge.to, ncur, nlast})
		}
	}
}

func (s *solver) adjust() {
	visited := make([]bool, s.n+1)
	for i := 1; i <= s.n; i++ {
		if s.vis[i] {
			s.dfs(i, visited)
		}
	}
	for s.eCnt < s.limit {
		k := 0
		minDelta := int(1e9)
		for i := 2; i <= s.n; i++ {
			if s.vis[i] || s.uCost[i] == 0 {
				continue
			}
			delta := s.uCost[i] - s.d[i]
			if delta <= minDelta {
				minDelta = delta
				k = i
			}
		}
		if k == 0 {
			break
		}
		pp := s.posEdge[k]
		if pp != 0 {
			s.g[pp].tip = true
			s.g[s.g[pp].rev].tip = true
		}
		s.vis[k] = true
		s.eCnt++
		visited = make([]bool, s.n+1)
		for i := 1; i <= s.n; i++ {
			if s.vis[i] {
				s.dfs(i, visited)
			}
		}
	}
}

func (s *solver) connected() bool {
	root := s.find(1)
	for i := 2; i <= s.n; i++ {
		if s.find(i) != root {
			return false
		}
	}
	return true
}

func (s *solver) solve() (int, bool) {
	s.buildMST()
	if !s.connected() {
		return -1, false
	}
	s.adjust()
	if s.eCnt != s.limit {
		return -1, false
	}

	totalCost := 0
	for i := 2; i <= s.n; i++ {
		if s.vis[i] {
			totalCost += s.uCost[i]
		}
	}

	used := make([]bool, len(s.g))
	for i := 1; i <= s.tot; i++ {
		if used[i] || s.g[i].tip {
			continue
		}
		totalCost += s.g[i].c
		used[i] = true
		used[s.g[i].rev] = true
	}
	return totalCost, true
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

func verifyUserOutput(tc testCase, userOutput string, expectedCost int, expectedPossible bool) error {
	lines := strings.Fields(userOutput)
	if !expectedPossible {
		if len(lines) == 1 && lines[0] == "-1" {
			return nil
		}
		return fmt.Errorf("expected -1, got %s", userOutput)
	}

	if len(lines) > 0 && lines[0] == "-1" {
		return fmt.Errorf("expected solution with cost %d, got -1", expectedCost)
	}

	if len(lines) == 0 {
		return fmt.Errorf("empty output")
	}

	numEdges, err := strconv.Atoi(lines[0])
	if err != nil {
		return fmt.Errorf("invalid edge count: %v", err)
	}

	if numEdges != tc.n-1 {
		return fmt.Errorf("expected %d edges, got %d", tc.n-1, numEdges)
	}

	if len(lines)-1 != numEdges {
		return fmt.Errorf("expected %d edge indices, found %d tokens", numEdges, len(lines)-1)
	}

	edgesByIdx := make(map[int]edgeInput)
	for i, e := range tc.edges {
		edgesByIdx[i+1] = e
	}

	adj := make([][]int, tc.n+1)
	userCost := 0
	degree1 := 0

	for i := 1; i < len(lines); i++ {
		idx, err := strconv.Atoi(lines[i])
		if err != nil {
			return fmt.Errorf("invalid edge index: %v", err)
		}
		e, ok := edgesByIdx[idx]
		if !ok {
			return fmt.Errorf("edge index %d out of bounds", idx)
		}

		userCost += e.c
		if e.u == 1 || e.v == 1 {
			degree1++
		}

		adj[e.u] = append(adj[e.u], e.v)
		adj[e.v] = append(adj[e.v], e.u)
	}

	visited := make(map[int]bool)
	var dfs func(int)
	dfs = func(u int) {
		visited[u] = true
		for _, v := range adj[u] {
			if !visited[v] {
				dfs(v)
			}
		}
	}
	dfs(1)
	if len(visited) != tc.n {
		return fmt.Errorf("graph not connected, visited %d/%d nodes", len(visited), tc.n)
	}

	if degree1 != tc.limit {
		return fmt.Errorf("degree of node 1 is %d, expected %d", degree1, tc.limit)
	}

	if userCost != expectedCost {
		return fmt.Errorf("user cost %d != expected cost %d", userCost, expectedCost)
	}

	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}

	var inputs []string
	for _, tc := range testcases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.limit))
		for _, e := range tc.edges {
			sb.WriteString(fmt.Sprintf("%d %d %d\n", e.u, e.v, e.c))
		}
		inputs = append(inputs, sb.String())
	}

	for idx, tc := range testcases {
		s := newSolver(tc)
		expectCost, expectPossible := s.solve()
		got, err := run(bin, inputs[idx])
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if err := verifyUserOutput(tc, got, expectCost, expectPossible); err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}