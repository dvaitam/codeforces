package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

// runProgram runs the binary with the given input and returns its output.
func runProgram(binary string, input string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, binary)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// hasStone checks if a statue occupies (r,c) at time t in the board.
func hasStone(board []string, r, c, t int) bool {
	pr := r - t
	if pr < 0 {
		return false
	}
	return board[pr][c] == 'S'
}

// solve runs BFS to see if Maria can reach Anna.
func solve(board []string) string {
	moves := [9][2]int{{0, 0}, {1, 0}, {-1, 0}, {0, 1}, {0, -1}, {1, 1}, {1, -1}, {-1, 1}, {-1, -1}}
	type state struct{ r, c, t int }
	visited := [8][8][9]bool{}
	q := []state{{7, 0, 0}}
	visited[7][0][0] = true
	for qi := 0; qi < len(q); qi++ {
		s := q[qi]
		r, c, t := s.r, s.c, s.t
		if hasStone(board, r, c, t) {
			continue
		}
		if r == 0 && c == 7 {
			return "WIN"
		}
		if t >= 8 {
			return "WIN"
		}
		for _, mv := range moves {
			nr, nc := r+mv[0], c+mv[1]
			nt := t + 1
			if nr < 0 || nr >= 8 || nc < 0 || nc >= 8 {
				continue
			}
			if hasStone(board, nr, nc, t) || hasStone(board, nr, nc, nt) {
				continue
			}
			if visited[nr][nc][nt] {
				continue
			}
			visited[nr][nc][nt] = true
			q = append(q, state{nr, nc, nt})
		}
	}
	return "LOSE"
}

// randomBoard generates a random board with statues.
func randomBoard(rng *rand.Rand) []string {
	grid := make([][]rune, 8)
	for i := 0; i < 8; i++ {
		grid[i] = []rune("........")
	}
	grid[7][0] = 'M'
	grid[0][7] = 'A'
	num := rng.Intn(7) // up to 6 statues
	for i := 0; i < num; i++ {
		r := rng.Intn(8)
		c := rng.Intn(8)
		if (r == 7 && c == 0) || (r == 0 && c == 7) || grid[r][c] == 'S' {
			i--
			continue
		}
		grid[r][c] = 'S'
	}
	res := make([]string, 8)
	for i := 0; i < 8; i++ {
		res[i] = string(grid[i])
	}
	return res
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: go run verifierA.go /path/to/binary")
		return
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(42))
	for i := 0; i < 100; i++ {
		board := randomBoard(rng)
		input := strings.Join(board, "\n") + "\n"
		expected := solve(board)
		output, err := runProgram(bin, input)
		if err != nil {
			fmt.Printf("Test %d runtime error: %v\n", i+1, err)
			return
		}
		if strings.TrimSpace(output) != expected {
			fmt.Printf("Test %d failed: expected %s got %s\n", i+1, expected, strings.TrimSpace(output))
			return
		}
	}
	fmt.Println("All tests passed")
}
