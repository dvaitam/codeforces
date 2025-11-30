package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	valMax = 100000
	logN   = 18
)

type node struct {
	left, right int
	sum         int
}

type solver struct {
	seg   []node
	roots []int
	adj   [][]int
	up    [][]int
	depth []int
	vals  []int
}

func (s *solver) update(prev, l, r, pos int) int {
	cur := len(s.seg)
	s.seg = append(s.seg, s.seg[prev])
	if l == r {
		s.seg[cur].sum++
		return cur
	}
	mid := (l + r) >> 1
	if pos <= mid {
		s.seg[cur].left = s.update(s.seg[prev].left, l, mid, pos)
	} else {
		s.seg[cur].right = s.update(s.seg[prev].right, mid+1, r, pos)
	}
	s.seg[cur].sum = s.seg[s.seg[cur].left].sum + s.seg[s.seg[cur].right].sum
	return cur
}

func (s *solver) dfsIterative(root int) {
	type frame struct{ v, p, i int }
	stack := []frame{{root, 0, 0}}
	s.depth[root] = 0
	s.roots[root] = s.update(0, 1, valMax, s.vals[root])
	for len(stack) > 0 {
		f := &stack[len(stack)-1]
		v := f.v
		if f.i < len(s.adj[v]) {
			to := s.adj[v][f.i]
			f.i++
			if to == f.p {
				continue
			}
			s.up[0][to] = v
			s.depth[to] = s.depth[v] + 1
			s.roots[to] = s.update(s.roots[v], 1, valMax, s.vals[to])
			stack = append(stack, frame{to, v, 0})
		} else {
			stack = stack[:len(stack)-1]
		}
	}
}

func (s *solver) lca(u, v int) int {
	if s.depth[u] < s.depth[v] {
		u, v = v, u
	}
	diff := s.depth[u] - s.depth[v]
	for i := 0; diff > 0; i++ {
		if diff&1 != 0 {
			u = s.up[i][u]
		}
		diff >>= 1
	}
	if u == v {
		return u
	}
	for i := logN - 1; i >= 0; i-- {
		if s.up[i][u] != s.up[i][v] {
			u = s.up[i][u]
			v = s.up[i][v]
		}
	}
	return s.up[0][u]
}

func (s *solver) getSum(a, b, c, d int) int {
	return s.seg[a].sum + s.seg[b].sum - s.seg[c].sum - s.seg[d].sum
}

func (s *solver) collect(a1, b1, c1, d1, a2, b2, c2, d2, l, r int, k *int, res *[]int) {
	if *k == 0 {
		return
	}
	diff := s.getSum(a1, b1, c1, d1) - s.getSum(a2, b2, c2, d2)
	if diff == 0 {
		return
	}
	if l == r {
		*res = append(*res, l)
		*k--
		return
	}
	mid := (l + r) >> 1
	la1, lb1, lc1, ld1 := s.seg[a1].left, s.seg[b1].left, s.seg[c1].left, s.seg[d1].left
	la2, lb2, lc2, ld2 := s.seg[a2].left, s.seg[b2].left, s.seg[c2].left, s.seg[d2].left
	if s.getSum(la1, lb1, lc1, ld1) != s.getSum(la2, lb2, lc2, ld2) {
		s.collect(la1, lb1, lc1, ld1, la2, lb2, lc2, ld2, l, mid, k, res)
		if *k == 0 {
			return
		}
	}
	ra1, rb1, rc1, rd1 := s.seg[a1].right, s.seg[b1].right, s.seg[c1].right, s.seg[d1].right
	ra2, rb2, rc2, rd2 := s.seg[a2].right, s.seg[b2].right, s.seg[c2].right, s.seg[d2].right
	if s.getSum(ra1, rb1, rc1, rd1) != s.getSum(ra2, rb2, rc2, rd2) {
		s.collect(ra1, rb1, rc1, rd1, ra2, rb2, rc2, rd2, mid+1, r, k, res)
	}
}

