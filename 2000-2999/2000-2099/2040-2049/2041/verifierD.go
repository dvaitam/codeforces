package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	refSource2041D = "2041D.go"
	refBinary2041D = "ref2041D.bin"
	totalTests     = 80
	maxCells       = 200000
)

type testCase struct {
	n, m int
	grid []string
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		return
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Println("failed to build reference:", err)
		return
	}
	defer os.Remove(ref)

	tests := generateTests()
	input := formatInput(tests)

	refOut, err := runProgram(ref, input)
	if err != nil {
		fmt.Printf("reference runtime error: %v\n", err)
		return
	}
	expected, err := parseOutput(refOut, len(tests))
	if err != nil {
		fmt.Printf("reference output parse error: %v\noutput:\n%s", err, refOut)
		return
	}

	candOut, err := runProgram(candidate, input)
	if err != nil {
		fmt.Printf("candidate runtime error: %v\n", err)
		return
	}
	got, err := parseOutput(candOut, len(tests))
	if err != nil {
		fmt.Printf("candidate output parse error: %v\noutput:\n%s", err, candOut)
		return
	}

	for i := range tests {
		if expected[i] != got[i] {
			fmt.Printf("Mismatch on case %d: expected %d, got %d\n", i+1, expected[i], got[i])
			fmt.Println("Input used:")
			fmt.Println(string(input))
			return
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}

func buildReference() (string, error) {
	cmd := exec.Command("go", "build", "-o", refBinary2041D, refSource2041D)
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("%v\n%s", err, string(out))
	}
	return filepath.Join(".", refBinary2041D), nil
}

func runProgram(path string, input []byte) (string, error) {
	cmd := exec.Command(path)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	return out.String(), err
}

func parseOutput(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d", t, len(fields))
	}
	res := make([]int, t)
	for i, tok := range fields {
		var val int
		if _, err := fmt.Sscanf(tok, "%d", &val); err != nil {
			return nil, fmt.Errorf("token %d: %v", i+1, err)
		}
		res[i] = val
	}
	return res, nil
}

func formatInput(tests []testCase) []byte {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d\n", len(tests))
	for _, tc := range tests {
		fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
		for _, row := range tc.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
	}
	return []byte(sb.String())
}

func generateTests() []testCase {
	rnd := rand.New(rand.NewSource(2041))
	var tests []testCase
	totalCells := 0

	// Small edge cases
	tests = append(tests, buildSimpleMaze(3, 3, []string{
		"###",
		"#S#",
		"#T#",
	}))
	totalCells += 9

	tests = append(tests, buildSimpleMaze(4, 4, []string{
		"####",
		"#S.#",
		"#..#",
		"#.T#",
	}))
	totalCells += 16

	for len(tests) < totalTests && totalCells < maxCells {
		capacity := maxCells - totalCells
		n := rnd.Intn(200) + 3
		m := rnd.Intn(200) + 3
		if n*m > capacity {
			n = max(3, capacity/m)
		}
		tc := randomMaze(n, m, rnd)
		tests = append(tests, tc)
		totalCells += n * m
	}
	return tests
}

func randomMaze(n, m int, rnd *rand.Rand) testCase {
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			grid[i][j] = '#'
		}
	}

	for i := 1; i < n-1; i++ {
		for j := 1; j < m-1; j++ {
			if rnd.Intn(100) < 70 {
				grid[i][j] = '.'
			} else {
				grid[i][j] = '#'
			}
		}
	}

	sx, sy := rnd.Intn(n-2)+1, rnd.Intn(m-2)+1
	tx, ty := rnd.Intn(n-2)+1, rnd.Intn(m-2)+1
	grid[sx][sy] = 'S'
	grid[tx][ty] = 'T'

	rows := make([]string, n)
	for i := 0; i < n; i++ {
		rows[i] = string(grid[i])
	}
	return testCase{n: n, m: m, grid: rows}
}

func buildSimpleMaze(n, m int, lines []string) testCase {
	return testCase{n: n, m: m, grid: append([]string(nil), lines...)}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
