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

const refSource = "2000-2999/2000-2099/2040-2049/2045/2045K.go"

type testCase struct {
	n int
	a []int
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierK.go /path/to/candidate")
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
	for i, tc := range tests {
		input := renderInput(tc)

		expOut, err := runWithInput(exec.Command(refBin), input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\noutput:\n%s\n", i+1, err, expOut)
			os.Exit(1)
		}
		gotOut, err := runWithInput(commandFor(candidate), input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate failed on test %d: %v\noutput:\n%s\ninput:\n%s", i+1, err, gotOut, input)
			os.Exit(1)
		}

		expVal, err := parseSingle(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", i+1, err, expOut)
			os.Exit(1)
		}
		gotVal, err := parseSingle(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d: %v\noutput:\n%s\ninput:\n%s", i+1, err, gotOut, input)
			os.Exit(1)
		}
		if expVal != gotVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d: expected %d got %d\ninput:\n%s", i+1, expVal, gotVal, input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2045K-ref-*")
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

func parseSingle(output string) (int64, error) {
	fields := strings.Fields(output)
	if len(fields) == 0 {
		return 0, fmt.Errorf("no output")
	}
	if len(fields) > 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.ParseInt(fields[0], 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse int: %v", err)
	}
	return val, nil
}

func renderInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(strconv.Itoa(tc.n))
	sb.WriteByte('\n')
	for i, v := range tc.a {
		if i > 0 {
			sb.WriteByte(' ')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('\n')
	return sb.String()
}

func generateTests() []testCase {
	tests := []testCase{
		{n: 3, a: []int{3, 3, 3}},                         // sample 1
		{n: 4, a: []int{2, 2, 4, 4}},                      // sample 2
		{n: 9, a: []int{4, 2, 6, 9, 7, 7, 7, 3, 3}},       // sample 3
		{n: 2, a: []int{1, 1}},                            // minimal
		{n: 6, a: []int{1, 2, 3, 4, 5, 6}},                // increasing
		{n: 8, a: []int{8, 8, 8, 8, 8, 8, 8, 8}},          // all equal big
		{n: 10, a: []int{2, 4, 6, 8, 10, 2, 4, 6, 8, 10}}, // even mix
	}

	totalN := 0
	for _, tc := range tests {
		totalN += tc.n
	}

	const maxTotal = 40000
	rng := rand.New(rand.NewSource(2045))

	for totalN < maxTotal {
		n := rng.Intn(3000) + 2
		if totalN+n > maxTotal {
			n = maxTotal - totalN
		}
		a := make([]int, n)
		for i := 0; i < n; i++ {
			a[i] = rng.Intn(n) + 1
		}
		tests = append(tests, testCase{n: n, a: a})
		totalN += n
		if len(tests) > 200 {
			break
		}
	}

	return tests
}
