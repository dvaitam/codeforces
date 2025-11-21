package main

import (
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
	d int
	s []int
	p []int
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-1046C-")
	if err != nil {
		return "", nil, err
	}
	path := filepath.Join(tmpDir, "oracleC")
	cmd := exec.Command("go", "build", "-o", path, "1046C.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() { os.RemoveAll(tmpDir) }
	return path, cleanup, nil
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
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.d))
	for i, v := range tc.s {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for i, v := range tc.p {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func parseAnswer(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer output, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	if val < 1 {
		return 0, fmt.Errorf("ranking must be positive, got %d", val)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 1, d: 1, s: []int{5}, p: []int{3}},
		{n: 2, d: 1, s: []int{10, 5}, p: []int{5, 0}},
		{n: 3, d: 2, s: []int{20, 15, 10}, p: []int{10, 5, 0}},
		{n: 5, d: 3, s: []int{30, 25, 20, 10, 0}, p: []int{10, 8, 5, 3, 1}},
	}
}

func randomSortedArray(rng *rand.Rand, n int, maxVal int) []int {
	arr := make([]int, n)
	prev := maxVal
	for i := 0; i < n; i++ {
		if prev < 0 {
			prev = 0
		}
		val := rng.Intn(prev + 1)
		arr[i] = val
		prev = val
	}
	return arr
}

func randomTest(rng *rand.Rand) testCase {
	n := rng.Intn(20) + 1
	if rng.Intn(4) == 0 {
		n = rng.Intn(200000-1) + 1
	}
	d := rng.Intn(n) + 1
	maxScore := 100_000_000
	s := make([]int, n)
	s[0] = maxScore
	for i := 1; i < n; i++ {
		if rng.Intn(3) == 0 {
			s[i] = s[i-1]
		} else {
			s[i] = max(0, s[i-1]-rng.Intn(1000)-1)
		}
	}
	p := make([]int, n)
	p[0] = rng.Intn(maxScore + 1)
	for i := 1; i < n; i++ {
		if rng.Intn(3) == 0 {
			p[i] = p[i-1]
		} else {
			p[i] = max(0, p[i-1]-rng.Intn(1000)-1)
		}
	}
	return testCase{n: n, d: d, s: s, p: p}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
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
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomTest(rng))
	}

	for idx, tc := range tests {
		input := buildInput(tc)

		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		exp, err := parseAnswer(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotOut, input)
			os.Exit(1)
		}

		if got != exp {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\ninput:\n%s\n", idx+1, exp, got, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
