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
	refSource   = "./2125A.go"
	maxTests    = 400
	totalLenLim = 10000
	alphabet    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

type testCase struct {
	s string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse candidate output:", err)
		os.Exit(1)
	}

	for i, tc := range tests {
		if err := validateAnswer(tc.s, candAns[i]); err != nil {
			fmt.Fprintf(os.Stderr, "test %d: %v\n", i+1, err)
			os.Exit(1)
		}
		if refAns[i] != "" && refAns[i] == candAns[i] {
			continue
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2125A-ref-*")
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

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%s\n", tc.s)
	}
	return b.String()
}

func parseOutput(out string, t int) ([]string, error) {
	lines := strings.Fields(out)
	if len(lines) != t {
		return nil, fmt.Errorf("expected %d strings, got %d", t, len(lines))
	}
	return lines, nil
}

func validateAnswer(original, cand string) error {
	if len(original) != len(cand) {
		return fmt.Errorf("length mismatch: expected %d got %d", len(original), len(cand))
	}
	if containsForbidden(cand) {
		return fmt.Errorf("contest remains difficult: contains FFT/NTT")
	}
	cntOrig := countLetters(original)
	cntCand := countLetters(cand)
	if cntOrig != cntCand {
		return fmt.Errorf("output is not a permutation of input")
	}
	return nil
}

func containsForbidden(s string) bool {
	return strings.Contains(s, "FFT") || strings.Contains(s, "NTT")
}

func countLetters(s string) [26]int {
	var cnt [26]int
	for i := 0; i < len(s); i++ {
		if s[i] < 'A' || s[i] > 'Z' {
			continue
		}
		cnt[s[i]-'A']++
	}
	return cnt
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	sumLen := 0

	add := func(s string) {
		if len(tests) >= maxTests || sumLen+len(s) > totalLenLim {
			return
		}
		tests = append(tests, testCase{s: s})
		sumLen += len(s)
	}

	add("FFT")
	add("NTT")
	add("A")
	add("TTT")
	add("FNFTT")
	add("FFTFFT")
	add("ABCD")

	for len(tests) < maxTests && sumLen < totalLenLim {
		length := rng.Intn(20) + 1
		if sumLen+length > totalLenLim {
			length = totalLenLim - sumLen
		}
		var b strings.Builder
		for i := 0; i < length; i++ {
			b.WriteByte(alphabet[rng.Intn(len(alphabet))])
		}
		add(b.String())
	}

	if len(tests) == 0 {
		tests = append(tests, testCase{s: "A"})
	}
	return tests
}
