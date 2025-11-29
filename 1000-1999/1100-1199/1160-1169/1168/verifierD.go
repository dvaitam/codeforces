package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type segTree struct {
	n   int
	max []int
	add []int
}

func newSegTree(n int, arr []int) *segTree {
	size := 1
	for size < n {
		size <<= 1
	}
	st := &segTree{n: n, max: make([]int, size<<1), add: make([]int, size<<1)}
	st.build(1, 0, n-1, arr)
	return st
}

func (st *segTree) build(idx, l, r int, arr []int) {
	if l == r {
		st.max[idx] = arr[l]
		return
	}
	mid := (l + r) >> 1
	st.build(idx<<1, l, mid, arr)
	st.build(idx<<1|1, mid+1, r, arr)
	if st.max[idx<<1] > st.max[idx<<1|1] {
		st.max[idx] = st.max[idx<<1]
	} else {
		st.max[idx] = st.max[idx<<1|1]
	}
}

func (st *segTree) apply(idx, v int) {
	st.max[idx] += v
	st.add[idx] += v
}

func (st *segTree) push(idx int) {
	if st.add[idx] != 0 {
		st.apply(idx<<1, st.add[idx])
		st.apply(idx<<1|1, st.add[idx])
		st.add[idx] = 0
	}
}

func (st *segTree) update(idx, l, r, ql, qr, v int) {
	if ql <= l && r <= qr {
		st.apply(idx, v)
		return
	}
	st.push(idx)
	mid := (l + r) >> 1
	if ql <= mid {
		st.update(idx<<1, l, mid, ql, qr, v)
	}
	if qr > mid {
		st.update(idx<<1|1, mid+1, r, ql, qr, v)
	}
	if st.max[idx<<1] > st.max[idx<<1|1] {
		st.max[idx] = st.max[idx<<1]
	} else {
		st.max[idx] = st.max[idx<<1|1]
	}
}

func (st *segTree) rangeAdd(l, r, v int) {
	st.update(1, 0, st.n-1, l, r, v)
}

func (st *segTree) Max() int {
	return st.max[1]
}

type solver struct {
	children [][]int
	edgeChar []byte
	lidx     []int
	ridx     []int
	leafArr  [][]int
	leafCnt  int
	depthBad bool
	depth    int
}

func newSolver(n int, parents []int, chars []byte) *solver {
	s := &solver{
		children: make([][]int, n+1),
		edgeChar: make([]byte, n+1),
		lidx:     make([]int, n+1),
		ridx:     make([]int, n+1),
		leafArr:  make([][]int, 26),
	}
	for i := 0; i < 26; i++ {
		s.leafArr[i] = make([]int, n)
	}
	copy(s.edgeChar, chars)
	for v := 2; v <= n; v++ {
		p := parents[v]
		s.children[p] = append(s.children[p], v)
	}
	return s
}

func (s *solver) dfs(v, d int, counts []int) {
	if v != 1 {
		c := s.edgeChar[v]
		if c != '?' {
			counts[c-'a']++
		}
	}
	if len(s.children[v]) == 0 {
		if s.leafCnt == 0 {
			s.depth = d
		} else if s.depth != d {
			s.depthBad = true
		}
		s.lidx[v] = s.leafCnt
		for i := 0; i < 26; i++ {
			s.leafArr[i][s.leafCnt] = counts[i]
		}
		s.leafCnt++
		s.ridx[v] = s.leafCnt - 1
	} else {
		s.lidx[v] = s.leafCnt
		for _, u := range s.children[v] {
			s.dfs(u, d+1, counts)
		}
		s.ridx[v] = s.leafCnt - 1
	}
	if v != 1 {
		c := s.edgeChar[v]
		if c != '?' {
			counts[c-'a']--
		}
	}
}

type query struct {
	v int
	c byte
}

