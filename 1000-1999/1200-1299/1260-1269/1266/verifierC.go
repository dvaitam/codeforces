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

type testCase struct {
	r int
	c int
}

type expectedResult struct {
	noSolution bool
	magnitude  int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	for i, tc := range tests {
		input := fmt.Sprintf("%d %d\n", tc.r, tc.c)

		refOut, err := runBinary(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d (%dx%d): %v\n", i+1, tc.r, tc.c, err)
			os.Exit(1)
		}
		exp, err := analyzeReferenceOutput(refOut, tc)
		if err != nil {
			fmt.Fprintf(os.Stderr, "unable to analyze reference output on test %d (%dx%d): %v\n", i+1, tc.r, tc.c, err)
			os.Exit(1)
		}

		got, err := runBinary(target, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "runtime error on test %d (%dx%d): %v\n", i+1, tc.r, tc.c, err)
			os.Exit(1)
		}

		if exp.noSolution {
			if strings.TrimSpace(got) != "0" {
				fmt.Fprintf(os.Stderr, "test %d (%dx%d): expected output 0 but got:\n%s\n", i+1, tc.r, tc.c, got)
				os.Exit(1)
			}
			continue
		}

		if err := validateCandidate(got, tc, exp.magnitude); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%dx%d) failed: %v\ncandidate output:\n%s\n", i+1, tc.r, tc.c, err, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1266C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean("1000-1999/1200-1299/1260-1269/1266/1266C.go"))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build failed: %v\n%s", err, out)
	}
	return tmp.Name(), nil
}

func runBinary(path, input string) (string, error) {
	cmd := commandFor(path)
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, stderr.String())
	}
	return stdout.String(), nil
}

func commandFor(path string) *exec.Cmd {
	if strings.HasSuffix(path, ".go") {
		return exec.Command("go", "run", path)
	}
	return exec.Command(path)
}

func analyzeReferenceOutput(out string, tc testCase) (expectedResult, error) {
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return expectedResult{}, fmt.Errorf("empty output")
	}
	if trimmed == "0" {
		return expectedResult{noSolution: true}, nil
	}
	mat, err := parseMatrix(out, tc.r, tc.c)
	if err != nil {
		return expectedResult{}, err
	}
	gcds, mag, err := computeStats(mat)
	if err != nil {
		return expectedResult{}, err
	}
	if err := ensureDistinct(gcds); err != nil {
		return expectedResult{}, err
	}
	return expectedResult{magnitude: mag}, nil
}

func validateCandidate(out string, tc testCase, expectedMag int) error {
	trimmed := strings.TrimSpace(out)
	if trimmed == "" {
		return fmt.Errorf("empty output")
	}
	if trimmed == "0" {
		return fmt.Errorf("printed 0 but a solution exists")
	}
	mat, err := parseMatrix(out, tc.r, tc.c)
	if err != nil {
		return err
	}
	gcds, mag, err := computeStats(mat)
	if err != nil {
		return err
	}
	if err := ensureDistinct(gcds); err != nil {
		return err
	}
	if mag != expectedMag {
		return fmt.Errorf("magnitude %d differs from expected %d", mag, expectedMag)
	}
	return nil
}

func parseMatrix(out string, r, c int) ([][]int, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	mat := make([][]int, r)
	for i := 0; i < r; i++ {
		mat[i] = make([]int, c)
		for j := 0; j < c; j++ {
			if _, err := fmt.Fscan(reader, &mat[i][j]); err != nil {
				return nil, fmt.Errorf("failed to read a[%d][%d]: %v", i+1, j+1, err)
			}
			if mat[i][j] < 1 || mat[i][j] > 1_000_000_000 {
				return nil, fmt.Errorf("a[%d][%d]=%d is out of bounds", i+1, j+1, mat[i][j])
			}
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err == nil {
		return nil, fmt.Errorf("unexpected extra token %q after matrix", extra)
	}
	return mat, nil
}

func computeStats(mat [][]int) ([]int, int, error) {
	if len(mat) == 0 || len(mat[0]) == 0 {
		return nil, 0, fmt.Errorf("empty matrix")
	}
	r := len(mat)
	c := len(mat[0])
	gcds := make([]int, 0, r+c)
	for i := 0; i < r; i++ {
		g := 0
		for j := 0; j < c; j++ {
			g = gcd(g, mat[i][j])
		}
		if g <= 0 {
			return nil, 0, fmt.Errorf("row %d has non-positive gcd %d", i+1, g)
		}
		gcds = append(gcds, g)
	}
	for j := 0; j < c; j++ {
		g := 0
		for i := 0; i < r; i++ {
			g = gcd(g, mat[i][j])
		}
		if g <= 0 {
			return nil, 0, fmt.Errorf("column %d has non-positive gcd %d", j+1, g)
		}
		gcds = append(gcds, g)
	}
	mag := 0
	for _, v := range gcds {
		if v > mag {
			mag = v
		}
	}
	return gcds, mag, nil
}

func ensureDistinct(vals []int) error {
	seen := make(map[int]struct{}, len(vals))
	for _, v := range vals {
		if _, ok := seen[v]; ok {
			return fmt.Errorf("GCD value %d appears multiple times", v)
		}
		seen[v] = struct{}{}
	}
	return nil
}

func gcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func buildTests() []testCase {
	tests := []testCase{
		{1, 1},
		{1, 2},
		{2, 1},
		{2, 2},
		{2, 3},
		{3, 2},
		{3, 3},
		{4, 5},
		{5, 4},
		{7, 1},
		{1, 7},
		{8, 8},
		{10, 15},
		{15, 10},
		{20, 20},
		{50, 3},
		{3, 50},
		{100, 1},
		{1, 100},
		{30, 40},
		{40, 30},
		{60, 60},
		{100, 200},
		{200, 100},
		{200, 200},
		{500, 1},
		{1, 500},
		{250, 500},
		{500, 250},
	}
	rng := rand.New(rand.NewSource(1266))
	for i := 0; i < 40; i++ {
		r := rng.Intn(60) + 1
		c := rng.Intn(60) + 1
		tests = append(tests, testCase{r: r, c: c})
	}
	return tests
}
