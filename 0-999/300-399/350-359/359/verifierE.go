package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
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
	return out.String(), nil
}

var dx = map[byte]int{'U': -1, 'D': 1, 'L': 0, 'R': 0}
var dy = map[byte]int{'U': 0, 'D': 0, 'L': -1, 'R': 1}

func validMove(n int, grid [][]int, x, y int, dir byte) bool {
	nx := x + dx[dir]
	ny := y + dy[dir]
	if nx < 1 || nx > n || ny < 1 || ny > n {
		return false
	}
	tx := nx
	ty := ny
	for tx >= 1 && tx <= n && ty >= 1 && ty <= n {
		if grid[tx][ty] == 1 {
			return true
		}
		tx += dx[dir]
		ty += dy[dir]
	}
	return false
}

func simulate(n, x0, y0 int, grid [][]int, actions string) error {
	if len(actions) > 3000000 {
		return fmt.Errorf("too many actions")
	}
	x, y := x0, y0
	for i := 0; i < len(actions); i++ {
		ch := actions[i]
		switch ch {
		case '1':
			if grid[x][y] == 1 {
				return fmt.Errorf("turn on when already on")
			}
			grid[x][y] = 1
		case '2':
			if grid[x][y] == 0 {
				return fmt.Errorf("turn off when already off")
			}
			grid[x][y] = 0
		case 'L', 'R', 'U', 'D':
			if !validMove(n, grid, x, y, ch) {
				return fmt.Errorf("invalid move %c", ch)
			}
			x += dx[ch]
			y += dy[ch]
		default:
			return fmt.Errorf("invalid character %c", ch)
		}
	}
	if x != x0 || y != y0 {
		return fmt.Errorf("did not return to start")
	}
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if grid[i][j] != 0 {
				return fmt.Errorf("light not off at %d %d", i, j)
			}
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, int, int, [][]int) {
	n := rng.Intn(5) + 2
	x0 := rng.Intn(n) + 1
	y0 := rng.Intn(n) + 1
	grid := make([][]int, n+1)
	for i := 1; i <= n; i++ {
		grid[i] = make([]int, n+1)
		for j := 1; j <= n; j++ {
			if rng.Intn(3) == 0 {
				grid[i][j] = 1
			}
		}
	}
	grid[x0][y0] = 0
	has := false
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if grid[i][j] == 1 {
				has = true
			}
		}
	}
	if !has {
		grid[rng.Intn(n)+1][rng.Intn(n)+1] = 1
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n, x0, y0))
	for i := 1; i <= n; i++ {
		for j := 1; j <= n; j++ {
			if j > 1 {
				sb.WriteByte(' ')
			}
			sb.WriteString(fmt.Sprintf("%d", grid[i][j]))
		}
		sb.WriteByte('\n')
	}
	return sb.String(), x0, y0, grid
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, x0, y0, grid := generateCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) == 0 {
			fmt.Fprintf(os.Stderr, "case %d failed: empty output\ninput:\n%s", i+1, in)
			os.Exit(1)
		}
		ans := strings.TrimSpace(lines[0])
		if ans == "NO" {
			continue
		}
		if ans != "YES" {
			fmt.Fprintf(os.Stderr, "case %d failed: expected YES or NO\ninput:\n%soutput:\n%s", i+1, in, out)
			os.Exit(1)
		}
		actions := ""
		if len(lines) > 1 {
			actions = strings.TrimSpace(lines[1])
		}
		gcopy := make([][]int, len(grid))
		for i := range grid {
			gcopy[i] = append([]int(nil), grid[i]...)
		}
		if err := simulate(len(grid)-1, x0, y0, gcopy, actions); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, err, in, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
