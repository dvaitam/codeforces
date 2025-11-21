package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func runCandidate(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, stderr.String())
	}
	return strings.TrimRight(out.String(), " \t\r\n"), nil
}

func expectedNonZero(i, j int) bool {
	return ((i/2)+(j/2))%2 == 0
}

func parseGrid(out string, size int) ([][]bool, error) {
	linesRaw := strings.Split(out, "\n")
	lines := make([]string, 0, len(linesRaw))
	for _, line := range linesRaw {
		line = strings.TrimRight(line, "\r")
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		lines = append(lines, line)
	}
	if len(lines) != size {
		return nil, fmt.Errorf("expected %d non-empty lines, got %d", size, len(lines))
	}
	grid := make([][]bool, size)
	for i := 0; i < size; i++ {
		row := lines[i]
		if len(row) != size {
			return nil, fmt.Errorf("line %d has length %d, expected %d", i+1, len(row), size)
		}
		grid[i] = make([]bool, size)
		for j := 0; j < size; j++ {
			switch row[j] {
			case 'X', 'x':
				grid[i][j] = true
			case '.':
				grid[i][j] = false
			default:
				return nil, fmt.Errorf("invalid character %q at row %d col %d", row[j], i+1, j+1)
			}
		}
	}
	return grid, nil
}

func verifyCase(bin string, n int) error {
	size := 1 << n
	input := fmt.Sprintf("%d\n", n)
	output, err := runCandidate(bin, input)
	if err != nil {
		return err
	}
	grid, err := parseGrid(output, size)
	if err != nil {
		return err
	}
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			exp := expectedNonZero(i, j)
			if grid[i][j] != exp {
				ch := '.'
				if grid[i][j] {
					ch = 'X'
				}
				want := '.'
				if exp {
					want = 'X'
				}
				return fmt.Errorf("N=%d mismatch at row %d col %d: got %c expected %c", n, i+1, j+1, ch, want)
			}
		}
	}
	return nil
}

func main() {
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "--" {
		args = args[1:]
	}
	if len(args) != 1 {
		fmt.Println("usage: go run verifierU2.go /path/to/binary")
		os.Exit(1)
	}
	bin := args[0]
	for _, n := range []int{2, 3, 4, 5} {
		if err := verifyCase(bin, n); err != nil {
			fmt.Fprintf(os.Stderr, "case N=%d failed: %v\n", n, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
