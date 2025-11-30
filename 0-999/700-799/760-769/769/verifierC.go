package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var rawTestcases = []string{
	"100",
	"1 1 1",
	"X",
	"",
	"3 3 3",
	"...",
	".**",
	"X*.",
	"",
	"4 3 0",
	".X.",
	".*.",
	"...",
	"...",
	"",
	"2 4 6",
	".*.*",
	".*X.",
	"",
	"2 4 4",
	"..*.",
	"*.X.",
	"",
	"2 3 2",
	".**",
	".X.",
	"",
	"4 3 3",
	"..*",
	"*..",
	"**.",
	"..X",
	"",
	"2 1 8",
	"X",
	".",
	"",
	"2 1 6",
	"X",
	".",
	"",
	"2 2 0",
	"X*",
	".*",
	"",
	"1 3 4",
	"X*.",
	"",
	"1 4 0",
	"*X.*",
	"",
	"1 3 5",
	"X..",
	"",
	"1 3 6",
	"*X.",
	"",
	"3 1 0",
	".",
	"X",
	"*",
	"",
	"4 3 2",
	".*.",
	".*.",
	".X.",
	"...",
	"",
	"1 2 3",
	"X.",
	"",
	"2 3 4",
	"...",
	".X.",
	"",
	"1 2 1",
	"X.",
	"",
	"2 1 3",
	"X",
	".",
	"",
	"4 2 3",
	"..",
	"..",
	"..",
	".X",
	"",
	"4 3 0",
	"..X",
	"*.*",
	"...",
	"***",
	"",
	"1 4 7",
	"*X..",
	"",
	"1 3 4",
	"X*.",
	"",
	"3 4 7",
	"..X.",
	".*.*",
	"...*",
	"",
	"4 3 5",
	"*..",
	".**",
	"**.",
	"X..",
	"",
	"3 4 1",
	"...*",
	"***.",
	"X*.*",
	"",
	"4 2 5",
	"..",
	".*",
	"*.",
	"X.",
	"",
	"3 2 2",
	".*",
	"X*",
	"..",
	"",
	"4 1 6",
	".",
	"X",
	".",
	".",
	"",
	"4 4 5",
	"**..",
	"..*.",
	"X..*",
	"*.*.",
	"",
	"2 4 3",
	"**.X",
	"*..*",
	"",
	"2 1 6",
	".",
	"X",
	"",
	"2 4 3",
	"..X*",
	".*..",
	"",
	"3 4 4",
	"..*.",
	"**..",
	"**.X",
	"",
	"2 4 6",
	"....",
	".*X.",
	"",
	"2 3 7",
	".X*",
	"...",
	"",
	"3 3 6",
	"...",
	"*..",
	"X..",
	"",
	"4 2 8",
	"*.",
	"*.",
	"X.",
	"..",
	"",
	"4 3 2",
	"*..",
	"**.",
	"..*",
	".X*",
	"",
	"3 2 5",
	"X.",
	"*.",
	".*",
	"",
	"3 3 8",
	"*..",
	"**.",
	"**X",
	"",
	"1 3 4",
	".X*",
	"",
	"3 4 7",
	".X.*",
	"...*",
	"*...",
	"",
	"4 3 7",
	".*.",
	".X.",
	"...",
	"...",
	"",
	"1 3 5",
	"X..",
	"",
	"4 3 2",
	"...",
	"**.",
	".*X",
	".*.",
	"",
	"1 4 6",
	".X..",
	"",
	"1 4 2",
	".X.*",
	"",
	"1 3 1",
	"X..",
	"",
	"1 3 7",
	".X.",
	"",
	"1 2 8",
	"X.",
	"",
	"1 3 8",
	"X..",
	"",
	"1 2 1",
	"X*",
	"",
	"2 4 5",
	".*.*",
	".X..",
	"",
	"2 4 2",
	"**..",
	".*X*",
	"",
	"2 4 2",
	"X.*.",
	"....",
	"",
	"3 4 8",
	"*X..",
	"....",
	".***",
	"",
	"4 3 3",
	"*.*",
	".X.",
	".*.",
	"..*",
	"",
	"3 1 3",
	".",
	"X",
	"*",
	"",
	"2 4 2",
	"**.X",
	"....",
	"",
	"4 3 4",
	"X..",
	".*.",
	"..*",
	"*.*",
	"",
	"1 1 0",
	"X",
	"",
	"1 3 8",
	".*X",
	"",
	"1 4 7",
	"..X.",
	"",
	"4 1 7",
	".",
	"*",
	"X",
	"*",
	"",
	"2 3 3",
	"**X",
	".*.",
	"",
	"2 2 8",
	".X",
	"..",
	"",
	"4 4 1",
	".*..",
	"..*.",
	"....",
	"...X",
	"",
	"3 4 1",
	"*...",
	"....",
	"..*X",
	"",
	"1 1 2",
	"X",
	"",
	"1 3 8",
	".X.",
	"",
	"3 1 5",
	"X",
	".",
	"*",
	"",
	"3 4 4",
	"..*.",
	".***",
	".X..",
	"",
	"4 4 6",
	"...*",
	"*X..",
	"*..*",
	"*.*.",
	"",
	"1 2 3",
	"X.",
	"",
	"4 1 8",
	".",
	".",
	"X",
	".",
	"",
	"4 4 2",
	".*.*",
	".**X",
	".**.",
	".**.",
	"",
	"3 2 2",
	"..",
	"..",
	".X",
	"",
	"3 2 5",
	"*.",
	"..",
	"X*",
	"",
	"4 3 4",
	"*.*",
	".X.",
	"...",
	"**.",
	"",
	"2 4 1",
	"..**",
	"..*X",
	"",
	"3 3 2",
	"X..",
	"..*",
	"*.*",
	"",
	"4 3 4",
	"...",
	"..*",
	"..X",
	".**",
	"",
	"1 4 6",
	"*..X",
	"",
	"3 4 2",
	"....",
	"..*.",
	"X*..",
	"",
	"4 4 7",
	"...*",
	"*..*",
	"X*..",
	".*..",
	"",
	"2 1 6",
	".",
	"X",
	"",
	"2 2 7",
	"X*",
	"*.",
	"",
	"1 4 3",
	".*X.",
	"",
	"1 3 0",
	".X.",
	"",
	"4 2 2",
	".*",
	".*",
	"**",
	"X.",
	"",
	"1 2 1",
	"X.",
	"",
	"3 4 3",
	"..X*",
	"*...",
	"..*.",
	"",
	"3 1 4",
	"X",
	"*",
	".",
	"",
	"2 3 2",
	"*.*",
	"*.X",
	"",
	"2 1 8",
	"X",
	"*",
	"",
	"4 4 5",
	".*.*",
	"...*",
	"***.",
	"*.X*",
	"",
	"2 2 5",
	"X.",
	"..",
	"",
	"1 3 6",
	".*X",
}

