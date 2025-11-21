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

const refSource = "0-999/600-699/640-649/648/648C.go"

type point struct {
	r, c int
}

type board struct {
	n, m  int
	rows  []string
	start point
	total int
}

type testCase struct {
	name  string
	input string
	board board
}

var moves = map[byte]point{
	'U': {-1, 0},
	'D': {1, 0},
	'L': {0, -1},
	'R': {0, 1},
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierC.go /path/to/binary")
		os.Exit(1)
	}

	refBin, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to build reference:", err)
		os.Exit(1)
	}
	defer os.Remove(refBin)

	tests := generateTests()
	candidate := os.Args[1]

	for idx, tc := range tests {
		expOut, err := runBinary(refBin, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "reference runtime error on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
		expPath, err := parsePath(expOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse reference output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, expOut)
			os.Exit(1)
		}
		if err := validatePath(expPath, tc.board); err != nil {
			fmt.Fprintf(os.Stderr, "reference produced invalid path on test %d (%s): %v\n", idx+1, tc.name, err)
			os.Exit(1)
		}

		gotOut, err := runCandidate(candidate, tc.input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "candidate runtime error on test %d (%s): %v\ninput:\n%soutput:\n%s\n", idx+1, tc.name, err, tc.input, gotOut)
			os.Exit(1)
		}
		gotPath, err := parsePath(gotOut)
		if err != nil {
			fmt.Fprintf(os.Stderr, "failed to parse candidate output on test %d (%s): %v\noutput:\n%s\n", idx+1, tc.name, err, gotOut)
			os.Exit(1)
		}
		if err := validatePath(gotPath, tc.board); err != nil {
			fmt.Fprintf(os.Stderr, "wrong answer on test %d (%s): %v\ninput:\n%s\n", idx+1, tc.name, err, tc.input)
			os.Exit(1)
		}
	}

	fmt.Printf("All %d tests passed.\n", len(tests))
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "648C-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	if out, err := cmd.CombinedOutput(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, string(out))
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

func runBinary(path, input string) (string, error) {
	cmd := exec.Command(path)
	return runWithInput(cmd, input)
}

