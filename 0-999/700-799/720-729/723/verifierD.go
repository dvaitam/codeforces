package main

import (
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"
)

type lake struct{ x, y, size int }

type node struct{ x, y int }

func expectedD(n, m, k int, grid []string) (int, []string) {
	g := make([][]rune, n)
	for i := range g {
		g[i] = []rune(grid[i])
	}
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	dx := []int{0, 0, 1, -1}
	dy := []int{1, -1, 0, 0}
	lakes := make([]lake, 0)
	var cells []node
	var flag bool
	var dfs func(x, y int)
	dfs = func(x, y int) {
		visited[x][y] = true
		cells = append(cells, node{x, y})
		if x == 0 || y == 0 || x == n-1 || y == m-1 {
			flag = false
		}
		for dir := 0; dir < 4; dir++ {
			nx, ny := x+dx[dir], y+dy[dir]
			if nx >= 0 && ny >= 0 && nx < n && ny < m && !visited[nx][ny] && g[nx][ny] == '.' {
				dfs(nx, ny)
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if g[i][j] == '.' && !visited[i][j] {
				cells = nil
				flag = true
				dfs(i, j)
				if flag {
					lakes = append(lakes, lake{cells[0].x, cells[0].y, len(cells)})
				}
			}
		}
	}
	sort.Slice(lakes, func(i, j int) bool { return lakes[i].size > lakes[j].size })
	var fill func(x, y int)
	fill = func(x, y int) {
		if g[x][y] == '*' {
			return
		}
		g[x][y] = '*'
		for dir := 0; dir < 4; dir++ {
			nx, ny := x+dx[dir], y+dy[dir]
			if nx >= 0 && ny >= 0 && nx < n && ny < m && g[nx][ny] == '.' {
				fill(nx, ny)
			}
		}
	}
	ans := 0
	for idx := k; idx < len(lakes); idx++ {
		ans += lakes[idx].size
		fill(lakes[idx].x, lakes[idx].y)
	}
	res := make([]string, n)
	for i := 0; i < n; i++ {
		res[i] = string(g[i])
	}
	return ans, res
}

func countLakes(n, m int, grid []string) int {
	g := make([][]rune, n)
	for i := range g {
		g[i] = []rune(grid[i])
	}
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	dx := []int{0, 0, 1, -1}
	dy := []int{1, -1, 0, 0}
	lakesCount := 0
	var flag bool
	var dfs func(x, y int)
	dfs = func(x, y int) {
		visited[x][y] = true
		if x == 0 || y == 0 || x == n-1 || y == m-1 {
			flag = false
		}
		for dir := 0; dir < 4; dir++ {
			nx, ny := x+dx[dir], y+dy[dir]
			if nx >= 0 && ny >= 0 && nx < n && ny < m && !visited[nx][ny] && g[nx][ny] == '.' {
				dfs(nx, ny)
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if g[i][j] == '.' && !visited[i][j] {
				flag = true
				dfs(i, j)
				if flag {
					lakesCount++
				}
			}
		}
	}
	return lakesCount
}

func runCase(exe, input string, expectedAns, n, m, k int, originalGrid []string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	if len(lines) != n+1 {
		return fmt.Errorf("expected %d lines of output, got %d", n+1, len(lines))
	}
	
	gotAns, err := strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return fmt.Errorf("invalid cost: %v", err)
	}
	if gotAns != expectedAns {
		return fmt.Errorf("expected cost %d, got %d", expectedAns, gotAns)
	}
	
	gotGrid := make([]string, n)
	for i := 0; i < n; i++ {
		gotGrid[i] = strings.TrimSpace(lines[i+1])
		if len(gotGrid[i]) != m {
			return fmt.Errorf("invalid grid line length at row %d", i)
		}
	}
	
	changes := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if originalGrid[i][j] != gotGrid[i][j] {
				if originalGrid[i][j] == '.' && gotGrid[i][j] == '*' {
					changes++
				} else {
					return fmt.Errorf("invalid change at (%d, %d)", i, j)
				}
			}
		}
	}
	
	if changes != expectedAns {
		return fmt.Errorf("expected %d changes based on output cost, but got %d changes in grid", expectedAns, changes)
	}
	
	gotLakes := countLakes(n, m, gotGrid)
	if gotLakes != k {
		return fmt.Errorf("expected %d lakes, got %d", k, gotLakes)
	}
	
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for caseNum := 0; caseNum < 100; caseNum++ {
		n := rng.Intn(20) + 3
		m := rng.Intn(20) + 3
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			b := make([]byte, m)
			for j := 0; j < m; j++ {
				if rng.Intn(3) == 0 {
					b[j] = '.'
				} else {
					b[j] = '*'
				}
			}
			grid[i] = string(b)
		}
		
		totalLakes := countLakes(n, m, grid)
		var k int
		if totalLakes > 0 {
			k = rng.Intn(totalLakes + 1)
		} else {
			k = 0
		}
		
		ans, _ := expectedD(n, m, k, grid)
		input := fmt.Sprintf("%d %d %d\n", n, m, k)
		for i := 0; i < n; i++ {
			input += grid[i] + "\n"
		}
		if err := runCase(exe, input, ans, n, m, k, grid); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
