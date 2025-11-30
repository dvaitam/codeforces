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

// Embedded copy of testcasesC.txt so the verifier is self-contained.
const testcasesRaw = `
4 2 2 4 0 4 0 2
3 4 2 2 0 4 0 1 2
1 2 4 2 1
2 3 5 5 5 2 2
4 3 3 2 5 1 0 2 4
1 1 3 3
4 1 3 3 3 3 0
1 1 1 0
2 4 1 3 4 0 3 4
4 1 1 1 3 1 1
3 3 2 3 0 4 2 4
2 3 3 4 4 3 4
3 3 1 0 0 4 5 0
2 4 1 1 2 5 5 0
4 1 0 3 5 2 0
3 2 5 5 5 4 5
3 2 5 1 0 4 2
3 2 2 5 3 2 4
2 2 0 4 4 2
3 1 1 3 1 1
1 2 1 3 4
2 2 5 1 1 3
3 2 5 3 0 4 0
3 4 3 2 0 1 4 5 1
4 4 4 2 5 0 2 1 4 3
2 3 2 3 0 1 0
3 2 2 3 5 5 5
2 3 2 5 1 2 0
3 1 5 5 1 2
1 4 5 0 0 1 0
1 2 5 2 0
1 3 5 3 1 1
4 4 1 2 2 1 5 2 5 5
4 4 0 3 2 4 4 5 5 5
4 1 4 0 3 3 1
1 2 4 5 4
2 1 2 1 1
2 1 1 5 5
2 1 3 4 0
4 2 4 4 0 2 2 5
4 1 5 0 3 0 2
3 2 3 1 4 2 1
4 3 2 3 3 0 2 4 2
4 1 4 4 4 2 5
1 4 3 5 0 3 2
4 1 0 2 5 0 2
3 2 4 4 2 3 2
2 1 4 2 1
3 2 1 2 5 0 4
1 2 2 2 2
3 3 3 3 4 3 1 0
2 1 3 1 2
1 4 5 5 2 5 4
2 1 4 3 2
2 2 0 5 5 1
1 4 3 2 2 2 0
4 3 2 4 4 4 2 2 1
4 2 0 4 1 5 5 4
4 3 3 0 4 3 5 4 1
1 3 4 1 5 1
3 4 0 5 5 1 4 1 2
2 4 0 4 3 1 5 5
4 1 1 1 3 0 3
4 3 2 3 2 4 2 0 2
1 4 0 2 1 3 3
4 4 5 0 4 3 2 4 1 4
3 3 0 3 3 4 1 0
1 3 2 0 0 1
1 4 0 2 0 2 3
2 3 3 5 1 4 1
4 3 4 0 1 1 1 3 5
1 1 4 5
4 2 2 4 0 0 4 1
3 2 2 2 5 5 2
3 4 1 3 4 1 0 5 4
2 1 2 0 1
4 2 3 4 1 5 2 5
2 4 4 0 4 0 4 2
1 1 0 3
3 4 2 3 4 2 5 0 1
3 1 1 5 1 0
3 2 0 3 5 0 5
3 4 0 4 5 1 2 0 3
1 4 2 4 4 2 5
3 4 3 1 0 3 2 1 3
1 3 0 0 0 2
1 2 3 4 5
1 3 3 4 5 5
4 4 0 0 4 3 2 0 5 4
1 1 2 2
1 4 0 1 4 5 2
1 1 3 2
2 2 1 1 1 5
2 3 0 3 5 5 3
1 2 1 5 2
3 2 1 2 1 1 3
4 3 4 1 3 4 0 1 4
1 4 5 1 1 0 0
3 1 0 0 0 4
2 2 5 3 0 3
`

type testCase struct {
	n int
	m int
	a []int
	b []int
}

type segTree struct {
	n    int
	min  []int
	lazy []int
}

func buildSeg(arr []int) *segTree {
	st := &segTree{n: len(arr), min: make([]int, 4*len(arr)), lazy: make([]int, 4*len(arr))}
	var build func(int, int, int)
	build = func(idx, l, r int) {
		if l+1 == r {
			st.min[idx] = arr[l]
			return
		}
		m := (l + r) / 2
		build(idx*2, l, m)
		build(idx*2+1, m, r)
		if st.min[idx*2] < st.min[idx*2+1] {
			st.min[idx] = st.min[idx*2]
		} else {
			st.min[idx] = st.min[idx*2+1]
		}
	}
	build(1, 0, st.n)
	return st
}

