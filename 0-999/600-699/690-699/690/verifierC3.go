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

const (
	refSource        = "690C3.go"
	tempOraclePrefix = "oracle-690C3-"
	randomTestsCount = 120
	maxRandomN       = 5000
)

type testCase struct {
	name    string
	n       int
	parents []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC3.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	oraclePath, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer cleanup()

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTests(randomTestsCount, rng)...)
	tests = append(tests, largeTests()...)

	for idx, tc := range tests {
		input := formatInput(tc)

		expOut, err := runProgram(oraclePath, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "oracle runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		exp, err := parseOutput(expOut, tc.n-1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse oracle output on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runProgram(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		got, err := parseOutput(gotOut, tc.n-1)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output parse error on test %d (%s): %v\noutput:\n%s", idx+1, tc.name, err, gotOut)
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}

		if len(exp) != len(got) {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d values, got %d\n", idx+1, tc.name, len(exp), len(got))
			fmt.Println("Input:")
			fmt.Print(input)
			os.Exit(1)
		}
		for i := range exp {
			if exp[i] != got[i] {
				fmt.Fprintf(os.Stderr, "test %d (%s) failed at position %d: expected %d got %d\n", idx+1, tc.name, i+1, exp[i], got[i])
				fmt.Println("Input:")
				fmt.Print(input)
				os.Exit(1)
			}
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("failed to determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", tempOraclePrefix)
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleC3")
	cmd := exec.Command("go", "build", "-o", outPath, refSource)
	cmd.Dir = dir
	if out, err := cmd.CombinedOutput(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("build oracle failed: %v\n%s", err, out)
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return outPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
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

func parseOutput(out string, expect int) ([]int, error) {
	if expect == 0 {
		if strings.TrimSpace(out) == "" {
			return nil, nil
		}
		return nil, fmt.Errorf("expected empty output but got %q", out)
	}
	fields := strings.Fields(out)
	if len(fields) != expect {
		return nil, fmt.Errorf("expected %d numbers, got %d", expect, len(fields))
	}
	res := make([]int, expect)
	for i, f := range fields {
		val, err := strconv.Atoi(f)
		if err != nil {
			return nil, fmt.Errorf("token %q is not an integer", f)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", tc.n)
	for i, p := range tc.parents {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(p))
	}
	if len(tc.parents) > 0 {
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "two_nodes", n: 2, parents: []int{1}},
		{name: "chain_5", n: 5, parents: []int{1, 2, 3, 4}},
		{name: "star_6", n: 6, parents: []int{1, 1, 1, 1, 1}},
		{name: "balanced", n: 9, parents: []int{1, 1, 2, 2, 3, 3, 4, 4}},
		{name: "alternating", n: 10, parents: []int{1, 2, 3, 4, 5, 1, 2, 3, 4}},
	}
}

func randomTests(count int, rng *rand.Rand) []testCase {
	tests := make([]testCase, 0, count)
	for i := 0; i < count; i++ {
		n := rng.Intn(maxRandomN-1) + 2
		parents := make([]int, n-1)
		for k := 2; k <= n; k++ {
			parents[k-2] = rng.Intn(k-1) + 1
		}
		tests = append(tests, testCase{
			name:    fmt.Sprintf("random_%d", i+1),
			n:       n,
			parents: parents,
		})
	}
	return tests
}

func largeTests() []testCase {
	chainParents := make([]int, 200000-1)
	for i := 0; i < len(chainParents); i++ {
		chainParents[i] = i + 1
	}
	starParents := make([]int, 200000-1)
	for i := range starParents {
		starParents[i] = 1
	}
	return []testCase{
		{name: "large_chain", n: 200000, parents: chainParents},
		{name: "large_star", n: 200000, parents: starParents},
	}
}
