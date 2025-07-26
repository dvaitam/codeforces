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

var dx = []int{-1, 0, 1, 0}
var dy = []int{0, 1, 0, -1}

func run(bin, input string) (string, error) {
	var cmd *exec.Cmd
	if strings.HasSuffix(bin, ".go") {
		cmd = exec.Command("go", "run", bin)
	} else {
		cmd = exec.Command(bin)
	}
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	var errb bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &errb
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("runtime error: %v\n%s", err, errb.String())
	}
	return strings.TrimSpace(out.String()), nil
}

func inBounds(x, y, n, m int) bool { return x >= 0 && x < n && y >= 0 && y < m }

func solveB(grid [][]byte, sx, sy, tx, ty int) string {
	n, m := len(grid), len(grid[0])
	dist := make([][][]int, n)
	for i := 0; i < n; i++ {
		dist[i] = make([][]int, m)
		for j := 0; j < m; j++ {
			dist[i][j] = []int{3, 3, 3, 3}
		}
	}
	type node struct{ x, y, d, t int }
	q := []node{}
	for d := 0; d < 4; d++ {
		dist[sx][sy][d] = 0
		q = append(q, node{sx, sy, d, 0})
	}
	for head := 0; head < len(q); head++ {
		cur := q[head]
		if cur.t > dist[cur.x][cur.y][cur.d] || cur.t > 2 {
			continue
		}
		if cur.x == tx && cur.y == ty {
			return "YES"
		}
		for nd := 0; nd < 4; nd++ {
			nx, ny := cur.x+dx[nd], cur.y+dy[nd]
			nt := cur.t
			if nd != cur.d {
				nt++
			}
			if !inBounds(nx, ny, n, m) || grid[nx][ny] == '*' || nt >= dist[nx][ny][nd] {
				continue
			}
			dist[nx][ny][nd] = nt
			q = append(q, node{nx, ny, nd, nt})
		}
	}
	return "NO"
}

func genCase(rng *rand.Rand) ([][]byte, int, int, int, int) {
	n := rng.Intn(5) + 2
	m := rng.Intn(5) + 2
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		row := make([]byte, m)
		for j := 0; j < m; j++ {
			if rng.Intn(4) == 0 {
				row[j] = '*'
			} else {
				row[j] = '.'
			}
		}
		grid[i] = row
	}
	sx, sy := rng.Intn(n), rng.Intn(m)
	tx, ty := rng.Intn(n), rng.Intn(m)
	if sx == tx && sy == ty {
		tx = (tx + 1) % n
	}
	grid[sx][sy] = 'S'
	grid[tx][ty] = 'T'
	return grid, sx, sy, tx, ty
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		grid, sx, sy, tx, ty := genCase(rng)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d\n", len(grid), len(grid[0]))
		for _, row := range grid {
			sb.WriteString(string(row))
			sb.WriteByte('\n')
		}
		expect := solveB(grid, sx, sy, tx, ty)
		got, err := run(bin, sb.String())
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\n", i+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != expect {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%sexpected %s got %s\n", i+1, sb.String(), expect, got)
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