func (st *segTree) push(idx int) {
	if st.lazy[idx] != 0 {
		for _, c := range []int{idx * 2, idx*2 + 1} {
			st.min[c] += st.lazy[idx]
			st.lazy[c] += st.lazy[idx]
		}
		st.lazy[idx] = 0
	}
}

func (st *segTree) rangeAdd(idx, l, r, ql, qr, val int) {
	if qr <= l || r <= ql {
		return
	}
	if ql <= l && r <= qr {
		st.min[idx] += val
		st.lazy[idx] += val
		return
	}
	st.push(idx)
	m := (l + r) / 2
	st.rangeAdd(idx*2, l, m, ql, qr, val)
	st.rangeAdd(idx*2+1, m, r, ql, qr, val)
	if st.min[idx*2] < st.min[idx*2+1] {
		st.min[idx] = st.min[idx*2]
	} else {
		st.min[idx] = st.min[idx*2+1]
	}
}

func (st *segTree) addSuffix(pos int, val int) {
	st.rangeAdd(1, 0, st.n, pos, st.n, val)
}

func (st *segTree) queryMin() int {
	return st.min[1]
}

type fenwick struct {
	n   int
	bit []int
}

func newFenwick(n int) *fenwick {
	return &fenwick{n: n, bit: make([]int, n+2)}
}

func (f *fenwick) add(idx, val int) {
	for idx <= f.n {
		f.bit[idx] += val
		idx += idx & -idx
	}
}

func (f *fenwick) sum(idx int) int {
	s := 0
	for idx > 0 {
		s += f.bit[idx]
		idx -= idx & -idx
	}
	return s
}

func solveCase(tc testCase) int {
	n := tc.n
	a := append([]int(nil), tc.a...)
	b := append([]int(nil), tc.b...)

	vals := make([]int, n)
	copy(vals, a)
	sort.Ints(vals)
	comp := make(map[int]int, len(vals))
	id := 1
	for _, v := range vals {
		if _, ok := comp[v]; !ok {
			comp[v] = id
			id++
		}
	}
	fw := newFenwick(len(comp) + 2)
	invA := 0
	for i := n - 1; i >= 0; i-- {
		c := comp[a[i]]
		invA += fw.sum(c - 1)
		fw.add(c, 1)
	}

	type pair struct{ val, idx int }
	pairs := make([]pair, n)
	for i, v := range a {
		pairs[i] = pair{v, i + 1}
	}
	sort.Slice(pairs, func(i, j int) bool { return pairs[i].val < pairs[j].val })

	arr := make([]int, n+1)
	for i := range arr {
		arr[i] = i
	}
	st := buildSeg(arr)

	lessPtr, eqPtr := 0, 0
	sort.Ints(b)
	ans := invA
	for _, x := range b {
		for lessPtr < n && pairs[lessPtr].val < x {
			st.addSuffix(pairs[lessPtr].idx, -1)
			lessPtr++
		}
		for eqPtr < n && pairs[eqPtr].val <= x {
			st.addSuffix(pairs[eqPtr].idx, -1)
			eqPtr++
		}
		ans += st.queryMin() + lessPtr
	}
	return ans
}

func parseTestcases() ([]testCase, error) {
	lines := strings.Split(strings.TrimSpace(testcasesRaw), "\n")
	cases := make([]testCase, 0, len(lines))
	for idx, line := range lines {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			return nil, fmt.Errorf("line %d: not enough numbers", idx+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse n: %v", idx+1, err)
		}
		m, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, fmt.Errorf("line %d: parse m: %v", idx+1, err)
		}
		expected := 2 + n + m
		if len(fields) != expected {
			return nil, fmt.Errorf("line %d: expected %d numbers got %d", idx+1, expected, len(fields))
		}
		tc := testCase{n: n, m: m, a: make([]int, n), b: make([]int, m)}
		for i := 0; i < n; i++ {
			v, err := strconv.Atoi(fields[2+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse a: %v", idx+1, err)
			}
			tc.a[i] = v
		}
		for i := 0; i < m; i++ {
			v, err := strconv.Atoi(fields[2+n+i])
			if err != nil {
				return nil, fmt.Errorf("line %d: parse b: %v", idx+1, err)
			}
			tc.b[i] = v
		}
		cases = append(cases, tc)
	}
	return cases, nil
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
		fmt.Println("usage: go run verifierC.go /path/to/binary")
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

	for i, tc := range cases {
		expected := solveCase(tc)

		var sb strings.Builder
		sb.WriteString("1\n")
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for idx, v := range tc.a {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for idx, v := range tc.b {
			if idx > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		got, err := runCandidate(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strconv.Itoa(expected) {
			fmt.Printf("case %d failed\nexpected: %d\ngot: %s\n", i+1, expected, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
