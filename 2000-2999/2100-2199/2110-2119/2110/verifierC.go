package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const refSource = "2000-2999/2100-2199/2110-2119/2110/2110C.go"

type testCase struct {
	n       int
	d       []int // original values, may contain -1/0/1
	l       []int
	r       []int
	comment string
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2110C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func deterministicTests() []testCase {
	// From statement and some custom small cases.
	t1 := testCase{
		n: 4,
		d: []int{0, -1, -1, 1},
		l: []int{0, 1, 2, 1},
		r: []int{4, 2, 4, 4},
	}
	t1.l[0], t1.r[0] = 0, 4

	t2 := testCase{
		n: 1,
		d: []int{-1},
		l: []int{0},
		r: []int{0},
	}

	t3 := testCase{
		n: 1,
		d: []int{1},
		l: []int{1},
		r: []int{1},
	}

	// Impossible case: require height 3 at step 1 (max reachable is 1).
	t4 := testCase{
		n:       3,
		d:       []int{-1, -1, -1},
		l:       []int{3, 0, 0},
		r:       []int{3, 3, 3},
		comment: "impossible high first obstacle",
	}

	// Simple increasing feasible.
	t5 := feasibleFromHeights([]int{0, 1, 2, 3, 3, 3}, []int{-1, -1, -1, -1, -1}, "simple climb")

	return []testCase{t1, t2, t3, t4, t5}
}

func feasibleFromHeights(h []int, dOrig []int, comment string) testCase {
	n := len(dOrig)
	if len(h) != n+1 {
		// Adjust to avoid panics; fallback to flat heights.
		h = make([]int, n+1)
		for i := 1; i <= n; i++ {
			h[i] = h[i-1]
		}
	}
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		l[i] = h[i+1]
		r[i] = h[i+1]
	}
	return testCase{n: n, d: append([]int(nil), dOrig...), l: l, r: r, comment: comment}
}

func randomFeasible(rng *rand.Rand, n int) testCase {
	dTrue := make([]int, n)
	h := make([]int, n+1)
	for i := 0; i < n; i++ {
		dTrue[i] = rng.Intn(2)
		h[i+1] = h[i] + dTrue[i]
	}
	l := make([]int, n)
	r := make([]int, n)
	for i := 0; i < n; i++ {
		leeway := rng.Intn(2)
		low := h[i+1] - leeway
		if low < 0 {
			low = 0
		}
		high := h[i+1] + leeway
		if high > n {
			high = n
		}
		l[i] = low
		r[i] = high
	}
	dOrig := make([]int, n)
	for i := 0; i < n; i++ {
		if rng.Intn(100) < 50 {
			dOrig[i] = -1
		} else {
			dOrig[i] = dTrue[i]
		}
	}
	return testCase{n: n, d: dOrig, l: l, r: r}
}

func randomImpossible(rng *rand.Rand, n int) testCase {
	tc := randomFeasible(rng, n)
	// Make the first obstacle impossible by requiring too high height.
	tc.l[0] = n
	tc.r[0] = n
	for i := 0; i < n; i++ {
		tc.d[i] = -1
	}
	return tc
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.d {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
		for i := 0; i < tc.n; i++ {
			sb.WriteString(fmt.Sprintf("%d %d\n", tc.l[i], tc.r[i]))
		}
	}
	return sb.String()
}

func parseSolutions(out string, tests []testCase) ([][]int, []bool, error) {
	fields := strings.Fields(out)
	pos := 0
	sols := make([][]int, len(tests))
	ok := make([]bool, len(tests))
	for i, tc := range tests {
		if pos >= len(fields) {
			return nil, nil, fmt.Errorf("test %d: missing output", i+1)
		}
		if fields[pos] == "-1" {
			ok[i] = false
			pos++
			continue
		}
		need := tc.n
		if pos+need > len(fields) {
			return nil, nil, fmt.Errorf("test %d: expected %d numbers, got %d", i+1, need, len(fields)-pos)
		}
		arr := make([]int, need)
		for j := 0; j < need; j++ {
			val, err := strconv.Atoi(fields[pos+j])
			if err != nil {
				return nil, nil, fmt.Errorf("test %d: invalid integer %q", i+1, fields[pos+j])
			}
			arr[j] = val
		}
		pos += need
		sols[i] = arr
		ok[i] = true
	}
	if pos != len(fields) {
		return nil, nil, fmt.Errorf("extra output detected after %d testcases", len(tests))
	}
	return sols, ok, nil
}

func validateSolution(tc testCase, sol []int) error {
	if len(sol) != tc.n {
		return fmt.Errorf("length mismatch: expected %d got %d", tc.n, len(sol))
	}
	h := 0
	for i := 0; i < tc.n; i++ {
		if sol[i] != 0 && sol[i] != 1 {
			return fmt.Errorf("d[%d]=%d not in {0,1}", i+1, sol[i])
		}
		if tc.d[i] != -1 && tc.d[i] != sol[i] {
			return fmt.Errorf("d[%d] must be %d", i+1, tc.d[i])
		}
		h += sol[i]
		if h < tc.l[i] || h > tc.r[i] {
			return fmt.Errorf("height %d at step %d outside [%d,%d]", h, i+1, tc.l[i], tc.r[i])
		}
	}
	return nil
}

func totalN(tests []testCase) int {
	s := 0
	for _, tc := range tests {
		s += tc.n
	}
	return s
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for i := 0; i < 30; i++ {
		n := rng.Intn(40) + 1
		tests = append(tests, randomFeasible(rng, n))
	}
	for i := 0; i < 20; i++ {
		n := rng.Intn(200) + 50
		tests = append(tests, randomFeasible(rng, n))
	}
	for i := 0; i < 5; i++ {
		n := rng.Intn(200) + 50
		tests = append(tests, randomImpossible(rng, n))
	}
	// Larger stress cases while respecting total n <= 2e5
	tests = append(tests, randomFeasible(rng, 50000))
	tests = append(tests, randomFeasible(rng, 40000))

	// Trim if needed to keep within limits.
	for totalN(tests) > 190000 {
		tests = tests[:len(tests)-1]
	}

	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	wantSol, wantOk, err := parseSolutions(wantOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	gotSol, gotOk, err := parseSolutions(gotOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		if !wantOk[i] {
			if gotOk[i] {
				fmt.Fprintf(os.Stderr, "test %d: reference says impossible, candidate provided solution\nn=%d\n", i+1, tc.n)
				os.Exit(1)
			}
			continue
		}
		if !gotOk[i] {
			fmt.Fprintf(os.Stderr, "test %d: solution exists but candidate output -1\nn=%d\n", i+1, tc.n)
			os.Exit(1)
		}
		if err := validateSolution(tc, gotSol[i]); err != nil {
			fmt.Fprintf(os.Stderr, "test %d invalid solution: %v\nn=%d\n", i+1, err, tc.n)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
