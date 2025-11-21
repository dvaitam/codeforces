package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type testCase struct {
	input string
	board []string
	n, m  int
}

type cell struct {
	x, y int
}

func buildReference() (string, error) {
	ref := "./ref97A.bin"
	cmd := exec.Command("go", "build", "-o", ref, "97A.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", fmt.Errorf("failed to build reference: %v\n%s", err, out)
	}
	return ref, nil
}

func runProgram(target, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(target, ".go") {
		cmd = exec.Command("go", "run", target)
	} else {
		cmd = exec.Command(target)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return out.String(), fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func parseOutput(out string, n, m int) (uint64, []string, error) {
	reader := bufio.NewReader(strings.NewReader(out))
	var ways uint64
	if _, err := fmt.Fscan(reader, &ways); err != nil {
		return 0, nil, fmt.Errorf("failed to read number of ways: %v", err)
	}
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		if _, err := fmt.Fscan(reader, &grid[i]); err != nil {
			return 0, nil, fmt.Errorf("failed to read row %d: %v", i+1, err)
		}
		if len(grid[i]) != m {
			return 0, nil, fmt.Errorf("row %d has length %d, expected %d", i+1, len(grid[i]), m)
		}
	}
	return ways, grid, nil
}

func chipLabel(idx int) byte {
	if idx < 26 {
		return byte('a' + idx)
	}
	if idx == 26 {
		return 'A'
	}
	if idx == 27 {
		return 'B'
	}
	return byte('C')
}

func buildBoard(rowsBlocks, colsBlocks int, rng *rand.Rand) []string {
	n := rowsBlocks * 2
	m := colsBlocks * 2
	board := make([][]byte, n)
	for i := range board {
		board[i] = make([]byte, m)
		for j := range board[i] {
			board[i][j] = '.'
		}
	}
	label := 0
	for rb := 0; rb < rowsBlocks; rb++ {
		for cb := 0; cb < colsBlocks; cb++ {
			x := rb * 2
			y := cb * 2
			if rng.Intn(2) == 0 {
				ch := chipLabel(label)
				label++
				board[x][y] = ch
				board[x+1][y] = ch
				ch = chipLabel(label)
				label++
				board[x][y+1] = ch
				board[x+1][y+1] = ch
			} else {
				ch := chipLabel(label)
				label++
				board[x][y] = ch
				board[x][y+1] = ch
				ch = chipLabel(label)
				label++
				board[x+1][y] = ch
				board[x+1][y+1] = ch
			}
		}
	}
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		lines[i] = string(board[i])
	}
	return lines
}

