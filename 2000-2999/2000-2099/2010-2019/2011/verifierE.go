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

const refSource = "2011E.go"
const totalLimit = 200000

type testCase struct {
	s string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected := tokenize(refOut)

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	got := tokenize(candOut)

	if len(expected) != len(got) {
		fmt.Fprintf(os.Stderr, "wrong number of tokens: expected %d got %d\n", len(expected), len(got))
		os.Exit(1)
	}
	for i := range expected {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "mismatch at token %d: expected %q got %q\n", i+1, expected[i], got[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2011E-ref-*")
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

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	err := cmd.Run()
	if err != nil {
		out.WriteString(errBuf.String())
		return out.String(), err
	}
	if errBuf.Len() > 0 {
		out.WriteString(errBuf.String())
	}
	return out.String(), nil
}

func tokenize(s string) []string {
	return strings.Fields(strings.TrimSpace(s))
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%s\n", tc.s)
	}
	return b.String()
}

func generateTests() []testCase {
	var tests []testCase
	total := 0

	add := func(s string) {
		if len(s) == 0 {
			return
		}
		if total+len(s) > totalLimit {
			return
		}
		tests = append(tests, testCase{s: s})
		total += len(s)
	}

	deterministic := []string{
		"R",
		"P",
		"S",
		"RS",
		"PR",
		"SP",
		"RPS",
		"SRP",
		"RSPPRSP",
		strings.Repeat("R", 20),
		strings.Repeat("S", 21),
		"RPRPRPRP",
		"SSRRPSRPPPSPPRSPPPPP", // sample
	}
	for _, s := range deterministic {
		add(s)
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for total < totalLimit {
		remaining := totalLimit - total
		maxLen := 20000
		if remaining < maxLen {
			maxLen = remaining
		}
		if maxLen == 0 {
			break
		}
		n := rng.Intn(maxLen) + 1
		b := make([]byte, n)
		for i := range b {
			r := rng.Intn(3)
			if r == 0 {
				b[i] = 'R'
			} else if r == 1 {
				b[i] = 'P'
			} else {
				b[i] = 'S'
			}
		}
		add(string(b))
	}

	return tests
}
