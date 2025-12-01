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

const refSource = "./2053H.go"

type testCase struct {
	n int
	w int64
	a []int64
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[len(os.Args)-1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	input := renderInput(tests)

	refOut, err := runWithInput(exec.Command(refBin), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference solution failed: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	candOut, err := runWithInput(commandFor(candidate), input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate failed: %v\noutput:\n%s\n", err, candOut)
		os.Exit(1)
	}

	expect, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse reference output: %v\n", err)
		os.Exit(1)
	}
	got, err := parseOutputs(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not parse candidate output: %v\n", err)
		os.Exit(1)
	}

	for i := range tests {
		if expect[i][0] != got[i][0] || expect[i][1] != got[i][1] {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d %d got %d %d\ninput:\n%s", i+1, expect[i][0], expect[i][1], got[i][0], got[i][1], formatSingleInput(tests[i]))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2053H-ref-*")
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

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutputs(output string, t int) ([][2]int64, error) {
	lines := strings.Fields(output)
	if len(lines) < 2*t {
		return nil, fmt.Errorf("expected at least %d numbers, got %d", 2*t, len(lines))
	}
	if len(lines) > 2*t {
		return nil, fmt.Errorf("extra output detected after %d numbers", 2*t)
	}
	res := make([][2]int64, t)
	for i := 0; i < t; i++ {
		a, err := strconv.ParseInt(lines[2*i], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number %d (%q): %v", 2*i+1, lines[2*i], err)
		}
		b, err := strconv.ParseInt(lines[2*i+1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number %d (%q): %v", 2*i+2, lines[2*i+1], err)
		}
		res[i] = [2]int64{a, b}
	}
	return res, nil
}

func renderInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(len(tests)))
	sb.WriteByte('\n')
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.w))
		for i, v := range tc.a {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.FormatInt(v, 10))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func formatSingleInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.w))
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.FormatInt(v, 10))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 5, w: 8, a: []int64{1, 2, 3, 4, 5}},                                   // sample 1
		{n: 7, w: 5, a: []int64{3, 1, 2, 3, 4, 1, 1}},                             // sample 2
		{n: 1, w: 10, a: []int64{5}},                                              // single element
		{n: 4, w: 2, a: []int64{1, 1, 1, 1}},                                      // all equal min w
		{n: 6, w: 10, a: []int64{10, 10, 10, 10, 10, 10}},                         // all at max
		{n: 6, w: 3, a: []int64{1, 2, 2, 2, 2, 3}},                                // long equal block
		{n: 8, w: 100, a: []int64{5, 5, 6, 6, 7, 8, 9, 9}},                        // multiple equal pairs
		{n: 10, w: 100000000, a: []int64{1, 100, 1, 100, 1, 100, 1, 100, 1, 100}}, // large w
	}

	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	const maxTotal = 150000
	rng := rand.New(rand.NewSource(2053))
	for totalN < maxTotal {
		n := rng.Intn(4000) + 1
		if totalN+n > maxTotal {
			n = maxTotal - totalN
		}
		w := int64(rng.Intn(1_000_000_000-1) + 2)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			a[i] = int64(rng.Intn(int(w)) + 1)
		}
		tests = append(tests, testCase{n: n, w: w, a: a})
		totalN += n
		if len(tests) > 400 {
			break
		}
	}

	// Add one large n test with small w to stress equal segments.
	if totalN < maxTotal {
		n := maxTotal - totalN
		if n > 200000 {
			n = 200000
		}
		w := int64(2)
		a := make([]int64, n)
		for i := 0; i < n; i++ {
			if i%2 == 0 {
				a[i] = 1
			} else {
				a[i] = 2
			}
		}
		tests = append(tests, testCase{n: n, w: w, a: a})
	}

	return tests
}
