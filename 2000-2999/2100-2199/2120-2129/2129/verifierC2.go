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
	refSource   = "2000-2999/2100-2199/2120-2129/2129/2129C2.go"
	randomTests = 80
	totalLen    = 50000 // keep total n moderate
	maxN        = 1000
	minN        = 2
)

type testCase struct {
	n int
	s string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC2.go /path/to/candidate")
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
	tmp, err := os.CreateTemp("", "2129C2-ref-*")
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

	// Deterministic cases.
	tests = append(tests,
		testCase{n: 2, s: "()"},
		testCase{n: 2, s: ")("},
		testCase{n: 4, s: "(())"},
		testCase{n: 6, s: "()(())"},
		testCase{n: 6, s: ")))((("},
	)

	total := 0
	for _, tc := range tests {
		total += tc.n
	}

	for i := 0; i < randomTests && total < totalLen; i++ {
		remain := totalLen - total
		maxPossible := maxN
		if maxPossible > remain {
			maxPossible = remain
		}
		n := rng.Intn(maxPossible-minN+1) + minN
		total += n
		tests = append(tests, randomCase(n, rng))
	}

	return tests
}

func randomCase(n int, rng *rand.Rand) testCase {
	var sb strings.Builder
	sb.Grow(n)
	pos := rng.Intn(n-1) + 1 // ensure both '(' and ')'
	for i := 0; i < n; i++ {
		if i < pos {
			if rng.Intn(2) == 0 {
				sb.WriteByte('(')
			} else {
				sb.WriteByte(')')
			}
		} else {
			if rng.Intn(2) == 0 {
				sb.WriteByte('(')
			} else {
				sb.WriteByte(')')
			}
		}
	}
	// ensure at least one of each
	str := sb.String()
	if !strings.Contains(str, "(") {
		str = "(" + str[1:]
	}
	if !strings.Contains(str, ")") {
		str = ")" + str[1:]
	}
	return testCase{n: n, s: str}
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
	copy(res, tokens[:t])
	return res
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
