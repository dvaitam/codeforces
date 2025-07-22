package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"strings"
	"time"
)

type mazeCase struct {
	n, m, k int
	grid    [][]byte
}

func generateCase(rng *rand.Rand) (string, mazeCase) {
	n := rng.Intn(6) + 1
	m := rng.Intn(6) + 1
	total := n * m
	open := rng.Intn(total) + 1
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := range grid[i] {
			grid[i][j] = '#'
		}
	}
	sx, sy := rng.Intn(n), rng.Intn(m)
	grid[sx][sy] = '.'
	open--
	cells := [][2]int{{sx, sy}}
	for open > 0 {
		c := cells[rng.Intn(len(cells))]
		dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
		d := dirs[rng.Intn(4)]
		nx, ny := c[0]+d[0], c[1]+d[1]
		if nx < 0 || nx >= n || ny < 0 || ny >= m || grid[nx][ny] == '.' {
			continue
		}
		grid[nx][ny] = '.'
		cells = append(cells, [2]int{nx, ny})
		open--
	}
	free := len(cells)
	k := rng.Intn(free)
	var b strings.Builder
	fmt.Fprintf(&b, "%d %d %d\n", n, m, k)
	for i := 0; i < n; i++ {
		b.Write(grid[i])
		b.WriteByte('\n')
	}
	return b.String(), mazeCase{n: n, m: m, k: k, grid: grid}
}

func verify(tc mazeCase, out string) error {
	scanner := bufio.NewScanner(strings.NewReader(out))
	res := make([]string, 0, tc.n)
	for scanner.Scan() && len(res) < tc.n {
		line := scanner.Text()
		if len(line) != tc.m {
			return fmt.Errorf("expected %d columns", tc.m)
		}
		res = append(res, line)
	}
	if len(res) != tc.n {
		return fmt.Errorf("expected %d rows", tc.n)
	}
	changed := 0
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			ch := res[i][j]
			if tc.grid[i][j] == '#' {
				if ch != '#' {
					return fmt.Errorf("cannot modify walls")
				}
			} else {
				if ch == '.' {
				} else if ch == 'X' {
					changed++
				} else {
					return fmt.Errorf("invalid char %c", ch)
				}
			}
		}
	}
	if changed != tc.k {
		return fmt.Errorf("expected %d X, got %d", tc.k, changed)
	}
	// connectivity
	var q [][2]int
	vis := make([][]bool, tc.n)
	for i := range vis {
		vis[i] = make([]bool, tc.m)
	}
	found := false
	for i := 0; i < tc.n && !found; i++ {
		for j := 0; j < tc.m; j++ {
			if res[i][j] == '.' {
				q = append(q, [2]int{i, j})
				vis[i][j] = true
				found = true
				break
			}
		}
	}
	if !found {
		return fmt.Errorf("no free cells left")
	}
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	idx := 0
	for idx < len(q) {
		x, y := q[idx][0], q[idx][1]
		idx++
		for _, d := range dirs {
			nx, ny := x+d[0], y+d[1]
			if nx < 0 || nx >= tc.n || ny < 0 || ny >= tc.m {
				continue
			}
			if vis[nx][ny] || res[nx][ny] != '.' {
				continue
			}
			vis[nx][ny] = true
			q = append(q, [2]int{nx, ny})
		}
	}
	total := 0
	for i := 0; i < tc.n; i++ {
		for j := 0; j < tc.m; j++ {
			if res[i][j] == '.' {
				total++
			}
		}
	}
	if total != len(q) {
		return fmt.Errorf("cells not connected")
	}
	return nil
}

func runCase(bin string, input string, tc mazeCase) error {
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
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	return verify(tc, out.String())
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierA.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		input, tc := generateCase(rng)
		if err := runCase(bin, input, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