type testCase struct {
	n, m, k int
	grid    [][]byte
	input   string
}

func parseCases() []testCase {
	lines := rawTestcases
	pos := 0
	nextLine := func() string {
		if pos >= len(lines) {
			return ""
		}
		s := lines[pos]
		pos++
		return strings.TrimSpace(s)
	}

	var cases []testCase
	t, _ := strconv.Atoi(nextLine())
	for c := 0; c < t; c++ {
		line := nextLine()
		for line == "" && pos < len(lines) {
			line = nextLine()
		}
		if line == "" {
			break
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			break
		}
		n, _ := strconv.Atoi(parts[0])
		m, _ := strconv.Atoi(parts[1])
		k, _ := strconv.Atoi(parts[2])
		grid := make([][]byte, n)
		var sb strings.Builder
		fmt.Fprintf(&sb, "%d %d %d\n", n, m, k)
		for i := 0; i < n; i++ {
			row := nextLine()
			grid[i] = []byte(row)
			sb.WriteString(row)
			sb.WriteByte('\n')
		}
		cases = append(cases, testCase{n: n, m: m, k: k, grid: grid, input: sb.String()})
	}
	return cases
}

type pair struct{ x, y int }

func solve(tc testCase) string {
	n, m, k := tc.n, tc.m, tc.k
	grid := make([][]byte, n)
	var sx, sy int
	for i := 0; i < n; i++ {
		grid[i] = make([]byte, m)
		copy(grid[i], tc.grid[i])
		for j := 0; j < m; j++ {
			if grid[i][j] == 'X' {
				sx, sy = i, j
				grid[i][j] = '.'
			}
		}
	}

	if k%2 == 1 {
		return "IMPOSSIBLE"
	}

	dist := make([][]int, n)
	for i := range dist {
		dist[i] = make([]int, m)
		for j := range dist[i] {
			dist[i][j] = -1
		}
	}
	q := make([]pair, 0, n*m)
	q = append(q, pair{sx, sy})
	dist[sx][sy] = 0
	for head := 0; head < len(q); head++ {
		p := q[head]
		d := dist[p.x][p.y] + 1
		if p.x+1 < n && grid[p.x+1][p.y] != '*' && dist[p.x+1][p.y] == -1 {
			dist[p.x+1][p.y] = d
			q = append(q, pair{p.x + 1, p.y})
		}
		if p.y-1 >= 0 && grid[p.x][p.y-1] != '*' && dist[p.x][p.y-1] == -1 {
			dist[p.x][p.y-1] = d
			q = append(q, pair{p.x, p.y - 1})
		}
		if p.y+1 < m && grid[p.x][p.y+1] != '*' && dist[p.x][p.y+1] == -1 {
			dist[p.x][p.y+1] = d
			q = append(q, pair{p.x, p.y + 1})
		}
		if p.x-1 >= 0 && grid[p.x-1][p.y] != '*' && dist[p.x-1][p.y] == -1 {
			dist[p.x-1][p.y] = d
			q = append(q, pair{p.x - 1, p.y})
		}
	}

	dirs := []byte{'D', 'L', 'R', 'U'}
	dx := []int{1, 0, 0, -1}
	dy := []int{0, -1, 1, 0}
	ans := make([]byte, 0, k)
	x, y := sx, sy
	for step := 0; step < k; step++ {
		rem := k - step - 1
		moved := false
		for i := 0; i < 4; i++ {
			nx, ny := x+dx[i], y+dy[i]
			if nx < 0 || nx >= n || ny < 0 || ny >= m || grid[nx][ny] == '*' {
				continue
			}
			d := dist[nx][ny]
			if d != -1 && d <= rem && (rem-d)%2 == 0 {
				ans = append(ans, dirs[i])
				x, y = nx, ny
				moved = true
				break
			}
		}
		if !moved {
			return "IMPOSSIBLE"
		}
	}
	return string(ans)
}

func run(bin, input string) (string, error) {
	cmd := exec.Command(bin)
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

func main() {
	if len(os.Args) != 2 {
		fmt.Println("usage: verifierC <solution-binary>")
		os.Exit(1)
	}
	bin := os.Args[1]
	cases := parseCases()
	for idx, tc := range cases {
		exp := solve(tc)
		got, err := run(bin, tc.input)
		if err != nil {
			fmt.Printf("case %d failed: %v\n", idx+1, err)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != strings.TrimSpace(exp) {
			fmt.Printf("case %d failed:\nexpected: %s\ngot: %s\n", idx+1, exp, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(cases))
}
