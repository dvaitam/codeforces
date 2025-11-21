package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const refSource = "0-999/500-599/530-539/538/538D.go"

type move struct {
	dx int
	dy int
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/candidate")
		os.Exit(1)
	}
	candidate := os.Args[1]

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		fail("failed to read input: %v", err)
	}

	n, grid, err := parseInput(input)
	if err != nil {
		fail("failed to parse input: %v", err)
	}
	boardSize := 2*n - 1

	refBin, err := buildReference()
	if err != nil {
		fail("failed to build reference: %v", err)
	}
	defer os.Remove(refBin)

	refOut, err := runProgram(refBin, input)
	if err != nil {
		fail("reference execution failed: %v", err)
	}
	refStatus, err := parseDecision(refOut)
	if err != nil {
		fail("failed to parse reference output: %v", err)
	}

	userOut, err := runProgram(candidate, input)
	if err != nil {
		fail("candidate execution failed: %v", err)
	}
	userStatus, board, err := parseCandidateOutput(userOut, boardSize)
	if err != nil {
		fail("failed to parse candidate output: %v", err)
	}

	if refStatus == "NO" {
		if userStatus != "NO" {
			fail("expected answer NO, but candidate produced YES")
		}
		fmt.Println("OK")
		return
	}

	if userStatus != "YES" {
		fail("expected answer YES, but candidate produced NO")
	}

	if err := validateBoard(board, grid); err != nil {
		fail("invalid board: %v", err)
	}

	fmt.Println("OK")
}

func buildReference() (string, error) {
	tmp, err := os.CreateTemp("", "538D-ref-*")
	if err != nil {
		return "", err
	}
	tmp.Close()

	cmd := exec.Command("go", "build", "-o", tmp.Name(), filepath.Clean(refSource))
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		os.Remove(tmp.Name())
		return "", fmt.Errorf("%v\n%s", err, out.String())
	}
	return tmp.Name(), nil
}

func runProgram(bin string, input []byte) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = bytes.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	return out.String(), err
}

func parseInput(data []byte) (int, [][]byte, error) {
	reader := bufio.NewReader(bytes.NewReader(data))
	var n int
	if _, err := fmt.Fscan(reader, &n); err != nil {
		return 0, nil, err
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		var row string
		if _, err := fmt.Fscan(reader, &row); err != nil {
			return 0, nil, err
		}
		if len(row) != n {
			return 0, nil, fmt.Errorf("row %d has length %d (expected %d)", i+1, len(row), n)
		}
		grid[i] = []byte(row)
	}
	return n, grid, nil
}

func parseDecision(out string) (string, error) {
	lines := splitLines(out)
	for _, line := range lines {
		trim := strings.TrimSpace(line)
		if trim == "" {
			continue
		}
		if trim == "YES" || trim == "NO" {
			return trim, nil
		}
		return "", fmt.Errorf("unexpected token %q", trim)
	}
	return "", fmt.Errorf("empty output")
}

func parseCandidateOutput(out string, boardSize int) (string, []string, error) {
	lines := splitLines(out)
	idx := 0
	for idx < len(lines) && strings.TrimSpace(lines[idx]) == "" {
		idx++
	}
	if idx == len(lines) {
		return "", nil, fmt.Errorf("no output")
	}
	first := strings.TrimSpace(lines[idx])
	idx++
	if first != "YES" && first != "NO" {
		return "", nil, fmt.Errorf("expected YES or NO, got %q", first)
	}
	if first == "NO" {
		for ; idx < len(lines); idx++ {
			if strings.TrimSpace(lines[idx]) != "" {
				return "", nil, fmt.Errorf("unexpected extra output after NO")
			}
		}
		return "NO", nil, nil
	}

	board := make([]string, 0, boardSize)
	for idx < len(lines) && len(board) < boardSize {
		line := strings.TrimSpace(lines[idx])
		idx++
		if line == "" {
			continue
		}
		board = append(board, line)
	}
	if len(board) != boardSize {
		return "", nil, fmt.Errorf("expected %d board lines, got %d", boardSize, len(board))
	}
	for ; idx < len(lines); idx++ {
		if strings.TrimSpace(lines[idx]) != "" {
			return "", nil, fmt.Errorf("unexpected extra output after board")
		}
	}
	return "YES", board, nil
}

func validateBoard(board []string, grid [][]byte) error {
	n := len(grid)
	if len(board) != 2*n-1 {
		return fmt.Errorf("invalid board size")
	}
	size := len(board)
	center := n - 1
	moves := make([]move, 0)
	oCount := 0
	for i := 0; i < size; i++ {
		if len(board[i]) != size {
			return fmt.Errorf("line %d has length %d (expected %d)", i+1, len(board[i]), size)
		}
		for j := 0; j < size; j++ {
			ch := board[i][j]
			if ch != '.' && ch != 'x' && ch != 'o' {
				return fmt.Errorf("invalid character %q", ch)
			}
			if ch == 'o' {
				if i != center || j != center {
					return fmt.Errorf("'o' must be only at the board center")
				}
				oCount++
			}
			if ch == 'x' {
				moves = append(moves, move{dx: j - center, dy: i - center})
			}
		}
	}
	if oCount != 1 {
		return fmt.Errorf("board must contain exactly one 'o' at the center")
	}

	nRows := len(grid)
	nCols := len(grid[0])
	pieces := make([][2]int, 0)
	for i := 0; i < nRows; i++ {
		for j := 0; j < nCols; j++ {
			if grid[i][j] == 'o' {
				pieces = append(pieces, [2]int{i, j})
			}
		}
	}

	attacked := make([][]bool, nRows)
	for i := 0; i < nRows; i++ {
		attacked[i] = make([]bool, nCols)
	}

	for _, p := range pieces {
		for _, mv := range moves {
			nr := p[0] + mv.dy
			nc := p[1] + mv.dx
			if nr < 0 || nr >= nRows || nc < 0 || nc >= nCols {
				continue
			}
			if grid[nr][nc] == 'o' {
				continue
			}
			attacked[nr][nc] = true
		}
	}

	for i := 0; i < nRows; i++ {
		for j := 0; j < nCols; j++ {
			if grid[i][j] == 'x' && !attacked[i][j] {
				return fmt.Errorf("cell (%d,%d) should be attacked but is not", i+1, j+1)
			}
			if grid[i][j] == '.' && attacked[i][j] {
				return fmt.Errorf("cell (%d,%d) should not be attacked", i+1, j+1)
			}
		}
	}
	return nil
}

func splitLines(out string) []string {
	lines := strings.Split(out, "\n")
	for i := range lines {
		lines[i] = strings.TrimRight(lines[i], "\r")
	}
	return lines
}

func fail(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(1)
}
