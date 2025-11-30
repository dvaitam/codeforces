package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

const testcasesB = `100
2 1
#
S
4 4
.#.#
#.##
..#.
.S..
4 2
.#
#.
#.
S.
4 3
.S.
#.#
#.#
###
1 4
#.#S
3 1
.
S
#
3 4
#.##
....
S..#
3 3
.#.
.S#
.##
2 4
##.#
#.S.
2 2
.S
..
1 4
.#S.
2 3
..#
S.#
3 4
##..
####
.#S.
3 2
..
.#
.S
1 2
.S
4 2
#S
#.
#.
..
3 1
S
#
.
4 3
..#
S..
#.#
..#
2 4
##.S
###.
2 2
##
S.
3 1
#
S
.
1 1
S
2 2
#S
##
3 3
.S#
..#
.#.
4 2
.#
.S
..
##
3 1
.
S
.
3 1
.
#
S
1 1
S
4 2
.S
..
.#
##
3 4
.#..
S.##
####
1 1
S
1 3
##S
3 2
#.
.#
S.
3 1
#
S
.
4 3
.S#
#.#
...
...
4 1
.
.
S
.
3 3
..#
.S.
...
3 3
#.S
...
#..
3 4
..#.
####
#S.#
3 2
.#
.S
.#
2 2
#S
##
2 1
S
.
2 3
..#
##S
2 4
...#
..S#
3 3
.#.
.#.
.S#
4 2
#S
##
..
.#
1 3
S.#
2 4
.##.
#..S
4 4
#.S#
..#.
..#.
.#..
3 3
#.#
##.
.S#
2 3
S..
#..
3 3
##.
.#S
###
4 3
##.
.#.
###
S.#
2 3
##S
##.
4 2
#.
##
S.
..
3 3
#.#
##.
#S.
3 1
.
#
S
1 4
.S.#
3 3
..S
##.
.##
1 2
#S
2 4
###.
#S..
4 3
#..
.#.
##.
S..
4 3
...
#..
.S#
##.
2 1
S
.
1 4
#.S.
1 2
S#
4 4
.###
.S##
##..
....
4 3
#.#
...
.#.
.#S
3 2
#S
..
..
1 3
.S.
1 4
.#S#
3 3
...
.#S
#.#
4 2
#.
#.
..
S.
2 4
S.##
.#..
1 4
#.#S
1 3
S#.
4 1
.
#
S
#
1 3
.S.
3 3
#..
#.#
.S.
3 4
.##.
..##
##S#
3 3
...
##S
#.#
3 4
#...
.#S#
.###
4 3
#..
.#.
...
#.S
1 3
S#.
1 1
S
4 1
#
.
S
.
4 2
#.
#.
##
.S
1 4
S##.
3 4
#S..
####
####
2 4
#S.#
.###
3 1
S
.
#
3 4
####
.#S#
.#..
2 3
###
S##
3 4
.###
...S
....
1 3
.S#
2 1
S
.
2 4
.#.#
S#.#
1 2
.S
4 3
.#.
.#.
.##
#S#
2 3
.##
S.#`

func solveCase(n, m int, grid []string) string {
	NM := n * m
	visited := make([]bool, NM)
	offX := make([]int, NM)
	offY := make([]int, NM)
	start := 0
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			if grid[i][j] == 'S' {
				start = i*m + j
			}
		}
	}
	dirs := [][2]int{{-1, 0}, {1, 0}, {0, -1}, {0, 1}}
	queue := []int{start}
	visited[start] = true
	for q := 0; q < len(queue); q++ {
		v := queue[q]
		i := v / m
		j := v % m
		bx := offX[v]
		by := offY[v]
		for _, d := range dirs {
			ni := i + d[0]
			nj := j + d[1]
			nx, ny := bx, by
			if ni < 0 {
				ni += n
				nx--
			} else if ni >= n {
				ni -= n
				nx++
			}
			if nj < 0 {
				nj += m
				ny--
			} else if nj >= m {
				nj -= m
				ny++
			}
			if grid[ni][nj] == '#' {
				continue
			}
			u := ni*m + nj
			if !visited[u] {
				visited[u] = true
				offX[u] = nx
				offY[u] = ny
				queue = append(queue, u)
			} else if offX[u] != nx || offY[u] != ny {
				return "Yes"
			}
		}
	}
	return "No"
}

type testCase struct {
	n, m int
	grid []string
}

func parseTests() ([]testCase, error) {
	reader := bufio.NewReader(strings.NewReader(testcasesB))
	var t int
	if _, err := fmt.Fscan(reader, &t); err != nil {
		return nil, err
	}
	tests := make([]testCase, t)
	for i := 0; i < t; i++ {
		var n, m int
		if _, err := fmt.Fscan(reader, &n, &m); err != nil {
			return nil, err
		}
		grid := make([]string, n)
		for r := 0; r < n; r++ {
			if _, err := fmt.Fscan(reader, &grid[r]); err != nil {
				return nil, err
			}
		}
		tests[i] = testCase{n: n, m: m, grid: grid}
	}
	return tests, nil
}

func buildInput(tc testCase) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "%d %d\n", tc.n, tc.m)
	for i, row := range tc.grid {
		if i > 0 {
			sb.WriteByte('\n')
		}
		sb.WriteString(row)
	}
	sb.WriteByte('\n')
	return sb.String()
}

func runCandidate(bin, input string) (string, error) {
	cmd := exec.Command(bin)
	cmd.Stdin = strings.NewReader(input)
	out, err := cmd.CombinedOutput()
	return strings.TrimSpace(string(out)), err
}

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintln(os.Stderr, "usage: verifierB /path/to/binary")
		os.Exit(1)
	}
	bin := os.Args[1]

	tests, err := parseTests()
	if err != nil {
		fmt.Fprintln(os.Stderr, "parse tests error:", err)
		os.Exit(1)
	}
	for i, tc := range tests {
		input := buildInput(tc)
		want := solveCase(tc.n, tc.m, tc.grid)
		got, err := runCandidate(bin, input)
		if err != nil {
			fmt.Fprintf(os.Stderr, "case %d runtime error: %v\n%s\n", i+1, err, got)
			os.Exit(1)
		}
		if strings.TrimSpace(got) != want {
			fmt.Fprintf(os.Stderr, "case %d failed\ninput:\n%s\nexpected: %s\ngot: %s\n", i+1, input, want, got)
			os.Exit(1)
		}
	}
	fmt.Printf("All %d tests passed\n", len(tests))
}
