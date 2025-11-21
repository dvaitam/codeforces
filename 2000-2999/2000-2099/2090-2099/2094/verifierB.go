package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type testCase struct {
	n int
	m int
	l int
	r int
}

// We still build the reference to ensure the inputs are sound, even though any valid
// output is accepted (multiple solutions).
func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2094B-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleB")
	cmd := exec.Command("go", "build", "-o", outPath, "2094B.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runBinary(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
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

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.Grow(len(tests) * 64)
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d %d\n", tc.n, tc.m, tc.l, tc.r))
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([][2]int, error) {
	res := make([][2]int, expected)
	r := bufio.NewReader(strings.NewReader(out))
	for i := 0; i < expected; i++ {
		if _, err := fmt.Fscan(r, &res[i][0], &res[i][1]); err != nil {
			return nil, fmt.Errorf("cannot read answer %d: %v", i+1, err)
		}
	}
	// ensure no extra tokens aside from whitespace
	var extra string
	if _, err := fmt.Fscan(r, &extra); err == nil {
		return nil, fmt.Errorf("extra output detected after %d answers", expected)
	}
	return res, nil
}

func validSegment(tc testCase, l2, r2 int) error {
	if r2-l2+1 != tc.m+1 {
		return fmt.Errorf("length mismatch (got [%d,%d])", l2, r2)
	}
	if l2 > 0 || r2 < 0 {
		return fmt.Errorf("segment must contain 0, got [%d,%d]", l2, r2)
	}
	if l2 < tc.l || r2 > tc.r {
		return fmt.Errorf("segment outside final bounds [%d,%d]", l2, r2)
	}

	L := -tc.l
	R := tc.r
	Lm := -l2
	Rm := r2
	if Lm < 0 || Rm < 0 {
		return fmt.Errorf("negative counts")
	}
	if Lm+Rm != tc.m {
		return fmt.Errorf("step count mismatch: left+right=%d, expected %d", Lm+Rm, tc.m)
	}
	if Lm > L || Rm > R {
		return fmt.Errorf("used more expansions than available (L:%d/%d, R:%d/%d)", Lm, L, Rm, R)
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 4, m: 2, l: -2, r: 2},
		{n: 4, m: 1, l: 0, r: 4},
		{n: 3, m: 3, l: -1, r: 2},
		{n: 9, m: 8, l: -6, r: 3},
		{n: 1, m: 1, l: -1, r: 0},
		{n: 1, m: 1, l: 0, r: 1},
		{n: 6, m: 3, l: -5, r: 1},  // heavily to the left
		{n: 6, m: 3, l: -1, r: 5},  // heavily to the right
		{n: 10, m: 5, l: -3, r: 7}, // mixed
	}
}

func randomTests(count int) []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(2000) + 1
		L := rng.Intn(n + 1)
		R := n - L
		l := -L
		r := R
		m := rng.Intn(n) + 1
		tests[i] = testCase{n: n, m: m, l: l, r: r}
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	tests = append(tests, randomTests(60)...)

	input := buildInput(tests)

	// Ensure the reference solution handles the input (sanity check).
	if _, err := runBinary(oracle, input); err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	actOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\ninput:\n%s", err, input)
		os.Exit(1)
	}

	answers, err := parseOutput(actOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "invalid output: %v\noutput:\n%s", err, actOut)
		os.Exit(1)
	}

	for i, tc := range tests {
		l2, r2 := answers[i][0], answers[i][1]
		if err := validSegment(tc, l2, r2); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v (n=%d m=%d l=%d r=%d output=%d %d)\ninput:\n%s",
				i+1, err, tc.n, tc.m, tc.l, tc.r, l2, r2, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
