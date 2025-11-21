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
	refSource    = "2000-2999/2000-2099/2060-2069/2060/2060F.go"
	mod          = 998244353
	randomTests  = 80
	totalKBudget = 95000 // keep sum of k across tests within limits and runtime
	maxK         = 100000
)

type testCase struct {
	k int
	n int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "2060F-ref-*")
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
	tests := make([]testCase, 0, randomTests+5)

	// Small deterministic coverage and the sample-like patterns.
	tests = append(tests,
		testCase{k: 2, n: 2},
		testCase{k: 4, n: 3},
		testCase{k: 10, n: 6},
		testCase{k: 1, n: 1},
		testCase{k: 3, n: 10},
	)

	usedK := 0
	for _, tc := range tests {
		usedK += tc.k
	}

	for i := 0; i < randomTests && usedK < totalKBudget; i++ {
		left := totalKBudget - usedK
		maxKHere := maxK
		if left < maxKHere {
			maxKHere = left
		}
		k := rng.Intn(maxKHere) + 1
		usedK += k

		// n can be as large as 9e8; test both small and large values.
		var n int
		if rng.Intn(3) == 0 {
			n = rng.Intn(10) + 1
		} else {
			n = rng.Intn(900_000_000) + 1
		}
		tests = append(tests, testCase{k: k, n: n})
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.k, tc.n))
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

func parseOutputs(out string, tests []testCase) ([]int, error) {
	tokens := strings.Fields(out)
	expected := 0
	for _, tc := range tests {
		expected += tc.k
	}
	if len(tokens) != expected {
		return nil, fmt.Errorf("expected %d tokens, got %d", expected, len(tokens))
	}
	res := make([]int, expected)
	for i, tok := range tokens {
		v, err := strconv.Atoi(tok)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at position %d: %v", i+1, err)
		}
		if v < 0 || v >= mod {
			return nil, fmt.Errorf("value out of range at position %d: %d", i+1, v)
		}
		res[i] = v
	}
	return res, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
