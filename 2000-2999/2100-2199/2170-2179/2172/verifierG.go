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

type testCase struct {
	name string
	s    string
	t    string
	n    int64
}

func buildOracle() (string, func(), error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", nil, fmt.Errorf("cannot determine verifier path")
	}
	dir := filepath.Dir(file)
	tmpDir, err := os.MkdirTemp("", "oracle-2172G-")
	if err != nil {
		return "", nil, err
	}
	outPath := filepath.Join(tmpDir, "oracleG")
	cmd := exec.Command("go", "build", "-o", outPath, "2172G.go")
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
	return strings.TrimSpace(stdout.String()), nil
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(tc.s)
		sb.WriteByte('\n')
		sb.WriteString(tc.t)
		sb.WriteByte('\n')
		sb.WriteString(strconv.FormatInt(tc.n, 10))
		sb.WriteByte('\n')
	}
	return sb.String()
}

func deterministicTests() []testCase {
	return []testCase{
		{name: "example1", s: "ABAB", t: "A", n: 1},
		{name: "example2", s: "ABABA", t: "AABAB", n: 1},
		{name: "simple_A", s: "AB", t: "AB", n: 2},
		{name: "simple_B", s: "ABABAB", t: "BBBB", n: 3},
		{name: "short_s", s: "AB", t: "BA", n: 10},
		{name: "long_t", s: "ABAB", t: strings.Repeat("AB", 5), n: 20},
	}
}

func randomS(rng *rand.Rand) string {
	l := rng.Intn(10) + 2
	for {
		var b strings.Builder
		for i := 0; i < l; i++ {
			if i == 0 {
				b.WriteByte('A')
			} else if i == l-1 {
				b.WriteByte('B')
			} else {
				if rng.Intn(2) == 0 {
					b.WriteByte('A')
				} else {
					b.WriteByte('B')
				}
			}
		}
		s := b.String()
		if isValidS(s) {
			return s
		}
	}
}

func isValidS(s string) bool {
	if s[0] != 'A' || s[len(s)-1] != 'B' {
		return false
	}
	for i := 0; i < len(s)-1; i++ {
		if s[i] == 'A' && s[i+1] == 'A' {
			return false
		}
		if i+2 < len(s) && s[i] == 'B' && s[i+1] == 'B' && s[i+2] == 'B' {
			return false
		}
	}
	return true
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]testCase, 0, 20)
	for i := 0; i < 10; i++ {
		s := randomS(rng)
		tLen := rng.Intn(20) + 1
		var sb strings.Builder
		for j := 0; j < tLen; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('A')
			} else {
				sb.WriteByte('B')
			}
		}
		n := rng.Int63n(1000) + 1
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_small_%d", i+1),
			s:    s,
			t:    sb.String(),
			n:    n,
		})
	}
	for i := 0; i < 10; i++ {
		s := randomS(rng)
		tLen := rng.Intn(100000) + 1
		var sb strings.Builder
		for j := 0; j < tLen; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('A')
			} else {
				sb.WriteByte('B')
			}
		}
		n := rng.Int63n(1_000_000_000) + 1
		tests = append(tests, testCase{
			name: fmt.Sprintf("random_large_%d", i+1),
			s:    s,
			t:    sb.String(),
			n:    n,
		})
	}
	return tests
}

func stressTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano() + 123))
	tests := make([]testCase, 0, 5)
	tests = append(tests, testCase{
		name: "stress_max_n",
		s:    "ABABABABAB",
		t:    strings.Repeat("AB", 50000),
		n:    1_000_000_000,
	})
	for i := 0; i < 4; i++ {
		s := randomS(rng)
		tLen := 100000
		var sb strings.Builder
		for j := 0; j < tLen; j++ {
			if rng.Intn(2) == 0 {
				sb.WriteByte('A')
			} else {
				sb.WriteByte('B')
			}
		}
		tests = append(tests, testCase{
			name: fmt.Sprintf("stress_%d", i+1),
			s:    s,
			t:    sb.String(),
			n:    1_000_000_000,
		})
	}
	return tests
}

func compareOutputs(expected, actual string, count int) error {
	exp := strings.Fields(expected)
	act := strings.Fields(actual)
	if len(exp) != count {
		return fmt.Errorf("oracle produced %d answers, expected %d", len(exp), count)
	}
	if len(act) != count {
		return fmt.Errorf("expected %d answers, got %d", count, len(act))
	}
	for i := 0; i < count; i++ {
		if exp[i] != act[i] {
			return fmt.Errorf("mismatch at case %d: expected %s got %s", i+1, exp[i], act[i])
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierG.go /path/to/binary")
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
	tests = append(tests, stressTests()...)
	input := buildInput(tests)

	expected, err := runBinary(oracle, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "oracle failed: %v\n", err)
		os.Exit(1)
	}
	actual, err := runBinary(target, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "target runtime error: %v\n", err)
		os.Exit(1)
	}
	if err := compareOutputs(expected, actual, len(tests)); err != nil {
		fmt.Fprintf(os.Stderr, "%v\nInput:\n%s\nExpected:\n%s\nActual:\n%s\n", err, input, expected, actual)
		os.Exit(1)
	}
	fmt.Println("All tests passed.")
}
