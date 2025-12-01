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
)

const refSource = "./506E.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		expOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expVal, err := parseModulo(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, expOut)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotVal, err := parseModulo(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, gotOut)
			os.Exit(1)
		}

		if gotVal != expVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): expected %d, got %d\ninput:\n%s\n", idx+1, tc.name, expVal, gotVal, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "506E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return strings.TrimSpace(stdout.String()), nil
}

func parseModulo(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q: %v", fields[0], err)
	}
	if val < 0 || val >= 10007 {
		val = ((val % 10007) + 10007) % 10007
	}
	return val, nil
}

func generateTests() []testCase {
	tests := []testCase{
		makeCase("sample1", "revive", 1),
		makeCase("sample2", "ad", 3),
		makeCase("needs-more-than-n", "abc", 1),
		makeCase("already-pal", "abba", 5),
		makeCase("single-char", "z", 10),
		makeCase("two-chars", "az", 2),
		makeCase("long-n", strings.Repeat("kitayuta", 10), 1_000_000_000),
		makeCase("repeat-a", strings.Repeat("a", 200), 999_999_937),
	}

	tests = append(tests, targetedCases()...)

	rng := rand.New(rand.NewSource(5060506))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomCase(rng, fmt.Sprintf("random-%d", i+1)))
	}
	return tests
}

func targetedCases() []testCase {
	var cases []testCase
	cases = append(cases, makeCase("alternating", repeatPattern("ab", 70), 200))
	cases = append(cases, makeCase("pal-odd", "abcba", 7))
	cases = append(cases, makeCase("pal-even", "abccba", 8))
	cases = append(cases, makeCase("max-length", randomLetters(200, 1), 123456789))
	return cases
}

func randomCase(rng *rand.Rand, name string) testCase {
	length := rng.Intn(200) + 1
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(byte('a' + rng.Intn(26)))
	}
	var n int64
	switch rng.Intn(4) {
	case 0:
		n = int64(rng.Intn(5) + 1)
	case 1:
		n = int64(rng.Intn(500) + 1)
	case 2:
		n = int64(rng.Int63n(1_000_000_000) + 1)
	default:
		n = int64(rng.Intn(200)) + int64(length)/2
		if n < 1 {
			n = 1
		}
	}
	return makeCase(name, sb.String(), n)
}

func makeCase(name, s string, n int64) testCase {
	return testCase{
		name:  name,
		input: fmt.Sprintf("%s\n%d\n", s, n),
	}
}

func repeatPattern(p string, count int) string {
	var sb strings.Builder
	for i := 0; i < count; i++ {
		sb.WriteString(p)
	}
	return sb.String()
}

func randomLetters(length int, seedOffset int64) string {
	rng := rand.New(rand.NewSource(100 + seedOffset))
	var sb strings.Builder
	for i := 0; i < length; i++ {
		sb.WriteByte(byte('a' + rng.Intn(26)))
	}
	return sb.String()
}