func runWithInput(cmd *exec.Cmd, input string) (string, error) {
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

func parsePath(out string) (string, error) {
	lines := strings.Split(out, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		return line, nil
	}
	return "", fmt.Errorf("output is empty")
}

func validatePath(path string, b board) error {
	if len(path) == 0 {
		return fmt.Errorf("path is empty")
	}
	if len(path) != b.total {
		return fmt.Errorf("path length %d does not match cycle length %d", len(path), b.total)
	}
	visited := make([][]bool, b.n)
	for i := range visited {
		visited[i] = make([]bool, b.m)
	}
	cur := b.start
	visited[cur.r][cur.c] = true
	visitedCount := 1
	for idx := 0; idx < len(path); idx++ {
		ch := path[idx]
		dir, ok := moves[ch]
		if !ok {
			return fmt.Errorf("invalid move %q at position %d", ch, idx+1)
		}
		nr := cur.r + dir.r
		nc := cur.c + dir.c
		if nr < 0 || nr >= b.n || nc < 0 || nc >= b.m {
			return fmt.Errorf("move %d goes out of bounds", idx+1)
		}
		cell := b.rows[nr][nc]
		if cell == '.' {
			return fmt.Errorf("move %d steps onto '.' cell", idx+1)
		}
		if nr == b.start.r && nc == b.start.c {
			if idx != len(path)-1 {
				return fmt.Errorf("returned to start before finishing at move %d", idx+1)
			}
			if visitedCount != b.total {
				return fmt.Errorf("returned to start without visiting all cells")
			}
			cur = point{nr, nc}
			continue
		}
		if visited[nr][nc] {
			return fmt.Errorf("revisited cell (%d,%d) at move %d", nr+1, nc+1, idx+1)
		}
		visited[nr][nc] = true
		visitedCount++
		cur = point{nr, nc}
	}
	if cur != b.start {
		return fmt.Errorf("final position (%d,%d) is not the start (%d,%d)", cur.r+1, cur.c+1, b.start.r+1, b.start.c+1)
	}
	if visitedCount != b.total {
		return fmt.Errorf("visited %d cells, expected %d", visitedCount, b.total)
	}
	return nil
}

func generateTests() []testCase {
	var tests []testCase
	tests = append(tests,
		rectangleCase("small-square", 3, 3, 0, 0, 2, 2, 0),
		rectangleCase("mid-rect", 6, 7, 1, 1, 4, 5, 5),
		rectangleCase("thin-vertical", 10, 5, 1, 1, 8, 3, 7),
		rectangleCase("thin-horizontal", 5, 12, 1, 2, 3, 10, 11),
		rectangleCase("max-border", 100, 100, 0, 0, 99, 99, 17),
		rectangleCase("offset-frame", 15, 20, 3, 4, 11, 15, 23),
	)

	rng := rand.New(rand.NewSource(6480648))
	for i := 0; i < 40; i++ {
		tests = append(tests, randomRectangleCase(rng, fmt.Sprintf("random-rect-%d", i+1)))
	}
	return tests
}

func rectangleCase(name string, n, m int, top, left, bottom, right, startIdx int) testCase {
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	cycle := buildRectangleCycle(top, left, bottom, right)
	if len(cycle) == 0 {
		panic("invalid rectangle cycle")
	}
	for _, p := range cycle {
		grid[p.r][p.c] = '*'
	}
	startIdx = ((startIdx % len(cycle)) + len(cycle)) % len(cycle)
	start := cycle[startIdx]
	grid[start.r][start.c] = 'S'
	rows := make([]string, n)
	for i := range rows {
		rows[i] = string(grid[i])
	}
	return makeTestCase(name, rows)
}

func buildRectangleCycle(top, left, bottom, right int) []point {
	if bottom <= top || right <= left {
		return nil
	}
	var cycle []point
	for c := left; c <= right; c++ {
		cycle = append(cycle, point{top, c})
	}
	for r := top + 1; r <= bottom-1; r++ {
		cycle = append(cycle, point{r, right})
	}
	for c := right; c >= left; c-- {
		cycle = append(cycle, point{bottom, c})
	}
	for r := bottom - 1; r >= top+1; r-- {
		cycle = append(cycle, point{r, left})
	}
	return cycle
}

func randomRectangleCase(rng *rand.Rand, name string) testCase {
	n := rng.Intn(98) + 3
	m := rng.Intn(98) + 3
	top := rng.Intn(n - 1)
	bottom := top + rng.Intn(n-top-1) + 1
	left := rng.Intn(m - 1)
	right := left + rng.Intn(m-left-1) + 1
	startIdx := rng.Intn(1000)
	return rectangleCase(name, n, m, top, left, bottom, right, startIdx)
}

func makeTestCase(name string, rows []string) testCase {
	b, err := newBoard(rows)
	if err != nil {
		panic(fmt.Sprintf("failed to build board %s: %v", name, err))
	}
	input := buildInput(rows)
	return testCase{name: name, board: b, input: input}
}

func newBoard(rows []string) (board, error) {
	n := len(rows)
	if n == 0 {
		return board{}, fmt.Errorf("empty grid")
	}
	m := len(rows[0])
	startCount := 0
	var start point
	total := 0
	for i, row := range rows {
		if len(row) != m {
			return board{}, fmt.Errorf("inconsistent row length")
		}
		for j := 0; j < m; j++ {
			switch row[j] {
			case 'S':
				start = point{i, j}
				startCount++
				total++
			case '*':
				total++
			case '.':
			default:
				return board{}, fmt.Errorf("invalid char %q at (%d,%d)", row[j], i, j)
			}
		}
	}
	if startCount != 1 {
		return board{}, fmt.Errorf("expected exactly one S, got %d", startCount)
	}
	return board{
		n:     n,
		m:     m,
		rows:  append([]string(nil), rows...),
		start: start,
		total: total,
	}, nil
}

func buildInput(rows []string) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", len(rows), len(rows[0]))
	for _, row := range rows {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}
