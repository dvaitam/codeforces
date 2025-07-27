package main

import (
	"bytes"
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type pair struct{ x, y int }

func runBinary(bin, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.CommandContext(ctx, "go", "run", bin)
	} else {
		cmd = exec.CommandContext(ctx, bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			return "", fmt.Errorf("time limit")
		}
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func generateBoard(rng *rand.Rand, n int) [][]byte {
	board := make([][]byte, n)
	for i := 0; i < n; i++ {
		board[i] = make([]byte, n)
		for j := 0; j < n; j++ {
			opts := []byte{'X'}
			if i > 0 {
				opts = append(opts, 'U')
			}
			if i+1 < n {
				opts = append(opts, 'D')
			}
			if j > 0 {
				opts = append(opts, 'L')
			}
			if j+1 < n {
				opts = append(opts, 'R')
			}
			board[i][j] = opts[rng.Intn(len(opts))]
		}
	}
	board[rng.Intn(n)][rng.Intn(n)] = 'X'
	return board
}

func follow(board [][]byte, r, c int) pair {
	n := len(board)
	visited := make(map[pair]bool)
	for step := 0; step < n*n+5; step++ {
		if board[r][c] == 'X' {
			return pair{r + 1, c + 1}
		}
		p := pair{r, c}
		if visited[p] {
			return pair{-1, -1}
		}
		visited[p] = true
		switch board[r][c] {
		case 'U':
			r--
		case 'D':
			r++
		case 'L':
			c--
		case 'R':
			c++
		}
	}
	return pair{-1, -1}
}

func computePairs(board [][]byte) [][]pair {
	n := len(board)
	res := make([][]pair, n)
	for i := 0; i < n; i++ {
		res[i] = make([]pair, n)
		for j := 0; j < n; j++ {
			res[i][j] = follow(board, i, j)
		}
	}
	return res
}

func pairsToInput(p [][]pair) string {
	n := len(p)
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if j > 0 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d %d", p[i][j].x, p[i][j].y))
		}
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseBoard(out string, n int) ([][]byte, error) {
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) < n+1 {
		return nil, fmt.Errorf("not enough lines")
	}
	if strings.TrimSpace(lines[0]) != "VALID" {
		return nil, fmt.Errorf("expected VALID")
	}
	board := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := strings.TrimSpace(lines[i+1])
		if len(row) != n {
			return nil, fmt.Errorf("row %d length", i)
		}
		for j := 0; j < n; j++ {
			ch := row[j]
			switch ch {
			case 'U', 'D', 'L', 'R', 'X':
			default:
				return nil, fmt.Errorf("invalid char")
			}
		}
		board[i] = []byte(row)
	}
	return board, nil
}

func check(candidate string, input string, expected [][]pair) error {
	out, err := runBinary(candidate, input)
	if err != nil {
		return err
	}
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) == 0 {
		return fmt.Errorf("no output")
	}
	if strings.TrimSpace(lines[0]) != "VALID" {
		return fmt.Errorf("candidate says INVALID")
	}
	n := len(expected)
	board, err := parseBoard(out, n)
	if err != nil {
		return err
	}
	pairs := computePairs(board)
	for i := 0; i < n; i++ {
		for j := 0; j < n; j++ {
			if pairs[i][j] != expected[i][j] {
				return fmt.Errorf("mismatch at %d %d", i, j)
			}
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, [][]pair) {
	n := rng.Intn(3) + 1
	board := generateBoard(rng, n)
	pairs := computePairs(board)
	input := pairsToInput(pairs)
	return input, pairs
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for t := 0; t < 100; t++ {
		input, pairs := generateCase(rng)
		if err := check(bin, input, pairs); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:%s", t+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
