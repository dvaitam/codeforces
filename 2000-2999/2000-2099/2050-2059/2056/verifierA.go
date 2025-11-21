package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

const refSource = "2056A.go"

type testCase struct {
	n     int
	m     int
	moves [][2]int
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := args[0]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := buildTests()
	input := buildInput(tests)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference failed: %v\n", err)
		os.Exit(1)
	}
	refAns, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse reference output: %v\n%s", err, refOut)
		os.Exit(1)
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate runtime error: %v\n", err)
		os.Exit(1)
	}
	candAns, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse candidate output: %v\n%s", err, candOut)
		os.Exit(1)
	}

	for i := range tests {
		if refAns[i] != candAns[i] {
			tc := tests[i]
			fmt.Fprintf(os.Stderr, "test %d mismatch: expected %d got %d\n", i+1, refAns[i], candAns[i])
			fmt.Fprintf(os.Stderr, "n=%d m=%d moves=%s\n", tc.n, tc.m, formatMoves(tc.moves))
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}

func buildReference() (string, error) {
	refPath, err := referencePath()
	if err != nil {
		return "", err
	}
	tmp, err := os.CreateTemp("", "ref_2056A_*.bin")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %v", err)
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), refPath)
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, string(out))
	}
	return tmp.Name(), nil
}

func referencePath() (string, error) {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "", fmt.Errorf("failed to locate verifier file path")
	}
	dir := filepath.Dir(file)
	return filepath.Join(dir, refSource), nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\nstdout:\n%s\nstderr:\n%s", err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String()), nil
}

func buildTests() []testCase {
	var tests []testCase
	add := func(n, m int, moves [][2]int) {
		tests = append(tests, testCase{n: n, m: m, moves: moves})
	}

	add(1, 2, [][2]int{{1, 1}})
	add(2, 3, [][2]int{{1, 1}, {1, 1}})
	add(4, 3, [][2]int{{1, 1}, {2, 2}, {2, 1}, {1, 2}})
	add(3, 4, [][2]int{{1, 3}, {3, 1}, {2, 2}})
	add(5, 5, [][2]int{{4, 4}, {1, 1}, {4, 4}, {1, 1}, {2, 2}})
	add(10, 100, uniformMoves(10, 100, 50, 75))
	add(100, 100, uniformMoves(100, 100, 99, 98))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for len(tests) < 250 {
		m := rng.Intn(99) + 2  // 2..100
		n := rng.Intn(100) + 1 // 1..100
		moves := make([][2]int, n)
		for i := 0; i < n; i++ {
			dx := rng.Intn(m-1) + 1
			dy := rng.Intn(m-1) + 1
			moves[i] = [2]int{dx, dy}
		}
		tests = append(tests, testCase{n: n, m: m, moves: moves})
	}
	return tests
}

func uniformMoves(n, m, dx, dy int) [][2]int {
	moves := make([][2]int, n)
	for i := 0; i < n; i++ {
		moves[i] = [2]int{dx, dy}
	}
	return moves
}

func buildInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		if len(tc.moves) != tc.n {
			panic("move count does not match n")
		}
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
		for _, mv := range tc.moves {
			sb.WriteString(fmt.Sprintf("%d %d\n", mv[0], mv[1]))
		}
	}
	return sb.String()
}

func parseOutput(out string, expected int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != expected {
		return nil, fmt.Errorf("expected %d outputs, got %d", expected, len(fields))
	}
	ans := make([]int, expected)
	for i, s := range fields {
		val, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("invalid integer %q", s)
		}
		ans[i] = val
	}
	return ans, nil
}

func formatMoves(moves [][2]int) string {
	parts := make([]string, len(moves))
	for i, mv := range moves {
		parts[i] = fmt.Sprintf("(%d,%d)", mv[0], mv[1])
	}
	return strings.Join(parts, " ")
}
