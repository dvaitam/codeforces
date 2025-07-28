package main

import (
	"bufio"
	"bytes"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

const numTestsE = 100

func prepareBinary(path string) (string, func(), error) {
	if strings.HasSuffix(path, ".go") {
		tmp := filepath.Join(os.TempDir(), "binE")
		cmd := exec.Command("go", "build", "-o", tmp, path)
		if out, err := cmd.CombinedOutput(); err != nil {
			return "", nil, fmt.Errorf("go build failed: %v: %s", err, out)
		}
		return tmp, func() { os.Remove(tmp) }, nil
	}
	return path, nil, nil
}

func prepareOracle() (string, func(), error) {
	tmp := filepath.Join(os.TempDir(), "oracleE")
	cmd := exec.Command("go", "build", "-o", tmp, "1799E.go")
	if out, err := cmd.CombinedOutput(); err != nil {
		return "", nil, fmt.Errorf("go build oracle failed: %v: %s", err, out)
	}
	return tmp, func() { os.Remove(tmp) }, nil
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err := cmd.Run()
	return strings.TrimSpace(buf.String()), err
}

type gridCase struct {
	n, m int
	grid [][]byte
}

func generateCase(rng *rand.Rand) gridCase {
	n := rng.Intn(3) + 3
	m := rng.Intn(3) + 3
	grid := make([][]byte, n)
	for i := range grid {
		grid[i] = make([]byte, m)
		for j := range grid[i] {
			grid[i][j] = '.'
		}
	}
	// place two cities as single cells
	r1, c1 := rng.Intn(n), rng.Intn(m)
	r2, c2 := rng.Intn(n), rng.Intn(m)
	for r1 == r2 && c1 == c2 {
		r2, c2 = rng.Intn(n), rng.Intn(m)
	}
	grid[r1][c1] = '#'
	grid[r2][c2] = '#'
	return gridCase{n, m, grid}
}

func formatInput(tc gridCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "1\n%d %d\n", tc.n, tc.m)
	for i := 0; i < tc.n; i++ {
		sb.Write(tc.grid[i])
		sb.WriteByte('\n')
	}
	return sb.String()
}

func parseGrid(out string, n, m int) ([][]byte, error) {
	scanner := bufio.NewScanner(strings.NewReader(out))
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		if !scanner.Scan() {
			return nil, fmt.Errorf("missing row")
		}
		line := scanner.Text()
		if len(line) != m {
			return nil, fmt.Errorf("bad row length")
		}
		grid[i] = []byte(line)
	}
	if scanner.Scan() {
		return nil, fmt.Errorf("extra output")
	}
	return grid, nil
}

func cost(orig, ans [][]byte) int {
	n := len(orig)
	m := len(orig[0])
	cnt := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if orig[i][j] == '.' && ans[i][j] == '#' {
				cnt++
			}
		}
	}
	return cnt
}

func connected(grid [][]byte) bool {
	n := len(grid)
	m := len(grid[0])
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	visited := make([][]bool, n)
	for i := range visited {
		visited[i] = make([]bool, m)
	}
	var q [][2]int
	found := false
	for i := 0; i < n && !found; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				q = append(q, [2]int{i, j})
				visited[i][j] = true
				found = true
				break
			}
		}
	}
	if !found {
		return true
	}
	idx := 0
	for idx < len(q) {
		x, y := q[idx][0], q[idx][1]
		idx++
		for _, d := range dirs {
			nx, ny := x+d[0], y+d[1]
			if nx >= 0 && nx < n && ny >= 0 && ny < m && grid[nx][ny] == '#' && !visited[nx][ny] {
				visited[nx][ny] = true
				q = append(q, [2]int{nx, ny})
			}
		}
	}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' && !visited[i][j] {
				return false
			}
		}
	}
	return true
}

func manhattanOK(grid [][]byte) bool {
	n := len(grid)
	m := len(grid[0])
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	cells := [][2]int{}
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == '#' {
				cells = append(cells, [2]int{i, j})
			}
		}
	}
	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, m)
	}
	for _, c := range cells {
		for i := 0; i < n; i++ {
			for j := 0; j < m; j++ {
				dist[i][j] = -1
			}
		}
		q := [][2]int{c}
		dist[c[0]][c[1]] = 0
		idx := 0
		for idx < len(q) {
			x, y := q[idx][0], q[idx][1]
			idx++
			for _, d := range dirs {
				nx, ny := x+d[0], y+d[1]
				if nx >= 0 && nx < n && ny >= 0 && ny < m && grid[nx][ny] == '#' && dist[nx][ny] == -1 {
					dist[nx][ny] = dist[x][y] + 1
					q = append(q, [2]int{nx, ny})
				}
			}
		}
		for _, c2 := range cells {
			md := abs(c[0]-c2[0]) + abs(c[1]-c2[1])
			if dist[c2[0]][c2[1]] != md {
				return false
			}
		}
	}
	return true
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func validate(orig, ans [][]byte) error {
	n := len(orig)
	m := len(orig[0])
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if orig[i][j] == '#' && ans[i][j] != '#' {
				return fmt.Errorf("removed filled cell")
			}
		}
	}
	if !connected(ans) {
		return fmt.Errorf("not connected")
	}
	if !manhattanOK(ans) {
		return fmt.Errorf("manhattan property failed")
	}
	return nil
}

func runCase(bin, oracle string, tc gridCase) error {
	input := formatInput(tc)
	oracleOut, err := run(oracle, input)
	if err != nil {
		return fmt.Errorf("oracle error: %v", err)
	}
	oracleGrid, err := parseGrid(oracleOut, tc.n, tc.m)
	if err != nil {
		return fmt.Errorf("oracle output bad: %v", err)
	}
	expectCost := cost(tc.grid, oracleGrid)
	gotOut, err := run(bin, input)
	if err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	gotGrid, err := parseGrid(gotOut, tc.n, tc.m)
	if err != nil {
		return err
	}
	if err := validate(tc.grid, gotGrid); err != nil {
		return err
	}
	gotCost := cost(tc.grid, gotGrid)
	if gotCost != expectCost {
		return fmt.Errorf("expected cost %d got %d", expectCost, gotCost)
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: go run verifierE.go /path/to/binary")
		return
	}
	binPath := os.Args[1]
	bin, cleanup, err := prepareBinary(binPath)
	if err != nil {
		fmt.Println("compile error:", err)
		return
	}
	if cleanup != nil {
		defer cleanup()
	}
	oracle, cleanOracle, err := prepareOracle()
	if err != nil {
		fmt.Println(err)
		return
	}
	defer cleanOracle()
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < numTestsE; i++ {
		tc := generateCase(rng)
		if err := runCase(bin, oracle, tc); err != nil {
			fmt.Printf("case %d failed: %v\ninput:\n%s", i+1, err, formatInput(tc))
			return
		}
	}
	fmt.Println("All tests passed")
}
