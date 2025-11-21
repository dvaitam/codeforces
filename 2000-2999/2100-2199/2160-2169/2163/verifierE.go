package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

type testCase struct {
	n    int
	grid []string
	conn int
}

type pair struct {
	r int
	c int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	tests := buildTests()
	firstInput := buildFirstInput(tests)

	refPath, err := buildReference()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to build reference: %v\n", err)
		os.Exit(1)
	}
	defer os.Remove(refPath)

	refPairs, err := executeFirst(refPath, firstInput, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference first run failed: %v\n", err)
		os.Exit(1)
	}
	expected, _, err := executeSecond(refPath, tests, refPairs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "reference second run failed: %v\n", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		if expected[i] != tc.conn {
			fmt.Fprintf(os.Stderr, "reference output mismatch on test %d: expected connectivity %d, reference produced %d\n", i+1, tc.conn, expected[i])
			os.Exit(1)
		}
	}

	candPairs, err := executeFirst(candidate, firstInput, tests)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate first run failed: %v\n", err)
		os.Exit(1)
	}
	candOutputs, candSecondInput, err := executeSecond(candidate, tests, candPairs)
	if err != nil {
		fmt.Fprintf(os.Stderr, "candidate second run failed: %v\nSecond run input:\n%s\n", err, candSecondInput)
		os.Exit(1)
	}

	for i := range tests {
		if candOutputs[i] != expected[i] {
			rowStr, colStr := fetchRowCol(tests[i], candPairs[i])
			fmt.Fprintf(os.Stderr, "Mismatch on test %d: expected %d, got %d\n", i+1, expected[i], candOutputs[i])
			fmt.Fprintf(os.Stderr, "Candidate chose row %d column %d\nRow: %s\nColumn: %s\n", candPairs[i].r, candPairs[i].c, rowStr, colStr)
			fmt.Fprintf(os.Stderr, "Second run input provided to candidate:\n%s\n", candSecondInput)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "2163E_ref_*")
	if err != nil {
		return "", err
	}
	path := tmp.Name()
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", path, "2163E.go")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(path)
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return path, nil
}

func buildFirstInput(tests []testCase) string {
	var sb strings.Builder
	sb.WriteString("first\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for _, tc := range tests {
		sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.conn))
		for _, row := range tc.grid {
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
	}
	return sb.String()
}

func executeFirst(target, firstInput string, tests []testCase) ([]pair, error) {
	out, err := runProgram(target, firstInput)
	if err != nil {
		return nil, fmt.Errorf("runtime error: %v\n%s", err, out)
	}
	pairs, err := parsePairs(out, tests)
	if err != nil {
		return nil, err
	}
	for i, tc := range tests {
		if pairs[i].r < 1 || pairs[i].r > tc.n || pairs[i].c < 1 || pairs[i].c > tc.n {
			return nil, fmt.Errorf("test %d: invalid row/column (%d,%d) for n=%d", i+1, pairs[i].r, pairs[i].c, tc.n)
		}
	}
	return pairs, nil
}

func executeSecond(target string, tests []testCase, pairs []pair) ([]int, string, error) {
	secondInput, err := buildSecondInput(tests, pairs)
	if err != nil {
		return nil, "", err
	}
	out, runErr := runProgram(target, secondInput)
	if runErr != nil {
		return nil, secondInput, fmt.Errorf("runtime error: %v\n%s", runErr, out)
	}
	results, err := parseConnectivityOutputs(out, len(tests))
	if err != nil {
		return nil, secondInput, err
	}
	return results, secondInput, nil
}

func buildSecondInput(tests []testCase, pairs []pair) (string, error) {
	var sb strings.Builder
	sb.WriteString("second\n")
	sb.WriteString(fmt.Sprintf("%d\n", len(tests)))
	for i, tc := range tests {
		r := pairs[i].r
		c := pairs[i].c
		if r < 1 || r > tc.n || c < 1 || c > tc.n {
			return "", fmt.Errorf("test %d: invalid row/column (%d,%d) for n=%d", i+1, r, c, tc.n)
		}
		row := tc.grid[r-1]
		colBytes := make([]byte, tc.n)
		for j := 0; j < tc.n; j++ {
			colBytes[j] = tc.grid[j][c-1]
		}
		sb.WriteString(fmt.Sprintf("%d\n", tc.n))
		sb.WriteString(row)
		sb.WriteByte('\n')
		sb.WriteString(string(colBytes))
		sb.WriteByte('\n')
	}
	return sb.String(), nil
}

