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
	t int
	x int
}

type testCase struct {
	n       int
	times   []int
	queries []query
}

func callerFile() (string, bool) {
	_, file, _, ok := runtime.Caller(0)
	return file, ok
}

func buildOracle() (string, func(), error) {
	_, file, ok := callerFile()
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2011I-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleI")
	cmd := exec.Command("go", "build", "-o", outPath, "2011I.go")
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
	sb.Grow(tc.n*16 + len(tc.queries)*16)
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.times {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(strconv.Itoa(q.t))
		sb.WriteByte(' ')
		sb.WriteString(strconv.Itoa(q.x))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseAnswers(out string, n int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != n {
		return nil, fmt.Errorf("expected %d answers, got %d", n, len(lines))
	}
	for i, s := range lines {
		lines[i] = strings.ToUpper(s)
		if lines[i] != "YES" && lines[i] != "NO" {
			return nil, fmt.Errorf("invalid answer %q", s)
		}
	}
	return lines, nil
}

func compareAnswers(exp, got []string) error {
	for i := range exp {
		if exp[i] != got[i] {
			return fmt.Errorf("at position %d expected %s, got %s", i+1, exp[i], got[i])
		}
	}
	return nil
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n:     3,
			times: []int{10, 15, 4},
			queries: []query{
				{t: 1, x: 1},
				{t: 2, x: 1},
				{t: 2, x: 2},
			},
		},
		{
			n:     4,
			times: []int{3, 7, 2, 5},
			queries: []query{
				{t: 1, x: 1},
				{t: 1, x: 2},
				{t: 3, x: 1},
				{t: 2, x: 3},
			},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 120)
	for len(tests) < cap(tests) {
		n := rng.Intn(80) + 8 // [8,87]
		times := make([]int, n)
		for i := range times {
			val := rng.Intn(1000) + 1
			if rng.Intn(7) == 0 {
				val = rng.Intn(1_000_000_000) + 1
			}
			times[i] = val
		}

		inQueue := make([]bool, n+1)
		inStack := make([]bool, n+1)
		queueOnly := func() []int {
			var res []int
			for i := 1; i <= n; i++ {
				if inQueue[i] && !inStack[i] {
					res = append(res, i)
				}
			}
			return res
		}
		notInQueue := func() []int {
			var res []int
			for i := 1; i <= n; i++ {
				if !inQueue[i] {
					res = append(res, i)
				}
			}
			return res
		}
		notInStack := func() []int {
			var res []int
			for i := 1; i <= n; i++ {
				if !inStack[i] {
					res = append(res, i)
				}
			}
			return res
		}

		queries := make([]query, 0, n)
		for i := 0; i < n; i++ {
			var options []int
			queuedOnly := queueOnly()
			if len(queuedOnly) > 0 {
				options = append(options, 3)
			}
			notQ := notInQueue()
			if len(notQ) > 0 {
				options = append(options, 1)
			}
			notS := notInStack()
			if len(notS) > 0 {
				options = append(options, 2)
			}
			if len(options) == 0 {
				// Fallback shouldn't happen; restart test generation.
				options = []int{1}
				notQ = []int{1}
			}
			t := options[rng.Intn(len(options))]
			switch t {
			case 1:
				choices := notQ
				x := choices[rng.Intn(len(choices))]
				inQueue[x] = true
				queries = append(queries, query{t: 1, x: x})
			case 2:
				choices := notInStack()
				x := choices[rng.Intn(len(choices))]
				inStack[x] = true
				queries = append(queries, query{t: 2, x: x})
			case 3:
				choices := queuedOnly
				x := choices[rng.Intn(len(choices))]
				inQueue[x] = false
				inStack[x] = true
				queries = append(queries, query{t: 3, x: x})
			}
		}

		tests = append(tests, testCase{n: n, times: times, queries: queries})
	}
	return tests
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierI.go /path/to/binary")
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
		expAns, err := parseAnswers(expOut, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle output invalid on test %d: %v\n%s", idx+1, err, expOut)
			os.Exit(1)
		}
		gotAns, err := parseAnswers(gotOut, tc.n)
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
