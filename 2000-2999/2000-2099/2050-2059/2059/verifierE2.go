package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	sumLimit = 300000
)

type testCase struct {
	n, m   int
	start  [][]int
	target [][]int
}

type operation struct {
	idx int
	val int
}

// ---------- embedded solver (from cf_t23_2059_E2.go) ----------

func solveAllCases(input string) string {
	scanner := bufio.NewScanner(strings.NewReader(input))
	scanner.Split(bufio.ScanWords)
	buf := make([]byte, 1024*1024)
	scanner.Buffer(buf, 10*1024*1024)

	readInt := func() int {
		scanner.Scan()
		val, _ := strconv.Atoi(scanner.Text())
		return val
	}

	t := readInt()
	out := &bytes.Buffer{}

	for tc := 0; tc < t; tc++ {
		n := readInt()
		m := readInt()

		A := make([]int, n*m)
		for i := 0; i < n*m; i++ {
			A[i] = readInt()
		}

		B := make([]int, n*m)
		posB := make([]int, 2*n*m+1)
		for i := range posB {
			posB[i] = -1
		}
		for i := 0; i < n*m; i++ {
			B[i] = readInt()
			posB[B[i]] = i
		}

		L := 0
		for L < n*m {
			p := posB[A[L]]
			if p != -1 {
				if L == 0 || p > posB[A[L-1]] {
					L++
				} else {
					break
				}
			} else {
				break
			}
		}

		check := func(K int) bool {
			isRet := make([]bool, n*m)
			for i := 0; i < K; i++ {
				isRet[posB[A[i]]] = true
			}
			c := 0
			for j := 0; j < n*m; j++ {
				if !isRet[j] {
					if c > (j/m)*m {
						return false
					}
				} else {
					c++
				}
			}
			return true
		}

		ansK := 0
		low := 0
		high := L
		for low <= high {
			mid := (low + high) / 2
			if check(mid) {
				ansK = mid
				low = mid + 1
			} else {
				high = mid - 1
			}
		}

		isRet := make([]bool, n*m)
		for i := 0; i < ansK; i++ {
			isRet[posB[A[i]]] = true
		}

		Ij := make([]int, 0, n*m-ansK)
		for j := 0; j < n*m; j++ {
			if !isRet[j] {
				Ij = append(Ij, j)
			}
		}

		numI := len(Ij)
		fmt.Fprintln(out, numI)

		if numI > 0 {
			type Op struct {
				i int
				x int
			}
			ansOps := make([]Op, numI)

			tree := make([]int, numI+1)
			for i := 1; i <= numI; i++ {
				tree[i] += 1
				nxt := i + (i & -i)
				if nxt <= numI {
					tree[nxt] += tree[i]
				}
			}

			for i := numI - 1; i >= 0; i-- {
				j := Ij[i]
				pi := i - (j % m)

				idx := 0
				rem := pi
				for step := 19; step >= 0; step-- {
					next := idx + (1 << step)
					if next <= numI && tree[next] <= rem {
						idx = next
						rem -= tree[next]
					}
				}
				pos := idx + 1

				for k := pos; k <= numI; k += k & -k {
					tree[k] -= 1
				}

				ansOps[pos-1] = Op{
					i: j/m + 1,
					x: B[j],
				}
			}

			for i := 0; i < numI; i++ {
				fmt.Fprintln(out, ansOps[i].i, ansOps[i].x)
			}
		}
	}
	return out.String()
}

// ---------- verifier infrastructure ----------

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE2.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	tests := generateTests()
	input := buildInput(tests)

	refOut := solveAllCases(input)

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refCounts, err := parseCounts(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	candOps, candCounts, err := parseCandidate(candOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if candCounts[i] != refCounts[i] {
			fmt.Fprintf(os.Stderr, "test %d: expected %d operations, got %d\n", i+1, refCounts[i], candCounts[i])
			os.Exit(1)
		}
		if err := simulate(tests[i], candOps[i]); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	switch filepath.Ext(path) {
	case ".go":
		return exec.Command("go", "run", path)
	case ".py":
		return exec.Command("python3", path)
	default:
		return exec.Command(path)
	}
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("%v\nstderr: %s", err, stderr.String())
	}
	return stdout.String(), nil
}

func parseCounts(output string, t int) ([]int, error) {
	tokens := strings.Fields(output)
	idx := 0
	counts := make([]int, t)
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		if idx >= len(tokens) {
			return nil, fmt.Errorf("missing number of operations for test %d", caseIdx+1)
		}
		q, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, fmt.Errorf("token %q is not integer", tokens[idx])
		}
		if q < 0 {
			return nil, fmt.Errorf("test %d: negative operation count", caseIdx+1)
		}
		idx++
		need := 2 * q
		if idx+need > len(tokens) {
			return nil, fmt.Errorf("test %d: insufficient tokens for operations", caseIdx+1)
		}
		for j := 0; j < need; j++ {
			if _, err := strconv.Atoi(tokens[idx]); err != nil {
				return nil, fmt.Errorf("test %d: invalid operation token %q", caseIdx+1, tokens[idx])
			}
			idx++
		}
		counts[caseIdx] = q
	}
	if idx != len(tokens) {
		return nil, fmt.Errorf("extra output detected starting at token %q", tokens[idx])
	}
	return counts, nil
}

