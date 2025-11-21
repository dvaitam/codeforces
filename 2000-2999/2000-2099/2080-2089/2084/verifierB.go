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
	refSource    = "2000-2999/2000-2099/2080-2089/2084/2084B.go"
	maxTests     = 500
	totalNLimit  = 10000
	randomTrials = 400
)

type testCase struct {
	n int
	a []int64
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
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
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse reference output:", err)
		os.Exit(1)
	}

	candOut, err := runCandidate(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to parse candidate output:", err)
		os.Exit(1)
	}

	for i := range tests {
		if candAns[i] != refAns[i] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %s, got %s\n", i+1, boolToStr(refAns[i]), boolToStr(candAns[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests).\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2084B-ref-*")
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
		fmt.Fprintf(&b, "%d\n", tc.n)
		for i, v := range tc.a {
			if i > 0 {
				b.WriteByte(' ')
			}
			fmt.Fprintf(&b, "%d", v)
		}
		b.WriteByte('\n')
	}
	return b.String()
}

func parseOutput(out string, t int) ([]bool, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d tokens, got %d", t, len(fields))
	}
	res := make([]bool, t)
	for i, f := range fields {
		val, err := parseBoolToken(f)
		if err != nil {
			return nil, fmt.Errorf("line %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func parseBoolToken(s string) (bool, error) {
	low := strings.ToLower(strings.TrimSpace(s))
	switch low {
	case "yes", "y", "true", "1":
		return true, nil
	case "no", "n", "false", "0":
		return false, nil
	default:
		return false, fmt.Errorf("invalid boolean token %q", s)
	}
}

func boolToStr(v bool) string {
	if v {
		return "Yes"
	}
	return "No"
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	sumN := 0

	add := func(tc testCase) {
		if len(tests) >= maxTests || sumN+tc.n > totalNLimit {
			return
		}
		tests = append(tests, tc)
		sumN += tc.n
	}

	// Deterministic edge and illustrative cases.
	add(testCase{n: 2, a: []int64{1, 1}})
	add(testCase{n: 2, a: []int64{2, 3}})
	add(testCase{n: 3, a: []int64{3, 2, 2}})
	add(testCase{n: 5, a: []int64{3, 4, 5, 6, 9}})
	add(testCase{n: 4, a: []int64{5, 10, 15, 20}})
	add(testCase{n: 4, a: []int64{7, 11, 13, 17}})

	for attempts := 0; attempts < randomTrials && len(tests) < maxTests && sumN < totalNLimit; attempts++ {
		n := rng.Intn(200) + 2
		if sumN+n > totalNLimit {
			n = totalNLimit - sumN
		}

		a := make([]int64, n)
		base := rng.Int63n(1_000_000_000) + 1
		for i := 0; i < n; i++ {
			mode := rng.Intn(6)
			switch mode {
			case 0:
				a[i] = 1
			case 1:
				a[i] = base
			case 2:
				a[i] = int64(rng.Intn(15)+1) * base
			case 3:
				a[i] = int64(rng.Intn(1_000_000_000)+1) * int64(rng.Intn(1_000_000_000)+1) % 1_000_000_000_000_000_000
				if a[i] == 0 {
					a[i] = 1
				}
			case 4:
				a[i] = (int64(rng.Intn(1_000_000)+1) << 30) + int64(rng.Intn(1_000_000)+1)
			default:
				a[i] = rng.Int63n(1_000_000_000_000_000_000) + 1
			}
		}

		add(testCase{n: n, a: a})
	}

	if len(tests) == 0 {
		tests = append(tests, testCase{n: 2, a: []int64{1, 1}})
	}
	return tests
}
