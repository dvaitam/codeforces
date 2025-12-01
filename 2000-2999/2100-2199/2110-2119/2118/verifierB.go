package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "./2118B.go"

type testCase struct {
	ns    []int
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	candidate := os.Args[1]
	tests := generateTests()

	for i, tc := range tests {
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
		if err := validateOutput(tc, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, tc.input, refOut)
			os.Exit(1)
		}

		got, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d: %v\ninput:\n%soutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}

		if err := validateOutput(tc, got); err != nil {
			fmt.Fprintf(os.Stderr, "wrong output on test %d: %v\ninput:\n%s\noutput:\n%s\n", i+1, err, tc.input, got)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2118B-ref-*")
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
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func runCandidate(path, input string) (string, error) {
	cmd := commandFor(path)
	return runWithInput(cmd, input)
}

func runProgram(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return strings.TrimSpace(out.String()), err
}

func generateTests() []testCase {
	rng := rand.New(rand.NewSource(21182118))
	var tests []testCase

	// Simple small cases.
	tests = append(tests, makeTest([]int{3}))
	tests = append(tests, makeTest([]int{4}))
	tests = append(tests, makeTest([]int{5, 3}))

	// Random moderate cases keeping sum n within limit.
	for i := 0; i < 20; i++ {
		t := rng.Intn(4) + 1
		var ns []int
		total := 0
		for j := 0; j < t; j++ {
			n := rng.Intn(40) + 3
			if total+n > 5000 {
				break
			}
			total += n
			ns = append(ns, n)
		}
		if len(ns) > 0 {
			tests = append(tests, makeTest(ns))
		}
	}

	// Larger structured cases.
	tests = append(tests, makeTest([]int{100, 120}))
	tests = append(tests, makeTest([]int{250, 250}))
	tests = append(tests, makeTest([]int{500}))

	// Stress near limits.
	tests = append(tests, makeTest([]int{1000, 800, 700}))
	tests = append(tests, makeTest([]int{5000}))

	return tests
}

func makeTest(ns []int) testCase {
	return testCase{ns: ns, input: buildInput(ns)}
}

func buildInput(ns []int) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(ns))
	for _, n := range ns {
		fmt.Fprintf(&b, "%d\n", n)
	}
	return b.String()
}

// validateOutput checks that the produced operations transform the matrix so each column is a permutation.
func validateOutput(tc testCase, out string) error {
	reader := bufio.NewReader(strings.NewReader(out))

	tOut, err := nextInt(reader)
	if err != nil {
		return fmt.Errorf("failed to read number of test cases: %v", err)
	}
	if tOut != len(tc.ns) {
		return fmt.Errorf("expected %d test cases, got %d", len(tc.ns), tOut)
	}

	for idx, n := range tc.ns {
		k, err := nextInt(reader)
		if err != nil {
			return fmt.Errorf("test %d: failed to read k: %v", idx+1, err)
		}
		if k < 0 || k > 2*n {
			return fmt.Errorf("test %d: k=%d out of bounds (0..%d)", idx+1, k, 2*n)
		}

		rows := make([][]int, n)

		for op := 0; op < k; op++ {
			i, err := nextInt(reader)
			if err != nil {
				return fmt.Errorf("test %d: failed to read op %d row: %v", idx+1, op+1, err)
			}
			l, err := nextInt(reader)
			if err != nil {
				return fmt.Errorf("test %d: failed to read op %d l: %v", idx+1, op+1, err)
			}
			r, err := nextInt(reader)
			if err != nil {
				return fmt.Errorf("test %d: failed to read op %d r: %v", idx+1, op+1, err)
			}
			if i < 1 || i > n || l < 1 || r < l || r > n {
				return fmt.Errorf("test %d: op %d has invalid indices i=%d l=%d r=%d", idx+1, op+1, i, l, r)
			}
			rowIdx := i - 1
			if rows[rowIdx] == nil {
				rows[rowIdx] = identityRow(n)
			}
			reverseSubarray(rows[rowIdx], l-1, r-1)
		}

		if err := checkColumns(rows, n); err != nil {
			return fmt.Errorf("test %d: %v", idx+1, err)
		}
	}

	// Ensure no extra integers in output (ignore trailing whitespace).
	if extra, err := nextInt(reader); err == nil {
		return fmt.Errorf("unexpected extra output after all test cases, next int=%d", extra)
	}

	return nil
}

func identityRow(n int) []int {
	row := make([]int, n)
	for i := range row {
		row[i] = i + 1
	}
	return row
}

func reverseSubarray(a []int, l, r int) {
	for l < r {
		a[l], a[r] = a[r], a[l]
		l++
		r--
	}
}

func checkColumns(rows [][]int, n int) error {
	visit := make([]int, n+1)
	iter := 0
	for col := 0; col < n; col++ {
		iter++
		for row := 0; row < n; row++ {
			val := col + 1
			if rows[row] != nil {
				val = rows[row][col]
			}
			if val < 1 || val > n {
				return fmt.Errorf("column %d row %d value %d out of range", col+1, row+1, val)
			}
			if visit[val] == iter {
				return fmt.Errorf("column %d has duplicate value %d", col+1, val)
			}
			visit[val] = iter
		}
	}
	return nil
}

func nextInt(r *bufio.Reader) (int, error) {
	sign, val, saw := 1, 0, false
	for {
		c, err := r.ReadByte()
		if err != nil {
			return 0, err
		}
		if c == '-' {
			sign = -1
			for {
				c, err = r.ReadByte()
				if err != nil {
					return 0, err
				}
				if c >= '0' && c <= '9' {
					break
				}
				if c > ' ' {
					return 0, fmt.Errorf("unexpected character %q", c)
				}
			}
		}
		if c >= '0' && c <= '9' {
			val = int(c - '0')
			saw = true
			for {
				c, err = r.ReadByte()
				if err != nil {
					if err.Error() == "EOF" {
						return sign * val, nil
					}
					return 0, err
				}
				if c < '0' || c > '9' {
					if err := r.UnreadByte(); err != nil {
						return 0, err
					}
					break
				}
				val = val*10 + int(c-'0')
			}
			return sign * val, nil
		}
		if c > ' ' {
			return 0, fmt.Errorf("unexpected character %q", c)
		}
		if saw {
			break
		}
	}
	return sign * val, nil
}