func solveCase(n int, parents []int, chars []byte, qs []query) []string {
	s := newSolver(n, parents, chars)
	counts := make([]int, 26)
	s.dfs(1, 0, counts)
	outputs := make([]string, 0, len(qs))
	if s.depthBad {
		for range qs {
			outputs = append(outputs, "Fou")
		}
		return outputs
	}

	L := s.leafCnt
	trees := make([]*segTree, 26)
	mval := make([]int, 26)
	sumM := 0
	var sumW int64
	for i := 0; i < 26; i++ {
		arr := s.leafArr[i][:L]
		t := newSegTree(L, arr)
		trees[i] = t
		mval[i] = t.Max()
		sumM += mval[i]
		sumW += int64(mval[i] * (i + 1))
	}

	const totalInd = 351

	for _, q := range qs {
		v := q.v
		cb := q.c
		old := s.edgeChar[v]
		if old != cb {
			if old != '?' {
				ci := int(old - 'a')
				l, r := s.lidx[v], s.ridx[v]
				trees[ci].rangeAdd(l, r, -1)
				newm := trees[ci].Max()
				if newm != mval[ci] {
					sumM += newm - mval[ci]
					sumW += int64((newm - mval[ci]) * (ci + 1))
					mval[ci] = newm
				}
			}
			if cb != '?' {
				ci := int(cb - 'a')
				l, r := s.lidx[v], s.ridx[v]
				trees[ci].rangeAdd(l, r, 1)
				newm := trees[ci].Max()
				if newm != mval[ci] {
					sumM += newm - mval[ci]
					sumW += int64((newm - mval[ci]) * (ci + 1))
					mval[ci] = newm
				}
			}
			s.edgeChar[v] = cb
		}
		if sumM > s.depth {
			outputs = append(outputs, "Fou")
		} else {
			ans := sumW + int64(s.depth-sumM)*totalInd
			outputs = append(outputs, fmt.Sprintf("Shi %d", ans))
		}
	}
	return outputs
}

