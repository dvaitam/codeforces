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
	refSource     = "./2101C.go"
	randomTests   = 120
	maxTotalN     = 200000
	maxCaseLength = 60000
)

type testCase struct {
	n int
	a []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	expectRaw, err := runProgram(exec.Command(refBin), input)
	if err != nil {
		fail("reference failed: %v\n%s", err, expectRaw)
	}
	expect, err := parseOutputs(expectRaw, len(tests))
	if err != nil {
		fail("could not parse reference output: %v\n%s", err, expectRaw)
	}

	gotRaw, err := runProgram(commandFor(candidate), input)
	if err != nil {
		fail("candidate failed: %v\n%s", err, gotRaw)
	}
	got, err := parseOutputs(gotRaw, len(tests))
	if err != nil {
		fail("could not parse candidate output: %v\n%s", err, gotRaw)
	}

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			tc := tests[i]
			fail("mismatch on test %d: expected %d, got %d (n=%d, a[:min(10)]%v)", i+1, expect[i], got[i], tc.n, preview(tc.a, 10))
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2101C-ref-*")
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

func buildTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, randomTests+5)

	totalN := 0
	add := func(a []int) {
		n := len(a)
		if n == 0 || totalN+n > maxTotalN {
			return
		}
		tests = append(tests, testCase{n: n, a: a})
		totalN += n
	}

	// Deterministic coverage including samples and edge cases.
	add([]int{1, 2, 1, 2})
	add([]int{2, 2})
	add([]int{1, 2, 1, 5, 1, 2, 2, 1, 1, 2})
	add([]int{1, 5, 2, 8, 4, 1, 4, 2})
	add([]int{1})
	add([]int{1, 1, 1, 1})
	add([]int{4, 3, 2, 1})

	for len(tests) < randomTests && totalN < maxTotalN {
		remain := maxTotalN - totalN
		maxLen := maxCaseLength
		if maxLen > remain {
			maxLen = remain
		}
		if maxLen < 1 {
			break
		}
		n := rng.Intn(maxLen) + 1

		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(n) + 1 // within [1, n]
		}
		add(a)
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(strconv.Itoa(tc.n))
		sb.WriteByte('\n')
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(v))
		}
		sb.WriteByte('\n')
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

func parseOutputs(out string, t int) ([]int64, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid integer at position %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res, nil
}

func preview(a []int, k int) []int {
	if k > len(a) {
		k = len(a)
	}
	cpy := make([]int, k)
	copy(cpy, a[:k])
	return cpy
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
