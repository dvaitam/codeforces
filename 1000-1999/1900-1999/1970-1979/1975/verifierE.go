package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

type expectedAnswer struct {
	value    bool
	testIdx  int
	queryIdx int
}

type info struct {
	u     int
	v     int
	valid bool
}

type solver struct {
	n     int
	adj   [][]int
	log   int
	depth []int
	up    [][]int
	tin   []int
	tout  []int
}

type segTree struct {
	n      int
	tree   []info
	colors []bool
	s      *solver
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	inputData, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read stdin: %v", err)
	}
	if len(inputData) == 0 {
		fail("empty input")
	}

	expected, err := computeAnswers(inputData)
	if err != nil {
		fail("failed to compute expected answers: %v", err)
	}

	output, err := runCandidate(candidate, inputData)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}

	if err := compareOutput(output, expected); err != nil {
		fail("%v", err)
	}

	fmt.Println("OK")
}

func computeAnswers(data []byte) ([]expectedAnswer, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("failed to read number of test cases: %v", err)
	}

	results := make([]expectedAnswer, 0)
	for tc := 1; tc <= t; tc++ {
		var n, q int
		if _, err := fmt.Fscan(reader, &n, &q); err != nil {
			return nil, fmt.Errorf("test %d: failed to read n and q: %v", tc, err)
		}
		colors := make([]bool, n+1)
		countBlack := 0
		for i := 1; i <= n; i++ {
			var c int
			if _, err := fmt.Fscan(reader, &c); err != nil {
				return nil, fmt.Errorf("test %d: failed to read initial colors: %v", tc, err)
			}
			if c != 0 {
				colors[i] = true
				countBlack++
			}
		}
		adj := make([][]int, n+1)
		for i := 0; i < n-1; i++ {
			var x, y int
			if _, err := fmt.Fscan(reader, &x, &y); err != nil {
				return nil, fmt.Errorf("test %d: failed to read edges: %v", tc, err)
			}
			adj[x] = append(adj[x], y)
			adj[y] = append(adj[y], x)
		}

		s := newSolver(n, adj)
		s.prepare(1)
		st := newSegTree(colors, s)
		current := countBlack
		for qi := 1; qi <= q; qi++ {
			var u int
			if _, err := fmt.Fscan(reader, &u); err != nil {
				return nil, fmt.Errorf("test %d: failed to read query %d: %v", tc, qi, err)
			}
			if colors[u] {
				colors[u] = false
				current--
			} else {
				colors[u] = true
				current++
			}
			st.update(u, colors[u])
			answer := s.evaluate(current, st.root())
			results = append(results, expectedAnswer{
				value:    answer,
				testIdx:  tc,
				queryIdx: qi,
			})
		}
	}
	return results, nil
}

