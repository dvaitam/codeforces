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
	refSource     = "2000-2999/2000-2099/2090-2099/2097/2097D.go"
	randomTests   = 120
	totalLenLimit = 200000 // keep well below 1e6 for speed
	maxNPerCase   = 50000
)

type testCase struct {
	n int
	s string
	t string
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

	expect := parseOutputs(expectRaw, len(tests))
	got := parseOutputs(gotRaw, len(tests))

	if len(expect) != len(got) {
		fail("output length mismatch: expected %d tokens, got %d", len(expect), len(got))
	}
	for i := range expect {
		if expect[i] != got[i] {
			fail("mismatch at test %d: expected %s, got %s", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("All %d test cases passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2097D-ref-*")
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

	// Deterministic coverage.
	tests = append(tests,
		testCase{n: 1, s: "0", t: "0"},
		testCase{n: 1, s: "1", t: "0"},
		testCase{n: 2, s: "00", t: "11"},
		testCase{n: 2, s: "10", t: "10"},
		testCase{n: 4, s: "0000", t: "0000"},
		testCase{n: 4, s: "0101", t: "1001"},
	)

	total := 0
	for _, tc := range tests {
		total += tc.n
	}

	for i := 0; i < randomTests && total < totalLenLimit; i++ {
		remain := totalLenLimit - total
		maxN := maxNPerCase
		if maxN > remain {
			maxN = remain
		}
		if maxN < 1 {
			break
		}
		n := rng.Intn(maxN) + 1
		total += n

		var sbS, sbT strings.Builder
		sbS.Grow(n)
		sbT.Grow(n)
		ensureOneS := rng.Intn(2) == 0
		ensureOneT := rng.Intn(2) == 0
		for j := 0; j < n; j++ {
			bs := byte('0' + rng.Intn(2))
			bt := byte('0' + rng.Intn(2))
			if ensureOneS && j == n-1 && !strings.Contains(sbS.String(), "1") {
				bs = '1'
			}
			if ensureOneT && j == n-1 && !strings.Contains(sbT.String(), "1") {
				bt = '1'
			}
			sbS.WriteByte(bs)
			sbT.WriteByte(bt)
		}
		tests = append(tests, testCase{n: n, s: sbS.String(), t: sbT.String()})
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
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		sb.WriteString(tc.t)
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

func parseOutputs(out string, t int) []string {
	tokens := strings.Fields(out)
	if len(tokens) != t {
		fail("expected %d tokens, got %d", t, len(tokens))
	}
	res := make([]string, t)
	for i, tok := range tokens {
		res[i] = strings.ToUpper(tok)
	}
	return res
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
