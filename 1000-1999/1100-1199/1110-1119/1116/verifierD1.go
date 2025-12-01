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
	"time"
)

const (
	refSource = "./1116D1.go"
	threshold = 1e-5
)

type testCase struct {
	name string
	n    int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		refOut, err := runBinary(refBin, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		refMask, err := analyzePattern(tc.n, refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, refOut)
			os.Exit(1)
		}

		candOut, err := runBinaryWithTarget(candidate, tc.n)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}
		candMask, err := analyzePattern(tc.n, candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, candOut)
			os.Exit(1)
		}

		if !equalPatterns(refMask, candMask) {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): pattern mismatch for N=%d\n", idx+1, tc.name, tc.n)
			fmt.Fprintf(os.Stderr, "expected:\n%s\n", renderPattern(tc.n, refMask))
			fmt.Fprintf(os.Stderr, "got:\n%s\n", renderPattern(tc.n, candMask))
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1116D1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func runBinary(bin string, n int) (string, error) {
	cmd := exec.Command(bin, strconv.Itoa(n))
	return runWithCmd(cmd)
}

func runBinaryWithTarget(path string, n int) (string, error) {
	cmd := commandFor(path, n)
	return runWithCmd(cmd)
}

func commandFor(path string, n int) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path, strconv.Itoa(n))
	}
	return exec.Command(path, strconv.Itoa(n))
}

func runWithCmd(cmd *exec.Cmd) (string, error) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func analyzePattern(n int, out string) ([][]bool, error) {
	values, err := parseFloats(out)
	if err != nil {
		return nil, err
	}
	size := 1 << n
	expected := size * size
	if len(values) != expected {
		return nil, fmt.Errorf("expected %d entries, got %d", expected, len(values))
	}
	mask := make([][]bool, size)
	for i := range mask {
		mask[i] = make([]bool, size)
	}
	idx := 0
	for col := 0; col < size; col++ {
		for row := 0; row < size; row++ {
			val := values[idx]
			mask[row][col] = val*val >= threshold
			idx++
		}
	}
	return mask, nil
}

func parseFloats(out string) ([]float64, error) {
	fields := strings.Fields(out)
	vals := make([]float64, len(fields))
	for i, f := range fields {
		v, err := strconv.ParseFloat(f, 64)
		if err != nil {
			return nil, fmt.Errorf("invalid float %q: %v", f, err)
		}
		vals[i] = v
	}
	return vals, nil
}

func equalPatterns(a, b [][]bool) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}

func renderPattern(n int, mask [][]bool) string {
	var sb strings.Builder
	size := 1 << n
	for row := 0; row < size; row++ {
		for col := 0; col < size; col++ {
			if mask[row][col] {
				sb.WriteByte('X')
			} else {
				sb.WriteByte('.')
			}
		}
		if row+1 < size {
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func generateTests() []testCase {
	tests := []testCase{
		{name: "min-n", n: 2},
		{name: "small-n", n: 3},
		{name: "medium-n", n: 4},
		{name: "max-n", n: 5},
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 10; i++ {
		tests = append(tests, testCase{
			name: fmt.Sprintf("random-%d", i+1),
			n:    rng.Intn(4) + 2,
		})
	}
	return tests
}
