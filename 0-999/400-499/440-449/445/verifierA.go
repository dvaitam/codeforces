package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
)

// validateOutput checks if out contains a valid coloring for board.
func validateOutput(board [][]byte, out string) error {
	n := len(board)
	if n == 0 {
		return fmt.Errorf("empty board")
	}
	m := len(board[0])
	scanner := bufio.NewScanner(strings.NewReader(out))
	lines := make([]string, 0, n)
	for scanner.Scan() {
		line := strings.TrimRight(scanner.Text(), "\r")
		if line != "" || len(lines) < n {
			lines = append(lines, line)
		}
	}
	if len(lines) < n {
		return fmt.Errorf("expected %d lines, got %d", n, len(lines))
	}
	if len(lines) > n {
		for _, extra := range lines[n:] {
			if strings.TrimSpace(extra) != "" {
				return fmt.Errorf("extra output line: %q", extra)
			}
		}
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		line := lines[i]
		if len(line) != m {
			return fmt.Errorf("line %d has length %d, want %d", i+1, len(line), m)
		}
		grid[i] = []byte(line)
		for j := 0; j < m; j++ {
			if board[i][j] == '-' {
				if grid[i][j] != '-' {
					return fmt.Errorf("cell (%d,%d) should be '-'", i+1, j+1)
				}
			} else {
				if grid[i][j] != 'B' && grid[i][j] != 'W' {
					return fmt.Errorf("cell (%d,%d) invalid char %q", i+1, j+1, grid[i][j])
				}
			}
		}
	}
	// check adjacency
	dirs := [][2]int{{1, 0}, {0, 1}}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '-' {
				continue
			}
			for _, d := range dirs {
				ni := i + d[0]
				nj := j + d[1]
				if ni < n && nj < m && grid[ni][nj] != '-' {
					if grid[ni][nj] == grid[i][j] {
						return fmt.Errorf("adjacent cells (%d,%d) and (%d,%d) share color %c", i+1, j+1, ni+1, nj+1, grid[i][j])
					}
				}
			}
		}
	}
	return nil
}

func randomBoard(rng *rand.Rand, n, m int) [][]byte {
	b := make([][]byte, n)
	for i := 0; i < n; i++ {
		b[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(4) == 0 {
				b[i][j] = '-'
			} else {
				b[i][j] = '.'
			}
		}
	}
	return b
}

func boardToString(b [][]byte) string {
	var sb strings.Builder
	n := len(b)
	m := len(b[0])
	sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
	for i := 0; i < n; i++ {
		sb.Write(b[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	// deterministic RNG for reproducibility
	rng := rand.New(rand.NewSource(1))
	var cases [][][]byte
	// some fixed edge cases
	cases = append(cases, [][]byte{{'.'}})
	cases = append(cases, [][]byte{{'-'}})
	cases = append(cases, [][]byte{{'.', '.'}, {'.', '.'}})
	cases = append(cases, [][]byte{{'.', '-'}, {'-', '.'}})
	cases = append(cases, [][]byte{{'.', '.', '.'}})
	// generate random boards until we have at least 100
	for len(cases) < 100 {
		n := rng.Intn(8) + 1 // 1..8
		m := rng.Intn(8) + 1
		cases = append(cases, randomBoard(rng, n, m))
	}

	for idx, b := range cases {
		tc := boardToString(b)
		cmd := exec.Command(bin)
		cmd.Stdin = strings.NewReader(tc)
		var out bytes.Buffer
		var errBuf bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &errBuf
		err := cmd.Run()
		if err != nil {
			fmt.Printf("case %d: runtime error: %v\nstderr: %s\n", idx+1, err, errBuf.String())
			os.Exit(1)
		}
		if err := validateOutput(b, out.String()); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%soutput:\n%s\n", idx+1, err, tc, out.String())
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
