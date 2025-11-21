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
	"runtime"
	"strings"
	"time"
)

type testCase struct {
	name  string
	input string
}

type caseData struct {
	n int
	m int
	g [][]byte
}

type op struct {
	cells [3][2]int
}

var problemDir string

func init() {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		panic("unable to determine verifier location")
	}
	problemDir = filepath.Dir(file)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA1.go /path/to/binary")
		os.Exit(1)
	}
	target := os.Args[1]

	refBin, err := buildReferenceBinary()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	for idx, tc := range tests {
		cases, err := parseInput(tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse generated test %q: %v\n", tc.name, err)
			os.Exit(1)
		}

		if _, err := runProgram(refBin, tc.input); err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}

		candOut, err := runProgram(target, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "solution runtime error on test %d (%s): %v\ninput:\n%s", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		if err := verifyOutput(cases, candOut); err != nil {
			fmt.Fprintf(os.Stderr, "test %d (%s) failed: %v\ninput:\n%soutput:\n%s", idx+1, tc.name, err, tc.input, candOut)
			os.Exit(1)
		}
	}

	fmt.Printf("Accepted (%d tests)\n", len(tests))
}

func buildReferenceBinary() (string, error) {
	tmp, err := os.CreateTemp("", "cf-1439A1-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()
	os.Remove(tmp.Name())

	cmd := exec.Command("go", "build", "-o", tmp.Name(), "1439A1.go")
	cmd.Dir = problemDir
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("go build error: %v\n%s", err, stderr.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func parseInput(input string) ([]caseData, error) {
	reader := bufio.NewReader(strings.NewReader(input))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, fmt.Errorf("read t: %w", err)
	}
	cases := make([]caseData, t)
	totalCells := 0
	for i := 0; i < t; i++ {
		var n, m int
		if _, err := fmt.Fscan(reader, &n, &m); err != nil {
			return nil, fmt.Errorf("case %d: read n,m: %w", i+1, err)
		}
		grid := make([][]byte, n)
		for r := 0; r < n; r++ {
			var line string
			if _, err := fmt.Fscan(reader, &line); err != nil {
				return nil, fmt.Errorf("case %d: read row %d: %w", i+1, r+1, err)
			}
			if len(line) != m {
				return nil, fmt.Errorf("case %d: row %d length mismatch", i+1, r+1)
			}
			grid[r] = []byte(line)
		}
		cases[i] = caseData{n: n, m: m, g: grid}
		totalCells += n * m
	}
	if totalCells > 20000 {
		return nil, fmt.Errorf("total cells exceeded 20000")
	}
	return cases, nil
}

func verifyOutput(cases []caseData, output string) error {
	reader := bufio.NewReader(strings.NewReader(output))
	for idx, tc := range cases {
		var k int
		if _, err := fmt.Fscan(reader, &k); err != nil {
			return fmt.Errorf("case %d: failed to read k: %v", idx+1, err)
		}
		if k < 0 || k > 3*tc.n*tc.m {
			return fmt.Errorf("case %d: k=%d exceeds 3nm=%d", idx+1, k, 3*tc.n*tc.m)
		}
		grid := cloneGrid(tc.g)
		for opIdx := 0; opIdx < k; opIdx++ {
			op, err := readOperation(reader, tc.n, tc.m)
			if err != nil {
				return fmt.Errorf("case %d: operation %d invalid: %v", idx+1, opIdx+1, err)
			}
			if err := applyOperation(grid, op); err != nil {
				return fmt.Errorf("case %d: operation %d invalid application: %v", idx+1, opIdx+1, err)
			}
		}
		if !allZero(grid) {
			return fmt.Errorf("case %d: grid not zeroed", idx+1)
		}
	}
	var extra string
	if _, err := fmt.Fscan(reader, &extra); err != nil {
		if err != io.EOF {
			return fmt.Errorf("failed to parse trailing output: %v", err)
		}
	} else {
		return fmt.Errorf("unexpected extra token %q after processing all test cases", extra)
	}
	return nil
}

func readOperation(reader *bufio.Reader, n, m int) (op, error) {
	var cells [3][2]int
	seen := make(map[[2]int]struct{})
	for i := 0; i < 3; i++ {
		var x, y int
		if _, err := fmt.Fscan(reader, &x, &y); err != nil {
			return op{}, fmt.Errorf("read cell %d: %v", i+1, err)
		}
		if x < 1 || x > n || y < 1 || y > m {
			return op{}, fmt.Errorf("cell (%d,%d) out of range", x, y)
		}
		coord := [2]int{x, y}
		if _, ok := seen[coord]; ok {
			return op{}, fmt.Errorf("duplicate cell (%d,%d)", x, y)
		}
		seen[coord] = struct{}{}
		cells[i] = coord
	}
	if !sameSquare(cells) {
		return op{}, fmt.Errorf("cells do not belong to same 2x2 square")
	}
	return op{cells: cells}, nil
}

func sameSquare(cells [3][2]int) bool {
	coords := make([][2]int, 0, 3)
	for _, c := range cells {
		coords = append(coords, c)
	}
	minX, maxX := coords[0][0], coords[0][0]
	minY, maxY := coords[0][1], coords[0][1]
	for _, c := range coords[1:] {
		if c[0] < minX {
			minX = c[0]
		}
		if c[0] > maxX {
			maxX = c[0]
		}
		if c[1] < minY {
			minY = c[1]
		}
		if c[1] > maxY {
			maxY = c[1]
		}
	}
	return maxX-minX <= 1 && maxY-minY <= 1
}

func applyOperation(grid [][]byte, operation op) error {
	for _, c := range operation.cells {
		x := c[0] - 1
		y := c[1] - 1
		if grid[x][y] == '0' {
			grid[x][y] = '1'
		} else {
			grid[x][y] = '0'
		}
	}
	return nil
}

func cloneGrid(g [][]byte) [][]byte {
	n := len(g)
	clone := make([][]byte, n)
	for i := range g {
		clone[i] = append([]byte(nil), g[i]...)
	}
	return clone
}

func allZero(grid [][]byte) bool {
	for _, row := range grid {
		for _, v := range row {
			if v != '0' {
				return false
			}
		}
	}
	return true
}

func generateTests() []testCase {
	var tests []testCase

	tests = append(tests, buildManualTest("simple", []caseData{
		{n: 2, m: 2, g: [][]byte{
			[]byte("10"),
			[]byte("11"),
		}},
		{n: 3, m: 3, g: [][]byte{
			[]byte("010"),
			[]byte("111"),
			[]byte("000"),
		}},
	}))

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	tests = append(tests, randomTest("random_small", 10, 5, 5, rng))
	tests = append(tests, randomTest("random_medium", 10, 20, 20, rng))
	tests = append(tests, randomTest("random_large", 5, 100, 100, rng))

	return tests
}

func buildManualTest(name string, cases []caseData) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", len(cases))
	for _, tc := range cases {
		fmt.Fprintf(&b, "%d %d\n", tc.n, tc.m)
		for _, row := range tc.g {
			fmt.Fprintf(&b, "%s\n", string(row))
		}
	}
	return testCase{name: name, input: b.String()}
}

func randomTest(name string, caseCount, maxN, maxM int, rng *rand.Rand) testCase {
	var b strings.Builder
	fmt.Fprintf(&b, "%d\n", caseCount)
	totalCells := 0
	for i := 0; i < caseCount; i++ {
		n := rng.Intn(maxN-1) + 2
		m := rng.Intn(maxM-1) + 2
		if totalCells+n*m > 20000 {
			remaining := 20000 - totalCells
			if remaining < 4 {
				n, m = 2, 2
			} else {
				m = min(m, remaining/n)
				if m < 2 {
					m = 2
				}
			}
		}
		totalCells += n * m
		fmt.Fprintf(&b, "%d %d\n", n, m)
		for r := 0; r < n; r++ {
			row := make([]byte, m)
			for c := 0; c < m; c++ {
				if rng.Intn(2) == 0 {
					row[c] = '0'
				} else {
					row[c] = '1'
				}
			}
			fmt.Fprintf(&b, "%s\n", string(row))
		}
	}
	return testCase{name: name, input: b.String()}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
