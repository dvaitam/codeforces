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
	refSource   = "2000-2999/2000-2099/2090-2099/2091/2091B.go"
	randomTests = 120
	totalNLimit = 200000 // follows problem constraint sum n <= 2e5
	maxNPerCase = 200000
	maxValue    = 1_000_000_000
)

type testCase struct {
	n int
	x int64
	a []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/candidate")
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

	expect := parseOutputs(expectRaw, len(tests))
	got := parseOutputs(gotRaw, len(tests))

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
	tmp, err := os.CreateTemp("", "2091B-ref-*")
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

	// Deterministic edge coverage.
	tests = append(tests,
		testCase{n: 1, x: 1, a: []int64{1}},
		testCase{n: 2, x: 10, a: []int64{1, 1}},
		testCase{n: 4, x: 4, a: []int64{5, 3, 3, 2}},
		testCase{n: 5, x: 3, a: []int64{5, 3, 2, 3, 2}},
		testCase{n: 3, x: 6, a: []int64{9, 1, 7}},
	)

	used := 0
	for _, tc := range tests {
		used += tc.n
	}

	for i := 0; i < randomTests && used < totalNLimit; i++ {
		remain := totalNLimit - used
		maxN := maxNPerCase
		if maxN > remain {
			maxN = remain
		}
		if maxN < 1 {
			break
		}
		n := rng.Intn(maxN) + 1
		used += n

		a := make([]int64, n)
		for j := 0; j < n; j++ {
			a[j] = rng.Int63n(maxValue) + 1
		}
		var x int64
		switch rng.Intn(4) {
		case 0:
			x = 1
		case 1:
			x = int64(rng.Intn(n) + 1)
		case 2:
			x = int64(maxValue/2 + rng.Int63n(maxValue/2))
		default:
			x = int64(rng.Intn(10) + 1)
		}
		tests = append(tests, testCase{n: n, x: x, a: a})
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.x))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
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

func parseOutputs(out string, t int) []int64 {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		fail("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]int64, t)
	for i, tok := range tokens {
		v, err := strconv.ParseInt(tok, 10, 64)
		if err != nil {
			fail("invalid integer at position %d: %v", i+1, err)
		}
		res[i] = v
	}
	return res
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
