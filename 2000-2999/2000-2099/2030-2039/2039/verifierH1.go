package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const referenceSource = "2000-2999/2000-2099/2030-2039/2039/2039H1.go"

type testCase struct {
	name  string
	input string
}

type caseData struct {
	n   int
	arr []int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierH1.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, cleanup, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	tests := buildTests()
	for idx, tc := range tests {
		meta, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse input for test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		_, _ = runProgram(refBin, tc.input) // ensure reference builds and runs

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := validateOutput(meta, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-2039H1-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref2039H1.bin")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSource)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
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

func parseInput(input string) ([]caseData, error) {
	reader := strings.NewReader(input)
	in := bufio.NewReader(reader)
	var t int
	if _, err := fmt.Fscan(in, &t); err != nil {
		return nil, fmt.Errorf("failed to read t: %v", err)
	}
	cases := make([]caseData, t)
	for i := 0; i < t; i++ {
		var n int
		if _, err := fmt.Fscan(in, &n); err != nil {
			return nil, fmt.Errorf("test %d: failed to read n: %v", i+1, err)
		}
		arr := make([]int, n)
		for j := 0; j < n; j++ {
			if _, err := fmt.Fscan(in, &arr[j]); err != nil {
				return nil, fmt.Errorf("test %d: failed to read a[%d]: %v", i+1, j+1, err)
			}
		}
		cases[i] = caseData{n: n, arr: arr}
	}
	return cases, nil
}

func validateOutput(cases []caseData, output string) error {
	reader := bufio.NewReader(strings.NewReader(output))
	for idx, cs := range cases {
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return fmt.Errorf("test %d: failed to read k: %v", idx+1, err)
		}
		if k < 0 || k > 2*cs.n+4 {
			return fmt.Errorf("test %d: k=%d out of range", idx+1, k)
		}
		arr := append([]int(nil), cs.arr...)
		for w := 0; w < k; w++ {
			var path string
			if _, err := fmt.Fscan(reader, &path); err != nil {
				return fmt.Errorf("test %d: failed to read path %d: %v", idx+1, w+1, err)
			}
			if err := applyPath(arr, cs.n, path); err != nil {
				return fmt.Errorf("test %d, path %d: %v", idx+1, w+1, err)
			}
		}
		if !isSorted(arr) {
			return fmt.Errorf("test %d: array not sorted after operations", idx+1)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return fmt.Errorf("extra output detected: %s", extra)
	} else if err != nil && err != io.EOF {
		return fmt.Errorf("failed to parse trailing output: %v", err)
	}
	return nil
}

func applyPath(arr []int, n int, path string) error {
	if len(path) != 2*n-2 {
		return fmt.Errorf("path length %d invalid", len(path))
	}
	rCount := strings.Count(path, "R")
	dCount := strings.Count(path, "D")
	if rCount != n-1 || dCount != n-1 {
		return fmt.Errorf("path must have %d R and %d D, got %d and %d", n-1, n-1, rCount, dCount)
	}
	x, y := 1, 1
	for i := 0; i < len(path); i++ {
		ch := path[i]
		switch ch {
		case 'R':
			y++
		case 'D':
			x++
		default:
			return fmt.Errorf("invalid character %q", ch)
		}
		if x < 1 || x > n || y < 1 || y > n {
			return fmt.Errorf("step %d out of bounds (%d,%d)", i+1, x, y)
		}
		if x != y {
			arr[x-1], arr[y-1] = arr[y-1], arr[x-1]
		}
	}
	if x != n || y != n {
		return fmt.Errorf("path does not end at (%d,%d)", n, n)
	}
	return nil
}

func isSorted(arr []int) bool {
	for i := 1; i < len(arr); i++ {
		if arr[i-1] > arr[i] {
			return false
		}
	}
	return true
}

func buildTests() []testCase {
	var tests []testCase
	tests = append(tests, testCase{
		name:  "sample",
		input: "3\n2\n1 2\n3\n2 1 3\n4\n3 2 3 4\n",
	})
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		tests = append(tests, testCase{name: fmt.Sprintf("random-%d", i+1), input: randomTest(rng)})
	}
	return tests
}

func randomTest(rng *rand.Rand) string {
	t := rng.Intn(3) + 1
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", t))
	for caseIdx := 0; caseIdx < t; caseIdx++ {
		n := rng.Intn(5) + 2
		sb.WriteString(fmt.Sprintf("%d\n", n))
		for i := 0; i < n; i++ {
			if i > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(strconv.Itoa(rng.Intn(n) + 1))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}
