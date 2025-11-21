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

type op struct {
	a, b, c, d int
	ch         byte
}

type testCase struct {
	n, m, k int
	grid    []string
	ops     []op
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-761F-")
	if err != nil {
		return "", nil, err
	}
	bin := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", bin, "761F.go")
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build oracle: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return bin, cleanup, nil
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
	sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.m, tc.k))
	for _, row := range tc.grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	for _, o := range tc.ops {
		sb.WriteString(fmt.Sprintf("%d %d %d %d %c\n", o.a, o.b, o.c, o.d, o.ch))
	}
	return sb.String()
}

func parseAnswer(out string) (int64, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer output, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	if val < 0 {
		return 0, fmt.Errorf("negative distance %d", val)
	}
	return val, nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 1, m: 1, k: 1,
			grid: []string{"a"},
			ops:  []op{{1, 1, 1, 1, 'b'}},
		},
		{
			n: 2, m: 2, k: 2,
			grid: []string{"ab", "cd"},
			ops: []op{
				{1, 1, 1, 2, 'z'},
				{2, 1, 2, 2, 'a'},
			},
		},
		{
			n: 3, m: 3, k: 3,
			grid: []string{"abc", "def", "ghi"},
			ops: []op{
				{1, 1, 3, 3, 'a'},
				{1, 2, 2, 3, 'm'},
				{2, 1, 3, 2, 'z'},
			},
		},
	}
}

func randomTest(rng *rand.Rand) testCase {
	// vary dimensions
	var n, m int
	switch rng.Intn(4) {
	case 0:
		n = rng.Intn(3) + 1
		m = rng.Intn(3) + 1
	case 1:
		n = rng.Intn(10) + 1
		m = rng.Intn(10) + 1
	default:
		n = rng.Intn(30) + 1
		m = rng.Intn(30) + 1
	}
	k := rng.Intn(30) + 1
	if rng.Intn(4) == 0 {
		k = rng.Intn(200) + 1
	}
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			row[j] = byte('a' + rng.Intn(26))
		}
		grid[i] = string(row)
	}
	ops := make([]op, k)
	for i := 0; i < k; i++ {
		a := rng.Intn(n) + 1
		b := rng.Intn(m) + 1
		c := rng.Intn(n-a+1) + a
		d := rng.Intn(m-b+1) + b
		ch := byte('a' + rng.Intn(26))
		ops[i] = op{a: a, b: b, c: c, d: d, ch: ch}
	}
	return testCase{n: n, m: m, k: k, grid: grid, ops: ops}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/binary")
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

		expStr, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		exp, err := parseAnswer(expStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid oracle output on test %d: %v\noutput:\n%s\n", idx+1, err, expStr)
			os.Exit(1)
		}

		gotStr, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		got, err := parseAnswer(gotStr)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\noutput:\n%s\ninput:\n%s\n", idx+1, err, gotStr, input)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %d got %d\ninput:\n%s\n", idx+1, exp, got, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
