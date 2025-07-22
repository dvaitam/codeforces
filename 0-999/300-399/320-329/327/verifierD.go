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
type Op struct {
	opt  byte
	x, y int
}

func solve(grid [][]byte, n, m int) []Op {
	mark := make([][]bool, n+2)
	vis := make([][]bool, n+2)
	for i := 0; i < n+2; i++ {
		mark[i] = make([]bool, m+2)
		vis[i] = make([]bool, m+2)
	}
	ans := make([]Op, 0, n*m*2)
	stack := make([]Point, 0, n*m)
	dx := [4]int{0, 0, 1, -1}
	dy := [4]int{1, -1, 0, 0}
	for i := 1; i <= n; i++ {
		for j := 1; j <= m; j++ {
			if mark[i][j] || grid[i][j] == '#' {
				continue
			}
			mark[i][j] = true
			stack = append(stack, Point{i, j})
			for len(stack) > 0 {
				a := stack[len(stack)-1]
				if vis[a.x][a.y] {
					stack = stack[:len(stack)-1]
					if len(stack) > 0 {
						ans = append(ans, Op{'D', a.x, a.y})
						ans = append(ans, Op{'R', a.x, a.y})
					}
				} else {
					vis[a.x][a.y] = true
					ans = append(ans, Op{'B', a.x, a.y})
					found := false
					for p := 0; p < 4; p++ {
						x := a.x + dx[p]
						y := a.y + dy[p]
						if x < 1 || x > n || y < 1 || y > m || mark[x][y] || grid[x][y] == '#' {
							continue
						}
						mark[x][y] = true
						stack = append(stack, Point{x, y})
						found = true
					}
					if !found {
						stack = stack[:len(stack)-1]
						if len(stack) > 0 {
							ans[len(ans)-1].opt = 'R'
						}
					}
				}
			}
		}
	}
	return ans
}

func generateCase(rng *rand.Rand) (int, int, [][]byte) {
	n := rng.Intn(3) + 1
	m := rng.Intn(3) + 1
	grid := make([][]byte, n+1)
	empty := false
	for i := 1; i <= n; i++ {
		row := make([]byte, m+1)
		for j := 1; j <= m; j++ {
			if rng.Intn(3) == 0 {
				row[j] = '#'
			} else {
				row[j] = '.'
				empty = true
			}
		}
		grid[i] = row
	}
	if !empty {
		grid[1][1] = '.'
	}
	return n, m, grid
}

func gridToStrings(grid [][]byte, n, m int) []string {
	lines := make([]string, n)
	for i := 1; i <= n; i++ {
		lines[i-1] = string(grid[i][1 : m+1])
	}
	return lines
}

func runCase(bin string, n, m int, grid [][]byte) error {
	lines := gridToStrings(grid, n, m)
	var input strings.Builder
	fmt.Fprintf(&input, "%d %d\n", n, m)
	for _, line := range lines {
		input.WriteString(line)
		input.WriteByte('\n')
	}
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input.String())
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	scanner := bufio.NewScanner(strings.NewReader(out.String()))
	if !scanner.Scan() {
		return fmt.Errorf("no output")
	}
	k, err := strconv.Atoi(strings.TrimSpace(scanner.Text()))
	if err != nil {
		return fmt.Errorf("bad k line")
	}
	var ops []Op
	for i := 0; i < k; i++ {
		if !scanner.Scan() {
			return fmt.Errorf("missing op %d", i+1)
		}
		f := strings.Fields(scanner.Text())
		if len(f) != 3 {
			return fmt.Errorf("bad op line %d", i+1)
		}
		opt := f[0][0]
		x, _ := strconv.Atoi(f[1])
		y, _ := strconv.Atoi(f[2])
		ops = append(ops, Op{opt, x, y})
	}
	if scanner.Err() != nil {
		return scanner.Err()
	}
	exp := solve(grid, n, m)
	if len(exp) != len(ops) {
		return fmt.Errorf("expected %d ops got %d", len(exp), len(ops))
	}
	for i := range exp {
		if exp[i] != ops[i] {
			return fmt.Errorf("op %d expected %c %d %d got %c %d %d", i+1, exp[i].opt, exp[i].x, exp[i].y, ops[i].opt, ops[i].x, ops[i].y)
		}
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierD.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		n, m, grid := generateCase(rng)
		if err := runCase(bin, n, m, grid); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