const testcasesRaw = `100
3 5 1 c 2 a 2 ? 3 b 2 ? 3 ? 2 b
3 5 1 a 1 b 2 c 2 c 3 ? 3 ? 3 b
4 1 1 b 2 b 2 ? 4 c
5 5 1 c 2 b 2 a 3 b 4 a 3 c 4 a 2 ? 5 a
4 1 1 b 1 c 2 ? 2 a
6 5 1 ? 2 c 3 b 1 c 1 a 2 a 3 ? 4 c 3 a 4 c
4 2 1 ? 2 ? 3 a 4 c 3 b
4 4 1 c 2 a 2 c 2 ? 4 b 2 c 3 c
4 5 1 ? 1 a 3 a 3 c 4 ? 3 c 2 c 2 c
4 5 1 c 2 a 1 b 3 b 4 c 2 c 2 ? 4 a
2 5 1 c 2 ? 2 a 2 b 2 c 2 a
2 5 1 c 2 c 2 a 2 b 2 c 2 ?
4 4 1 ? 2 a 2 b 2 a 3 ? 4 b 2 ?
6 3 1 b 1 c 1 b 1 a 5 b 5 a 2 ? 2 b
6 3 1 a 2 a 3 a 3 b 3 ? 2 c 3 b 2 a
3 2 1 b 1 ? 3 a 3 b
4 5 1 a 2 c 1 a 2 a 2 a 2 ? 2 a 4 ?
4 2 1 a 2 ? 3 ? 4 c 3 c
3 3 1 a 1 a 3 a 2 a 3 ?
6 5 1 a 2 a 2 ? 3 ? 4 ? 2 b 3 c 6 a 5 b 5 b
2 3 1 c 2 ? 2 ? 2 c
6 5 1 a 2 b 1 ? 1 a 5 a 5 b 2 c 2 a 2 ? 4 c
2 1 1 a 2 a
6 3 1 a 1 b 3 b 4 ? 3 c 6 ? 4 ? 2 ?
6 2 1 b 2 ? 1 ? 2 b 1 ? 5 ? 6 b
3 3 1 b 2 b 3 ? 3 b 3 a
4 4 1 b 1 c 1 c 3 b 3 ? 4 b 3 a
5 1 1 a 2 b 1 a 2 c 3 b
4 2 1 a 2 b 1 c 2 ? 2 a
2 1 1 c 2 c
5 5 1 a 1 c 2 ? 4 ? 2 b 5 ? 3 c 2 c 2 ?
2 4 1 a 2 c 2 c 2 a 2 a
5 5 1 a 2 b 3 b 2 c 5 a 2 b 4 ? 4 b 5 ?
4 2 1 a 1 ? 2 a 3 c 3 a
3 1 1 c 1 c 3 ?
5 5 1 c 2 a 1 c 1 a 2 b 3 ? 5 b 3 ? 3 b
6 2 1 c 2 c 1 ? 4 ? 3 c 5 c 5 a
6 2 1 a 1 ? 1 ? 2 b 5 b 2 ? 2 c
3 2 1 a 1 a 3 b 2 ?
6 1 1 c 2 c 3 b 1 a 1 c 3 a
3 4 1 ? 2 a 2 ? 2 b 2 ? 3 ?
2 5 1 b 2 c 2 ? 2 ? 2 b 2 a
4 3 1 ? 2 ? 1 a 2 c 3 b 3 b
3 2 1 b 2 b 2 b 2 b
5 4 1 a 1 ? 2 c 4 a 2 b 5 a 5 ? 2 b
3 1 1 ? 1 b 3 a
4 1 1 c 2 ? 2 ? 2 c
2 3 1 c 2 a 2 a 2 a
2 1 1 b 2 ?
2 2 1 b 2 c 2 c
2 4 1 ? 2 c 2 ? 2 ? 2 ?
3 4 1 b 2 ? 3 b 3 ? 3 c 2 ?
5 1 1 b 1 ? 1 a 3 b 2 a
5 4 1 c 1 a 1 ? 1 c 4 a 5 a 3 b 2 b
3 5 1 a 1 c 2 b 3 c 2 ? 2 a 2 a
6 1 1 b 2 a 1 ? 4 b 5 a 5 c
5 5 1 a 2 ? 2 b 1 a 5 c 2 c 4 b 5 a 4 ?
2 5 1 b 2 a 2 b 2 ? 2 c 2 c
3 2 1 ? 2 c 2 ? 3 ?
5 5 1 c 1 b 1 a 2 c 2 ? 2 b 3 c 5 ? 2 c
5 2 1 b 1 ? 2 ? 4 a 3 b 5 c
5 3 1 ? 2 ? 2 b 3 c 4 a 5 c 4 c
3 5 1 a 1 b 2 b 2 a 2 ? 3 c 2 ?
3 2 1 c 2 c 2 a 3 c
4 2 1 ? 2 c 1 a 3 a 4 b
6 2 1 ? 2 b 2 a 2 a 1 ? 2 ? 5 ?
4 5 1 c 1 a 3 b 4 a 3 c 3 ? 3 c 3 c
6 2 1 b 1 ? 2 b 4 c 1 c 6 a 5 b
4 3 1 ? 1 c 1 ? 4 ? 2 c 3 b
5 4 1 a 2 b 1 c 3 ? 4 a 3 ? 4 ? 2 ?
3 2 1 a 1 a 3 a 3 c
5 4 1 c 1 b 1 ? 2 b 3 b 5 ? 2 b 5 a
3 2 1 b 1 b 2 b 2 ?
2 4 1 b 2 c 2 b 2 b 2 b
5 3 1 b 1 c 1 ? 2 a 4 ? 2 c 5 a
4 5 1 ? 2 c 2 b 4 ? 3 b 4 c 4 b 3 b
4 3 1 b 1 a 2 b 2 a 2 a 3 c
6 3 1 ? 2 c 1 b 3 ? 5 c 4 c 6 c 4 c
4 3 1 a 1 c 3 b 3 ? 3 ? 2 ?
4 5 1 ? 1 a 2 ? 4 c 2 b 3 ? 2 a 3 b
2 5 1 a 2 ? 2 a 2 c 2 ? 2 c
2 5 1 b 2 b 2 c 2 a 2 a 2 c
5 3 1 a 1 c 3 ? 3 b 5 c 2 ? 3 ?
3 3 1 b 1 b 2 ? 2 c 3 a
5 2 1 c 1 ? 1 ? 3 ? 4 a 2 a
5 3 1 a 2 ? 2 a 3 ? 3 c 3 ? 2 a
6 1 1 a 1 b 2 ? 3 c 4 b 3 c
4 5 1 c 2 c 3 b 3 a 4 ? 3 b 3 b 2 a
5 1 1 b 2 b 1 c 4 b 5 a
3 4 1 ? 1 b 3 c 3 c 2 c 2 b
3 3 1 b 1 b 3 b 2 c 2 b
6 4 1 b 2 ? 2 a 1 ? 4 a 6 b 4 ? 5 b 5 ?
4 4 1 b 1 ? 2 b 4 c 2 c 4 b 3 ?
6 2 1 c 2 a 3 c 2 a 5 c 6 c 6 c
5 3 1 ? 2 ? 2 c 2 b 4 b 5 a 4 b
4 4 1 ? 2 b 2 c 2 a 4 ? 4 b 2 b
2 1 1 ? 2 ?
4 1 1 c 2 b 2 a 2 a
5 3 1 b 2 c 3 b 3 b 3 b 4 c 4 a
2 2 1 b 2 c 2 c`

