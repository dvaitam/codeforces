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

type testCase struct {
	grid []string
}

func buildCase(grid []string) testCase {
	return testCase{grid: grid}
}

func generateCase(rng *rand.Rand) testCase {
	n := rng.Intn(5) + 3 // 3..7
	grid := make([]string, n)
	for i := 0; i < n; i++ {
		row := make([]byte, n)
		for j := 0; j < n; j++ {
			if i == 0 && j == 0 {
				row[j] = 'S'
			} else if i == n-1 && j == n-1 {
				row[j] = 'F'
			} else {
				if rng.Intn(2) == 0 {
					row[j] = '0'
				} else {
					row[j] = '1'
				}
			}
		}
		grid[i] = string(row)
	}
	return buildCase(grid)
}

func reachable(grid [][]byte, digit byte) bool {
	n := len(grid)
	type pair struct{ r, c int }
	q := []pair{{0, 0}}
	vis := make([][]bool, n)
	for i := range vis {
		vis[i] = make([]bool, n)
	}
	vis[0][0] = true
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.r == n-1 && cur.c == n-1 {
			return true
		}
		for _, d := range dirs {
			nr := cur.r + d[0]
			nc := cur.c + d[1]
			if nr >= 0 && nr < n && nc >= 0 && nc < n && !vis[nr][nc] {
				ch := grid[nr][nc]
				if nr == 0 && nc == 0 || nr == n-1 && nc == n-1 || ch == digit {
					vis[nr][nc] = true
					q = append(q, pair{nr, nc})
				}
			}
		}
	}
	return false
}

func applyOps(grid [][]byte, ops [][2]int) error {
	n := len(grid)
	seen := make(map[[2]int]bool)
	for _, op := range ops {
		r, c := op[0], op[1]
		if r <= 0 || r > n || c <= 0 || c > n {
			return fmt.Errorf("invalid coordinates")
		}
		if (r == 1 && c == 1) || (r == n && c == n) {
			return fmt.Errorf("cannot modify start or end")
		}
		key := [2]int{r, c}
		if seen[key] {
			return fmt.Errorf("duplicate cell")
		}
		seen[key] = true
		r--
		c--
		if grid[r][c] == '0' {
			grid[r][c] = '1'
		} else if grid[r][c] == '1' {
			grid[r][c] = '0'
		} else {
			return fmt.Errorf("invalid cell value")
		}
	}
	return nil
}

func runCase(bin string, tc testCase) error {
	var sb strings.Builder
	n := len(tc.grid)
	sb.WriteString("1\n")
	sb.WriteString(fmt.Sprintf("%d\n", n))
	for _, row := range tc.grid {
		sb.WriteString(row + "\n")
	}
	input := sb.String()
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v\n%s", err, out.String())
	}
	fields := strings.Fields(strings.TrimSpace(out.String()))
	if len(fields) == 0 {
		return fmt.Errorf("no output")
	}
	c, err := strconv.Atoi(fields[0])
	if err != nil {
		return fmt.Errorf("invalid first number: %v", err)
	}
	if c < 0 || c > 2 {
		return fmt.Errorf("operations count out of range")
	}
	if len(fields) != 1+2*c {
		return fmt.Errorf("expected %d coordinates got %d", 2*c, len(fields)-1)
	}
	ops := make([][2]int, c)
	idx := 1
	for i := 0; i < c; i++ {
		r, err1 := strconv.Atoi(fields[idx])
		s, err2 := strconv.Atoi(fields[idx+1])
		if err1 != nil || err2 != nil {
			return fmt.Errorf("bad coordinates")
		}
		ops[i] = [2]int{r, s}
		idx += 2
	}
	grid := make([][]byte, n)
	for i := 0; i < n; i++ {
		grid[i] = []byte(tc.grid[i])
	}
	if err := applyOps(grid, ops); err != nil {
		return err
	}
	if reachable(grid, '0') || reachable(grid, '1') {
		return fmt.Errorf("path still exists")
	}
	return nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierB.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	var cases []testCase
	// simple deterministic case
	cases = append(cases, buildCase([]string{"S0", "0F"}))
	for i := 0; i < 100; i++ {
		cases = append(cases, generateCase(rng))
	}

	for idx, tc := range cases {
		if err := runCase(bin, tc); err != nil {
			var sb strings.Builder
			n := len(tc.grid)
			sb.WriteString("1\n")
			sb.WriteString(fmt.Sprintf("%d\n", n))
			for _, row := range tc.grid {
				sb.WriteString(row + "\n")
			}
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", idx+1, err, sb.String())
			os.Exit(1)
		}
	}
	fmt.Println("All tests passed")
}
