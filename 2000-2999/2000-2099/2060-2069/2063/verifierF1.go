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
	refSource   = "2000-2999/2000-2099/2060-2069/2063/2063F1.go"
	mod         = 998244353
	randomTests = 120
	totalBudget = 5000 // sum of n across all tests
	maxNPerCase = 300  // keep individual cases moderate for speed
)

type testCase struct {
	n     int
	pairs [][2]int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF1.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "2063F1-ref-*")
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

	// Deterministic coverage with small structures.
	tests = append(tests,
		makeFromString("()"),
		makeFromString("(())"),
		makeFromString("()()"),
		makeFromString("(()())"),
		makeFromString("((()))"),
	)

	used := 0
	for _, tc := range tests {
		used += tc.n
	}

	for i := 0; i < randomTests && used < totalBudget; i++ {
		remain := totalBudget - used
		maxN := maxNPerCase
		if maxN > remain {
			maxN = remain
		}
		n := rng.Intn(maxN) + 1
		seq := randomBalancedSequence(2*n, rng)
		tc := makeFromString(seq)
		tests = append(tests, tc)
		used += tc.n
	}

	return tests
}

func makeFromString(s string) testCase {
	stack := make([]int, 0, len(s)/2)
	pairs := make([][2]int, 0, len(s)/2)
	for i := 0; i < len(s); i++ {
		if s[i] == '(' {
			stack = append(stack, i+1)
		} else {
			open := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			pairs = append(pairs, [2]int{open, i + 1})
		}
	}
	// Shuffle clues to avoid relying on any fixed order.
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := len(pairs) - 1; i > 0; i-- {
		j := rng.Intn(i + 1)
		pairs[i], pairs[j] = pairs[j], pairs[i]
	}
	return testCase{n: len(pairs), pairs: pairs}
}

func randomBalancedSequence(length int, rng *rand.Rand) string {
	open := length / 2  // remaining '('
	close := length / 2 // remaining ')'
	balance := 0        // current '(' - ')'
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		// If we must place '(', if all ')' would break balance.
		if open > 0 && (balance == close) {
			sb.WriteByte('(')
			open--
			balance++
			continue
		}
		if close > 0 && balance == 0 {
			sb.WriteByte('(')
			open--
			balance++
			continue
		}
		if open == 0 {
			sb.WriteByte(')')
			close--
			balance--
			continue
		}
		if close == 0 {
			sb.WriteByte('(')
			open--
			balance++
			continue
		}
		if rng.Intn(2) == 0 {
			sb.WriteByte('(')
			open--
			balance++
		} else {
			sb.WriteByte(')')
			close--
			balance--
		}
	}
	return sb.String()
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for _, p := range tc.pairs {
			sb.WriteString(fmt.Sprintf("%d %d\n", p[0], p[1]))
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

func parseOutputs(out string, tests []testCase) ([]int, error) {
	tokens := strings.Fields(out)
	expected := 0
	for _, tc := range tests {
		expected += tc.n + 1
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