var testcases = mustParseTestcases(testcasesRaw)

func mustParseTestcases(raw string) []struct {
	n   int
	q   int
	par []int
	ch  []byte
	qs  []query
} {
	lines := strings.Split(strings.TrimSpace(raw), "\n")
	start := 0
	firstFields := strings.Fields(lines[0])
	if len(firstFields) == 1 {
		start = 1
	}
	res := make([]struct {
		n   int
		q   int
		par []int
		ch  []byte
		qs  []query
	}, 0, len(lines)-start)
	for idx := start; idx < len(lines); idx++ {
		line := lines[idx]
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			panic(fmt.Sprintf("line %d: too short", idx+1))
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			panic(fmt.Sprintf("line %d: bad n: %v", idx+1, err))
		}
		q, err := strconv.Atoi(fields[1])
		if err != nil {
			panic(fmt.Sprintf("line %d: bad q: %v", idx+1, err))
		}
		expectedLen := 2 + 2*(n-1) + 2*q
		if len(fields) != expectedLen {
			panic(fmt.Sprintf("line %d: expected %d fields got %d", idx+1, expectedLen, len(fields)))
		}
		parents := make([]int, n+1)
		chars := make([]byte, n+1)
		pos := 2
		for v := 2; v <= n; v++ {
			p, _ := strconv.Atoi(fields[pos])
			parents[v] = p
			pos++
			cf := fields[pos]
			if len(cf) != 1 {
				panic(fmt.Sprintf("line %d: invalid char %q", idx+1, cf))
			}
			chars[v] = cf[0]
			pos++
		}
		qs := make([]query, q)
		for i := 0; i < q; i++ {
			v, _ := strconv.Atoi(fields[pos])
			pos++
			cf := fields[pos]
			if len(cf) != 1 {
				panic(fmt.Sprintf("line %d: invalid query char %q", idx+1, cf))
			}
			pos++
			qs[i] = query{v: v, c: cf[0]}
		}
		res = append(res, struct {
			n   int
			q   int
			par []int
			ch  []byte
			qs  []query
		}{n: n, q: q, par: parents, ch: chars, qs: qs})
	}
	return res
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
	return out.String(), nil
}

func parseCandidateOutput(out string) []string {
	lines := []string{}
	scanner := bufio.NewScanner(strings.NewReader(out))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			lines = append(lines, line)
		}
	}
	return lines
}

func buildInput(tc struct {
	n   int
	q   int
	par []int
	ch  []byte
	qs  []query
}) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
	for v := 2; v <= tc.n; v++ {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.par[v], int(tc.ch[v])))
	}
	for _, q := range tc.qs {
		sb.WriteString(fmt.Sprintf("%d %d\n", q.v, int(q.c)))
	}
	return sb.String()
}

func checkCase(bin string, idx int, tc struct {
	n   int
	q   int
	par []int
	ch  []byte
	qs  []query
}) error {
	input := buildInput(tc)
	expected := solveCase(tc.n, tc.par, tc.ch, tc.qs)
	out, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	got := parseCandidateOutput(out)
	if len(got) != len(expected) {
		return fmt.Errorf("expected %d outputs, got %d", len(expected), len(got))
	}
	for i := range expected {
		if got[i] != expected[i] {
			return fmt.Errorf("output %d: expected %s got %s", i+1, expected[i], got[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	for i, tc := range testcases {
		if err := checkCase(bin, i, tc); err != nil {
			fmt.Printf("case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