func parsePairs(out string, tests []testCase) ([]pair, error) {
	reader := strings.NewReader(out)
	pairs := make([]pair, len(tests))
	for i := range tests {
		if _, err := fmt.Fscan(reader, &pairs[i].r, &pairs[i].c); err != nil {
			return nil, fmt.Errorf("failed to read row/col for test %d: %v\noutput:\n%s", i+1, err, out)
		}
	}
	return pairs, nil
}

func parseConnectivityOutputs(out string, t int) ([]int, error) {
	fields := strings.Fields(out)
	if len(fields) != t {
		return nil, fmt.Errorf("expected %d outputs, got %d\noutput:\n%s", t, len(fields), out)
	}
	res := make([]int, t)
	for i, tok := range fields {
		if tok == "0" {
			res[i] = 0
		} else if tok == "1" {
			res[i] = 1
		} else {
			return nil, fmt.Errorf("invalid connectivity value %q on line %d", tok, i+1)
		}
	}
	return res, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		return stdout.String() + stderr.String(), err
	}
	return stdout.String(), nil
}

func buildTests() []testCase {
	const limit = 2_000_000
	total := 0
	var tests []testCase
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	add := func(grid []string) {
		n := len(grid)
		conn := computeConnectivity(grid)
		tests = append(tests, testCase{n: n, grid: grid, conn: conn})
		total += n * n
	}

	add([]string{"11", "10"})
	add([]string{"10", "01"})
	add([]string{"101", "010", "101"})
	add([]string{"100", "000", "001"})

	for total+16 <= limit && len(tests) < 20 {
		n := rng.Intn(4) + 2
		if total+n*n > limit {
			break
		}
		prob := 0.3 + rng.Float64()*0.3
		add(randomBinaryGrid(rng, n, prob))
	}

	for total+400 <= limit && len(tests) < 40 {
		n := rng.Intn(6) + 5
		if total+n*n > limit {
			break
		}
		prob := 0.25 + rng.Float64()*0.5
		add(randomBinaryGrid(rng, n, prob))
	}

	for total+2500 <= limit && len(tests) < 60 {
		n := rng.Intn(10) + 10
		if total+n*n > limit {
			break
		}
		prob := 0.25 + rng.Float64()*0.5
		add(randomBinaryGrid(rng, n, prob))
	}

	sizes := []int{25, 40, 60, 80, 120, 160, 220, 280}
	for _, n := range sizes {
		if total+n*n > limit {
			continue
		}
		add(randomBinaryGrid(rng, n, 0.45))
	}

	if total+400*400 <= limit {
		add(randomBinaryGrid(rng, 400, 0.5))
	}
	if total+500*500 <= limit {
		add(randomBinaryGrid(rng, 500, 0.4))
	}
	if total+600*600 <= limit {
		add(randomBinaryGrid(rng, 600, 0.5))
	}
	return tests
}

func randomBinaryGrid(rng *rand.Rand, n int, prob float64) []string {
	for {
		grid := make([]string, n)
		ones := 0
		for i := 0; i < n; i++ {
			row := make([]byte, n)
			for j := 0; j < n; j++ {
				if rng.Float64() < prob {
					row[j] = '1'
					ones++
				} else {
					row[j] = '0'
				}
			}
			grid[i] = string(row)
		}
		if ones > 0 {
			return grid
		}
	}
}

type cell struct {
	r int
	c int
}

func computeConnectivity(grid []string) int {
	n := len(grid)
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, n)
	}
	totalOnes := 0
	var queue []cell
	found := false
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if grid[i][j] == '1' {
				totalOnes++
				if !found {
					queue = append(queue, cell{i, j})
					visited[i][j] = true
					found = true
				}
			}
		}
	}
	if !found {
		return 0
	}
	reached := 0
	for head := 0; head < len(queue); head++ {
		cur := queue[head]
		reached++
		dirs := [][2]int{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
		for _, d := range dirs {
			nr := cur.r + d[0]
			nc := cur.c + d[1]
			if nr < 0 || nr >= n || nc < 0 || nc >= n {
				continue
			}
			if visited[nr][nc] || grid[nr][nc] != '1' {
				continue
			}
			visited[nr][nc] = true
			queue = append(queue, cell{nr, nc})
		}
	}
	if reached == totalOnes {
		return 1
	}
	return 0
}

func fetchRowCol(tc testCase, p pair) (string, string) {
	row := tc.grid[p.r-1]
	colBytes := make([]byte, tc.n)
	for i := 0; i < tc.n; i++ {
		colBytes[i] = tc.grid[i][p.c-1]
	}
	return row, string(colBytes)
}
