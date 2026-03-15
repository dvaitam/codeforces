package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

// Embedded reference solver for 2053F.

type refPair struct {
	c   int
	cnt int
}

type refItem struct {
	val   int64
	color int
}

func refHeapPush(h []refItem, x refItem) []refItem {
	h = append(h, x)
	i := len(h) - 1
	for i > 0 {
		p := (i - 1) >> 1
		if h[p].val >= h[i].val {
			break
		}
		h[p], h[i] = h[i], h[p]
		i = p
	}
	return h
}

func refHeapPop(h []refItem) []refItem {
	n := len(h) - 1
	if n == 0 {
		return h[:0]
	}
	h[0] = h[n]
	h = h[:n]
	i := 0
	for {
		l := i*2 + 1
		if l >= n {
			break
		}
		r := l + 1
		j := l
		if r < n && h[r].val > h[l].val {
			j = r
		}
		if h[i].val >= h[j].val {
			break
		}
		h[i], h[j] = h[j], h[i]
		i = j
	}
	return h
}

func solveCase(n, m, k int, grid [][]int) int64 {
	counts := make([][]refPair, n)
	blanks := make([]int64, n)

	for i := 0; i < n; i++ {
		vals := make([]int, 0, m)
		var b int64
		for j := 0; j < m; j++ {
			x := grid[i][j]
			if x == -1 {
				b++
			} else {
				vals = append(vals, x)
			}
		}
		blanks[i] = b
		if len(vals) == 0 {
			continue
		}
		sort.Ints(vals)
		pairs := make([]refPair, 0, len(vals))
		cur := vals[0]
		cnt := 1
		for j := 1; j < len(vals); j++ {
			if vals[j] == cur {
				cnt++
			} else {
				pairs = append(pairs, refPair{c: cur, cnt: cnt})
				cur = vals[j]
				cnt = 1
			}
		}
		pairs = append(pairs, refPair{c: cur, cnt: cnt})
		counts[i] = pairs
	}

	var base int64
	for i := 0; i+1 < n; i++ {
		a := counts[i]
		bv := counts[i+1]
		p, q := 0, 0
		for p < len(a) && q < len(bv) {
			if a[p].c == bv[q].c {
				base += int64(a[p].cnt) * int64(bv[q].cnt)
				p++
				q++
			} else if a[p].c < bv[q].c {
				p++
			} else {
				q++
			}
		}
	}

	raw := make([]int64, k+1)
	h := make([]refItem, 0)
	var lazy int64
	var d int64

	currentMax := func() int64 {
		for len(h) > 0 {
			top := h[0]
			if raw[top.color] != top.val || top.val <= lazy {
				h = refHeapPop(h)
				continue
			}
			return top.val - lazy
		}
		return 0
	}

	addColor := func(c int, inc int64) {
		cur := raw[c]
		if cur < lazy {
			cur = lazy
		}
		cur += inc
		raw[c] = cur
		h = refHeapPush(h, refItem{val: cur, color: c})
	}

	addUnary := func(i int) {
		r := blanks[i]
		if r == 0 {
			return
		}
		if i > 0 {
			for _, p := range counts[i-1] {
				addColor(p.c, r*int64(p.cnt))
			}
		}
		if i+1 < n {
			for _, p := range counts[i+1] {
				addColor(p.c, r*int64(p.cnt))
			}
		}
	}

	addUnary(0)
	for i := 1; i < n; i++ {
		mx := currentMax()
		bv := blanks[i-1] * blanks[i]
		if mx > bv {
			d += mx
			lazy += mx - bv
		} else {
			d += bv
		}
		addUnary(i)
	}

	return base + d + currentMax()
}

type matrixTest struct {
	n, m, k int
	data    [][]int
}

