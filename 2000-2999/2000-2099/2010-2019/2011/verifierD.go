package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

const referenceSolutionRel = "2000-2999/2000-2099/2010-2019/2011/2011D.go"

var referenceSolutionPath string

func init() {
	referenceSolutionPath = referenceSolutionRel
	if _, file, _, ok := runtime.Caller(0); ok {
		dir := filepath.Dir(file)
		candidate := filepath.Join(dir, "2011D.go")
		if _, err := os.Stat(candidate); err == nil {
			referenceSolutionPath = candidate
			return
		}
	}
	if abs, err := filepath.Abs(referenceSolutionRel); err == nil {
		if _, err := os.Stat(abs); err == nil {
			referenceSolutionPath = abs
		}
	}
}

type testCase struct {
	n int
	h int64
	b int64
	g [2]string
}

func deterministicTests() []testCase {
	return []testCase{
		{
			n: 2, h: 3, b: 7,
			g: [2]string{"S.", ".."},
		},
		{
			n: 2, h: 3, b: 7,
			g: [2]string{"S.", ".W"},
		},
		{
			n: 2, h: 7, b: 3,
			g: [2]string{"S.", ".W"},
		},
		{
			n: 4, h: 999999999, b: 1000000000,
			g: [2]string{"W.S.", "W..W"},
		},
		{
			n: 6, h: 6, b: 7,
			g: [2]string{"W....W", "W.S..W"},
		},
		{
			n: 6, h: 6, b: 7,
			g: [2]string{"W...WW", "..S..W"},
		},
	}
}

func randomTests() []testCase {
	rng := rand.New(rand.NewSource(20250305))
	var tests []testCase
	totalN := 0
	for totalN < 200000 && len(tests) < 80 {
		n := rng.Intn(40) + 2
		if totalN+n > 200000 {
			n = 200000 - totalN
		}
		if n < 2 {
			break
		}
		h := int64(rng.Intn(1_000_000_000) + 1)
		b := int64(rng.Intn(1_000_000_000) + 1)
		if h == b {
			h++
		}
		grid := [2][]rune{make([]rune, n), make([]rune, n)}
		sRow := rng.Intn(2)
		sCol := rng.Intn(n)
		for r := 0; r < 2; r++ {
			for c := 0; c < n; c++ {
				grid[r][c] = '.'
			}
		}
		grid[sRow][sCol] = 'S'
		for r := 0; r < 2; r++ {
			for c := 0; c < n; c++ {
				if r == sRow && c == sCol {
					continue
				}
				if abs(r-sRow)+abs(c-sCol) == 1 {
					continue
				}
				if rng.Intn(3) == 0 {
					grid[r][c] = 'W'
				}
			}
		}
		tests = append(tests, testCase{
			n: n,
			h: h,
			b: b,
			g: [2]string{string(grid[0]), string(grid[1])},
		})
		totalN += n
	}
	return tests
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func formatInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d %d\n", tc.n, tc.h, tc.b))
		sb.WriteString(tc.g[0])
		sb.WriteByte('\n')
		sb.WriteString(tc.g[1])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func runProgram(path, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(path, ".go") {
		cmd = exec.Command("go", "run", path)
	} else {
		cmd = exec.Command(path)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func buildReferenceBinary() (string, func(), error) {
	if referenceSolutionPath == "" {
		return "", nil, fmt.Errorf("reference solution path not set")
	}
	if _, err := os.Stat(referenceSolutionPath); err != nil {
		return "", nil, fmt.Errorf("reference solution not found: %v", err)
	}
	tmpDir, err := os.MkdirTemp("", "2011D-ref")
	if err != nil {
		return "", nil, err
	}
	binPath := filepath.Join(tmpDir, "ref_2011D")
	cmd := exec.Command("go", "build", "-o", binPath, referenceSolutionPath)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.RemoveAll(tmpDir)
		return "", nil, fmt.Errorf("failed to build reference: %v\n%s", err, out.String())
	}
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	return binPath, cleanup, nil
}

func parseOutputs(out string, count int) ([]int64, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	scanner.Split(bufio.ScanWords)
	results := make([]int64, 0, count)
	for scanner.Scan() {
		var val int64
		if _, err := fmt.Sscan(scanner.Text(), &val); err != nil {
			return nil, fmt.Errorf("failed to parse %q: %v", scanner.Text(), err)
		}
		results = append(results, val)
	}
	if len(results) != count {
		return nil, fmt.Errorf("expected %d numbers, got %d", count, len(results))
	}
	return results, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests := append(deterministicTests(), randomTests()...)
	input := formatInput(tests)

	refBin, cleanup, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer cleanup()

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference runtime error: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}
	expected, err := parseOutputs(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\noutput:\n%s\n", err, refOut)
		os.Exit(1)
	}

	userOut, err := runProgram(bin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "runtime error: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}
	got, err := parseOutputs(userOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse participant output: %v\noutput:\n%s\n", err, userOut)
		os.Exit(1)
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d, got %d\nn=%d h=%d b=%d\ngrid:\n%s\n%s\n", i+1, expected[i], got[i], tests[i].n, tests[i].h, tests[i].b, tests[i].g[0], tests[i].g[1])
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
