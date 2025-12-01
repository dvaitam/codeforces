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

type query struct {
	l int
	r int
}

type testCase struct {
	n       int
	arr     []int
	queries []query
}

func callerFile() (string, bool) {
	_, file, _, ok := runtime.Caller(0)
	return file, ok
}

func buildOracle() (string, func(), error) {
	file, ok := callerFile()
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2043F-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleF")
	cmd := exec.Command("go", "build", "-o", outPath, "2043F.go")
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

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.Grow(tc.n*4 + len(tc.queries)*8 + 32)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte(' ')
	sb.WriteString(strconv.Itoa(len(tc.queries)))
	sb.WriteByte('\n')
	for i, v := range tc.arr {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(strconv.Itoa(q.l))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(q.r))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseLines(out string, q int) ([][2]int64, error) {
	lines := strings.Fields(out)
	// Each query produces either 1 token (-1) or 2 tokens.
	res := make([][2]int64, 0, q)
	for i := 0; i < len(lines); {
		if len(res) >= q {
			return nil, fmt.Errorf("too many answers")
		}
		first, err := strconv.ParseInt(lines[i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", lines[i])
		}
		if first == -1 {
			res = append(res, [2]int64{-1, 0})
			i++
		} else {
			if i+1 >= len(lines) {
				return nil, fmt.Errorf("line %d missing second number", len(res)+1)
			}
			second, err := strconv.ParseInt(lines[i+1], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid integer %q", lines[i+1])
			}
			res = append(res, [2]int64{first, second})
			i += 2
		}
	}
	if len(res) != q {
		return nil, fmt.Errorf("expected %d answers, got %d", q, len(res))
	}
	return res, nil
}

func compareAnswers(exp, got [][2]int64) error {
	if len(exp) != len(got) {
		return fmt.Errorf("answer count mismatch")
	}
	for i := range exp {
		if exp[i][0] != got[i][0] || exp[i][1] != got[i][1] {
			return fmt.Errorf("at query %d expected (%d,%d) got (%d,%d)", i+1, exp[i][0], exp[i][1], got[i][0], got[i][1])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n:   9,
			arr: []int{0, 1, 2, 1, 3, 4, 5, 6, 0},
			queries: []query{
				{l: 1, r: 5},
				{l: 2, r: 5},
				{l: 3, r: 5},
				{l: 4, r: 5},
				{l: 1, r: 9},
			},
		},
		{
			n:   3,
			arr: []int{1, 1, 1},
			queries: []query{
				{l: 1, r: 3},
				{l: 2, r: 3},
			},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		n := rng.Intn(180) + 20
		if rng.Intn(5) == 0 {
			n = rng.Intn(800) + 200
		}
		arr := make([]int, n)
		for i := range arr {
			arr[i] = rng.Intn(51)
		}
		q := rng.Intn(120) + 10
		if rng.Intn(5) == 0 {
			q = rng.Intn(600) + 100
		}
		queries := make([]query, q)
		for i := 0; i < q; i++ {
			l := rng.Intn(n) + 1
			r := rng.Intn(n-l+1) + l
			queries[i] = query{l: l, r: r}
		}
		tests = append(tests, testCase{n: n, arr: arr, queries: queries})
	}
	return tests
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
	for idx, tc := range tests {
		input := buildInput(tc)
		expOut, err := runBinary(oracle, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle failed on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		gotOut, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "target runtime error on test %d: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
		expAns, err := parseLines(expOut, len(tc.queries))
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n%s", idx+1, err, expOut)
			os.Exit(1)
		}
		gotAns, err := parseLines(gotOut, len(tc.queries))
		if err != nil {
			fmt.Fprintf(os.Stderr, "target output invalid on test %d: %v\n%s", idx+1, err, gotOut)
			os.Exit(1)
		}
		if err := compareAnswers(expAns, gotAns); err != nil {
			fmt.Fprintf(os.Stderr, "test %d mismatch: %v\ninput:\n%s", idx+1, err, input)
			os.Exit(1)
		}
	}

	fmt.Println("All tests passed.")
}