func runCandidate(path string, input []byte) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = bytes.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func compareOutput(out string, expected []expectedAnswer) error {
	reader := bufio.NewReader(strings.NewReader(out))
	for _, ans := range expected {
		var token string
		if _, err := fmt.Fscan(reader, &token); err != nil {
			return fmt.Errorf("not enough answers; expected output for test %d query %d", ans.testIdx, ans.queryIdx)
		}
		val, ok := parseToken(token)
		if !ok {
			return fmt.Errorf("invalid answer %q at test %d query %d", token, ans.testIdx, ans.queryIdx)
		}
		if val != ans.value {
			want := "No"
			if ans.value {
				want = "Yes"
			}
			return fmt.Errorf("wrong answer at test %d query %d: expected %s, got %s", ans.testIdx, ans.queryIdx, want, token)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("extra output detected, first unexpected token %q", extra)
	}
	return nil
}

func parseToken(token string) (bool, bool) {
	t := strings.ToLower(token)
	switch t {
	case "yes":
		return true, true
	case "no":
		return false, true
	default:
		return false, false
	}
}

func newSolver(n int, adj [][]int) *solver {
	log := 1
	for (1 << log) <= n {
		log++
	}
	return &solver{
		n:   n,
		adj: adj,
		log: log,
	}
}

func (s *solver) prepare(root int) {
	s.depth = make([]int, s.n+1)
	s.tin = make([]int, s.n+1)
	s.tout = make([]int, s.n+1)
	s.up = make([][]int, s.log)
	for i := range s.up {
		s.up[i] = make([]int, s.n+1)
	}
	parent := make([]int, s.n+1)

	type frame struct {
		node   int
		parent int
		idx    int
	}
	stack := []frame{{node: root, parent: root}}
	timer := 0
	for len(stack) > 0 {
		top := &stack[len(stack)-1]
		if top.idx == 0 {
			timer++
			s.tin[top.node] = timer
			parent[top.node] = top.parent
			if top.node == top.parent {
				s.depth[top.node] = 0
			} else {
				s.depth[top.node] = s.depth[top.parent] + 1
			}
		}
		if top.idx < len(s.adj[top.node]) {
			nxt := s.adj[top.node][top.idx]
			top.idx++
			if nxt == top.parent {
				continue
			}
			stack = append(stack, frame{node: nxt, parent: top.node})
			continue
		}
		timer++
		s.tout[top.node] = timer
		stack = stack[:len(stack)-1]
	}

	for v := 1; v <= s.n; v++ {
		if parent[v] == 0 {
			parent[v] = v
		}
		s.up[0][v] = parent[v]
	}
	for k := 1; k < s.log; k++ {
		for v := 1; v <= s.n; v++ {
			s.up[k][v] = s.up[k-1][s.up[k-1][v]]
		}
	}
}

func (s *solver) lca(u, v int) int {
	if u == v {
		return u
	}
	if s.depth[u] < s.depth[v] {
		u, v = v, u
	}
	diff := s.depth[u] - s.depth[v]
	for i := s.log - 1; i >= 0; i-- {
		if diff&(1<<uint(i)) != 0 {
			u = s.up[i][u]
		}
	}
	if u == v {
		return u
	}
	for i := s.log - 1; i >= 0; i-- {
		if s.up[i][u] != s.up[i][v] {
			u = s.up[i][u]
			v = s.up[i][v]
		}
	}
	return s.up[0][u]
}

func (s *solver) dist(u, v int) int {
	l := s.lca(u, v)
	return s.depth[u] + s.depth[v] - 2*s.depth[l]
}

func (s *solver) isAncestor(u, v int) bool {
	return s.tin[u] <= s.tin[v] && s.tout[v] <= s.tout[u]
}

func (s *solver) isOnPath(a, b, x int) bool {
	if a == 0 || b == 0 || x == 0 {
		return false
	}
	l := s.lca(a, b)
	if !s.isAncestor(l, x) {
		return false
	}
	return s.isAncestor(x, a) || s.isAncestor(x, b)
}

func (s *solver) addNode(state info, x int) info {
	if !state.valid {
		return state
	}
	if state.u == 0 {
		state.u = x
		state.v = x
		state.valid = true
		return state
	}
	if x == state.u || x == state.v {
		return state
	}
	if s.isOnPath(state.u, state.v, x) {
		return state
	}
	if s.isOnPath(state.v, x, state.u) {
		state.u = x
		return state
	}
	if s.isOnPath(state.u, x, state.v) {
		state.v = x
		return state
	}
	state.valid = false
	return state
}

func (s *solver) mergeInfo(left, right info) info {
	if !left.valid || !right.valid {
		return info{valid: false}
	}
	if left.u == 0 {
		return right
	}
	if right.u == 0 {
		return left
	}
	res := left
	res = s.addNode(res, right.u)
	if !res.valid {
		return res
	}
	res = s.addNode(res, right.v)
	return res
}

func (s *solver) evaluate(count int, state info) bool {
	if count == 0 {
		return false
	}
	if !state.valid || state.u == 0 {
		return false
	}
	pathLen := s.dist(state.u, state.v) + 1
	return pathLen == count
}

func newSegTree(colors []bool, s *solver) *segTree {
	n := len(colors) - 1
	tree := make([]info, 4*n)
	st := &segTree{
		n:      n,
		tree:   tree,
		colors: colors,
		s:      s,
	}
	st.build(1, 1, n)
	return st
}

func (st *segTree) build(node, l, r int) {
	if l == r {
		if st.colors[l] {
			st.tree[node] = info{u: l, v: l, valid: true}
		} else {
			st.tree[node] = info{u: 0, v: 0, valid: true}
		}
		return
	}
	mid := (l + r) >> 1
	st.build(node<<1, l, mid)
	st.build(node<<1|1, mid+1, r)
	st.tree[node] = st.s.mergeInfo(st.tree[node<<1], st.tree[node<<1|1])
}

func (st *segTree) update(pos int, black bool) {
	st.updateRec(1, 1, st.n, pos, black)
}

func (st *segTree) updateRec(node, l, r, pos int, black bool) {
	if l == r {
		if black {
			st.tree[node] = info{u: l, v: l, valid: true}
		} else {
			st.tree[node] = info{u: 0, v: 0, valid: true}
		}
		return
	}
	mid := (l + r) >> 1
	if pos <= mid {
		st.updateRec(node<<1, l, mid, pos, black)
	} else {
		st.updateRec(node<<1|1, mid+1, r, pos, black)
	}
	st.tree[node] = st.s.mergeInfo(st.tree[node<<1], st.tree[node<<1|1])
}

func (st *segTree) root() info {
	if st.n == 0 {
		return info{u: 0, v: 0, valid: true}
	}
	return st.tree[1]
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
