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

var n, m, q int
var grid [][]byte
var dist []int

func idx(r, c int) int { return r*m + c }

func inBounds(r, c int) bool { return r >= 0 && r < n && c >= 0 && c < m }

func computeDist() {
	dist = make([]int, n*m)
	for i := range dist {
		dist[i] = -1
	}
	queue := make([]int, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'v' {
				id := idx(i, j)
				dist[id] = 0
				queue = append(queue, id)
			}
		}
	}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		r, c := v/m, v%m
		d := dist[v] + 1
		if r > 0 {
			u := v - m
			if dist[u] == -1 {
				dist[u] = d
				queue = append(queue, u)
			}
		}
		if r+1 < n {
			u := v + m
			if dist[u] == -1 {
				dist[u] = d
				queue = append(queue, u)
			}
		}
		if c > 0 {
			u := v - 1
			if dist[u] == -1 {
				dist[u] = d
				queue = append(queue, u)
			}
		}
		if c+1 < m {
			u := v + 1
			if dist[u] == -1 {
				dist[u] = d
				queue = append(queue, u)
			}
		}
	}
}

func buildComponent(start, t int) []bool {
	comp := make([]bool, n*m)
	if dist[start] < t {
		return comp
	}
	queue := []int{start}
	comp[start] = true
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		r, c := v/m, v%m
		if r > 0 {
			u := v - m
			if !comp[u] && grid[r-1][c] != '#' && dist[u] >= t {
				comp[u] = true
				queue = append(queue, u)
			}
		}
		if r+1 < n {
			u := v + m
			if !comp[u] && grid[r+1][c] != '#' && dist[u] >= t {
				comp[u] = true
				queue = append(queue, u)
			}
		}
		if c > 0 {
			u := v - 1
			if !comp[u] && grid[r][c-1] != '#' && dist[u] >= t {
				comp[u] = true
				queue = append(queue, u)
			}
		}
		if c+1 < m {
			u := v + 1
			if !comp[u] && grid[r][c+1] != '#' && dist[u] >= t {
				comp[u] = true
				queue = append(queue, u)
			}
		}
	}
	return comp
}

func islandToBorder(comp []bool) bool {
	visited := make([]bool, n*m)
	queue := make([]int, 0)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				id := idx(i, j)
				visited[id] = true
				queue = append(queue, id)
			}
		}
	}
	dirs := [8][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}, {-1, -1}, {-1, 1}, {1, -1}, {1, 1}}
	for head := 0; head < len(queue); head++ {
		v := queue[head]
		r, c := v/m, v%m
		if r == 0 || r == n-1 || c == 0 || c == m-1 {
			return true
		}
		for _, d := range dirs {
			nr, nc := r+d[0], c+d[1]
			if !inBounds(nr, nc) {
				continue
			}
			ni := idx(nr, nc)
			if comp[ni] || visited[ni] {
				continue
			}
			visited[ni] = true
			queue = append(queue, ni)
		}
	}
	return false
}

func check(start, t int) bool {
	if dist[start] < t {
		return false
	}
	comp := buildComponent(start, t)
	if islandToBorder(comp) {
		return false
	}
	return true
}

func solveQuery(x, y int) int {
	start := idx(x-1, y-1)
	left, right := 0, dist[start]
	ans := 0
	for left <= right {
		mid := (left + right) / 2
		if check(start, mid) {
			ans = mid
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return ans
}

func solveCaseF1(g [][]byte, queries [][2]int) []int {
	grid = g
	n = len(g)
	m = len(g[0])
	computeDist()
	res := make([]int, len(queries))
	for i, q := range queries {
		res[i] = solveQuery(q[0], q[1])
	}
	return res
}

func genGrid(rng *rand.Rand) ([][]byte, [][2]int) {
	n = rng.Intn(5) + 3
	m = rng.Intn(5) + 3
	g := make([][]byte, n)
	for i := 0; i < n; i++ {
		g[i] = make([]byte, m)
		for j := 0; j < m; j++ {
			g[i][j] = '.'
		}
	}
	ir := rng.Intn(n-2) + 1
	ic := rng.Intn(m-2) + 1
	g[ir][ic] = '#'
	vr := rng.Intn(n)
	vc := rng.Intn(m)
	for vr == ir && vc == ic {
		vr = rng.Intn(n)
		vc = rng.Intn(m)
	}
	g[vr][vc] = 'v'
	qcnt := rng.Intn(20) + 1
	qs := make([][2]int, qcnt)
	for i := 0; i < qcnt; i++ {
		r, c := rng.Intn(n), rng.Intn(m)
		for g[r][c] == '#' {
			r, c = rng.Intn(n), rng.Intn(m)
		}
		qs[i] = [2]int{r + 1, c + 1}
	}
	return g, qs
}

func formatInput(g [][]byte, qs [][2]int) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d %d\n", len(g), len(g[0]), len(qs)))
	for i := 0; i < len(g); i++ {
		sb.WriteString(string(g[i]))
		sb.WriteByte('\n')
	}
	for i, q := range qs {
		sb.WriteString(fmt.Sprintf("%d %d", q[0], q[1]))
		if i+1 < len(qs) {
			sb.WriteByte('\n')
		}
	}
	sb.WriteByte('\n')
	return sb.String()
}

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
	return strings.TrimSpace(out.String()), nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierF2.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 100; i++ {
		g, qs := genGrid(rng)
		expected := solveCaseF1(g, qs)
		input := formatInput(g, qs)
		out, err := run(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, input)
			os.Exit(1)
		}
		fields := strings.Fields(out)
		if len(fields) != len(expected) {
			fmt.Fprintf(os.Stderr, "case %d failed: expected %d numbers got %d\ninput:\n%s", i+1, len(expected), len(fields), input)
			os.Exit(1)
		}
		for j, f := range fields {
			var v int
			if _, err := fmt.Sscan(f, &v); err != nil || v != expected[j] {
				fmt.Fprintf(os.Stderr, "case %d failed at position %d expected %d got %s\ninput:\n%s", i+1, j+1, expected[j], f, input)
				os.Exit(1)
			}
		}
	}
	fmt.Println("All tests passed")
}
