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

const (
	refSourceD   = "1000-1999/1700-1799/1770-1779/1773/1773D.go"
	answerCap    = 1000000
	maxRandomLen = 80
)

type testCase struct {
	n    int
	m    int
	grid []string
}

func (tc testCase) input() string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
	for _, row := range tc.grid {
		b.WriteString(row)
		b.WriteByte('\n')
	}
	return b.String()
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		input := tc.input()
		refOut, err := runProgram(refBin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference failed on test %d: %v\ninput:\n%s\n", idx+1, err, input)
			os.Exit(1)
		}
		refVal, err := parseAnswer(refOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d: %v\noutput:\n%s\n", idx+1, err, refOut)
			os.Exit(1)
		}

		candOut, err := runCandidate(candidate, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d: %v\ninput:\n%s\nstdout/stderr:\n%s\n", idx+1, err, input, candOut)
			os.Exit(1)
		}
		candVal, err := parseAnswer(candOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "invalid candidate output on test %d: %v\noutput:\n%s\n", idx+1, err, candOut)
			os.Exit(1)
		}

		if candVal != refVal {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d\ninput:\n%sreference: %d\ncandidate: %d\n", idx+1, input, refVal, candVal)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "1773D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSourceD))
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
	return out.String(), err
}

func parseAnswer(out string) (int, error) {
	fields := strings.Fields(out)
	if len(fields) == 0 {
		return 0, fmt.Errorf("empty output")
	}
	val, err := strconv.Atoi(fields[0])
	if err != nil {
		return 0, fmt.Errorf("invalid integer %q", fields[0])
	}
	if val < 0 || val > answerCap {
		return 0, fmt.Errorf("answer %d out of range", val)
	}
	return val, nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests, manualTest([]string{
		"...#..",
		"......",
		"#...##",
	}))
	tests = append(tests, manualTest([]string{
		"..",
		"..",
	}))
	tests = append(tests, manualTest([]string{
		"#.",
		"#.",
	}))

	rng := rand.New(rand.NewSource(1773))
	tests = append(tests, tileableBoard(1, 8, 4, rng))   // thin row
	tests = append(tests, tileableBoard(2, 50, 30, rng)) // long corridor
	tests = append(tests, tileableBoard(6, 6, 18, rng))  // dense square
	tests = append(tests, tileableBoard(20, 20, 150, rng))
	tests = append(tests, tileableBoard(40, 52, 40*52/2, rng)) // forces cap

	for len(tests) < maxRandomLen {
		n := rng.Intn(20) + 1
		m := rng.Intn(20) + 1
		if n*m%2 == 1 {
			if rng.Intn(2) == 0 {
				if n < 1000 {
					n++
				} else {
					m++
				}
			} else {
				if m < 1000 {
					m++
				} else {
					n++
				}
			}
		}
		if n*m%2 == 1 {
			m++
		}
		if n > 1000 {
			n = 1000
		}
		if m > 1000 {
			m = 1000
		}
		totalDominoes := (n * m) / 2
		keep := rng.Intn(totalDominoes) + 1
		tests = append(tests, tileableBoard(n, m, keep, rng))
	}
	return tests
}

func manualTest(grid []string) testCase {
	n := len(grid)
	m := len(grid[0])
	return testCase{
		n:    n,
		m:    m,
		grid: append([]string(nil), grid...),
	}
}

type domino struct {
	a cell
	b cell
}

type cell struct {
	r int
	c int
}

func tileableBoard(n, m, keep int, rng *rand.Rand) testCase {
	if n*m%2 == 1 {
		panic("tileableBoard requires even cell count")
	}
	board := make([][]byte, n)
	for i := range board {
		row := make([]byte, m)
		for j := range row {
			row[j] = '.'
		}
		board[i] = row
	}

	totalDominoes := (n * m) / 2
	if keep < 1 {
		keep = 1
	}
	if keep > totalDominoes {
		keep = totalDominoes
	}
	var dominos []domino
	if m%2 == 0 {
		for i := 0; i < n; i++ {
			for j := 0; j < m; j += 2 {
				dominos = append(dominos, domino{cell{i, j}, cell{i, j + 1}})
			}
		}
	} else {
		for j := 0; j < m; j++ {
			for i := 0; i < n; i += 2 {
				dominos = append(dominos, domino{cell{i, j}, cell{i + 1, j}})
			}
		}
	}

	remove := totalDominoes - keep
	if remove > 0 {
		order := rng.Perm(len(dominos))
		for i := 0; i < remove; i++ {
			d := dominos[order[i]]
			board[d.a.r][d.a.c] = '#'
			board[d.b.r][d.b.c] = '#'
		}
	}

	grid := make([]string, n)
	for i := range board {
		grid[i] = string(board[i])
	}
	return testCase{n: n, m: m, grid: grid}
}
