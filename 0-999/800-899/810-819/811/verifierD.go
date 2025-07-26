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

type point struct{ x, y int }

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
	return strings.TrimSpace(out.String()), nil
}

var dirs = []struct {
	dx, dy int
	c      byte
}{
	{-1, 0, 'U'},
	{1, 0, 'D'},
	{0, -1, 'L'},
	{0, 1, 'R'},
}

func bfs(grid []string, n, m int) ([]byte, bool, int, int) {
	parent := make([][]point, n)
	move := make([][]byte, n)
	used := make([][]bool, n)
	for i := range parent {
		parent[i] = make([]point, m)
		move[i] = make([]byte, m)
		used[i] = make([]bool, m)
	}
	q := []point{{0, 0}}
	used[0][0] = true
	var fx, fy int
	found := false
	for len(q) > 0 && !found {
		p := q[0]
		q = q[1:]
		if grid[p.x][p.y] == 'F' {
			fx, fy = p.x, p.y
			found = true
			break
		}
		for _, d := range dirs {
			nx, ny := p.x+d.dx, p.y+d.dy
			if nx < 0 || nx >= n || ny < 0 || ny >= m {
				continue
			}
			if grid[nx][ny] == '*' || used[nx][ny] {
				continue
			}
			used[nx][ny] = true
			parent[nx][ny] = p
			move[nx][ny] = d.c
			q = append(q, point{nx, ny})
		}
	}
	if !found {
		return nil, false, 0, 0
	}
	path := []byte{}
	x, y := fx, fy
	for x != 0 || y != 0 {
		path = append(path, move[x][y])
		p := parent[x][y]
		x, y = p.x, p.y
	}
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
	return path, true, fx, fy
}

func solveCase(n, m int, swapLR, swapUD int, grid []string) string {
	path, ok, _, _ := bfs(grid, n, m)
	if !ok {
		return ""
	}
	for i, c := range path {
		if c == 'L' || c == 'R' {
			if swapLR == 1 {
				if c == 'L' {
					path[i] = 'R'
				} else {
					path[i] = 'L'
				}
			}
		} else {
			if swapUD == 1 {
				if c == 'U' {
					path[i] = 'D'
				} else {
					path[i] = 'U'
				}
			}
		}
	}
	var sb strings.Builder
	for i, c := range path {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteByte(c)
	}
	return sb.String()
}

func genGrid(rng *rand.Rand, n, m int) ([]string, int, int) {
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			grid[i][j] = '.'
		}
	}
	fx := rng.Intn(n)
	fy := rng.Intn(m)
	grid[fx][fy] = 'F'
	path := []point{{0, 0}}
	visited := map[point]bool{{0, 0}: true}
	for len(path) == 1 || path[len(path)-1] != (point{fx, fy}) {
		cur := path[len(path)-1]
		if cur == (point{fx, fy}) {
			break
		}
		var opts []point
		for _, d := range dirs {
			nx, ny := cur.x+d.dx, cur.y+d.dy
			if nx < 0 || nx >= n || ny < 0 || ny >= m {
				continue
			}
			p := point{nx, ny}
			if !visited[p] {
				opts = append(opts, p)
			}
		}
		if len(opts) == 0 {
			break
		}
		next := opts[rng.Intn(len(opts))]
		visited[next] = true
		path = append(path, next)
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if (i == 0 && j == 0) || (i == fx && j == fy) {
				continue
			}
			if visited[point{i, j}] {
				continue
			}
			if rng.Intn(5) == 0 {
				grid[i][j] = '*'
			} else {
				grid[i][j] = '.'
			}
		}
	}
	lines := make([]string, n)
	for i := 0; i < n; i++ {
		lines[i] = string(grid[i])
	}
	return lines, fx, fy
}

func genCase(rng *rand.Rand) (string, string) {
	n := rng.Intn(5) + 2
	m := rng.Intn(5) + 2
	grid, _, _ := genGrid(rng, n, m)
	swapLR := rng.Intn(2)
	swapUD := rng.Intn(2)
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d %d %d\n", n, m, swapLR, swapUD)
	for _, row := range grid {
		sb.WriteString(row)
		sb.WriteByte('\n')
	}
	expected := solveCase(n, m, swapLR, swapUD, grid)
	return sb.String(), expected
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		in, exp := genCase(rng)
		out, err := runCandidate(bin, in)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, in)
			os.Exit(1)
		}
		if strings.TrimSpace(out) != strings.TrimSpace(exp) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %s got %s\ninput:\n%s", i+1, exp, out, in)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