type testCase struct {
	input    string
	caseCnt  int
	matrices []matrixTest
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := generateTests()
	for idx, tc := range tests {
		// Compute reference answers using embedded solver
		refVals := make([]int64, tc.caseCnt)
		for i, mt := range tc.matrices {
			refVals[i] = solveCase(mt.n, mt.m, mt.k, mt.data)
		}

		gotOut, err := runProgram(bin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVals, err := parseOutputs(gotOut, tc.caseCnt)
		if err != nil {
			fmt.Fprintf(os.Stderr, "participant output invalid on test %d: %v\noutput:\n%s\n", idx+1, err, gotOut)
			os.Exit(1)
		}

		for caseIdx := 0; caseIdx < tc.caseCnt; caseIdx++ {
			if refVals[caseIdx] != gotVals[caseIdx] {
				fmt.Fprintf(os.Stderr, "test %d case %d mismatch: expected %d got %d\ninput:\n%sparticipant output:\n%s\n",
					idx+1, caseIdx+1, refVals[caseIdx], gotVals[caseIdx], tc.input, gotOut)
				os.Exit(1)
			}
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	switch {
	case strings.HasSuffix(path, ".go"):
		cmd = exec.Command("go", "run", path)
	default:
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseOutputs(out string, expected int) ([]int64, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d integers got %d", expected, len(fields))
	}
	vals := make([]int64, expected)
	for i, token := range fields {
		v, err := strconv.ParseInt(token, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", token)
		}
		vals[i] = v
	}
	return vals, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTests()...)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(rng, 60, 8)...)
	tests = append(tests, randomTests(rng, 60, 60)...)
	tests = append(tests, randomTests(rng, 40, 200)...)
	tests = append(tests, stressTests()...)
	return tests
}

func manualTests() []testCase {
	cases := []matrixTest{
		{
			n: 3, m: 3, k: 3,
			data: [][]int{
				{1, 2, 2},
				{3, 1, 3},
				{3, 2, 1},
			},
		},
		{
			n: 2, m: 3, k: 3,
			data: [][]int{
				{-1, 3, 3},
				{2, 2, -1},
			},
		},
	}
	return []testCase{makeTestCase(cases)}
}

func randomTests(rng *rand.Rand, batches int, maxNM int) []testCase {
	const limit = 600000
	var tests []testCase
	for b := 0; b < batches; b++ {
		caseCnt := rng.Intn(3) + 1
		var cases []matrixTest
		sumNM := 0
		for len(cases) < caseCnt {
			n := rng.Intn(maxNM-1) + 2
			m := rng.Intn(maxNM-1) + 2
			if n*m > maxNM || sumNM+n*m > limit {
				break
			}
			k := rng.Intn(n*m) + 1
			mat := make([][]int, n)
			for i := 0; i < n; i++ {
				row := make([]int, m)
				for j := 0; j < m; j++ {
					if rng.Intn(5) == 0 {
						row[j] = -1
					} else {
						row[j] = rng.Intn(k) + 1
					}
				}
				mat[i] = row
			}
			cases = append(cases, matrixTest{n: n, m: m, k: k, data: mat})
			sumNM += n * m
		}
		if len(cases) == 0 {
			cases = append(cases, matrixTest{
				n: 2, m: 2, k: 4,
				data: [][]int{
					{rng.Intn(4) + 1, -1},
					{-1, rng.Intn(4) + 1},
				},
			})
		}
		tests = append(tests, makeTestCase(cases))
	}
	return tests
}

func stressTests() []testCase {
	return []testCase{
		makeTestCase([]matrixTest{
			fullMatrix(300, 2, 10),
		}),
		makeTestCase([]matrixTest{
			stripedMatrix(200, 3, 20),
		}),
	}
}

func fullMatrix(n, m, k int) matrixTest {
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			row[j] = (i+j)%k + 1
		}
		mat[i] = row
	}
	return matrixTest{n: n, m: m, k: k, data: mat}
}

func stripedMatrix(n, m, k int) matrixTest {
	mat := make([][]int, n)
	for i := 0; i < n; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			if j%2 == 0 {
				row[j] = -1
			} else {
				row[j] = (i + j) % k
				if row[j] == 0 {
					row[j] = k
				}
			}
		}
		mat[i] = row
	}
	return matrixTest{n: n, m: m, k: k, data: mat}
}

func makeTestCase(cases []matrixTest) testCase {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(cases)))
	sb.WriteByte('\n')
	for _, c := range cases {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", c.n, c.m, c.k))
		for i := 0; i < c.n; i++ {
			for j := 0; j < c.m; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				sb.WriteString(strconv.Itoa(c.data[i][j]))
			}
			sb.WriteByte('\n')
		}
	}
	return testCase{input: sb.String(), caseCnt: len(cases), matrices: cases}
}
