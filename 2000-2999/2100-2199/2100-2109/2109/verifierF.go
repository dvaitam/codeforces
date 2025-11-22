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
	r int
	k int
	a [][]int
	c []string
}

func callerDir() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("cannot determine caller")
	}
	return filepath.Dir(file), nil
}

func buildOracle() (string, func(), error) {
	dir, err := callerDir()
	if err != nil {
		return "", nil, err
	}
	tmpDir, err := os.MkdirTemp("", "oracle-2109F-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", outPath, "2109F.go")
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

func parseOutput(out string, t int) ([][2]int, error) {
	fields := strings.Fields(out)
	if len(fields) != 2*t {
		return nil, fmt.Errorf("expected %d tokens, got %d", 2*t, len(fields))
	}
	ans := make([][2]int, t)
	for i := 0; i < t; i++ {
		dm, err := strconv.Atoi(fields[2*i])
		if err != nil {
			return nil, fmt.Errorf("token %d invalid int %q: %v", 2*i+1, fields[2*i], err)
		}
		df, err := strconv.Atoi(fields[2*i+1])
		if err != nil {
			return nil, fmt.Errorf("token %d invalid int %q: %v", 2*i+2, fields[2*i+1], err)
		}
		ans[i] = [2]int{dm, df}
	}
	return ans, nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d %d\n", tc.n, tc.r, tc.k)
		for i := 0; i < tc.n; i++ {
			for j := 0; j < tc.n; j++ {
				if j > 0 {
					sb.WriteByte(' ')
				}
				fmt.Fprintf(&sb, "%d", tc.a[i][j])
			}
			sb.WriteByte('\n')
		}
		for i := 0; i < tc.n; i++ {
			sb.WriteString(tc.c[i])
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func deterministicTests() []testCase {
	// Small handcrafted cases including samples.
	return []testCase{
		{
			n: 2, r: 1, k: 30,
			a: [][]int{{2, 2}, {1, 1}},
			c: []string{"10", "11"},
		},
		{
			n: 3, r: 3, k: 5,
			a: [][]int{{9, 2, 2}, {2, 3, 2}, {2, 2, 2}},
			c: []string{"111", "110", "100"},
		},
		{
			n: 2, r: 2, k: 0,
			a: [][]int{{1, 1}, {1, 1}},
			c: []string{"00", "00"},
		},
	}
}

func randomTests() []testCase {
	const limit = 90000
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 40)
	sumSq := 0

	add := func(tc testCase) {
		if sumSq+tc.n*tc.n > limit {
			return
		}
		tests = append(tests, tc)
		sumSq += tc.n * tc.n
	}

	// Variety of small cases.
	for len(tests) < 10 {
		add(randomCase(rng, 2, 8))
	}
	for len(tests) < 20 {
		add(randomCase(rng, 5, 20))
	}
	for len(tests) < 28 {
		add(randomCase(rng, 20, 80))
	}
	if sumSq+300*300 <= limit {
		add(randomCase(rng, 300, 300))
	} else if sumSq+200*200 <= limit {
		add(randomCase(rng, 200, 200))
	}

	return tests
}

func randomCase(rng *rand.Rand, lo, hi int) testCase {
	n := lo
	if hi > lo {
		n = rng.Intn(hi-lo+1) + lo
	}
	r := rng.Intn(n) + 1
	k := rng.Intn(1_000_001)
	a := make([][]int, n)
	for i := 0; i < n; i++ {
		a[i] = make([]int, n)
		for j := 0; j < n; j++ {
			a[i][j] = rng.Intn(1_000_000) + 1
		}
	}
	c := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				row[j] = '0'
			} else {
				row[j] = '1'
			}
		}
		c[i] = string(row)
	}
	return testCase{n: n, r: r, k: k, a: a, c: c}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierF.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	oracle, cleanup, err := buildOracle()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := append(deterministicTests(), randomTests()...)
	input := buildInput(tests)

	expOut, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle runtime error: %v\noutput:\n%s\n", err, expOut)
		os.Exit(1)
	}
	gotOut, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\noutput:\n%s\n", err, gotOut)
		os.Exit(1)
	}

	expAns, err := parseOutput(expOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse oracle output: %v\n", err)
		os.Exit(1)
	}
	gotAns, err := parseOutput(gotOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse target output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if expAns[i] != gotAns[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d (n=%d): expected (%d %d), got (%d %d)\n", i+1, tests[i].n, expAns[i][0], expAns[i][1], gotAns[i][0], gotAns[i][1])
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}