func solveCase(ints []int) (string, error) {
	if len(ints) < 1 {
		return "", fmt.Errorf("empty test")
	}
	idx := 0
	n := ints[idx]
	idx++
	if n <= 0 {
		return "", fmt.Errorf("n must be positive")
	}
	if len(ints) < idx+n {
		return "", fmt.Errorf("need %d values, have %d", n, len(ints)-idx)
	}
	vals := make([]int, n+1)
	for i := 1; i <= n; i++ {
		vals[i] = ints[idx]
		idx++
	}

	if len(ints) < idx+2*(n-1) {
		return "", fmt.Errorf("need %d edge numbers, have %d", 2*(n-1), len(ints)-idx)
	}
	adj := make([][]int, n+1)
	for i := 0; i < n-1; i++ {
		u := ints[idx]
		v := ints[idx+1]
		idx += 2
		adj[u] = append(adj[u], v)
		adj[v] = append(adj[v], u)
	}

	if len(ints) <= idx {
		return "", fmt.Errorf("missing q")
	}
	q := ints[idx]
	idx++
	if q < 0 {
		return "", fmt.Errorf("negative q")
	}
	if len(ints) != idx+5*q {
		return "", fmt.Errorf("expected %d query numbers, got %d (n=%d idx=%d len=%d)", 5*q, len(ints)-idx, n, idx, len(ints))
	}

	s := solver{
		seg:   make([]node, 1),
		roots: make([]int, n+1),
		adj:   adj,
		up:    make([][]int, logN),
		depth: make([]int, n+1),
		vals:  vals,
	}
	for i := range s.up {
		s.up[i] = make([]int, n+1)
	}

	s.dfsIterative(1)
	for i := 1; i < logN; i++ {
		for v := 1; v <= n; v++ {
			s.up[i][v] = s.up[i-1][s.up[i-1][v]]
		}
	}

	var out strings.Builder
	for ; q > 0; q-- {
		u1 := ints[idx]
		v1 := ints[idx+1]
		u2 := ints[idx+2]
		v2 := ints[idx+3]
		k := ints[idx+4]
		idx += 5
		l1 := s.lca(u1, v1)
		l2 := s.lca(u2, v2)
		p1 := s.up[0][l1]
		p2 := s.up[0][l2]
		res := make([]int, 0, k)
		kval := k
		s.collect(s.roots[u1], s.roots[v1], s.roots[l1], s.roots[p1], s.roots[u2], s.roots[v2], s.roots[l2], s.roots[p2], 1, valMax, &kval, &res)
		out.WriteString(strconv.Itoa(len(res)))
		for _, v := range res {
			out.WriteByte(' ')
			out.WriteString(strconv.Itoa(v))
		}
		out.WriteByte('\n')
	}
	return strings.TrimRight(out.String(), "\n"), nil
}

func parseLine(line string) ([]int, error) {
	fields := strings.Fields(line)
	ints := make([]int, len(fields))
	for i, f := range fields {
		x, err := strconv.Atoi(f)
		if err != nil {
			return nil, err
		}
		ints[i] = x
	}
	return ints, nil
}

func buildInput(ints []int) string {
	var sb strings.Builder
	idx := 0
	n := ints[idx]
	idx++
	sb.WriteString(strconv.Itoa(n))
	sb.WriteByte('\n')
	for i := 0; i < n; i++ {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(ints[idx+i]))
	}
	sb.WriteByte('\n')
	idx += n
	for i := 0; i < n-1; i++ {
		sb.WriteString(strconv.Itoa(ints[idx]))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(ints[idx+1]))
		sb.WriteByte('\n')
		idx += 2
	}
	q := ints[idx]
	idx++
	sb.WriteString(strconv.Itoa(q))
	sb.WriteByte('\n')
	for i := 0; i < q; i++ {
		sb.WriteString(fmt.Sprintf("%d %d %d %d %d\n", ints[idx], ints[idx+1], ints[idx+2], ints[idx+3], ints[idx+4]))
		idx += 5
	}
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

