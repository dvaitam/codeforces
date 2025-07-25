package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
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

func runCase(exe, input, exp string) error {
	cmd := exec.Command(exe)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	got := strings.TrimSpace(out.String())
	exp = strings.TrimSpace(exp)
	if got != exp {
		return fmt.Errorf("expected %q got %q", exp, got)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	exe := os.Args[1]
	f, err := os.Open("testcasesD.txt")
	if err != nil {
		fmt.Println("could not read testcasesD.txt:", err)
		os.Exit(1)
	}
	defer f.Close()
	scan := bufio.NewScanner(f)
	scan.Split(bufio.ScanLines)
	if !scan.Scan() {
		fmt.Println("invalid test file")
		os.Exit(1)
	}
	t, _ := strconv.Atoi(strings.TrimSpace(scan.Text()))
	for caseNum := 0; caseNum < t; caseNum++ {
		if !scan.Scan() {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		parts := strings.Fields(scan.Text())
		if len(parts) != 3 {
			fmt.Println("bad test file")
			os.Exit(1)
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		k, _ := strconv.Atoi(parts[2])
		grid := make([]string, n)
		for i := 0; i < n; i++ {
			if !scan.Scan() {
				fmt.Println("bad test file")
				os.Exit(1)
			}
			grid[i] = scan.Text()
		}
		ans, g := expectedD(n, m, k, grid)
		input := fmt.Sprintf("%d %d %d\n", n, m, k)
		for i := 0; i < n; i++ {
			input += grid[i] + "\n"
		}
		exp := fmt.Sprintf("%d\n", ans)
		for i := 0; i < n; i++ {
			exp += g[i] + "\n"
		}
		if err := runCase(exe, input, exp); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", caseNum+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