func parseCandidate(output string, tests []testCase) ([][]operation, []int, error) {
	tokens := strings.Fields(output)
	idx := 0
	opsPerTest := make([][]operation, len(tests))
	counts := make([]int, len(tests))
	for caseIdx, tc := range tests {
		if idx >= len(tokens) {
			return nil, nil, fmt.Errorf("missing operation count for test %d", caseIdx+1)
		}
		q, err := strconv.Atoi(tokens[idx])
		if err != nil {
			return nil, nil, fmt.Errorf("test %d: token %q is not integer", caseIdx+1, tokens[idx])
		}
		if q < 0 {
			return nil, nil, fmt.Errorf("test %d: negative operation count", caseIdx+1)
		}
		idx++
		if idx+2*q > len(tokens) {
			return nil, nil, fmt.Errorf("test %d: insufficient tokens for %d operations", caseIdx+1, q)
		}
		ops := make([]operation, q)
		for j := 0; j < q; j++ {
			iVal, err1 := strconv.Atoi(tokens[idx])
			xVal, err2 := strconv.Atoi(tokens[idx+1])
			if err1 != nil || err2 != nil {
				return nil, nil, fmt.Errorf("test %d: invalid operation tokens %q %q", caseIdx+1, tokens[idx], tokens[idx+1])
			}
			idx += 2
			ops[j] = operation{idx: iVal, val: xVal}
		}
		opsPerTest[caseIdx] = ops
		counts[caseIdx] = q
		if err := validateOperationValues(tc, ops); err != nil {
			return nil, nil, fmt.Errorf("test %d: %v", caseIdx+1, err)
		}
	}
	if idx != len(tokens) {
		return nil, nil, fmt.Errorf("extra output detected starting at token %q", tokens[idx])
	}
	return opsPerTest, counts, nil
}

func validateOperationValues(tc testCase, ops []operation) error {
	limit := 2 * tc.n * tc.m
	for idx, op := range ops {
		if op.idx < 1 || op.idx > tc.n {
			return fmt.Errorf("operation %d has invalid index %d", idx+1, op.idx)
		}
		if op.val < 1 || op.val > limit {
			return fmt.Errorf("operation %d has invalid value %d", idx+1, op.val)
		}
	}
	return nil
}

func simulate(tc testCase, ops []operation) error {
	arr := copyMatrix(tc.start)
	n, m := tc.n, tc.m
	for _, op := range ops {
		x := op.val
		for k := op.idx - 1; k < n; k++ {
			row := arr[k]
			last := row[m-1]
			for pos := m - 1; pos > 0; pos-- {
				row[pos] = row[pos-1]
			}
			row[0] = x
			x = last
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if arr[i][j] != tc.target[i][j] {
				return fmt.Errorf("final configuration mismatch at row %d column %d", i+1, j+1)
			}
		}
	}
	return nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(2059))
	var tests []testCase
	total := 0
	addCase := func(tc testCase) {
		need := tc.n * tc.m
		if need == 0 {
			return
		}
		if total+need > sumLimit {
			return
		}
		tests = append(tests, tc)
		total += need
	}

	addCase(randomCase(rng, 1, 1, 0))
	addCase(randomCase(rng, 1, 3, 2))
	addCase(randomCase(rng, 2, 2, 3))
	addCase(randomCase(rng, 2, 4, 5))
	addCase(randomCase(rng, 3, 3, 6))

	for total < sumLimit && len(tests) < 80 {
		remain := sumLimit - total
		maxN := min(8, remain)
		if maxN == 0 {
			break
		}
		n := rng.Intn(maxN) + 1
		maxM := remain / n
		if maxM == 0 {
			n = 1
			maxM = remain
		}
		m := rng.Intn(min(80, maxM)) + 1
		opCount := rng.Intn(2*m + 3)
		addCase(randomCase(rng, n, m, opCount))
	}

	for total < sumLimit {
		remain := sumLimit - total
		if remain < 1000 {
			break
		}
		n := min(50, remain/10)
		if n == 0 {
			break
		}
		m := min(100, remain/n)
		if m == 0 {
			break
		}
		opCount := rng.Intn(n + m + 5)
		addCase(randomCase(rng, n, m, opCount))
	}

	return tests
}

func randomCase(rng *rand.Rand, n, m, operations int) testCase {
	if n < 1 {
		n = 1
	}
	if m < 1 {
		m = 1
	}
	total := n * m
	start := make([][]int, n)
	val := 1
	for i := 0; i < n; i++ {
		row := make([]int, m)
		for j := 0; j < m; j++ {
			row[j] = val
			val++
		}
		start[i] = row
	}
	target := copyMatrix(start)
	pool := make([]int, 0, total)
	for v := total + 1; v <= 2*total; v++ {
		pool = append(pool, v)
	}
	for op := 0; op < operations; op++ {
		if len(pool) == 0 {
			break
		}
		i := rng.Intn(n) + 1
		idx := rng.Intn(len(pool))
		x := pool[idx]
		pool[idx] = pool[len(pool)-1]
		pool = pool[:len(pool)-1]
		cur := x
		for k := i - 1; k < n; k++ {
			row := target[k]
			last := row[m-1]
			for pos := m - 1; pos > 0; pos-- {
				row[pos] = row[pos-1]
			}
			row[0] = cur
			cur = last
		}
		pool = append(pool, cur)
	}
	return testCase{
		n:      n,
		m:      m,
		start:  copyMatrix(start),
		target: target,
	}
}

func copyMatrix(src [][]int) [][]int {
	dst := make([][]int, len(src))
	for i, row := range src {
		tmp := make([]int, len(row))
		copy(tmp, row)
		dst[i] = tmp
	}
	return dst
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
		for i := 0; i < tc.n; i++ {
			writeRow(&b, tc.start[i])
		}
		for i := 0; i < tc.n; i++ {
			writeRow(&b, tc.target[i])
		}
	}
	return b.String()
}

func writeRow(b *strings.Builder, row []int) {
	for i, v := range row {
		if i > 0 {
			b.WriteByte(' ')
		}
		fmt.Fprintf(b, "%d", v)
	}
	b.WriteByte('\n')
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