const testcasesF2 = `3 2 1 3 1 2 1 3 1 3 3 2 3 2
3 5 4 5 1 2 2 3 1 3 3 3 1 2
2 5 5 1 2 2 2 1 1 1 2 2 2 2 2 1
3 4 2 5 1 2 2 3 1 1 1 3 1 1
1 1 2 1 1 1 1 1 1 1 1 1 2
2 4 1 1 2 3 1 2 1 2 1 1 1 1 2 2 2 2 2 2 2
3 3 5 3 1 2 2 3 2 3 2 2 2 3 1 1 3 2 2
4 4 2 2 4 1 2 1 3 3 4 1 1 3 2 2 3
4 1 4 2 4 1 2 1 3 1 4 3 1 1 1 2 1 2 3 3 1 1 4 4 2 3 3
5 2 5 2 5 5 1 2 2 3 1 4 2 5 2 1 2 2 2 3 3 3 3 2 2
6 3 5 1 1 4 1 1 2 1 3 2 4 4 5 2 6 2 4 5 1 3 1 4 5 4 3 1
4 1 3 3 2 1 2 2 3 2 4 1 2 2 3 4 3
1 4 2 1 1 1 1 1 1 1 1 1 2
3 5 2 4 1 2 1 3 3 2 1 2 1 2 1 3 1 1 2 3 2 3 3 1
4 3 2 5 3 1 2 1 3 1 4 1 1 2 4 2 1
4 3 3 3 5 1 2 2 3 3 4 3 1 3 4 4 1 4 1 4 3 1 3 3 2 2 3
6 3 2 1 3 4 2 1 2 2 3 3 4 2 5 4 6 1 6 1 2 4 2
1 1 3 1 1 1 1 3 1 1 1 1 3 1 1 1 1 3
1 4 1 1 1 1 1 2
6 2 5 3 5 5 1 1 2 2 3 3 4 1 5 5 6 2 6 3 5 2 3 3 1 5 4 2
1 3 2 1 1 1 1 2 1 1 1 1 2
2 4 5 1 2 2 1 1 2 1 3 1 1 1 1 2
4 5 2 2 5 1 2 1 3 1 4 2 4 3 2 4 1 3 4 2 3 3
2 1 4 1 2 3 1 1 2 2 3 1 1 1 1 1 1 2 2 2 3
3 2 2 5 1 2 2 3 1 2 1 3 2 1
4 2 4 4 2 1 2 2 3 1 4 3 2 2 4 4 1 1 1 2 2 2 1 3 1 2 3
6 5 5 2 5 4 3 1 2 2 3 1 4 1 5 5 6 1 6 4 6 5 2
4 1 2 5 3 1 2 2 3 3 4 1 3 4 3 3 1
3 1 3 3 1 2 2 3 1 1 1 1 3 2
1 5 2 1 1 1 1 3 1 1 1 1 3
2 5 3 1 2 1 2 2 1 1 1
3 3 2 4 1 2 1 3 3 3 3 1 1 2 1 2 2 2 2 1 1 2 3 2
5 2 1 4 4 4 1 2 2 3 2 4 2 5 2 1 4 1 4 3 2 1 1 4 1
3 1 4 5 1 2 1 3 3 1 2 3 2 3 2 1 1 1 3 3 1 1 3 3
2 3 3 1 2 3 1 2 1 1 2 1 1 2 1 1 2 1 1 1 3
6 4 5 5 4 3 3 1 2 1 3 1 4 4 5 3 6 3 6 3 2 5 2 3 6 6 6 1 4 2 1 6 1
4 3 3 3 4 1 2 2 3 3 4 2 2 3 3 3 3 4 2 1 4 3
5 5 3 4 5 1 1 2 2 3 1 4 4 5 1 1 2 1 5 2
2 3 2 1 2 1 2 2 2 1 1
4 3 4 2 2 1 2 1 3 1 4 1 2 1 1 3 1
6 2 2 4 3 3 5 1 2 2 3 3 4 4 5 5 6 3 1 2 2 6 1 2 1 2 1 1 5 5 1 5 3
5 4 4 4 4 1 1 2 1 3 1 4 3 5 3 3 4 5 5 1 4 5 5 2 2 3 5 3 1 2
5 3 5 3 4 4 1 2 2 3 2 4 1 5 3 4 3 1 1 3 3 5 1 5 1 5 1 1 1 2
6 4 1 5 1 4 2 1 2 2 3 2 4 2 5 4 6 3 4 4 3 3 1 5 5 5 3 1 3 6 1 2 2
4 3 3 3 5 1 2 1 3 3 4 2 4 3 4 1 1 2 1 3 3 2
3 1 4 3 1 2 2 3 2 1 3 3 2 2 2 3 3 2 1
3 1 4 5 1 2 2 3 3 2 3 2 2 1 2 1 1 2 1 1 3 1 2 2
4 5 1 3 1 1 2 2 3 2 4 1 2 4 1 4 3
1 3 1 1 1 1 1 2
5 5 3 1 3 4 1 2 2 3 2 4 2 5 1 5 2 4 5 3
1 4 2 1 1 1 1 1 1 1 1 1 3
4 4 2 4 3 1 2 1 3 2 4 3 3 3 3 2 2 1 1 1 2 3 3 2 1 1 3
2 5 2 1 2 1 1 1 2 2 1
3 3 3 5 1 2 1 3 3 1 2 2 2 1 3 3 3 2 2 1 2 1 2 3
3 2 2 5 1 2 1 3 3 1 2 3 3 3 3 3 1 3 1 2 2 2 2 2
2 4 4 1 2 1 2 1 2 1 2
4 1 3 3 3 1 2 1 3 2 4 1 2 2 4 4 3
5 2 1 3 2 1 1 2 2 3 3 4 1 5 1 4 3 5 2 2
1 5 1 1 1 1 1 2
4 1 3 2 2 1 2 1 3 3 4 3 4 2 4 2 3 2 3 1 4 2 4 3 2 3 3
1 4 3 1 1 1 1 1 1 1 1 1 1 1 1 1 1 1
4 5 4 3 4 1 2 1 3 3 4 3 4 4 3 3 1 3 4 3 4 3 2 4 2 2 2
4 1 4 2 5 1 2 2 3 1 4 2 3 2 2 2 2 2 2 3 2 3
5 2 1 1 2 3 1 2 2 3 1 4 1 5 2 5 3 1 4 2 3 5 5 5 3
3 1 5 5 1 2 1 3 2 3 2 3 2 2 3 3 1 1 2
4 1 5 3 5 1 2 1 3 1 4 1 2 4 3 1 2
3 5 5 5 1 2 1 3 1 3 3 1 1 2
6 5 4 2 1 1 1 1 2 1 3 2 4 1 5 4 6 1 3 1 1 6 1
3 2 2 2 1 2 2 3 1 1 1 3 2 2
6 2 5 1 5 3 3 1 2 1 3 2 4 2 5 5 6 1 5 6 1 3 2
6 1 1 1 3 3 1 1 2 1 3 1 4 3 5 5 6 1 1 4 4 2 2
2 2 2 1 2 1 1 2 2 1 2
2 2 5 1 2 2 1 2 1 2 2 1 2 2 2 1
4 2 5 2 5 1 2 2 3 2 4 2 4 2 2 4 1 2 3 3 1 2
4 2 4 2 3 1 2 1 3 2 4 1 4 1 1 1 1
4 5 3 4 3 1 2 2 3 1 4 2 3 1 3 1 1 3 1 1 1 3
5 1 2 5 5 4 1 2 1 3 1 4 2 5 1 2 5 2 2 2
5 3 4 2 1 5 1 2 1 3 1 4 2 5 2 4 1 2 2 1 2 4 4 4 2
5 2 5 1 5 4 1 2 2 3 3 4 2 5 1 3 4 5 4 3
4 5 4 3 5 1 2 1 3 2 4 3 2 3 3 2 3 4 3 3 4 2 3 2 1 3 2
2 1 2 1 2 3 1 2 1 2 2 1 2 2 2 1 2 2 1 1 3
1 5 2 1 1 1 1 1 1 1 1 1 2
3 4 4 1 1 2 2 3 3 3 3 1 1 1 1 3 3 1 1 3 2 2 1 3
3 4 2 1 1 2 1 3 3 1 1 1 1 3 1 2 2 2 3 3 1 2 2 3
6 2 1 2 1 5 4 1 2 1 3 2 4 1 5 4 6 3 5 2 3 4 1 4 4 4 2 1 6 1 5 5 1
5 2 5 4 3 3 1 2 1 3 2 4 3 5 1 5 5 2 1 2
1 4 2 1 1 1 1 3 1 1 1 1 2
5 5 5 1 1 2 1 2 2 3 2 4 4 5 3 5 3 3 5 2 3 2 1 4 3 1 1 2 2 2
4 2 5 3 3 1 2 2 3 3 4 3 1 1 4 1 3 4 4 1 1 3 4 2 4 2 2
5 1 3 4 3 3 1 2 1 3 3 4 4 5 3 3 5 2 5 3 1 2 4 5 2 2 5 5 1 3
2 3 5 1 2 2 1 2 2 2 1 1 1 2 1 2
3 1 1 4 1 2 1 3 3 1 3 1 2 2 1 2 2 3 2 3 1 3 1 2
2 1 4 1 2 2 2 1 1 1 2 1 1 2 1 1
4 2 1 1 2 1 2 2 3 2 4 3 1 2 3 3 3 1 2 3 4 2 1 3 1 4 2
2 5 5 1 2 1 2 2 1 2 2
6 1 1 1 5 4 2 1 2 1 3 2 4 4 5 5 6 2 2 2 1 5 2 1 4 6 6 2
6 2 1 5 5 3 4 1 2 2 3 3 4 4 5 1 6 1 3 3 5 3 3
4 5 1 4 1 1 2 2 3 3 4 2 3 4 1 3 1 4 4 3 4 3
2 4 3 1 2 2 1 1 2 1 1 1 2 2 2 1
1 4 2 1 1 1 1 1 1 1 1 1 2
`

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierF2 /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	lines := strings.Split(testcasesF2, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		ints, err := parseLine(line)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d parse error: %v\n", i+1, err)
			os.Exit(1)
		}
		input := buildInput(ints)
		want, err := solveCase(ints)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d solve error: %v\n", i+1, err)
			os.Exit(1)
		}
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", i+1, err, got)
			os.Exit(1)
		}
		if got != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected:\n%s\ngot:\n%s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(lines))
}
