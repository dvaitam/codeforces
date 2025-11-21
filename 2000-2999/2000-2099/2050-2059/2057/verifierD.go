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

const (
	refSource     = "2000-2999/2000-2099/2050-2059/2057/2057D.go"
	randomTests   = 80
	totalBudget   = 40000 // limit for sum of n+q across all generated tests
	maxValue      = 1_000_000_000
	maxCaseLength = 500 // to keep runs fast while covering variety
)

type testCase struct {
	n       int
	q       int
	a       []int64
	updates [][2]int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	expectRaw, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference failed: %v\n%s", err, expectRaw)
	}
	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}

	expect, err := parseOutputs(expectRaw, tests)
	if err != nil {
		fail("could not parse reference output: %v", err)
	}
	got, err := parseOutputs(gotRaw, tests)
	if err != nil {
		fail("could not parse candidate output: %v", err)
	}

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			fail("mismatch at output %d: expected %d, got %d", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2057D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, randomTests+4)

	// Deterministic edge coverage.
	tests = append(tests, testCase{
		n: 1, q: 3,
		a: []int64{5},
		updates: [][2]int64{
			{1, 5}, {1, 1}, {1, 10},
		},
	})
	tests = append(tests, testCase{
		n: 2, q: 2,
		a: []int64{1, 10},
		updates: [][2]int64{
			{1, 10}, {2, 2},
		},
	})
	tests = append(tests, testCase{
		n: 5, q: 4,
		a: []int64{1, 2, 3, 4, 5},
		updates: [][2]int64{
			{3, 1}, {5, 10}, {2, 8}, {4, 4},
		},
	})
	tests = append(tests, testCase{
		n: 5, q: 4,
		a: []int64{9, 7, 7, 7, 9},
		updates: [][2]int64{
			{3, 1}, {5, 1}, {2, 9}, {4, 9},
		},
	})

	used := 0
	for _, tc := range tests {
		used += tc.n + tc.q
	}
	remaining := totalBudget - used

	for i := 0; i < randomTests && remaining > 0; i++ {
		maxLen := maxCaseLength
		if maxLen > remaining-2 { // leave room for q at least 1
			maxLen = remaining - 2
		}
		if maxLen < 1 {
			break
		}
		n := rng.Intn(maxLen) + 1
		remaining -= n

		maxQ := maxCaseLength
		if maxQ > remaining {
			maxQ = remaining
		}
		if maxQ < 1 {
			break
		}
		q := rng.Intn(maxQ) + 1
		remaining -= q

		a := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = 1 + rng.Int63n(maxValue)
		}
		updates := make([][2]int64, q)
		for j := 0; j < q; j++ {
			pos := rng.Intn(n) + 1
			val := 1 + rng.Int63n(maxValue)
			updates[j] = [2]int64{int64(pos), val}
		}
		tests = append(tests, testCase{n: n, q: q, a: a, updates: updates})
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.q))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
		for _, up := range tc.updates {
			sb.WriteString(fmt.Sprintf("%d %d\n", up[0], up[1]))
		}
	}
	return sb.String()
}

func runProgram(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("stderr not empty")
	}
	return out.String(), nil
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

func parseOutputs(out string, tests []testCase) ([]int64, error) {
	tokens := strings.Fields(out)
	expectCount := 0
	for _, tc := range tests {
		expectCount += tc.q + 1
	}
	if len(tokens) != expectCount {
		return nil, fmt.Errorf("expected %d tokens, got %d", expectCount, len(tokens))
	}
	res := make([]int64, expectCount)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at position %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
