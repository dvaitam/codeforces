package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const (
	randomCases = 80
	maxALen     = 1000
	maxBLen     = 1000
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierE /path/to/candidate")
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

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refAnswers := strings.Fields(refOut)

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	candAnswers := strings.Fields(candOut)

	numTests := len(tests)
	if len(refAnswers) != numTests {
		fail("reference produced %d answers, expected %d", len(refAnswers), numTests)
	}
	if len(candAnswers) != numTests {
		fail("candidate produced %d answers, expected %d", len(candAnswers), numTests)
	}

	for i := 0; i < numTests; i++ {
		if refAnswers[i] != candAnswers[i] {
			fail("test %d: expected %s got %s\n  a=%s\n  b=%s\n  c=%s",
				i+1, refAnswers[i], candAnswers[i],
				tests[i][0], tests[i][1], tests[i][2])
		}
	}

	fmt.Printf("All %d tests passed.\n", numTests)
}

// Each test is [3]string{a, b, c}.
type testCase [3]string

func buildReference() (string, error) {
	refSource := os.Getenv("REFERENCE_SOURCE_PATH")
	if refSource == "" {
		fail("REFERENCE_SOURCE_PATH environment variable not set")
	}

	tmp, err := os.CreateTemp("", "2050E-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func randString(rng *rand.Rand, length int) string {
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(byte('a' + rng.Intn(26)))
	}
	return sb.String()
}

// Generate c by interleaving a and b, then optionally changing some characters.
func interleave(rng *rand.Rand, a, b string, changes int) string {
	c := make([]byte, len(a)+len(b))
	ia, ib := 0, 0
	for i := 0; i < len(c); i++ {
		if ia >= len(a) {
			c[i] = b[ib]
			ib++
		} else if ib >= len(b) {
			c[i] = a[ia]
			ia++
		} else if rng.Intn(2) == 0 {
			c[i] = a[ia]
			ia++
		} else {
			c[i] = b[ib]
			ib++
		}
	}
	// Apply random changes.
	for ch := 0; ch < changes; ch++ {
		pos := rng.Intn(len(c))
		newChar := byte('a' + rng.Intn(26))
		c[pos] = newChar
	}
	return string(c)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	// Sample tests from the problem statement.
	samples := []testCase{
		{"abc", "ba", "abcba"},
		{"a", "b", "ab"},
		{"baa", "aab", "baaaab"},
		{"xxx", "yyy", "xyxyxy"},
	}
	tests = append(tests, samples...)

	// Edge cases.
	// Single character strings.
	tests = append(tests, testCase{"a", "a", "aa"})
	tests = append(tests, testCase{"a", "b", "ba"})
	tests = append(tests, testCase{"a", "a", "bb"})

	// Same characters.
	tests = append(tests, testCase{"aaa", "aaa", "aaaaaa"})

	// Completely different.
	tests = append(tests, testCase{"aaa", "bbb", "ababab"})

	// Max changes needed.
	tests = append(tests, testCase{"abc", "def", "zzzzzz"})

	// Keep total |a| and |b| sums within 2000 each for the batch.
	totalA, totalB := 0, 0
	for _, t := range tests {
		totalA += len(t[0])
		totalB += len(t[1])
	}

	for i := 0; i < randomCases; i++ {
		remainA := maxALen*2 - totalA
		remainB := maxBLen*2 - totalB
		if remainA < 2 || remainB < 2 {
			break
		}

		lenA := rng.Intn(min(remainA, 50)) + 1
		lenB := rng.Intn(min(remainB, 50)) + 1

		a := randString(rng, lenA)
		b := randString(rng, lenB)

		// Vary number of changes: sometimes 0, sometimes a few.
		changes := rng.Intn(lenA + lenB + 1)
		if rng.Intn(3) == 0 {
			changes = 0 // perfect interleaving, no changes
		}
		c := interleave(rng, a, b, changes)

		tests = append(tests, testCase{a, b, c})
		totalA += lenA
		totalB += lenB
	}

	return tests
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, t := range tests {
		sb.WriteString(t[0])
		sb.WriteByte('\n')
		sb.WriteString(t[1])
		sb.WriteByte('\n')
		sb.WriteString(t[2])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\nstderr: %s", err, errBuf.String())
	}
	return out.String(), nil
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\nstderr: %s", err, errBuf.String())
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

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
