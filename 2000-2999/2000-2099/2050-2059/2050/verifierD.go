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
	refSource   = "2050D.go"
	randomCases = 120
	maxTotalLen = 200000
)

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

	expectRaw, err := runProgram(refBin, input)
	if err != nil {
		fail("reference failed: %v\n%s", err, expectRaw)
	}
	gotRaw, err := runCandidate(candidate, input)
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

	for i := range tests {
		if expect[i] != got[i] {
			fail("mismatch on test %d\nexpected: %s\n   found: %s", i+1, expect[i], got[i])
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2050D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	source := filepath.Join(".", refSource)
	cmd := exec.Command("go", "build", "-o", tmp.Name(), source)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func generateTests() []string {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests := make([]string, 0, randomCases+4)

	// Include sample tests to catch formatting issues quickly.
	samples := []string{"19", "1709", "115555", "51476", "9876543210", "5891917899"}
	tests = append(tests, samples...)

	// Add a few deterministic edge patterns.
	tests = append(tests, "1")                      // single digit
	tests = append(tests, "90")                     // leading high digit with zero
	tests = append(tests, strings.Repeat("9", 200)) // long uniform string

	remaining := maxTotalLen
	for _, s := range tests {
		remaining -= len(s)
	}

	for i := 0; i < randomCases && remaining > 0; i++ {
		// Keep lengths varied while respecting the total length limit.
		maxLen := remaining
		if maxLen > 500 { // cap for diversity and speed
			maxLen = 500
		}
		length := rng.Intn(maxLen) + 1
		remaining -= length

		var sb strings.Builder
		sb.Grow(length)
		sb.WriteByte(byte('1' + rng.Intn(9))) // leading digit cannot be zero
		for j := 1; j < length; j++ {
			sb.WriteByte(byte('0' + rng.Intn(10)))
		}
		tests = append(tests, sb.String())
	}

	return tests
}

func buildInput(tests []string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, s := range tests {
		sb.WriteString(s)
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
		return errBuf.String(), err
	}
	if errBuf.Len() > 0 {
		return errBuf.String(), fmt.Errorf("stderr not empty")
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

func parseOutputs(out string, t int) ([]string, error) {
	tokens := strings.Fields(out)
	if len(tokens) < t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(tokens))
	}
	if len(tokens) > t {
		return nil, fmt.Errorf("too many output tokens: expected %d, got %d", t, len(tokens))
	}
	res := make([]string, t)
	copy(res, tokens[:t])
	return res, nil
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
