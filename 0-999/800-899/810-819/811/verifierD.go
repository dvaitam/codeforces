package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type point struct{ x, y int }

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

func genValidGrid(rng *rand.Rand) (int, int, []string) {
	for {
		n := rng.Intn(10) + 2
		m := rng.Intn(10) + 2
		grid, _, _ := genGrid(rng, n, m)
		_, ok, _, _ := bfs(grid, n, m)
		if ok {
			return n, m, grid
		}
	}
}

func runCandidateInteractive(bin string, n, m, swapLR, swapUD int, grid []string) error {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	
	if err := cmd.Start(); err != nil {
		return err
	}
	
	defer func() {
		cmd.Process.Kill()
		cmd.Wait()
	}()
	
	// Send initial input
	fmt.Fprintf(stdin, "%d %d\n", n, m)
	for _, row := range grid {
		fmt.Fprintf(stdin, "%s\n", row)
	}
	
	scanner := bufio.NewScanner(stdout)
	
	x, y := 0, 0
	
	for {
		if !scanner.Scan() {
			break // EOF or error
		}
		move := strings.TrimSpace(scanner.Text())
		if move == "" {
			continue
		}
		
		// Map the candidate's move using swapLR / swapUD
		actualMove := move
		if move == "L" || move == "R" {
			if swapLR == 1 {
				if move == "L" { actualMove = "R" } else { actualMove = "L" }
			}
		} else if move == "U" || move == "D" {
			if swapUD == 1 {
				if move == "U" { actualMove = "D" } else { actualMove = "U" }
			}
		}
		
		nx, ny := x, y
		switch actualMove {
		case "U": nx--
		case "D": nx++
		case "L": ny--
		case "R": ny++
		}
		
		if nx >= 0 && nx < n && ny >= 0 && ny < m && grid[nx][ny] != '*' {
			x, y = nx, ny
		}
		
		if grid[x][y] == 'F' {
			fmt.Fprintf(stdin, "%d %d\n", x+1, y+1)
			break
		}
		
		fmt.Fprintf(stdin, "%d %d\n", x+1, y+1)
	}
	
	if grid[x][y] != 'F' {
		return fmt.Errorf("did not reach finish, ended at %d %d", x+1, y+1)
	}
	
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, grid := genValidGrid(rng)
		swapLR := rng.Intn(2)
		swapUD := rng.Intn(2)
		
		err := runCandidateInteractive(bin, n, m, swapLR, swapUD, grid)
		if err != nil {
			var sb strings.Builder
			sb.WriteString(fmt.Sprintf("%d %d\n", n, m))
			for _, row := range grid {
				sb.WriteString(row)
				sb.WriteByte('\n')
			}
			
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
