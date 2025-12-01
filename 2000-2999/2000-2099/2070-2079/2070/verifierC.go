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
	refSource   = "./2070C.go"
	randomTests = 120
	totalNLimit = 200000 // keep within problem sum constraint of 3e5
	maxNPerCase = 5000
	maxPenalty  = 1_000_000_000
)

type testCase struct {
	n int
	k int
	s string
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

	expect, err := parseOutputs(expectRaw, len(tests))
	if err != nil {
		fail("could not parse reference output: %v", err)
	}
	got, err := parseOutputs(gotRaw, len(tests))
	if err != nil {
		fail("could not parse candidate output: %v", err)
	}

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			fail("mismatch at test %d: expected %d, got %d", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2070C-ref-*")
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
	tests := make([]testCase, 0, randomTests+6)

	// Deterministic coverage: small sizes and corner k values.
	tests = append(tests,
		testCase{n: 1, k: 0, s: "R", a: []int{1}},
		testCase{n: 1, k: 1, s: "B", a: []int{2}},
		testCase{n: 4, k: 1, s: "BRBR", a: []int{9, 3, 5, 4}},
		testCase{n: 4, k: 2, s: "BRBR", a: []int{9, 3, 5, 4}},
		testCase{n: 5, k: 5, s: "RRRRR", a: []int{5, 3, 1, 2, 4}},
		testCase{n: 6, k: 0, s: "BBBBBB", a: []int{1, 2, 3, 4, 5, 6}},
	)

	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	for i := 0; i < randomTests && totalN < totalNLimit; i++ {
		remain := totalNLimit - totalN
		maxN := maxNPerCase
		if maxN > remain {
			maxN = remain
		}
		if maxN < 1 {
			break
		}
		n := rng.Intn(maxN) + 1
		k := rng.Intn(n + 1)
		var sb strings.Builder
		sb.Grow(n)
		for j := 0; j < n; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('R')
			} else {
				sb.WriteByte('B')
			}
		}
		a := make([]int, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Intn(maxPenalty) + 1
		}
		// Make a few structured patterns to stress segment grouping.
		if i%5 == 0 && n >= 4 {
			sb.Reset()
			for j := 0; j < n; j++ {
				if (j/2)%2 == 0 {
					sb.WriteByte('R')
				} else {
					sb.WriteByte('B')
				}
			}
		}
		tests = append(tests, testCase{n: n, k: k, s: sb.String(), a: a})
		totalN += n
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.k))
		sb.WriteString(tc.s)
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

func parseOutputs(out string, t int) ([]int, error) {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]int, t)
	for i, tok := range tokens {
		v, err := strconv.Atoi(tok)
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
