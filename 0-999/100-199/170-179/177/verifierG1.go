package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	refSource          = "./177G1.go"
	maxTotalPatternLen = 100000
)

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierG1.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fatal("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		exp, err := runProgram(refBin, tc.input)
		if err != nil {
			fatal("reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fatal("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		if normalize(got) != normalize(exp) {
			fatal("wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s", i+1, tc.input, exp, got)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fatal(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "177G1-ref-*")
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
		return "", fmt.Errorf("build reference failed: %v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", filepath.Clean(bin))
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func normalize(s string) string {
	trimmed := strings.TrimSpace(s)
	if trimmed == "" {
		return ""
	}
	lines := strings.Split(trimmed, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	return strings.Join(lines, "\n")
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(20240529))
	var tests []testCase
	tests = append(tests, makeTestCase(1, []string{"a"}))
	tests = append(tests, makeTestCase(2, []string{"a", "b", "ab"}))
	tests = append(tests, makeTestCase(5, []string{"ba", "ab", "aba", "bab"}))
	tests = append(tests, makeTestCase(7, []string{"bbbb", "aaaaa", "abba"}))
	tests = append(tests, makeLongPatternTest())
	for i := 0; i < 50; i++ {
		tests = append(tests, randomTest(rng, 30, 80, false))
	}
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTest(rng, 80, 250, true))
	}
	for i := 0; i < 10; i++ {
		tests = append(tests, randomTest(rng, 5, 20000, true))
	}
	tests = append(tests, makeLargeMTest(rng))
	tests = append(tests, makeHugeKTest(rng))
	return tests
}

func randomTest(rng *rand.Rand, maxM, maxLen int, allowHugeK bool) testCase {
	if maxM < 1 {
		maxM = 1
	}
	if maxLen < 1 {
		maxLen = 1
	}
	m := rng.Intn(maxM) + 1
	patterns := make([]string, m)
	total := 0
	for i := 0; i < m; i++ {
		remainingSlots := m - i - 1
		capacity := maxTotalPatternLen - total - remainingSlots
		if capacity < 1 {
			capacity = 1
		}
		maxForNow := maxLen
		if maxForNow > capacity {
			maxForNow = capacity
		}
		if maxForNow < 1 {
			maxForNow = 1
		}
		length := rng.Intn(maxForNow) + 1
		patterns[i] = randomPattern(rng, length)
		total += length
	}
	k := randomK(rng, allowHugeK)
	return makeTestCase(k, patterns)
}

func randomK(rng *rand.Rand, allowHuge bool) int64 {
	switch rng.Intn(5) {
	case 0:
		return int64(rng.Intn(20) + 1)
	case 1:
		return int64(rng.Intn(500) + 1)
	case 2:
		return int64(rng.Intn(100000) + 1)
	case 3:
		return int64(rng.Intn(1000000) + 1)
	default:
		if allowHuge {
			return rng.Int63n(1_000_000_000_000_000_000) + 1
		}
		return int64(rng.Intn(5000000) + 1)
	}
}

func randomPattern(rng *rand.Rand, length int) string {
	if length < 1 {
		length = 1
	}
	b := make([]byte, length)
	for i := range b {
		if rng.Intn(2) == 0 {
			b[i] = 'a'
		} else {
			b[i] = 'b'
		}
	}
	return string(b)
}

func makeLargeMTest(rng *rand.Rand) testCase {
	m := 10000
	patterns := make([]string, m)
	total := 0
	for i := 0; i < m; i++ {
		remaining := maxTotalPatternLen - total - (m - i - 1)
		if remaining < 1 {
			remaining = 1
		}
		length := 1 + rng.Intn(10)
		if length > remaining {
			length = remaining
		}
		patterns[i] = randomPattern(rng, length)
		total += length
	}
	return makeTestCase(987654321012345678, patterns)
}

func makeLongPatternTest() testCase {
	pattern := strings.Repeat("a", maxTotalPatternLen)
	return makeTestCase(12, []string{pattern})
}

func makeHugeKTest(rng *rand.Rand) testCase {
	patterns := []string{
		"a",
		"b",
		randomPattern(rng, 3),
		strings.Repeat("ab", 500),
		strings.Repeat("b", 4000),
	}
	return makeTestCase(1_000_000_000_000_000_000, patterns)
}

func makeTestCase(k int64, patterns []string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", k, len(patterns))
	for _, p := range patterns {
		sb.WriteString(p)
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String()}
}
