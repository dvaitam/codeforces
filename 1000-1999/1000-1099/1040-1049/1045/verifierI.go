package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const refSource = "1000-1999/1000-1099/1040-1049/1045/1045I.go"

type testCase struct {
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierI.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for i, tc := range tests {
		want, err := runProgram(refBin, tc.input)
		if err != nil {
			fail("reference runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		got, err := runProgram(candidate, tc.input)
		if err != nil {
			fail("candidate runtime error on test %d: %v\ninput:\n%s", i+1, err, tc.input)
		}
		if normalize(got) != normalize(want) {
			fail("wrong answer on test %d\ninput:\n%s\nexpected:\n%s\ngot:\n%s", i+1, tc.input, want, got)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1045I-ref-*")
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
	return strings.TrimSpace(s)
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	tests = append(tests, makeTest([]string{"a"}))
	tests = append(tests, makeTest([]string{"a", "b"}))
	tests = append(tests, makeTest([]string{"ab", "ba"}))
	tests = append(tests, makeTest([]string{"ab", "abc", "c"}))
	tests = append(tests, makeTest([]string{"abc", "bca", "cab", "abc"}))
	tests = append(tests, makeTest([]string{"aaaa", "bbbb", "cccc", "dddd"}))
	tests = append(tests, makeTest([]string{"aa", "bb", "cc", "abc", "cba"}))

	for size := 1; size <= 6; size++ {
		tests = append(tests, randomTest(rng, size, 5))
	}
	for i := 0; i < 30; i++ {
		tests = append(tests, randomTest(rng, rng.Intn(60)+1, 10))
	}
	for i := 0; i < 20; i++ {
		tests = append(tests, randomTest(rng, rng.Intn(600)+50, 40))
	}
	tests = append(tests, randomTest(rng, 2000, 100))
	tests = append(tests, randomTest(rng, 5000, 200))
	tests = append(tests, randomTest(rng, 100000, 200))
	return tests
}

func makeTest(words []string) testCase {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(words))
	for _, w := range words {
		sb.WriteString(w)
		sb.WriteByte('\n')
	}
	return testCase{input: sb.String()}
}

func randomTest(rng *rand.Rand, n, maxLen int) testCase {
	if maxLen <= 0 {
		maxLen = 1
	}
	words := make([]string, 0, n)
	totalChars := 0
	for len(words) < n && totalChars < 950000 {
		remaining := 950000 - totalChars
		if remaining <= 0 {
			break
		}
		length := rng.Intn(maxLen) + 1
		if length > remaining {
			length = remaining
		}
		b := make([]byte, length)
		for j := 0; j < length; j++ {
			b[j] = byte('a' + rng.Intn(26))
		}
		words = append(words, string(b))
		totalChars += length
	}
	for len(words) < n {
		words = append(words, "a")
	}
	return makeTest(words)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
