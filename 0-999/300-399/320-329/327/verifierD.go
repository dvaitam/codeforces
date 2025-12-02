package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Point struct{ x, y int }

const (
	Empty = iota
	Blue
	Red
	Wall
)

func calculateMaxScore(n, m int, grid [][]byte) int {
	visited := make([][]bool, n+1)
	for i := range visited {
		visited[i] = make([]bool, m+1)
	}

	dx := []int{0, 0, 1, -1}
	dy := []int{1, -1, 0, 0}

	totalScore := 0

	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if grid[i][j] == '.' && !visited[i][j] {
				// Found a new component
				size := 0
				q := []Point{{i, j}}
				visited[i][j] = true
				size++

				for len(q) > 0 {
					curr := q[0]
					q = q[1:]

					for k := 0; k < 4; k++ {
						nx, ny := curr.x+dx[k], curr.y+dy[k]
						if nx >= 1 && nx <= n && ny >= 1 && ny <= m && grid[nx][ny] == '.' && !visited[nx][ny] {
							visited[nx][ny] = true
							size++
							q = append(q, Point{nx, ny})
						}
					}
				}

				// Score for component of size S:
				// 1 Blue + (S-1) Red
				// 100 + (S-1)*200
				if size > 0 {
					totalScore += 100 + (size-1)*200
				}
			}
		}
	}
	return totalScore
}

func runCase(bin string, n, m int, grid [][]byte) error {
	// Prepare input string
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, m)
	for i := 1; i <= n; i++ {
		input.Write(grid[i][1 : m+1])
		input.WriteByte('\n')
	}

	// Run user binary
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr // Let user stderr pass through for debugging
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}

	// Parse output
	scanner := bufio.NewScanner(&out)
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	kStr := strings.TrimSpace(scanner.Text())
	k, err := strconv.Atoi(kStr)
	if err != nil {
		return fmt.Errorf("invalid k: %v", err)
	}

	// Simulation Board
	// 0: Empty, 1: Blue, 2: Red
	board := make([][]int, n+1)
	for i := range board {
		board[i] = make([]int, m+1)
	}

	dx := []int{0, 0, 1, -1}
	dy := []int{1, -1, 0, 0}

	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("expected %d ops, got %d", k, i)
		}
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) != 3 {
			return fmt.Errorf("invalid op format at line %d: %s", i+2, line)
		}

		opType := parts[0]
		r, err1 := strconv.Atoi(parts[1])
		c, err2 := strconv.Atoi(parts[2])

		if err1 != nil || err2 != nil {
			return fmt.Errorf("invalid coords at line %d: %s", i+2, line)
		}

		if r < 1 || r > n || c < 1 || c > m {
			return fmt.Errorf("coords out of bounds at line %d: %d %d", i+2, r, c)
		}

		if grid[r][c] == '#' {
			return fmt.Errorf("operation on wall at %d %d", r, c)
		}

		switch opType {
		case "B":
			if board[r][c] != Empty {
				return fmt.Errorf("building Blue on non-empty cell at %d %d", r, c)
			}
			board[r][c] = Blue
		case "R":
			if board[r][c] != Empty {
				return fmt.Errorf("building Red on non-empty cell at %d %d", r, c)
			}
			// Check neighbors for Blue
			hasBlue := false
			for d := 0; d < 4; d++ {
				nr, nc := r+dx[d], c+dy[d]
				if nr >= 1 && nr <= n && nc >= 1 && nc <= m && board[nr][nc] == Blue {
					hasBlue = true
					break
				}
			}
			if !hasBlue {
				return fmt.Errorf("building Red at %d %d without Blue neighbor", r, c)
			}
			board[r][c] = Red
		case "D":
			if board[r][c] == Empty {
				return fmt.Errorf("destroying empty cell at %d %d", r, c)
			}
			board[r][c] = Empty
		default:
			return fmt.Errorf("unknown operation type %s", opType)
		}
	}

	// Calculate final score
	userScore := 0
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if board[i][j] == Blue {
				userScore += 100
			} else if board[i][j] == Red {
				userScore += 200
			}
		}
	}

	expectedScore := calculateMaxScore(n, m, grid)
	if userScore != expectedScore {
		return fmt.Errorf("score mismatch: expected %d, got %d", expectedScore, userScore)
	}

	return nil
}

func generateCase(rng *rand.Rand) (int, int, [][]byte) {
	n := rng.Intn(10) + 1 // Increased size slightly for better coverage
	m := rng.Intn(10) + 1
	grid := make([][]byte, n+1)
	hasEmpty := false
	for i := 1; i <= n; i++ {
		grid[i] = make([]byte, m+1)
		for j := 1; j <= m; j++ {
			if rng.Intn(3) == 0 { // 33% chance of wall
				grid[i][j] = '#'
			} else {
				grid[i][j] = '.'
				hasEmpty = true
			}
		}
	}
	// Ensure at least one empty cell to avoid trivial 0 score cases constantly
	if !hasEmpty {
		grid[1][1] = '.'
	}
	return n, m, grid
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "usage: go run verifierD.go /path/to/binary\n")
		os.Exit(1)
	}
	bin := os.Args[1]
	
	// Compile the solution if it's a .go file (convenience)
	if strings.HasSuffix(bin, ".go") {
		cmd := exec.Command("go", "build", "-o", "solution", bin)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "failed to build solution: %v\n", err)
			os.Exit(1)
		}
		bin = "./solution"
		defer os.Remove("./solution")
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 50; i++ {
		n, m, grid := generateCase(rng)
		if err := runCase(bin, n, m, grid); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}