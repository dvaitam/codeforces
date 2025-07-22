package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type pt struct{ x, y, d int }

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return out.String(), nil
}

func bfs(grid [][]byte, dist, sx, sy int) [][]bool {
	n := len(grid)
	m := len(grid[0])
	vis := make([][]bool, n)
	for i := range vis {
		vis[i] = make([]bool, m)
	}
	if sx < 0 || sy < 0 || sx >= n || sy >= m || grid[sx][sy] == 'X' {
		return vis
	}
	q := []pt{{sx, sy, 0}}
	vis[sx][sy] = true
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for head := 0; head < len(q); head++ {
		p := q[head]
		if p.d == dist {
			continue
		}
		nd := p.d + 1
		for _, d := range dirs {
			x := p.x + d[0]
			y := p.y + d[1]
			if x >= 0 && y >= 0 && x < n && y < m && grid[x][y] != 'X' && !vis[x][y] {
				vis[x][y] = true
				q = append(q, pt{x, y, nd})
			}
		}
	}
	return vis
}

func hasSolution(grid [][]byte, dist int) bool {
	// use official solver to check
	n := len(grid)
	m := len(grid[0])
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n-2, m-2, dist))
	for i := 1; i < n-1; i++ {
		sb.WriteString(string(grid[i][1:m-1]) + "\n")
	}
	out, err := run("254D.go", sb.String())
	if err != nil {
		return false
	}
	out = strings.TrimSpace(out)
	return out != "-1"
}

func verifyOut(grid [][]byte, dist int, out string, expectPossible bool) error {
	out = strings.TrimSpace(out)
	if !expectPossible {
		if out != "-1" {
			return fmt.Errorf("expected -1")
		}
		return nil
	}
	fields := strings.Fields(out)
	if len(fields) != 4 {
		return fmt.Errorf("expected four integers")
	}
	vals := make([]int, 4)
	for i, f := range fields {
		v, err := strconv.Atoi(f)
		if err != nil {
			return fmt.Errorf("bad int")
		}
		vals[i] = v
	}
	n := len(grid)
	m := len(grid[0])
	x1, y1, x2, y2 := vals[0]-1, vals[1]-1, vals[2]-1, vals[3]-1
	if x1 < 0 || x1 >= n-2 || y1 < 0 || y1 >= m-2 || x2 < 0 || x2 >= n-2 || y2 < 0 || y2 >= m-2 {
		return fmt.Errorf("coords out of range")
	}
	if x1 == x2 && y1 == y2 {
		return fmt.Errorf("same cell")
	}
	grid2 := make([][]byte, n-2)
	for i := 1; i < n-1; i++ {
		grid2[i-1] = grid[i][1 : m-1]
	}
	cover1 := bfs(grid2, dist, x1, y1)
	cover2 := bfs(grid2, dist, x2, y2)
	for i := 0; i < len(grid2); i++ {
		for j := 0; j < len(grid2[0]); j++ {
			if grid2[i][j] == 'R' && !cover1[i][j] && !cover2[i][j] {
				return fmt.Errorf("rats remain")
			}
		}
	}
	return nil
}

func generateCase(rng *rand.Rand) (string, [][]byte, int) {
	n := rng.Intn(6) + 4
	m := rng.Intn(6) + 4
	dist := rng.Intn(3) + 1
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if i == 0 || j == 0 || i == n-1 || j == m-1 {
				grid[i][j] = 'X'
			} else {
				r := rng.Intn(5)
				switch r {
				case 0:
					grid[i][j] = 'R'
				case 1:
					grid[i][j] = 'X'
				default:
					grid[i][j] = '.'
				}
			}
		}
	}
	// ensure at least one rat and two empty cells
	rats := 0
	empties := 0
	for i := 1; i < n-1; i++ {
		for j := 1; j < m-1; j++ {
			if grid[i][j] == 'R' {
				rats++
			}
			if grid[i][j] != 'X' {
				empties++
			}
		}
	}
	if rats == 0 {
		grid[1][1] = 'R'
		rats = 1
	}
	if empties < 2 {
		grid[1][2] = '.'
		grid[2][1] = '.'
		empties = 2
	}
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", n-2, m-2, dist))
	for i := 1; i < n-1; i++ {
		sb.WriteString(string(grid[i][1 : m-1]))
		sb.WriteByte('\n')
	}
	return sb.String(), grid, dist
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, grid, dist := generateCase(rng)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		expect := hasSolution(grid, dist)
		if e := verifyOut(grid, dist, out, expect); e != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%soutput:\n%s", i+1, e, input, out)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
