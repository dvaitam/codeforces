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
	refSource   = "./2120A.go"
	targetTests = 160
	maxTests    = 500
	maxDim      = 100
)

type testCase struct {
	l1 int
	b1 int
	l2 int
	b2 int
	l3 int
	b3 int
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
	refAns, err := parseAnswers(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseAnswers(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	if len(refAns) != len(candAns) {
		fmt.Fprintf(os.Stderr, "answer count mismatch: expected %d, got %d\n", len(refAns), len(candAns))
		os.Exit(1)
	}
	for i := range refAns {
		if refAns[i] != candAns[i] {
			fmt.Fprintf(os.Stderr, "mismatch on test %d: expected %v, got %v\n", i+1, refAns[i], candAns[i])
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2120A-ref-*")
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

func parseAnswers(out string, t int) ([]bool, error) {
	lines := strings.Fields(out)
	if len(lines) < t {
		return nil, fmt.Errorf("expected %d answers, got %d tokens", t, len(lines))
	}
	ans := make([]bool, 0, t)
	for _, tok := range lines {
		low := strings.ToLower(tok)
		if low == "yes" {
			ans = append(ans, true)
		} else if low == "no" {
			ans = append(ans, false)
		}
		if len(ans) == t {
			break
		}
	}
	if len(ans) != t {
		return nil, fmt.Errorf("expected %d answers, parsed %d", t, len(ans))
	}
	return ans, nil
}

func buildInput(tests []testCase) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&b, "%d %d %d %d %d %d\n", tc.l1, tc.b1, tc.l2, tc.b2, tc.l3, tc.b3)
	}
	return b.String()
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase

	add := func(tc testCase) {
		if len(tests) >= maxTests {
			return
		}
		tests = append(tests, tc)
	}

	// Sample-inspired cases.
	add(testCase{100, 100, 10, 10, 1, 1})
	add(testCase{5, 3, 5, 1, 5, 1})
	add(testCase{2, 3, 1, 2, 1, 1})
	add(testCase{8, 5, 3, 5, 3, 3})
	add(testCase{3, 3, 3, 2, 1, 1})

	for len(tests) < targetTests {
		// generate sorted dimensions to match constraints l3<=l2<=l1 etc.
		l := randomSorted(rng)
		b := randomSorted(rng)
		add(testCase{
			l1: l[2], b1: b[2],
			l2: l[1], b2: b[1],
			l3: l[0], b3: b[0],
		})
	}

	return tests
}

func randomSorted(rng *rand.Rand) [3]int {
	arr := [3]int{}
	for i := 0; i < 3; i++ {
		arr[i] = rng.Intn(maxDim) + 1
	}
	// simple bubble sort for 3 elements
	for i := 0; i < 3; i++ {
		for j := i + 1; j < 3; j++ {
			if arr[i] > arr[j] {
				arr[i], arr[j] = arr[j], arr[i]
			}
		}
	}
	return arr
}
