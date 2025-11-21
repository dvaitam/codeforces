package main

import (
	"bufio"
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

const refSource = "2000-2999/2100-2199/2120-2129/2129/2129C3.go"

type testCase struct {
	n int
	s string
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2129C3-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		abs, err := filepath.Abs(bin)
		if err != nil {
			return "", err
		}
		cmd = exec.Command("go", "run", abs)
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
	return stdout.String(), nil
}

func deterministicTests() []testCase {
	return []testCase{
		{n: 2, s: "()"},
		{n: 3, s: "(()"},
		{n: 4, s: "())("},
		{n: 5, s: "((())"},
	}
}

func randomBracket(rng *rand.Rand, n int) string {
	for {
		b := make([]byte, n)
		hasL, hasR := false, false
		for i := 0; i < n; i++ {
			if rng.Intn(2) == 0 {
				b[i] = '('
				hasL = true
			} else {
				b[i] = ')'
				hasR = true
			}
		}
		if hasL && hasR {
			return string(b)
		}
	}
}

func randomTest(rng *rand.Rand, n int) testCase {
	return testCase{n: n, s: randomBracket(rng, n)}
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

func parseReference(out string, t int) ([]string, error) {
	sc := bufio.NewScanner(strings.NewReader(out))
	sc.Split(bufio.ScanWords)
	res := make([]string, t)
	for i := 0; i < t; i++ {
		if !sc.Scan() {
			return nil, fmt.Errorf("reference: missing output for test %d", i+1)
		}
		res[i] = sc.Text()
	}
	if sc.Scan() {
		return nil, fmt.Errorf("reference: extra output detected")
	}
	return res, nil
}

func isBracketString(s string, n int) bool {
	if len(s) != n {
		return false
	}
	for i := 0; i < n; i++ {
		if s[i] != '(' && s[i] != ')' {
			return false
		}
	}
	return true
}

func parseCandidate(out string, tests []testCase) ([]string, error) {
	fields := strings.Fields(out)
	res := make([]string, len(tests))
	idx := 0
	for _, tok := range fields {
		if idx >= len(tests) {
			break
		}
		n := tests[idx].n
		if tok == "?" || tok == "!" {
			continue
		}
		if isBracketString(tok, n) {
			res[idx] = tok
			idx++
		}
	}
	if idx != len(tests) {
		return nil, fmt.Errorf("candidate: expected %d bracket strings, got %d", len(tests), idx)
	}
	return res, nil
}

func totalN(tests []testCase) int {
	sum := 0
	for _, tc := range tests {
		sum += tc.n
	}
	return sum
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC3.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := deterministicTests()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 30 && totalN(tests) < 15000 {
		n := rng.Intn(50) + 2
		tests = append(tests, randomTest(rng, n))
	}
	for len(tests) < 50 && totalN(tests) < 25000 {
		n := rng.Intn(1000) + 100
		tests = append(tests, randomTest(rng, n))
	}

	input := buildInput(tests)

	wantOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\n", err)
		os.Exit(1)
	}
	want, err := parseReference(wantOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n", err)
		os.Exit(1)
	}

	gotOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	got, err := parseCandidate(gotOut, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if want[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d failed:\nexpected: %s\ngot:      %s\nn=%d\n", i+1, want[i], got[i], tests[i].n)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed\n", len(tests))
}
