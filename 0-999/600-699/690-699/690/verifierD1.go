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

const refSource = "0-999/600-699/690-699/690/690D1.go"

type testCase struct {
	name  string
	input string
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD1.go /path/to/binary")
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
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}
		exp, err := parseOutput(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, refOut)
			os.Exit(1)
		}

		candOut, err := runProgram(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		got, err := parseOutput(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate output invalid on test %d (%s): %v\ninput:\n%soutput:\n%s",
				idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
		if got != exp {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: expected %d got %d\ninput:\n%sreference output:\n%s\ncandidate output:\n%s",
				idx+1, tc.name, exp, got, tc.input, refOut, candOut)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, func(), error) {
	dir, err := os.MkdirTemp("", "cf-690D1-ref-")
	if err != nil {
		return "", nil, fmt.Errorf("failed to create temp dir: %v", err)
	}
	binPath := filepath.Join(dir, "ref690D1.bin")
	cmd := exec.Command("go", "build", "-o", binPath, refSource)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.RemoveAll(dir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, stderr.String())
	}
	cleanup := func() { _ = os.RemoveAll(dir) }
	return binPath, cleanup, nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) != 1 {
		return 0, fmt.Errorf("expected single integer, got %d tokens", len(fields))
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	if val < 0 {
		return 0, fmt.Errorf("negative segment count %d", val)
	}
	return val, nil
}

func buildTests() []testCase {
	tests := []testCase{
		formatManual("single_column_full", []string{"B"}),
		formatManual("single_segment", []string{
			".B.",
			".B.",
			"BBB",
		}),
		formatManual("multiple_segments", []string{
			"B.B",
			"B.B",
			"B.B",
		}),
		formatManual("no_bricks", []string{
			"...",
			"...",
			"...",
		}),
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 200; i++ {
		tests = append(tests, randomCase(rng, i))
	}
	tests = append(tests, stressCase())
	return tests
}

func formatManual(name string, grid []string) testCase {
	R := len(grid)
	C := len(grid[0])
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", R, C)
	for i := 0; i < R; i++ {
		sb.WriteString(grid[i])
		sb.WriteByte('\n')
	}
	return testCase{name: name, input: sb.String()}
}

func randomCase(rng *rand.Rand, idx int) testCase {
	R := rng.Intn(100) + 1
	C := rng.Intn(100) + 1
	grid := make([][]byte, R)
	for i := range grid {
		grid[i] = make([]byte, C)
		for j := 0; j < C; j++ {
			if rng.Intn(3) == 0 {
				grid[i][j] = 'B'
			} else {
				grid[i][j] = '.'
			}
		}
	}
	// enforce property: bricks stacked from bottom
	for c := 0; c < C; c++ {
		height := rng.Intn(R + 1)
		for r := 0; r < R; r++ {
			if r < height {
				grid[R-1-r][c] = 'B'
			} else {
				grid[R-1-r][c] = '.'
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", R, C)
	for i := 0; i < R; i++ {
		sb.Write(grid[i])
		sb.WriteByte('\n')
	}
	return testCase{
		name:  fmt.Sprintf("random_%d", idx+1),
		input: sb.String(),
	}
}

func stressCase() testCase {
	R, C := 100, 100
	grid := make([][]byte, R)
	for i := range grid {
		grid[i] = make([]byte, C)
		for j := 0; j < C; j++ {
			if j%3 == 0 {
				grid[i][j] = 'B'
			} else {
				grid[i][j] = '.'
			}
		}
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", R, C)
	for i := 0; i < R; i++ {
		sb.Write(grid[i])
		sb.WriteByte('\n')
	}
	return testCase{name: "stress_pattern", input: sb.String()}
}
