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
6 17 48 23 45 42 34 -5 2 -2 5 -5 -3
2 24 31 -2 1
6 7 37 16 1 14 27 -1 -3 1 -3 -4 -3
6 40 29 9 47 1 45 -2 -2 -3 -3 -1 0
3 35 44 41 -2 -3 -2
5 20 2 24 27 11 -3 -1 -4 0 -1
6 38 1 39 44 22 5 -1 0 -1 2 0 -3
5 31 46 12 4 17 -5 0 1 -5 3
5 24 25 38 1 29 -5 -3 4 -2 -4
3 30 23 33 0 3 -1
5 7 38 48 24 19 -5 1 -4 -2 0
6 40 24 10 22 18 35 -4 -1 5 0 -1 -3
2 41 10 -1 2
3 47 4 6 4 3 1
2 16 48 4 0
4 30 42 27 10 -5 5 -5 2
4 14 9 47 37 -3 5 1 -4
3 28 24 10 -5 1 -1
3 30 40 11 3 2 2
7 47 21 31 18 19 49 26 -3 -4 1 3 -3 5 2
4 12 6 32 18 3 3 3 0
2 23 45 4 5
2 49 20 0 3
7 43 18 32 17 45 19 22 5 -3 4 -5 2 3 -1
4 43 18 30 19 3 5 5 0
4 18 42 23 48 1 0 -3 2
4 22 34 10 11 -2 0 2 -1
7 6 47 43 27 11 40 38 3 5 1 -1 4 3 5
4 47 2 13 11 4 2 4 5
3 15 49 44 -3 5 -5
5 15 11 4 9 8 0 -3 2 -2 3
2 27 30 0 1
7 40 5 38 14 16 24 1 0 1 -1 1 -4 3 0
2 36 40 -1 -4
4 35 33 22 38 -1 0 -3 1
5 37 42 35 24 30 -3 -3 4 1 4
5 13 9 39 6 23 5 -5 1 -4 0
6 40 35 10 21 41 37 1 1 1 -2 2 -1
5 46 25 11 39 17 -1 2 -1 1 -5
4 20 32 19 10 2 -5 -4 5
6 29 16 19 3 9 26 -5 2 3 3 -1 -2
5 3 16 32 18 10 -1 -1 2 4 2
6 42 39 8 2 9 20 -1 3 0 4 -1 3
2 30 23 0 5
7 38 9 3 1 17 36 30 5 -4 5 3 -2 -5 1
5 39 37 45 46 41 5 2 1 2 1
7 47 13 19 30 5 20 1 1 4 -1 5 2 -1 -3
3 31 45 36 2 0 3
3 28 38 35 -5 -4 -2
4 6 5 43 2 0 1 -4 1
7 32 4 8 47 15 40 42 -4 -3 -1 2 -3 -3 4
3 27 11 5 4 -2 -5
6 7 43 25 5 18 4 4 4 -4 1 4 -3
2 28 6 0 5
6 32 49 23 42 24 4 -3 -1 -3 4 5 5
6 19 36 48 40 15 17 -4 3 -2 -1 -1 3
3 16 24 30 1 -3 -3
7 2 42 22 6 37 43 3 -4 -4 3 4 2 -2 1
5 31 21 7 34 2 3 1 -5 -3 1
7 15 48 8 6 43 32 14 -3 4 1 0 -2 -1 0
6 46 23 25 47 9 48 5 -1 5 1 0 3
2 38 37 -2 -3
5 5 7 3 12 13 -2 -5 2 2 5
4 1 28 31 20 4 1 0 2
5 7 13 10 42 11 -4 0 1 2 -3
6 17 8 18 11 19 44 -2 -5 2 -5 0 0
4 4 45 2 44 2 2 -3 -4
7 21 19 30 44 16 11 3 -2 -5 4 -2 5 -4 4
5 46 24 44 20 12 2 0 -1 -4 2
3 16 12 44 -2 -5 5
6 28 47 18 1 30 4 2 5 1 -3 -5 -5
6 34 37 23 7 5 16 2 -4 2 -5 5 -2
7 4 32 26 49 46 17 27 2 -1 5 -5 -5 -2 -3
7 33 41 26 13 35 15 6 0 -4 -4 3 -3 4 -4
7 14 40 2 29 35 23 32 1 4 1 3 5 -3 -5
4 45 29 15 36 -4 -3 -1 2
2 17 26 0 -3
3 35 6 31 -2 -4 1
6 43 16 18 4 32 48 -5 -2 -1 0 -4 -4
5 21 25 49 36 1 -1 5 -3 4 -5
5 10 45 42 2 47 -5 1 -1 -4 1
6 9 21 5 20 11 49 -1 2 5 0 -1 -5
2 35 49 3 -5
2 8 4 -3 0
5 17 40 2 15 34 1 -3 0 -2 0
4 1 11 38 40 -3 -4 -1 3
6 11 41 46 10 29 28 1 -4 0 0 2 -2
5 39 28 11 29 4 5 -3 3 2 2
6 1 7 9 27 44 4 -2 -5 -5 4 -2 2
7 23 26 30 14 13 18 45 -1 3 3 0 -2 0 0
3 43 4 40 2 1 0
3 10 16 32 1 -5 -4
3 26 7 49 -1 -2 0
5 44 31 35 40 6 -2 1 -5 -4 3
6 38 40 12 25 2 21 5 -5 4 -1 3 1
6 31 35 46 45 26 29 4 3 5 -5 3 -3
5 24 38 2 44 34 -2 -1 -1 -4 -5
2 23 40 4 1
4 1 18 46 16 -4 -5 -4 5
4 23 45 48 39 4 3 -2 2
`

type testCase struct {
	n  int
	xs []int64
	vs []int
}

type BIT struct {
	n    int
	tree []int64
}

func newBIT(n int) *BIT {
	return &BIT{n: n, tree: make([]int64, n+1)}
}

func (b *BIT) Update(i int, v int64) {
	for i <= b.n {
		b.tree[i] += v
		i += i & -i
	}
}

func (b *BIT) Query(i int) int64 {
	var s int64
	for i > 0 {
		s += b.tree[i]
		i -= i & -i
	}
	return s
}

func solveCase(tc testCase) int64 {
	idx := make([]int, tc.n)
	for i := 0; i < tc.n; i++ {
		idx[i] = i
	}
	sort.Slice(idx, func(i, j int) bool {
		return tc.xs[idx[i]] < tc.xs[idx[j]]
	})
	sortedX := make([]int64, tc.n)
	sortedV := make([]int, tc.n)
	for i, id := range idx {
		sortedX[i] = tc.xs[id]
		sortedV[i] = tc.vs[id]
	}
	uniq := make([]int, tc.n)
	copy(uniq, sortedV)
	sort.Ints(uniq)
	comp := make(map[int]int, len(uniq))
	m := 0
	for _, v := range uniq {
		if _, ok := comp[v]; !ok {
			m++
			comp[v] = m
		}
	}
	bitCnt := newBIT(m)
	bitSum := newBIT(m)
	var res int64
	for i := 0; i < tc.n; i++ {
		ci := comp[sortedV[i]]
		cnt := bitCnt.Query(ci)
		sumX := bitSum.Query(ci)
		res += cnt*sortedX[i] - sumX
		bitCnt.Update(ci, 1)
		bitSum.Update(ci, sortedX[i])
	}
	return res
}

func parseTestcases() ([]testCase, error) {
	rawLines := strings.Split(strings.TrimSpace(testcasesData), "\n")
	lines := make([]string, 0, len(rawLines))
	for _, ln := range rawLines {
		ln = strings.TrimSpace(ln)
		if ln != "" {
			lines = append(lines, ln)
		}
	}
	if len(lines) == 0 {
		return nil, fmt.Errorf("no testcases")
	}
	t := 0
	startIdx := 0
	if fields := strings.Fields(lines[0]); len(fields) == 1 {
		if v, err := strconv.Atoi(fields[0]); err == nil {
			t = v
			startIdx = 1
		}
	}
	if t == 0 {
		t = len(lines)
		startIdx = 0
	}
	if startIdx+t != len(lines) {
		return nil, fmt.Errorf("testcase count mismatch: declared %d actual %d", t, len(lines)-startIdx)
	}

	cases := make([]testCase, 0, t)
	idx := startIdx
	for i := 0; i < t; i++ {
		if idx >= len(lines) {
			return nil, fmt.Errorf("unexpected EOF at case %d", i+1)
		}
		fields := strings.Fields(lines[idx])
		idx++
		if len(fields) < 1 {
			return nil, fmt.Errorf("bad line at case %d", i+1)
		}
		n, err := strconv.Atoi(fields[0])
		if err != nil {
			return nil, err
		}
		if len(fields) != 1+2*n {
			return nil, fmt.Errorf("case %d expected %d numbers got %d", i+1, 1+2*n, len(fields))
		}
		xs := make([]int64, n)
		vs := make([]int, n)
		for j := 0; j < n; j++ {
			v, err := strconv.ParseInt(fields[1+j], 10, 64)
			if err != nil {
				return nil, err
			}
			xs[j] = v
		}
		for j := 0; j < n; j++ {
			v, err := strconv.Atoi(fields[1+n+j])
			if err != nil {
				return nil, err
			}
			vs[j] = v
		}
		cases = append(cases, testCase{n: n, xs: xs, vs: vs})
	}
	return cases, nil
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
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	testcases, err := parseTestcases()
	if err != nil {
		fmt.Println("failed to parse testcases:", err)
		os.Exit(1)
	}

	for idx, tc := range testcases {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		for i, v := range tc.xs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for i, v := range tc.vs {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')

		want := strconv.FormatInt(solveCase(tc), 10)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Printf("case %d failed: expected %s got %s\n", idx+1, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(testcases))
}
