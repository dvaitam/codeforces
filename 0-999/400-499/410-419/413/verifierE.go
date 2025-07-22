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

type testCase struct {
	n, m    int
	grid    [2]string
	queries [][2]int
}

func bfs(grid [2]string, n int, s, t int) int {
	sr := (s - 1) / n
	sc := (s - 1) % n
	tr := (t - 1) / n
	tc := (t - 1) % n
	if grid[sr][sc] == 'X' || grid[tr][tc] == 'X' {
		return -1
	}
	type node struct{ r, c int }
	dist := make([][]int, 2)
	for i := 0; i < 2; i++ {
		dist[i] = make([]int, n)
		for j := 0; j < n; j++ {
			dist[i][j] = -1
		}
	}
	q := []node{{sr, sc}}
	dist[sr][sc] = 0
	dirs := [][2]int{{1, 0}, {-1, 0}, {0, 1}, {0, -1}}
	for len(q) > 0 {
		cur := q[0]
		q = q[1:]
		if cur.r == tr && cur.c == tc {
			return dist[cur.r][cur.c]
		}
		for _, d := range dirs {
			nr, nc := cur.r+d[0], cur.c+d[1]
			if nr < 0 || nr >= 2 || nc < 0 || nc >= n {
				continue
			}
			if grid[nr][nc] == 'X' {
				continue
			}
			if dist[nr][nc] == -1 {
				dist[nr][nc] = dist[cur.r][cur.c] + 1
				q = append(q, node{nr, nc})
			}
		}
	}
	return -1
}

func expected(tc testCase) []int {
	res := make([]int, len(tc.queries))
	for i, q := range tc.queries {
		res[i] = bfs(tc.grid, tc.n, q[0], q[1])
	}
	return res
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("%d %d\n", tc.n, tc.m))
	sb.WriteString(tc.grid[0])
	sb.WriteByte('\n')
	sb.WriteString(tc.grid[1])
	sb.WriteByte('\n')
	for _, q := range tc.queries {
		sb.WriteString(fmt.Sprintf("%d %d\n", q[0], q[1]))
	}
	return sb.String()
}

func runCase(bin string, tc testCase) error {
	input := buildInput(tc)
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("runtime error: %v", err)
	}
	lines := strings.Split(strings.TrimSpace(out.String()), "\n")
	exp := expected(tc)
	if len(lines) != len(exp) {
		return fmt.Errorf("expected %d lines got %d", len(exp), len(lines))
	}
	for i, line := range lines {
		var v int
		if _, err := fmt.Sscan(strings.TrimSpace(line), &v); err != nil {
			return fmt.Errorf("bad output: %v", err)
		}
		if v != exp[i] {
			return fmt.Errorf("expected %v got %v", exp, lines)
		}
	}
	return nil
}

func randomCase(rng *rand.Rand) testCase {
	n := rng.Intn(8) + 1
	m := rng.Intn(5) + 1
	var g [2]string
	for r := 0; r < 2; r++ {
		b := make([]byte, n)
		for i := 0; i < n; i++ {
			if rng.Intn(4) == 0 {
				b[i] = 'X'
			} else {
				b[i] = '.'
			}
		}
		g[r] = string(b)
	}
	// ensure start and end cells open
	queries := make([][2]int, m)
	for i := 0; i < m; i++ {
		for {
			v := rng.Intn(2*n) + 1
			u := rng.Intn(2*n) + 1
			sr := (v - 1) / n
			sc := (v - 1) % n
			tr := (u - 1) / n
			tc := (u - 1) % n
			if g[sr][sc] == '.' && g[tr][tc] == '.' {
				queries[i] = [2]int{v, u}
				break
			}
		}
	}
	return testCase{n: n, m: m, grid: g, queries: queries}
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: go run verifierE.go /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var tests []testCase
	// trivial deterministic
	tests = append(tests, testCase{n: 1, m: 1, grid: [2]string{".", "."}, queries: [][2]int{{1, 2}}})
	for len(tests) < 100 {
		tests = append(tests, randomCase(rng))
	}
	for i, tc := range tests {
		if err := runCase(bin, tc); err != nil {
			fmt.Fprintf(os.Stderr, "case %d failed: %v\ninput:\n%s", i+1, err, buildInput(tc))
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