func formatInput(board []string) string {
	var sb strings.Builder
	n := len(board)
	m := len(board[0])
	fmt.Fprintf(&sb, "%d %d\n", n, m)
	for _, row := range board {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	return sb.String()
}

func generateTests() []testCase {
	var tests []testCase
	deterministicSeeds := []int64{1, 2, 3}
	for _, seed := range deterministicSeeds {
		tc := makeTest(seed, 2, 7)
		tests = append(tests, tc)
	}
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	pairs := [][2]int{{1, 14}, {2, 7}, {7, 2}, {14, 1}}
	for len(tests) < 60 {
		p := pairs[rng.Intn(len(pairs))]
		tc := makeTest(rng.Int63(), p[0], p[1])
		tests = append(tests, tc)
	}
	return tests
}

func makeTest(seed int64, rowsBlocks, colsBlocks int) testCase {
	rng := rand.New(rand.NewSource(seed))
	board := buildBoard(rowsBlocks, colsBlocks, rng)
	return testCase{
		input: formatInput(board),
		board: board,
		n:     len(board),
		m:     len(board[0]),
	}
}

func verifySquares(board, grid []string) error {
	n := len(board)
	m := len(board[0])
	covered := make([][]bool, n)
	for i := range covered {
		covered[i] = make([]bool, m)
	}
	squares := 0
	for x := 0; x < n; x++ {
		for y := 0; y < m; y++ {
			if board[x][y] == '.' {
				if grid[x][y] != '.' {
					return fmt.Errorf("cell (%d,%d) should be '.'", x+1, y+1)
				}
				continue
			}
			if covered[x][y] {
				continue
			}
			if x+1 >= n || y+1 >= m {
				return fmt.Errorf("cell (%d,%d) cannot start 2x2 square", x+1, y+1)
			}
			d := grid[x][y]
			coords := [][2]int{{x, y}, {x + 1, y}, {x, y + 1}, {x + 1, y + 1}}
			for _, c := range coords {
				if board[c[0]][c[1]] == '.' {
					return fmt.Errorf("2x2 square at (%d,%d) hits empty cell", x+1, y+1)
				}
				if covered[c[0]][c[1]] {
					return fmt.Errorf("cell (%d,%d) covered twice", c[0]+1, c[1]+1)
				}
				if grid[c[0]][c[1]] != d {
					return fmt.Errorf("2x2 square at (%d,%d) not uniform", x+1, y+1)
				}
			}
			for _, c := range coords {
				covered[c[0]][c[1]] = true
			}
			squares++
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if board[i][j] != '.' && !covered[i][j] {
				return fmt.Errorf("cell (%d,%d) not covered by 2x2 square", i+1, j+1)
			}
		}
	}
	if squares != 14 {
		return fmt.Errorf("expected 14 squares, got %d", squares)
	}
	return nil
}

func checkDominoes(board, grid []string) error {
	n := len(board)
	m := len(board[0])
	chips := make(map[byte][]cell)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if board[i][j] == '.' {
				continue
			}
			chips[board[i][j]] = append(chips[board[i][j]], cell{i, j})
		}
	}
	if len(chips) != 28 {
		return fmt.Errorf("expected 28 chips, got %d", len(chips))
	}
	used := make(map[int]bool)
	for ch, cells := range chips {
		if len(cells) != 2 {
			return fmt.Errorf("chip %c occupies %d cells", ch, len(cells))
		}
		a := grid[cells[0].x][cells[0].y]
		b := grid[cells[1].x][cells[1].y]
		if a < '0' || a > '6' || b < '0' || b > '6' {
			return fmt.Errorf("chip %c uses invalid digits", ch)
		}
		key := encodePair(int(a-'0'), int(b-'0'))
		if used[key] {
			return fmt.Errorf("domino %d-%d reused", min(int(a-'0'), int(b-'0')), max(int(a-'0'), int(b-'0')))
		}
		used[key] = true
	}
	if len(used) != len(chips) {
		return fmt.Errorf("expected %d dominoes, got %d", len(chips), len(used))
	}
	return nil
}

func encodePair(a, b int) int {
	if a > b {
		a, b = b, a
	}
	return a*7 + b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func validateGrid(board, grid []string) error {
	n := len(board)
	m := len(board[0])
	if len(grid) != n {
		return fmt.Errorf("expected %d rows in answer, got %d", n, len(grid))
	}
	for i := 0; i < n; i++ {
		if len(grid[i]) != m {
			return fmt.Errorf("row %d has length %d, expected %d", i+1, len(grid[i]), m)
		}
		for j := 0; j < m; j++ {
			if board[i][j] == '.' {
				if grid[i][j] != '.' {
					return fmt.Errorf("cell (%d,%d) should be '.'", i+1, j+1)
				}
			} else {
				if grid[i][j] < '0' || grid[i][j] > '6' {
					return fmt.Errorf("cell (%d,%d) must contain digit 0-6", i+1, j+1)
				}
			}
		}
	}
	if err := verifySquares(board, grid); err != nil {
		return err
	}
	if err := checkDominoes(board, grid); err != nil {
		return err
	}
	return nil
}

func verifyCase(candidate, ref string, tc testCase) error {
	refOut, err := runProgram(ref, tc.input)
	if err != nil {
		return fmt.Errorf("reference error: %v", err)
	}
	refCount, _, err := parseOutput(refOut, tc.n, tc.m)
	if err != nil {
		return fmt.Errorf("invalid reference output: %v", err)
	}

	candOut, err := runProgram(candidate, tc.input)
	if err != nil {
		return err
	}
	ways, grid, err := parseOutput(candOut, tc.n, tc.m)
	if err != nil {
		return err
	}
	if ways != refCount {
		return fmt.Errorf("expected %d ways, got %d", refCount, ways)
	}
	if err := validateGrid(tc.board, grid); err != nil {
		return err
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	candidate := os.Args[1]

	ref, err := buildReference()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer os.Remove(ref)

	tests := generateTests()
	for i, tc := range tests {
		if err := verifyCase(candidate, ref, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, tc.input)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed.\n", len(tests))
}
