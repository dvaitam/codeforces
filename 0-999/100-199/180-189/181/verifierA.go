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

const refSource = "./181A.go"

type testCase struct {
	n, m                   int
	missingRow, missingCol int
	input                  string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
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
		refOut, err := runProgram(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc.input, refOut)
			os.Exit(1)
		}
		if err := checkOutput(tc, refOut); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid output on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%soutput:\n%s", idx+1, err, tc.input, candOut)
			os.Exit(1)
		}
		if err := checkOutput(tc, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d failed: %v\ninput:\n%scandidate output:\n%s", idx+1, err, tc.input, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-181A-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref181A.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() {
		_ = os.RemoveAll(dir)
	}
	return binPath, cleanup, nil
}

func buildTests() []testCase {
	tests := make([]testCase, 0, 220)
	tests = append(tests,
		newTestCase(2, 2, 1, 2, 1, 2, 0),
		newTestCase(2, 100, 1, 2, 1, 100, 3),
		newTestCase(100, 2, 1, 100, 1, 2, 1),
		newTestCase(100, 100, 5, 80, 10, 90, 2),
	)
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 200 {
		tests = append(tests, generateCase(rng))
	}
	return tests
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(99) + 2
	m := rng.Intn(99) + 2
	r1 := rng.Intn(n) + 1
	r2 := rng.Intn(n) + 1
	for r2 == r1 {
		r2 = rng.Intn(n) + 1
	}
	c1 := rng.Intn(m) + 1
	c2 := rng.Intn(m) + 1
	for c2 == c1 {
		c2 = rng.Intn(m) + 1
	}
	missing := rng.Intn(4)
	return newTestCase(n, m, r1, r2, c1, c2, missing)
}

func newTestCase(n, m, r1, r2, c1, c2, missing int) testCase {
	coords := [4][2]int{
		{r1, c1},
		{r1, c2},
		{r2, c1},
		{r2, c2},
	}
	grid := make([][]byte, n)
	for i := range grid {
		row := make([]byte, m)
		for j := range row {
			row[j] = '.'
		}
		grid[i] = row
	}
	for idx, pos := range coords {
		if idx == missing {
			continue
		}
		grid[pos[0]-1][pos[1]-1] = '*'
	}
	var sb strings.Builder
	sb.Grow(n*(m+1) + 20)
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for _, row := range grid {
		sb.Write(row)
		sb.WriteByte('\n')
	}
	missingRow := coords[missing][0]
	missingCol := coords[missing][1]
	return testCase{
		n:          n,
		m:          m,
		missingRow: missingRow,
		missingCol: missingCol,
		input:      sb.String(),
	}
}

func runProgram(bin string, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func checkOutput(tc testCase, output string) error {
	fields := strings.Fields(output)
	if len(fields) != 2 {
		return fmt.Errorf("expected 2 integers, got %d tokens", len(fields))
	}
	row, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid row value %q", fields[0])
	}
	col, err := strconv.Atoi(fields[1])
	if err != nil {
		return fmt.Errorf("invalid column value %q", fields[1])
	}
	if row < 1 || row > tc.n {
		return fmt.Errorf("row %d out of range [1,%d]", row, tc.n)
	}
	if col < 1 || col > tc.m {
		return fmt.Errorf("column %d out of range [1,%d]", col, tc.m)
	}
	if row != tc.missingRow || col != tc.missingCol {
		return fmt.Errorf("expected (%d %d) got (%d %d)", tc.missingRow, tc.missingCol, row, col)
	}
	return nil
}
